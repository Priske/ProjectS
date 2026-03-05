package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

// makeFormationSection stays small: it just opens the editor modal.
func (ps *PlayScreen) makeFormationSection(g core.Game) core.Widget {
	widgets := []core.Widget{}

	createFormationBtn := GUI.MakeButton(0, 0, 240, 50, "Create formation", func() {
		ps.resetFormationDraft()
		ps.openFormationEditorModal(g)
	})
	widgets = append(widgets, createFormationBtn)
	for i := range g.LocalPlayer().Formations {

		index := i
		f := g.LocalPlayer().Formations[index]

		name := f.Name
		if name == "" {
			name = "Unnamed"
		}
		btn := GUI.MakeButton(0, 0, 240, 44, name, func() {

			ps.drag.Active = true
			ps.drag.Source = interaction.DragFromFormation

			ps.drag.Payload = &g.LocalPlayer().Formations[index]

			ps.drag.GrabOffX = 0
			ps.drag.GrabOffY = 0
		})

		const iconZone = 40

		hover := &GUI.HoverPopup{
			Target: btn,

			OffsetX: 8,
			OffsetY: 0,

			PopupW: 3*32 + 16,
			PopupH: 5*32 + 16,

			ClampToScreen: true,

			DrawPopup: func(dst *ebiten.Image, px, py int) {
				ps.drawFormationPreview(dst, g, index, px+8, py+8, 32)
			},
		}

		widgets = append(widgets, btn, hover)
	}

	return GUI.MakeCollapsible(0, 0, 240, 50, "Formations . . . ", widgets)
}

// Opens the modal and wires the formation editor UI.
func (ps *PlayScreen) openFormationEditorModal(g core.Game) {
	ps.ensureFormationEditorState(g)

	mw, mh := 420, 360
	mx, my := (core.VirtualW-mw)/2, (core.VirtualH-mh)/2
	padding := 14
	cell := 66

	gridX := mx + padding
	gridY := my + padding

	// Left: formation grid (3x5)
	formationGrid := ps.makeFormationGrid(g, gridX, gridY, cell)
	ps.formationGrid = formationGrid
	// Right-top: category grid
	categoryGridX := gridX + (3 * cell) + padding
	categoryGridY := gridY
	categoryGrid := ps.makeUnitCategoryGrid(g, categoryGridX, categoryGridY)
	unitOptionsGridX := categoryGridX
	unitOptionsGridY := categoryGridY + 48 + padding

	unitOptionsGrid := ps.makeUnitOptionsGrid(g, unitOptionsGridX, unitOptionsGridY)

	// Buttons
	saveBtn, closeBtn := ps.makeFormationModalButtons(g, mx, my, mw, mh, padding)

	ps.modal = GUI.MakeModal(mx, my, mw, mh, []core.Widget{
		formationGrid,
		categoryGrid,
		unitOptionsGrid,
		saveBtn,
		closeBtn,
	})
}

// Ensures editor state exists and sets sensible defaults.
func (ps *PlayScreen) ensureFormationEditorState(g core.Game) {
	if ps.formationWants == nil {
		ps.formationWants = map[core.Pos]core.UnitType{}
	}

	// Pick a valid default brush from actual owned units (guaranteed drawable)
	if len(g.LocalPlayer().Units) > 0 {
		ps.formationBrushUnitType = g.LocalPlayer().Units[0].Type
	} else {
		// fallback: only if you truly have no units
		ps.formationBrushUnitType = core.Soldier
	}

	// Defaults for category UI
	ps.selectedUnitCategory = core.Attack
	ps.availableUnitTypesForCategory = unitTypesFor(ps.selectedUnitCategory)
}

// Builds the left formation editor grid.
func (ps *PlayScreen) makeFormationGrid(g core.Game, x, y, cell int) *GUI.GridField {
	grid := GUI.MakeGridField(x, y, 3, 5, cell)
	grid.ShowGrid = true

	grid.OnCellRightClick = func(cx, cy int) {

		delete(ps.formationWants, core.Pos{X: cx, Y: cy})
	}

	grid.Get = func(cx, cy int) any {
		ut, ok := ps.formationWants[core.Pos{X: cx, Y: cy}]
		if !ok {
			return nil
		}
		return ut
	}

	grid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut := payload.(core.UnitType)
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, size)
	}

	return grid
}

// Builds the category grid on the right side.
// (Next step: use selected category to populate a unit-options grid.)
func (ps *PlayScreen) makeUnitCategoryGrid(g core.Game, x, y int) *GUI.GridField {
	unitCategories := []core.UnitCategory{
		core.Attack,
		core.Defense,
		core.Support,
	}

	grid := GUI.MakeGridField(x, y, len(unitCategories), 1, 48)
	grid.ShowGrid = true

	grid.Get = func(cx, cy int) any {
		return unitCategories[cx]
	}

	grid.OnCellClick = func(cx, cy int) {
		ps.selectedUnitCategory = unitCategories[cx]
		ps.availableUnitTypesForCategory = unitTypesFor(ps.selectedUnitCategory)
	}

	return grid
}

// Builds modal Save/Close buttons.
func (ps *PlayScreen) makeFormationModalButtons(g core.Game, mx, my, mw, mh, padding int) (saveBtn, closeBtn core.Widget) {
	btnW, btnH := 90, 40
	gap := 10
	btnY := my + mh - padding - btnH

	closeBtn = GUI.MakeButton(mx+mw-padding-btnW, btnY, btnW, btnH, "Close", func() {
		if ps.modal != nil {

			ps.modal.Close()
		}
	})

	saveBtn = GUI.MakeButton(mx+mw-padding-2*btnW-gap, btnY, btnW, btnH, "Save", func() {
		// build formation from the current draft
		newFormation := core.Formation{
			Name:  "My Formation", // later: prompt for name
			W:     3,
			H:     5,
			Wants: copyFormationWants(ps.formationWants),
		}

		g.LocalPlayer().Formations = append(g.LocalPlayer().Formations, newFormation)

		ps.resetFormationDraft()
		if ps.modal != nil {
			ps.modal.Close()
		}
		ps.rebuildLeftSidebar(g)
	})

	return saveBtn, closeBtn
}

func (ps *PlayScreen) openUnitSelectionModal(g core.Game) {

	mw, mh := 320, 220
	mx := (core.VirtualW - mw) / 2
	my := (core.VirtualH - mh) / 2

	padding := 10
	cell := 64

	gridX := mx + padding
	gridY := my + padding

	unitGrid := GUI.MakeGridField(
		gridX,
		gridY,
		len(ps.availableUnitTypesForCategory),
		1,
		cell,
	)
	unitGrid.ShowGrid = true

	unitGrid.Get = func(cx, cy int) any {
		if cx < 0 || cx >= len(ps.availableUnitTypesForCategory) {
			return nil
		}
		return ps.availableUnitTypesForCategory[cx]
	}

	unitGrid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut := payload.(core.UnitType)
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, size)
	}

	// Clicking an option sets the brush and closes the popup
	unitGrid.OnCellClick = func(cx, cy int) {
		if cx < 0 || cx >= len(ps.availableUnitTypesForCategory) {
			return
		}
		ps.formationBrushUnitType = ps.availableUnitTypesForCategory[cx]
		if ps.modal != nil {
			ps.modal.Close()
		}
	}

	closeBtn := GUI.MakeButton(
		mx+mw-90,
		my+mh-40,
		80,
		30,
		"Close",
		func() {
			if ps.modal != nil {
				ps.modal.Close()
			}
		},
	)

	ps.modal = GUI.MakeModal(
		mx,
		my,
		mw,
		mh,
		[]core.Widget{
			unitGrid,
			closeBtn,
		},
	)
}

func (ps *PlayScreen) makeUnitOptionsGrid(g core.Game, x, y int) *GUI.GridField {
	const columns = 3
	const rows = 4
	const cellSize = 48

	grid := GUI.MakeGridField(x, y, columns, rows, cellSize)
	grid.ShowGrid = true

	grid.Get = func(cx, cy int) any {
		index := cy*columns + cx
		if index < 0 || index >= len(ps.availableUnitTypesForCategory) {
			return nil
		}
		return ps.availableUnitTypesForCategory[index]
	}

	grid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut, ok := payload.(core.UnitType)
		if !ok {
			return
		}
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, size)
	}

	// IMPORTANT: do NOT reference ps.modal here (it might still be nil while building)
	grid.OnBeginDrag = func(cx, cy int, payload any) {
		ut, ok := payload.(core.UnitType)
		if !ok {
			return
		}

		ps.drag.Active = true
		ps.drag.Source = interaction.DragFromFormationPalette
		ps.drag.Payload = ut

		// center sprite under cursor
		ps.drag.GrabOffX = grid.Cell / 2
		ps.drag.GrabOffY = grid.Cell / 2

		// IMPORTANT: set mouse immediately so it draws same frame

		ps.drag.MX = g.Input().MX
		ps.drag.MY = g.Input().MY
	}
	return grid
}

func (ps *PlayScreen) resetFormationDraft() {
	ps.formationWants = make(map[core.Pos]core.UnitType) // new empty map
	// optional: reset selected category / options
	ps.selectedUnitCategory = core.Attack
	ps.availableUnitTypesForCategory = nil
	ps.formationGrid = nil
}

func copyFormationWants(src map[core.Pos]core.UnitType) map[core.Pos]core.UnitType {
	dst := make(map[core.Pos]core.UnitType, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (ps *PlayScreen) rebuildLeftSidebar(g core.Game) {
	ps.widgets = nil
	ps.widgets = append(ps.widgets,
		ps.makeOptionsSidebar(g), // or whatever you call it
	)
}

func (ps *PlayScreen) drawFormationPreview(dst *ebiten.Image, g core.Game, formationIndex int, x, y, cell int) {
	f := g.LocalPlayer().Formations[formationIndex]

	// draw grid lines (optional)
	// then draw units
	for pos, ut := range f.Wants {
		px := x + pos.X*cell
		py := y + pos.Y*cell
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, cell)
	}
}

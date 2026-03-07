package screens

import (
	"fmt"

	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

func (ps *PlayScreen) makeUnitOptionsGrid(g core.Game, x, y int) *GUI.GridField {
	const columns = 3
	const rows = 4
	const cellSize = 48

	grid := GUI.MakeGridField(x, y, columns, rows, cellSize)
	grid.ShowGrid = true

	grid.Get = func(cx, cy int) any {
		index := cy*columns + cx
		if index < 0 || index >= len(ps.formation.availableUnitTypesForCategory) {
			return nil
		}
		return ps.formation.availableUnitTypesForCategory[index]
	}

	grid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut, ok := payload.(core.UnitType)
		if !ok {
			return
		}
		drawUnitImage(dst, g.Assets(), ut, px, py, size)
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
	ps.formation.formationGrid = formationGrid
	// Right-top: category grid
	categoryGridX := gridX + (3 * cell) + padding
	categoryGridY := gridY
	categoryGrid := ps.makeUnitCategoryGrid(g, categoryGridX, categoryGridY)
	unitOptionsGridX := categoryGridX
	unitOptionsGridY := categoryGridY + 48 + padding

	unitOptionsGrid := ps.makeUnitOptionsGrid(g, unitOptionsGridX, unitOptionsGridY)

	// Buttons
	saveBtn, closeBtn := ps.makeFormationModalButtons(g, mx, my, mw, mh, padding)

	ps.ui.modal = GUI.MakeModal(mx, my, mw, mh, []core.Widget{
		formationGrid,
		categoryGrid,
		unitOptionsGrid,
		saveBtn,
		closeBtn,
	})
}

// Ensures editor state exists and sets sensible defaults.
func (ps *PlayScreen) ensureFormationEditorState(g core.Game) {
	if ps.formation.formationWants == nil {
		ps.formation.formationWants = map[core.Pos]core.UnitType{}
	}

	// Pick a valid default brush from actual owned units (guaranteed drawable)
	if len(g.LocalPlayer().Units) > 0 {
		ps.formation.formationBrushUnitType = g.LocalPlayer().Units[0].Type
	} else {
		// fallback: only if you truly have no units
		ps.formation.formationBrushUnitType = core.Soldier
	}

	// Defaults for category UI
	ps.formation.selectedUnitCategory = core.Attack
	ps.formation.availableUnitTypesForCategory = unitTypesFor(ps.formation.selectedUnitCategory)
}

// Builds the left formation editor grid.
func (ps *PlayScreen) makeFormationGrid(g core.Game, x, y, cell int) *GUI.GridField {
	grid := GUI.MakeGridField(x, y, 3, 5, cell)
	grid.ShowGrid = true

	grid.OnCellRightClick = func(cx, cy int) {

		delete(ps.formation.formationWants, core.Pos{X: cx, Y: cy})
	}

	grid.Get = func(cx, cy int) any {
		ut, ok := ps.formation.formationWants[core.Pos{X: cx, Y: cy}]
		if !ok {
			return nil
		}
		return ut
	}

	grid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut := payload.(core.UnitType)
		drawUnitImage(dst, g.Assets(), ut, px, py, size)
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
		ps.formation.selectedUnitCategory = unitCategories[cx]
		ps.formation.availableUnitTypesForCategory = unitTypesFor(ps.formation.selectedUnitCategory)
	}

	return grid
}

// Builds modal Save/Close buttons.
func (ps *PlayScreen) makeFormationModalButtons(g core.Game, mx, my, mw, mh, padding int) (saveBtn, closeBtn core.Widget) {
	btnW, btnH := 90, 40
	gap := 10
	btnY := my + mh - padding - btnH

	closeBtn = GUI.MakeButton(mx+mw-padding-btnW, btnY, btnW, btnH, "Close", func() {
		if ps.ui.modal != nil {

			ps.ui.modal.Close()
		}
	})

	saveBtn = GUI.MakeButton(mx+mw-padding-2*btnW-gap, btnY, btnW, btnH, "Save", func() {
		fmt.Println("SAVE CLICKED")
		ps.openFormationNamePrompt(g)
	})

	return saveBtn, closeBtn
}
func (ps *PlayScreen) openFormationNamePrompt(g core.Game) {
	px, py := 340, 180
	pw, ph := 320, 140
	maxlen := 50
	placeholder := "Formation name"
	fmt.Println("OPEN NAME PROMPT")
	nameField := GUI.MakeTextField(px+20, py+20, 280, 36, maxlen, placeholder)

	confirmBtn := GUI.MakeButton(px+20, py+80, 120, 36, "Confirm", func() {
		name := nameField.Text
		if name == "" {
			name = "Unnamed"
		}

		newFormation := core.Formation{
			Name:  name,
			W:     3,
			H:     5,
			Wants: copyFormationWants(ps.formation.formationWants),
		}

		g.LocalPlayer().Formations = append(g.LocalPlayer().Formations, newFormation)

		ps.ui.overlay = nil
		ps.resetFormationDraft()
		if ps.ui.modal != nil {
			ps.ui.modal.Close()
		}
		// TEMPORARY until formation list becomes truly dynamic:
		ps.rebuildLeftSidebar(g)
	})

	cancelBtn := GUI.MakeButton(px+160, py+80, 120, 36, "Cancel", func() {
		ps.ui.overlay = nil
	})

	ps.ui.overlay = GUI.MakeModal(px, py, pw, ph, []core.Widget{
		nameField,
		confirmBtn,
		cancelBtn,
	})
}

func (ps *PlayScreen) resetFormationDraft() {
	ps.formation.formationWants = make(map[core.Pos]core.UnitType) // new empty map
	// optional: reset selected category / options
	ps.formation.selectedUnitCategory = core.Attack
	ps.formation.availableUnitTypesForCategory = nil
	ps.formation.formationGrid = nil
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
		len(ps.formation.availableUnitTypesForCategory),
		1,
		cell,
	)
	unitGrid.ShowGrid = true

	unitGrid.Get = func(cx, cy int) any {
		if cx < 0 || cx >= len(ps.formation.availableUnitTypesForCategory) {
			return nil
		}
		return ps.formation.availableUnitTypesForCategory[cx]
	}

	unitGrid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut := payload.(core.UnitType)
		drawUnitImage(dst, g.Assets(), ut, px, py, size)
	}

	// Clicking an option sets the brush and closes the popup
	unitGrid.OnCellClick = func(cx, cy int) {
		if cx < 0 || cx >= len(ps.formation.availableUnitTypesForCategory) {
			return
		}
		ps.formation.formationBrushUnitType = ps.formation.availableUnitTypesForCategory[cx]
		if ps.ui.modal != nil {
			ps.ui.modal.Close()
		}
	}

	closeBtn := GUI.MakeButton(
		mx+mw-90,
		my+mh-40,
		80,
		30,
		"Close",
		func() {
			if ps.ui.modal != nil {
				ps.ui.modal.Close()
			}
		},
	)

	ps.ui.modal = GUI.MakeModal(
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

func (ps *PlayScreen) makeCategoryBar(g core.Game, x, y int) *GUI.GridField {
	categories := []core.UnitCategory{core.Attack, core.Defense, core.Support}

	gf := GUI.MakeGridField(x, y, len(categories), 1, 48)
	gf.ShowGrid = true

	gf.Get = func(cx, cy int) any { return categories[cx] }

	gf.OnCellClick = func(cx, cy int) {
		cat := categories[cx]
		ps.openUnitPickerModal(g, cat) // opens popup of unit types for that category
	}

	gf.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		// draw icon/text for category here
	}

	return gf
}
func (ps *PlayScreen) openUnitPickerModal(g core.Game, category core.UnitCategory) {
	types := unitTypesFor(category)
	cols := 4
	rows := (len(types) + cols - 1) / cols

	pw, ph := 320, 240
	px, py := (core.VirtualW-pw)/2, (core.VirtualH-ph)/2

	picker := GUI.MakeGridField(px+16, py+16, cols, rows, 48)
	picker.ShowGrid = true

	picker.Get = func(cx, cy int) any {
		i := cy*cols + cx
		if i < 0 || i >= len(types) {
			return nil
		}
		return types[i]
	}

	picker.OnCellClick = func(cx, cy int) {
		i := cy*cols + cx
		if i < 0 || i >= len(types) {
			return
		}
		ps.formation.formationBrushUnitType = types[i] // set brush
		if ps.ui.modal != nil {
			ps.ui.modal.Close()
		}
	}

	picker.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut := payload.(core.UnitType)
		drawUnitImage(dst, g.Assets(), ut, px, py, 60)
	}

	closeBtn := GUI.MakeButton(px+pw-110, py+ph-60, 90, 40, "Close", func() {
		if ps.ui.modal != nil {
			ps.ui.modal.Close()
		}
	})

	ps.ui.modal = GUI.MakeModal(px, py, pw, ph, []core.Widget{picker, closeBtn})
}

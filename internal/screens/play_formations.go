package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

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
		ps.formationBrushUnitType = types[i] // set brush
		if ps.modal != nil {
			ps.modal.Close()
		}
	}

	picker.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut := payload.(core.UnitType)
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, 60)
	}

	closeBtn := GUI.MakeButton(px+pw-110, py+ph-60, 90, 40, "Close", func() {
		if ps.modal != nil {
			ps.modal.Close()
		}
	})

	ps.modal = GUI.MakeModal(px, py, pw, ph, []core.Widget{picker, closeBtn})
}

func (ps *PlayScreen) handleFormationDrop(g core.Game, mx, my int) (bool, string) {
	defer func() { ps.drag.Active = false }()

	if ps.drag.Source != interaction.DragFromFormationPalette {
		return false, "not formation drag"
	}

	ut, ok := ps.drag.Payload.(core.UnitType)
	if !ok {
		return false, "bad payload"
	}

	// Convert mouse to formation cell
	// You already know formation grid position: gridX/gridY/cell, 3x5.
	// BEST: store a pointer on ps when you create it:
	// ps.formationGrid = formationGrid
	gf := ps.formationGrid
	if gf == nil {
		return false, "no formation grid"
	}

	cx, cy, ok := gf.MouseToCell(mx, my) // currently unexported
	if !ok {
		return false, "drop outside formation"
	}

	ps.formationWants[core.Pos{X: cx, Y: cy}] = ut
	return true, "placed in formation"
}

func (ps *PlayScreen) formationFits(cx, cy int) bool {

	if cx < 0 || cx+3 > 3 { // formation width
		return false
	}

	if cy < 0 || cy+5 > 10 { // formation height
		return false
	}

	return true
}
func (ps *PlayScreen) deployFormation(g core.Game, f *core.Formation, cx, cy int) {

	board := g.Board()

	available := append([]*core.Unit{}, ps.unPlacedUnits...)

	for pos, ut := range f.Wants {

		unitIndex := -1

		for i, u := range available {
			if u.Type == ut {
				unitIndex = i
				break
			}
		}

		if unitIndex == -1 {
			continue // no unit available of that type
		}

		unit := available[unitIndex]
		available = append(available[:unitIndex], available[unitIndex+1:]...)

		bx := cx + pos.X
		by := cy + pos.Y

		board.Location[by][bx].Unit = unit
	}

	ps.unPlacedUnits = available
}

package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

// makeFormationSection stays small: it just opens the editor modal.
func (ps *PlayScreen) makeFormationSection(g core.Game) core.Widget {
	createFormationBtn := GUI.MakeButton(0, 0, 240, 50, "Create formation", func() {
		ps.openFormationEditorModal(g)
	})

	return GUI.MakeCollapsible(0, 0, 240, 50, "Formations . . . ", []core.Widget{createFormationBtn})
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

	// Right-top: category grid
	categoryGridX := gridX + (3 * cell) + padding
	categoryGridY := gridY
	categoryGrid := ps.makeUnitCategoryGrid(categoryGridX, categoryGridY)

	// Buttons
	saveBtn, closeBtn := ps.makeFormationModalButtons(mx, my, mw, mh, padding)

	ps.modal = GUI.MakeModal(mx, my, mw, mh, []core.Widget{
		formationGrid,
		categoryGrid,
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

	grid.OnCellClick = func(cx, cy int) {
		pos := core.Pos{X: cx, Y: cy}

		// Eraser brush clears
		if ps.formationBrushUnitType == core.UnitNone {
			delete(ps.formationWants, pos)
			return
		}

		// Toggle same type off
		if ps.formationWants[pos] == ps.formationBrushUnitType {
			delete(ps.formationWants, pos)
			return
		}

		ps.formationWants[pos] = ps.formationBrushUnitType
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
func (ps *PlayScreen) makeUnitCategoryGrid(x, y int) *GUI.GridField {
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

	// Optional: implement DrawCell to show text/icons for categories.
	// grid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) { ... }

	return grid
}

// Builds modal Save/Close buttons.
func (ps *PlayScreen) makeFormationModalButtons(mx, my, mw, mh, padding int) (saveBtn, closeBtn core.Widget) {
	btnW, btnH := 90, 40
	gap := 10
	btnY := my + mh - padding - btnH

	closeBtn = GUI.MakeButton(mx+mw-padding-btnW, btnY, btnW, btnH, "Close", func() {
		if ps.modal != nil {
			ps.modal.Close()
		}
	})

	saveBtn = GUI.MakeButton(mx+mw-padding-2*btnW-gap, btnY, btnW, btnH, "Save", func() {
		if ps.modal != nil {
			ps.modal.Close()
		}
		// TODO: save ps.formationWants -> core.Formation
	})

	return saveBtn, closeBtn
}

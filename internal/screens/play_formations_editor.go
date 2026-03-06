package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

type Editor struct {
	g  core.Game
	ps *PlayScreen

	formationGrid   *GUI.GridField
	unitOptionsGrid *GUI.GridField

	wants            map[core.Pos]core.UnitType
	selectedCategory core.UnitCategory
	availableTypes   []core.UnitType
	brushType        core.UnitType

	nameField *GUI.TextField
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

func (ps *PlayScreen) resetFormationDraft() {
	ps.formationWants = make(map[core.Pos]core.UnitType) // new empty map
	// optional: reset selected category / options
	ps.selectedUnitCategory = core.Attack
	ps.availableUnitTypesForCategory = nil
	ps.formationGrid = nil
}

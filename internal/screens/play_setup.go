package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) updateSetupWidgets(g core.Game) {

	if len(ps.setup.unPlacedUnits) == 0 && !ps.setup.readyAdded {
		ps.setup.readyWidget = ps.makeReadyButton(g)
		ps.ui.widgets = append(ps.ui.widgets, ps.setup.readyWidget)
		ps.setup.readyAdded = true
	}

	if len(ps.setup.unPlacedUnits) != 0 && ps.setup.readyAdded {
		ps.ui.widgets = removeWidget(ps.ui.widgets, ps.setup.readyWidget)
		ps.setup.readyWidget = nil
		ps.setup.readyAdded = false
	}

}
func (ps *PlayScreen) resetSetupState(g core.Game) {
	board := g.Board()
	board.ClearUnits()

	ps.setup.unPlacedUnits = append([]*core.Unit(nil), g.LocalPlayer().Units...)
	ps.setup.setupMode = true

	if ps.setup.readyAdded {
		ps.ui.widgets = removeWidget(ps.ui.widgets, ps.setup.readyWidget)
		ps.setup.readyWidget = nil
		ps.setup.readyAdded = false
	}
}

func (ps *PlayScreen) confirmSetup(g core.Game) {
	ps.setup.setupMode = false
	ps.enterBattle(g)

}

func (ps *PlayScreen) canPlaceSetupUnit(toX int) bool {
	return toX < 3
}

func (ps *PlayScreen) canMoveSetupUnit(fromX, toX int) bool {
	return fromX < 3 && toX < 3
}
func (ps *PlayScreen) validateSetupPlacement(toX int) (bool, string) {
	if !ps.canPlaceSetupUnit(toX) {
		return false, "place only in first 3 columns"
	}
	return true, ""
}

func (ps *PlayScreen) validateSetupMove(fromX, toX int) (bool, string) {
	if !ps.canMoveSetupUnit(fromX, toX) {
		return false, "can't move outside placement zone"
	}
	return true, ""
}

func (ps *PlayScreen) tryReturnUnitToReserve(g core.Game, mx, my int) (bool, string) {
	if ps.drag.Source != interaction.DragFromBoard {
		return false, ""
	}
	if !ps.mouseOverReserve(mx, my) {
		return false, ""
	}

	board := g.Board()
	src := board.TilePtr(ps.drag.FromX, ps.drag.FromY)
	if src == nil || src.Unit == nil {
		return false, "no src unit"
	}

	u := src.Unit
	src.Unit = nil
	ps.setup.unPlacedUnits = append(ps.setup.unPlacedUnits, u)
	return true, "returned"
}

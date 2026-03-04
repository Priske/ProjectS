package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeRightSidebar(g core.Game) {

	x := 999
	y := 40
	//padding :=10
	reset := GUI.MakeButton(x, y, 240, 50, "reset", func() {
		board := g.Board()
		board.ClearUnits()

		// restore reserve list back to the player's full unit list
		ps.unPlacedUnits = append([]*core.Unit(nil), g.LocalPlayer().Units...)
		ps.setupMode = true

		if ps.readyAdded {
			ps.widgets = removeWidget(ps.widgets, ps.readyWidget)
			ps.readyWidget = nil
		}
		ps.readyAdded = false
		// if you track it
	})
	ps.widgets = append(ps.widgets, reset)

}

func (ps *PlayScreen) makeReadyButton(g core.Game) core.Widget {
	x := 999
	y := 400
	ready := GUI.MakeButton(x, y, 240, 50, "Ready", func() {

		ps.setupMode = false
		//spawn enemies
	})
	return ready
}

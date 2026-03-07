package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeRightSidebar(g core.Game) core.Widget {

	x := 999
	y := 40
	reset := GUI.MakeButton(x, y, 240, 50, "reset", func() {
		ps.resetSetupState(g)
	})
	return reset

}

func (ps *PlayScreen) makeReadyButton(g core.Game) core.Widget {
	x := 999
	y := 400
	ready := GUI.MakeButton(x, y, 240, 50, "Ready", func() {
		ps.confirmSetup(g)
	})
	return ready
}

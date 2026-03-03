package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeRightSidebar(g core.Game) core.Widget {

	x := 999
	y := 40

	exit := GUI.MakeButton(x, y, 240, 50, "Ready", func() {

	})

	return exit
}

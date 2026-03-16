package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeOptionsSidebar(g core.Game, contentW int) core.Widget {
	childW := contentW - 20
	if childW < 0 {
		childW = 0
	}

	placeUnit := ps.makePlaceUnitSection(g, childW)
	formation := ps.makeFormationSection(g, childW)

	exit := GUI.MakeButton(0, 0, childW, 50, "Save & Quit", func() {
		g.SetScreen(NewMenuScreen(g))
	})

	widgetsSidebar := []core.Widget{
		placeUnit,
		formation,
		exit,
	}

	return GUI.MakeCollapsible(0, 0, contentW, 50, "Options . . .", widgetsSidebar)
}
func (ps *PlayScreen) makePlaceUnitSection(g core.Game, contentW int) core.Widget {
	widgets := []core.Widget{}
	ps.reserve.grid = ps.makeUnitsGrid(g)
	grid := ps.reserve.grid

	widgets = append(widgets, grid)
	return GUI.MakeCollapsible(0, 0, contentW, 50, "Place Unit . . .", widgets)
}

func (ps *PlayScreen) makeSetupLeftPanel(g core.Game) core.Widget {
	x := 20
	y := 40
	w := 260
	h := 640

	padding := 10
	contentW := w - padding*2
	if contentW < 0 {
		contentW = 0
	}

	options := ps.makeOptionsSidebar(g, contentW)

	return GUI.MakePanel(x, y, w, h, "Setup", []core.Widget{
		options,
	})
}

package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeOptionsSidebar(g core.Game) core.Widget {
	panelW := 260
	headerH := 44
	x := 20
	y := 40

	placeUnit := ps.makePlaceUnitSection(g)
	formation := ps.makeFormationSection(g)
	exit := GUI.MakeButton(0, 0, 240, 50, "Save & Quit", func() {
		g.SetScreen(NewMenuScreen(g))
	})

	widgetsSidebar := []core.Widget{placeUnit, formation, exit}
	return GUI.MakeCollapsible(x, y, panelW, headerH, "Options . . .", widgetsSidebar)
}

func (ps *PlayScreen) makePlaceUnitSection(g core.Game) core.Widget {
	ps.reserve.grid = ps.makeUnitsGrid(g)
	grid := ps.reserve.grid
	return GUI.MakeCollapsible(0, 0, 240, 50, "Place Unit . . .", []core.Widget{
		grid,
	})
}

func (ps *PlayScreen) rebuildLeftSidebar(g core.Game) {
	ps.ui.widgets = nil
	ps.ui.widgets = append(ps.ui.widgets,
		ps.makeOptionsSidebar(g), // or whatever you call it
	)
}

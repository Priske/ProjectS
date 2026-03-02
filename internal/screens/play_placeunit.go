package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makePlaceUnitSection(g core.Game) core.Widget {
	grid := ps.makeUnitsGrid(g) // returns core.Widget

	return GUI.MakeCollapsible(0, 0, 240, 50, "Place Unit . . .", []core.Widget{
		grid,
	})
}

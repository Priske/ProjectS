package screens

import (
	"math"

	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeUnitsGrid(g core.Game) core.Widget {
	// Pick the player you actually want here.
	// Your current code loops all players but only uses the last one.
	players := g.Players()
	if len(players) == 0 {
		// return an empty grid or a label widget if you have one
		grid := GUI.MakeGridField(0, 0, 1, 1, 48)
		grid.ShowGrid = true
		return grid
	}

	p := players[0] // TODO: current player index

	unitCount := len(p.Units)
	cols := 5
	if unitCount < cols {
		cols = unitCount
		if cols == 0 {
			cols = 1
		}
	}
	rows := int(math.Ceil(float64(unitCount) / float64(5)))
	if rows == 0 {
		rows = 1
	}

	// IMPORTANT: child widgets should start at 0,0 and let collapsible layout them
	grid := GUI.MakeGridField(0, 0, cols, rows, 48)
	grid.ShowGrid = true

	// Next step: grid.Get + grid.DrawCell, but you said you’re handling that next.
	return grid
}

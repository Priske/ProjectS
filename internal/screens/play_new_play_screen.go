package screens

import (
	"github.com/Priske/ProjectS/internal/core"
)

func NewPlayScreen(g core.Game) *PlayScreen {
	ps := &PlayScreen{}

	if localPlayer := g.LocalPlayer(); localPlayer != nil {
		ps.unPlacedUnits = make([]*core.Unit, len(localPlayer.Units))
		copy(ps.unPlacedUnits, localPlayer.Units)
	}
	ps.setupMode = true
	options := ps.makeOptionsSidebar(g)

	ps.widgets = []core.Widget{options}

	ps.makeRightSidebar(g)
	ps.formationBrushUnitType = core.Soldier
	return ps
}

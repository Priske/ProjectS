package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
)

type PlayScreen struct {
	formation FormationEditorState
	setup     SetupState
	ui        PlayUI
	reserve   ReserveState

	drag interaction.DragState
}

func (ps *PlayScreen) buildUI(g core.Game) {

	ps.ui.widgets = []core.Widget{
		ps.makeOptionsSidebar(g),
		ps.makeRightSidebar(g),
	}

}

func NewPlayScreen(g core.Game) *PlayScreen {
	ps := &PlayScreen{}

	ps.initSetupState(g)
	ps.initFormationState(g)
	ps.buildUI(g)

	return ps
}

func (ps *PlayScreen) initSetupState(g core.Game) {
	if localPlayer := g.LocalPlayer(); localPlayer != nil {
		ps.setup.unPlacedUnits = make([]*core.Unit, len(localPlayer.Units))
		copy(ps.setup.unPlacedUnits, localPlayer.Units)
	}
	ps.setup.setupMode = true
}

func (ps *PlayScreen) initFormationState(g core.Game) {
	ps.formation.formationBrushUnitType = core.Soldier
	ps.formation.formationWants = map[core.Pos]core.UnitType{}
}

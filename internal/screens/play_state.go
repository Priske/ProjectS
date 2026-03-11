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
	battle    BattleState

	drag interaction.DragState
}

func (ps *PlayScreen) buildSetupUI(g core.Game) {

	ps.ui.widgets = []core.Widget{
		ps.makeOptionsSidebar(g),
		ps.makeRightSidebarSetup(g),
	}

}
func (ps *PlayScreen) buildBattleUI(g core.Game) {
	right := ps.makeBattleRightSidebar(g)

	ps.ui.widgets = []core.Widget{
		right,
	}
}
func (ps *PlayScreen) swapAndResetUI(build func(core.Game), g core.Game) {
	ps.ui.widgets = nil
	ps.ui.modal = nil
	ps.ui.overlay = nil
	build(g)
}
func NewPlayScreen(g core.Game) *PlayScreen {
	ps := &PlayScreen{}
	ps.enterSetup(g)
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

func (ps *PlayScreen) enterSetup(g core.Game) {
	ps.initSetupState(g)
	ps.initFormationState(g)
	ps.swapAndResetUI(ps.buildSetupUI, g)
}

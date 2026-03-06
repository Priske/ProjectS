package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

type PlayScreen struct {
	widgets                []core.Widget
	drag                   interaction.DragState
	lastDrop               string
	modal                  *GUI.Modal
	unPlacedUnits          []*core.Unit
	readyAdded             bool
	setupMode              bool
	reserveGrid            *GUI.GridField
	readyWidget            core.Widget
	formationGrid          *GUI.GridField
	unitOptionsGrid        *GUI.GridField
	nameFormationTextField *GUI.TextField

	formationWants                map[core.Pos]core.UnitType
	selectedUnitCategory          core.UnitCategory
	availableUnitTypesForCategory []core.UnitType
	formationBrushUnitType        core.UnitType
}

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

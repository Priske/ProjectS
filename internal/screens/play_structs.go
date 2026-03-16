package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

type PlayUI struct {
	widgets []core.Widget
	modal   *GUI.Modal
	overlay core.Widget // generic popup/tooltip/etc

	setupLeftPanel      core.Widget
	setupActionsPanel   core.Widget
	setupInventoryPanel core.Widget

	lastDrop string
}

type SetupState struct {
	setupMode     bool
	unPlacedUnits []*core.Unit
	readyAdded    bool
	readyWidget   core.Widget

	Selected *core.Unit
}

type FormationEditorState struct {
	formationGrid          *GUI.GridField
	unitOptionsGrid        *GUI.GridField
	nameFormationTextField *GUI.TextField

	formationWants                map[core.Pos]core.UnitType
	selectedUnitCategory          core.UnitCategory
	availableUnitTypesForCategory []core.UnitType
	formationBrushUnitType        core.UnitType
}

type ReserveState struct {
	grid *GUI.GridField
}

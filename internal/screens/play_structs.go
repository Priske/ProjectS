package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

type PlayUI struct {
	widgets  []core.Widget
	modal    *GUI.Modal
	overlay  core.Widget // generic popup/tooltip/etc
	lastDrop string
}

type SetupState struct {
	setupMode     bool
	unPlacedUnits []*core.Unit
	readyAdded    bool
	readyWidget   core.Widget
}

type FormationEditorState struct {
	formationGrid   *GUI.GridField
	unitOptionsGrid *GUI.GridField

	draftWants       map[core.Pos]core.UnitType
	selectedCategory core.UnitCategory
	availableTypes   []core.UnitType
	brushType        core.UnitType

	// optional: if you keep a textfield around for reuse
	nameTextField *GUI.TextField
}

package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

type Editor struct {
	g  core.Game
	ps *PlayScreen

	formationGrid   *GUI.GridField
	unitOptionsGrid *GUI.GridField

	wants            map[core.Pos]core.UnitType
	selectedCategory core.UnitCategory
	availableTypes   []core.UnitType
	brushType        core.UnitType

	nameField *GUI.TextField
}

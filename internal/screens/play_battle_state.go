package screens

import (
	"github.com/Priske/ProjectS/internal/core"
)

type BattleState struct {
	Active bool

	Turn      TurnState
	Selected  *core.Unit
	SelectedX int
	SelectedY int

	ActionMenuOpen bool
	ActionMenuX    int
	ActionMenuY    int
	SelectedAction *core.UnitAction

	Log []string
}

type boardUnitRef struct {
	U *core.Unit
	X int
	Y int
}

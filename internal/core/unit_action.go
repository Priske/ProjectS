package core

type UnitAction struct {
	ID          string
	Name        string
	Kind        ActionKind
	Range       int
	Power       int
	UsesPerTurn int
	Description string
}

type ActionKind int

const (
	ActionMove ActionKind = iota
	ActionAttack
	ActionSkill
	ActionSupport
	ActionWait
)

func (k ActionKind) String() string {
	switch k {
	case ActionMove:
		return "Move"
	case ActionAttack:
		return "Attack"
	case ActionSkill:
		return "Skill"
	case ActionSupport:
		return "Support"
	case ActionWait:
		return "Wait"
	default:
		return "Unknown"
	}
}

func defaultMeleeActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{
		{
			ID:          "move",
			Name:        "Move",
			Kind:        ActionMove,
			Range:       moveRange,
			Power:       0,
			UsesPerTurn: 1,
			Description: "Move across the battlefield",
		},
		{
			ID:          "basic_attack",
			Name:        "Attack",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic melee attack",
		},
		{
			ID:          "wait",
			Name:        "Wait",
			Kind:        ActionWait,
			Range:       0,
			Power:       0,
			UsesPerTurn: 1,
			Description: "End this unit's actions",
		},
	}
}

func defaultSoldierActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{
		{
			ID:          "move",
			Name:        "Move",
			Kind:        ActionMove,
			Range:       moveRange,
			Power:       0,
			UsesPerTurn: 1,
			Description: "Move across the battlefield",
		},
		{
			ID:          "basic_attack",
			Name:        "Attack",
			Kind:        ActionAttack,
			Range:       2,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic ranged attack",
		},
		{
			ID:          "wait",
			Name:        "Wait",
			Kind:        ActionWait,
			Range:       0,
			Power:       0,
			UsesPerTurn: 1,
			Description: "End this unit's actions",
		},
	}

}

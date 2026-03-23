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
			ID:          "basic_attack",
			Name:        "Pistol",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic melee attack",
		},
	}
}
func defaultShieldActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{
		{
			ID:          "basic_attack",
			Name:        "Shield_Bash",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Shield_Bash",
		},
		{
			ID:          "defend",
			Name:        "defend",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "defend",
		},
	}
}
func defaultCultistMeleeActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "basic_attack",
			Name:        "Stab",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic melee attack",
		},
	}
}

func defaultRatBroodLordActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "spawn_rats",
			Name:        "Spawn",
			Kind:        ActionSupport,
			Range:       0,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Spawn rats",
		},
	}
}
func defaultSoldierActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "basic_attack",
			Name:        "Rifle",
			Kind:        ActionAttack,
			Range:       2,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic ranged attack",
		},
	}

}
func defaultMedicActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "basic_attack",
			Name:        "Pistol",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic ranged attack",
		},
		{
			ID:          "heal",
			Name:        "heal",
			Kind:        ActionSupport,
			Range:       1,
			Power:       2,
			UsesPerTurn: attackUses,
			Description: "Basic healing action",
		},
	}

}

func defaulSniperActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "basic_attack",
			Name:        "Sniper Rifle",
			Kind:        ActionAttack,
			Range:       4,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic sniper rifle attack",
		},
	}
}
func defaulMedicActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "basic_attack",
			Name:        "Pistol",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Basic melee attack",
		},
	}
}

func defaultCultistLordActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "basic_attack",
			Name:        "fire bolt",
			Kind:        ActionAttack,
			Range:       3,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Ranged fire bolt",
		},
	}
}

func defaultCultistShieldActions(moveRange int, attackUses int) []UnitAction {
	return []UnitAction{

		{
			ID:          "basic_attack",
			Name:        "Shield bash",
			Kind:        ActionAttack,
			Range:       1,
			Power:       0,
			UsesPerTurn: attackUses,
			Description: "Shield bash",
		},
	}
}

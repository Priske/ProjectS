package core

type UnitAction struct {
	Name        string
	Kind        ActionKind
	Range       int
	Power       int
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

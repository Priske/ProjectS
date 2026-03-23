package screens

import "github.com/Priske/ProjectS/internal/core"

type TurnSide int

const (
	TurnPlayer TurnSide = iota
	TurnEnemy
)

type UnitTurnState struct {
	RemainingMoveActions int
	RemainingActions     int

	ActionUses map[string]int
}

type TurnState struct {
	Side  TurnSide
	Round int
	Units map[int]UnitTurnState
}

func makeUnitTurnState(u *core.Unit) UnitTurnState {
	return UnitTurnState{
		RemainingMoveActions: u.MoveActionsPerTurn,
		RemainingActions:     u.AttackActionsPerTurn,
		ActionUses:           make(map[string]int),
	}
}

func (t *TurnState) ensureUnit(u *core.Unit) {
	if u == nil {
		return
	}

	if t.Units == nil {
		t.Units = make(map[int]UnitTurnState)
	}

	if _, ok := t.Units[u.UnitId]; !ok {
		t.Units[u.UnitId] = makeUnitTurnState(u)
	}
}

func (t *TurnState) ensureUnitState(u *core.Unit) UnitTurnState {
	if u == nil {
		return UnitTurnState{}
	}

	t.ensureUnit(u)
	return t.Units[u.UnitId]
}

func (t *TurnState) CanMove(u *core.Unit) bool {
	if u == nil {
		return false
	}
	s := t.ensureUnitState(u)
	return s.RemainingMoveActions > 0
}

func (t *TurnState) CanUseTurnAction(u *core.Unit) bool {
	if u == nil {
		return false
	}
	s := t.ensureUnitState(u)
	return s.RemainingActions > 0
}

func (t *TurnState) MarkMoved(u *core.Unit) {
	if u == nil {
		return
	}

	t.ensureUnit(u)
	s := t.Units[u.UnitId]
	if s.RemainingMoveActions > 0 {
		s.RemainingMoveActions--
	}
	t.Units[u.UnitId] = s
}

func (t *TurnState) MarkActionUsed(u *core.Unit) {
	if u == nil {
		return
	}

	t.ensureUnit(u)
	s := t.Units[u.UnitId]
	if s.RemainingActions > 0 {
		s.RemainingActions--
	}
	t.Units[u.UnitId] = s
}

func (t *TurnState) CanUseNamedAction(u *core.Unit, actionID string, maxUses int) bool {
	if u == nil {
		return false
	}

	s := t.ensureUnitState(u)
	return s.ActionUses[actionID] < maxUses
}

func (t *TurnState) MarkNamedActionUsed(u *core.Unit, actionID string) {
	if u == nil {
		return
	}

	t.ensureUnit(u)
	s := t.Units[u.UnitId]
	s.ActionUses[actionID]++
	t.Units[u.UnitId] = s
}

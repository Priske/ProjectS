package screens

import (
	"fmt"

	"github.com/Priske/ProjectS/internal/core"
)

func manhattan(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	if dx < 0 {
		dx = -dx
	}

	dy := y1 - y2
	if dy < 0 {
		dy = -dy
	}

	return dx + dy
}

func (ps *PlayScreen) canUseAction(u *core.Unit, a *core.UnitAction) bool {
	if u == nil || a == nil {
		return false
	}
	if ps.battle.Turn.Side != TurnPlayer {
		return false
	}

	switch a.Kind {
	case core.ActionMove:
		return ps.battle.Turn.CanMove(u) &&
			ps.battle.Turn.CanUseNamedAction(u, a.ID, a.UsesPerTurn)

	case core.ActionAttack:
		return ps.battle.Turn.CanAttack(u) &&
			ps.battle.Turn.CanUseNamedAction(u, a.ID, a.UsesPerTurn)

	case core.ActionSkill, core.ActionSupport, core.ActionWait:
		return ps.battle.Turn.CanUseNamedAction(u, a.ID, a.UsesPerTurn)
	}

	return false
}

func (ps *PlayScreen) tryActionMenuClick(g core.Game, mx, my int) bool {
	if !ps.battle.ActionMenuOpen || ps.battle.Selected == nil {
		return false
	}

	u := ps.battle.Selected
	x := ps.battle.ActionMenuX
	y := ps.battle.ActionMenuY
	w := 140
	rowH := 22
	h := len(u.Actions)*rowH + 8

	if mx < x || mx >= x+w || my < y || my >= y+h {
		return false
	}

	index := (my - (y + 4)) / rowH
	if index < 0 || index >= len(u.Actions) {
		return true
	}

	a := &u.Actions[index]
	if !ps.canUseAction(u, a) {
		return true
	}

	ps.battle.SelectedAction = a
	ps.battle.ActionMenuOpen = false
	ps.addBattleLog(fmt.Sprintf("Selected action: %s", a.Name))
	return true
}

func actionDamage(attacker *core.Unit, action *core.UnitAction) int {
	if attacker == nil || action == nil {
		return 0
	}

	damage := attacker.TotalAttackPower() + action.Power
	if damage < 0 {
		damage = 0
	}

	return damage
}

func firstAttackAction(u *core.Unit) *core.UnitAction {
	if u == nil {
		return nil
	}

	actions := u.AllActions()
	for i := range actions {
		if actions[i].Kind == core.ActionAttack {
			a := actions[i]
			return &a
		}
	}

	return nil
}

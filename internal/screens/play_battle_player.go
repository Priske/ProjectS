package screens

import (
	"fmt"

	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) trySelectUnit(g core.Game, mx, my int) bool {
	if !ps.battle.Active || ps.battle.Turn.Side != TurnPlayer {
		return false
	}

	cx, cy, ok := mouseToCell(g, mx, my)
	if !ok {
		return false
	}

	u := g.Board().Location[cy][cx].Unit
	if u == nil {
		return false
	}

	if u.Playerid != 1 {
		return false
	}

	ps.battle.Selected = u
	ps.battle.SelectedX = cx
	ps.battle.SelectedY = cy
	ps.battle.SelectedAction = firstAttackAction(u)

	return true
}
func (ps *PlayScreen) tryMoveUnit(g core.Game, u *core.Unit, fromX, fromY, toX, toY int) bool {
	if u == nil {
		ps.addBattleLog("Move blocked: no unit")
		return false
	}
	if u.Playerid != 1 {
		ps.addBattleLog("Move blocked: not your unit")
		return false
	}
	if ps.battle.Turn.Side != TurnPlayer {
		ps.addBattleLog("Move blocked: not player turn")
		return false
	}
	if !ps.battle.Turn.CanMove(u) {
		ps.addBattleLog("Move blocked: unit already moved")
		return false
	}

	dx := abs(toX - fromX)
	dy := abs(toY - fromY)
	dist := dx + dy

	if dist == 0 {
		ps.addBattleLog("Move blocked: same tile")
		return false
	}
	if dist > u.TotalMoveRange() {
		ps.addBattleLog("Move blocked: out of range")
		return false
	}

	board := g.Board()

	if fromY < 0 || fromY >= len(board.Location) || fromX < 0 || fromX >= len(board.Location[fromY]) {
		ps.addBattleLog("Move blocked: invalid source")
		return false
	}
	if toY < 0 || toY >= len(board.Location) || toX < 0 || toX >= len(board.Location[toY]) {
		ps.addBattleLog("Move blocked: invalid destination")
		return false
	}

	src := &board.Location[fromY][fromX]
	dst := &board.Location[toY][toX]

	if src.Unit == nil {
		ps.addBattleLog("Move blocked: source empty")
		return false
	}
	if src.Unit != u {
		ps.addBattleLog("Move blocked: dragged unit mismatch")
		return false
	}
	if dst.Unit != nil {
		ps.addBattleLog("Move blocked: destination occupied")
		return false
	}

	dst.Unit = src.Unit
	src.Unit = nil

	ps.battle.Selected = dst.Unit
	ps.battle.SelectedX = toX
	ps.battle.SelectedY = toY
	ps.battle.Turn.MarkMoved(dst.Unit)

	ps.addBattleLog(fmt.Sprintf("Unit moved to (%d,%d)", toX, toY))
	return true
}
func (ps *PlayScreen) tryAttackWithSelectedAction(g core.Game, action *core.UnitAction, targetX, targetY int) bool {
	if ps.battle.Selected == nil || action == nil {
		return false
	}
	if ps.battle.Turn.Side != TurnPlayer {
		return false
	}
	if action.Kind != core.ActionAttack {
		return false
	}

	attacker := ps.battle.Selected

	if !ps.battle.Turn.CanAttack(attacker) {
		return false
	}
	if !ps.battle.Turn.CanUseNamedAction(attacker, action.ID, action.UsesPerTurn) {
		return false
	}

	fromX := ps.battle.SelectedX
	fromY := ps.battle.SelectedY

	dx := abs(targetX - fromX)
	dy := abs(targetY - fromY)
	dist := dx + dy
	if dist == 0 || dist > action.Range {
		return false
	}

	board := g.Board()

	if fromY < 0 || fromY >= len(board.Location) || fromX < 0 || fromX >= len(board.Location[fromY]) {
		return false
	}
	if targetY < 0 || targetY >= len(board.Location) || targetX < 0 || targetX >= len(board.Location[targetY]) {
		return false
	}

	src := &board.Location[fromY][fromX]
	dst := &board.Location[targetY][targetX]

	if src.Unit == nil || dst.Unit == nil {
		return false
	}

	defender := dst.Unit
	if attacker.Playerid == defender.Playerid {
		return false
	}

	damage := actionDamage(attacker, action)

	defender.CurrentHealth -= damage
	if defender.CurrentHealth < 0 {
		defender.CurrentHealth = 0
	}

	defender.BattleStats.DamageTaken += damage
	attacker.BattleStats.DamageDealt += damage

	ps.battle.Turn.MarkAttacked(attacker)
	ps.battle.Turn.MarkNamedActionUsed(attacker, action.ID)

	ps.addBattleLog(fmt.Sprintf("%s used %s on %s", attacker.Name, action.Name, defender.Name))

	if defender.CurrentHealth <= 0 {
		dst.Unit = nil
		attacker.BattleStats.Kills++
		ps.addBattleLog("Enemy defeated")
	}

	ps.resolveBattleResult(g)
	return true
}

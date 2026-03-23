package screens

import (
	"fmt"

	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) useAction(g core.Game, user *core.Unit, action *core.UnitAction, tx, ty int) bool {
	if user == nil || action == nil {
		return false
	}

	switch action.ID {
	case "basic_attack":
		return ps.useBasicAttack(g, user, action, tx, ty)
	case "heal":
		return ps.useHeal(g, user, action, tx, ty)
	case "spawn_rats":
		return ps.useSpawnRats(g, user, action, tx, ty)
	case "defend":
		return ps.useDefend(g, user, action)
	}

	return false
}

func (ps *PlayScreen) useHeal(g core.Game, user *core.Unit, action *core.UnitAction, tx, ty int) bool {
	if user == nil || action == nil {
		return false
	}
	if action.Kind != core.ActionSupport {
		return false
	}

	board := g.Board()

	ux, uy := ps.findUnitOnBoard(g, user)
	if ux == -1 || uy == -1 {
		return false
	}

	if ty < 0 || ty >= len(board.Location) || tx < 0 || tx >= len(board.Location[ty]) {
		return false
	}

	target := board.Location[ty][tx].Unit
	if target == nil {
		return false
	}
	if target.Playerid != user.Playerid {
		return false
	}

	dist := abs(tx-ux) + abs(ty-uy)
	if dist == 0 || dist > action.Range {
		return false
	}

	if target.CurrentHealth >= target.MaxHealth {
		return false
	}

	heal := action.Power
	if heal <= 0 {
		heal = 1
	}

	target.CurrentHealth += heal
	if target.CurrentHealth > target.MaxHealth {
		target.CurrentHealth = target.MaxHealth
	}

	ps.addBattleLog(fmt.Sprintf("%s used %s on %s", user.Name, action.Name, target.Name))
	return true
}

func (ps *PlayScreen) useBasicAttack(g core.Game, user *core.Unit, action *core.UnitAction, tx, ty int) bool {
	if user == nil || action == nil {
		return false
	}
	if action.Kind != core.ActionAttack {
		return false
	}

	board := g.Board()

	ux, uy := ps.findUnitOnBoard(g, user)
	if ux == -1 || uy == -1 {
		return false
	}

	if ty < 0 || ty >= len(board.Location) || tx < 0 || tx >= len(board.Location[ty]) {
		return false
	}

	target := board.Location[ty][tx].Unit
	if target == nil {
		return false
	}
	if target.Playerid == user.Playerid {
		return false
	}

	dist := abs(tx-ux) + abs(ty-uy)
	if dist == 0 || dist > action.Range {
		return false
	}

	damage := actionDamage(user, action)

	target.CurrentHealth -= damage
	if target.CurrentHealth < 0 {
		target.CurrentHealth = 0
	}

	target.BattleStats.DamageTaken += damage
	user.BattleStats.DamageDealt += damage

	ps.addBattleLog(fmt.Sprintf("%s used %s on %s", user.Name, action.Name, target.Name))

	if target.CurrentHealth <= 0 {
		board.Location[ty][tx].Unit = nil
		user.BattleStats.Kills++
		ps.removeUnitFromRoster(target, g)
		ps.addBattleLog(target.Name + " was defeated")
	}

	ps.resolveBattleResult(g)
	return true
}

func (ps *PlayScreen) useSpawnRats(g core.Game, user *core.Unit, action *core.UnitAction, tx, ty int) bool {
	ux, uy := ps.findUnitOnBoard(g, user)
	if ux == -1 {
		return false
	}

	rat := core.MakeNewEnemyRatKnife(user.Playerid, g.NewUnitID())
	if !ps.spawnUnitNear(g, ux, uy, rat) {
		return false
	}

	ps.addBattleLog(user.Name + " spawned a rat")
	return true
}

func (ps *PlayScreen) spawnUnitNear(g core.Game, sourceX, sourceY int, u *core.Unit) bool {
	board := g.Board()

	dirs := [][2]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}

	for _, d := range dirs {
		nx := sourceX + d[0]
		ny := sourceY + d[1]

		tile := board.TilePtr(nx, ny)
		if tile == nil || tile.Unit != nil {
			continue
		}

		tile.Unit = u
		return true
	}

	return false
}

func (ps *PlayScreen) useDefend(g core.Game, user *core.Unit, action *core.UnitAction) bool {
	if user == nil {
		return false
	}

	user.Defending = true
	ps.addBattleLog(user.Name + " is defending")
	return true
}

func (ps *PlayScreen) canUseActionOnCell(g core.Game, action *core.UnitAction, targetX, targetY int) bool {
	if ps.battle.Selected == nil || action == nil {
		return false
	}

	user := ps.battle.Selected
	ux := ps.battle.SelectedX
	uy := ps.battle.SelectedY

	dist := abs(targetX-ux) + abs(targetY-uy)
	if dist == 0 || dist > action.Range {
		return false
	}

	board := g.Board()
	if targetY < 0 || targetY >= len(board.Location) || targetX < 0 || targetX >= len(board.Location[targetY]) {
		return false
	}

	target := board.Location[targetY][targetX].Unit

	switch action.ID {
	case "basic_attack":
		return target != nil && target.Playerid != user.Playerid
	case "heal":
		return target != nil && target.Playerid == user.Playerid && target.CurrentHealth < target.MaxHealth
	case "spawn_rats":
		return true
	case "defend":
		return true
	}

	return false
}

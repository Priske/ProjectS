package screens

import (
	"strconv"

	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) enemyUnitsOnBoard(g core.Game) []boardUnitRef {
	board := g.Board()
	out := []boardUnitRef{}

	for y := range board.Location {
		for x := range board.Location[y] {
			u := board.Location[y][x].Unit
			if u == nil {
				continue
			}
			if u.Playerid != -1 {
				continue
			}
			out = append(out, boardUnitRef{
				U: u,
				X: x,
				Y: y,
			})
		}
	}

	return out
}
func (ps *PlayScreen) enemyTryAttackAdjacent(g core.Game, ex, ey int) bool {
	board := g.Board()
	src := board.TilePtr(ex, ey)
	if src == nil || src.Unit == nil {
		return false
	}

	dirs := [][2]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}

	bestScore := -1
	var bestDst *core.Tile

	for _, d := range dirs {
		tx := ex + d[0]
		ty := ey + d[1]

		dst := board.TilePtr(tx, ty)
		if dst == nil || dst.Unit == nil {
			continue
		}
		if dst.Unit.Playerid == src.Unit.Playerid {
			continue
		}

		score := attackScore(src.Unit, dst.Unit)
		if score > bestScore {
			bestScore = score
			bestDst = dst
		}
	}

	if bestDst == nil {
		return false
	}

	target := bestDst.Unit
	damage := src.Unit.TotalAttackPower()

	target.CurrentHealth -= damage
	if target.CurrentHealth < 0 {
		target.CurrentHealth = 0
	}

	target.BattleStats.DamageTaken += damage
	src.Unit.BattleStats.DamageDealt += damage

	ps.addBattleLog("Enemy attacked " + target.Name)

	if target.CurrentHealth <= 0 {
		src.Unit.BattleStats.Kills++
		bestDst.Unit = nil
		ps.addBattleLog(target.Name + " was defeated")
	}

	ps.resolveBattleResult(g)
	return true
}

func (ps *PlayScreen) enemyTryMoveTowardPlayer(g core.Game, ex, ey int) bool {
	board := g.Board()
	src := board.TilePtr(ex, ey)
	if src == nil || src.Unit == nil {
		return false
	}

	players := ps.playerUnitsOnBoard(g)
	if len(players) == 0 {
		return false
	}

	target := ps.chooseEnemyTarget(players)

	dirs := [][2]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}

	bestScore := -999999
	bestX, bestY := -1, -1

	for _, d := range dirs {
		nx := ex + d[0]
		ny := ey + d[1]

		dst := board.TilePtr(nx, ny)
		if dst == nil || dst.Unit != nil {
			continue
		}

		score := -manhattan(nx, ny, target.X, target.Y)

		// Bonus if move ends next to a player
		for _, p := range players {
			if manhattan(nx, ny, p.X, p.Y) == 1 {
				score += 10
			}
		}

		if score > bestScore {
			bestScore = score
			bestX, bestY = nx, ny
		}
	}

	if bestX == -1 {
		return false
	}

	board.TilePtr(bestX, bestY).Unit = src.Unit
	src.Unit = nil
	ps.addBattleLog("Enemy moved")
	return true
}

func (ps *PlayScreen) playerUnitsOnBoard(g core.Game) []boardUnitRef {
	board := g.Board()
	out := []boardUnitRef{}

	for y := range board.Location {
		for x := range board.Location[y] {
			u := board.Location[y][x].Unit
			if u == nil {
				continue
			}
			if u.Playerid == -1 {
				continue
			}
			out = append(out, boardUnitRef{
				U: u,
				X: x,
				Y: y,
			})
		}
	}

	return out
}

func attackScore(attacker, target *core.Unit) int {
	score := 0

	if attacker == nil || target == nil {
		return score
	}

	if attacker.TotalAttackPower() >= target.CurrentHealth {
		score += 100
	}

	if isFlagUnit(target) {
		score += 1000
	}

	return score
}

func moveScore(nx, ny int, target boardUnitRef) int {
	return -manhattan(nx, ny, target.X, target.Y)
}

func (ps *PlayScreen) chooseEnemyTarget(players []boardUnitRef) boardUnitRef {
	best := players[0]
	bestScore := -999999

	for _, p := range players {
		score := 0

		if isFlagUnit(p.U) {
			score += 1000
		}

		score += 20 - p.U.CurrentHealth

		if score > bestScore {
			bestScore = score
			best = p
		}
	}

	return best
}
func (ps *PlayScreen) runEnemyTurn(g core.Game) {
	enemies := ps.enemyUnitsOnBoard(g)

	for _, e := range enemies {
		board := g.Board()
		tile := board.TilePtr(e.X, e.Y)
		if tile == nil || tile.Unit == nil || tile.Unit.Playerid != -1 {
			continue
		}

		unit := tile.Unit

		if ps.enemyTryBestAttack(g, e.X, e.Y) {
			if ps.checkBattleResult(g) != BattleOngoing {
				return
			}
			continue
		}

		if ps.enemyTryMoveTowardPlayer(g, e.X, e.Y) {
			newX, newY, ok := ps.findUnitPosition(g, unit)
			if ok {
				ps.enemyTryBestAttack(g, newX, newY)
				if ps.checkBattleResult(g) != BattleOngoing {
					return
				}
			}
		}
	}
	ps.battle.Turn.Side = TurnPlayer

	ps.battle.Turn.Round++
	ps.resetUnitsForTurn(g, 0) // player id

}

func (ps *PlayScreen) findUnitPosition(g core.Game, u *core.Unit) (int, int, bool) {
	board := g.Board()

	for y := range board.Location {
		for x := range board.Location[y] {
			if board.Location[y][x].Unit == u {
				return x, y, true
			}
		}
	}

	return 0, 0, false
}

func getUsableAttackActions(u *core.Unit) []core.UnitAction {
	out := []core.UnitAction{}

	if u == nil {
		return out
	}

	for _, a := range u.AllActions() {
		if a.Kind != core.ActionAttack {
			continue
		}
		out = append(out, a)
	}

	return out
}
func (ps *PlayScreen) enemyTryBestAttack(g core.Game, ex, ey int) bool {
	board := g.Board()
	src := board.TilePtr(ex, ey)
	if src == nil || src.Unit == nil {
		return false
	}

	attacker := src.Unit
	actions := getUsableAttackActions(attacker)
	if len(actions) == 0 {
		return false
	}

	bestScore := -1
	var bestAction *core.UnitAction
	var bestTargetTile *core.Tile

	for _, action := range actions {
		for y := range board.Location {
			for x := range board.Location[y] {
				dst := board.TilePtr(x, y)
				if dst == nil || dst.Unit == nil {
					continue
				}
				if dst.Unit.Playerid == attacker.Playerid {
					continue
				}

				dist := manhattan(ex, ey, x, y)
				if dist > action.Range {
					continue
				}

				score := enemyAttackActionScore(attacker, &action, dst.Unit)
				if score > bestScore {
					a := action
					bestScore = score
					bestAction = &a
					bestTargetTile = dst
				}
			}
		}
	}

	if bestAction == nil || bestTargetTile == nil || bestTargetTile.Unit == nil {
		return false
	}

	target := bestTargetTile.Unit
	damage := actionDamage(attacker, bestAction)

	ps.addBattleLog(target.Name + " HP before " + strconv.Itoa(target.CurrentHealth))
	ps.addBattleLog(target.Name + " DT before " + strconv.Itoa(target.BattleStats.DamageTaken))

	target.CurrentHealth -= damage
	if target.CurrentHealth < 0 {
		target.CurrentHealth = 0
	}

	target.BattleStats.DamageTaken += damage
	attacker.BattleStats.DamageDealt += damage

	ps.addBattleLog(target.Name + " HP now " + strconv.Itoa(target.CurrentHealth))
	ps.addBattleLog(target.Name + " DT now " + strconv.Itoa(target.BattleStats.DamageTaken))
	ps.addBattleLog(attacker.Name + " DD now " + strconv.Itoa(attacker.BattleStats.DamageDealt))
	ps.addBattleLog(attacker.Name + " used " + bestAction.Name + " on " + target.Name)

	if target.CurrentHealth <= 0 {
		attacker.BattleStats.Kills++
		ps.killUnit(target, bestTargetTile, g)
		ps.addBattleLog(target.Name + " was defeated")

	}

	ps.resolveBattleResult(g)
	return true
}
func (ps *PlayScreen) removeUnitFromRoster(target *core.Unit, g core.Game) {
	if target == nil {
		return
	}

	for _, p := range g.Players() {
		if p == nil || p.Playerid != target.Playerid {
			continue
		}

		for i, u := range p.Units {
			if u == target || u.UnitId == target.UnitId {
				p.Units = append(p.Units[:i], p.Units[i+1:]...)
				return
			}
		}
	}
}
func (ps *PlayScreen) killUnit(target *core.Unit, tile *core.Tile, g core.Game) {
	if target == nil {
		return
	}

	tile.Unit = nil
	ps.removeUnitFromRoster(target, g)

	if ps.battle.Selected == target {
		ps.battle.Selected = nil
		ps.battle.SelectedAction = nil
		ps.battle.ActionMenuOpen = false
	}

	if ps.setup.Selected == target {
		ps.setup.Selected = nil
	}

	ps.addBattleLog(target.Name + " was defeated")
}
func enemyAttackActionScore(attacker *core.Unit, action *core.UnitAction, target *core.Unit) int {
	score := 0

	if attacker == nil || action == nil || target == nil {
		return score
	}

	damage := actionDamage(attacker, action)

	if isFlagUnit(target) {
		score += 1000
	}

	if damage >= target.CurrentHealth {
		score += 500
	}

	score += damage * 10
	score += 20 - target.CurrentHealth

	return score
}
func (ps *PlayScreen) resetUnitsForTurn(g core.Game, playerID int) {
	board := g.Board()

	for y := range board.Location {
		for x := range board.Location[y] {
			u := board.Location[y][x].Unit
			if u == nil || u.Playerid != playerID {
				continue
			}

			ps.battle.Turn.Units[u.UnitId] = UnitTurnState{
				RemainingMoveActions:   u.TotalMoveActionsPerTurn(),
				RemainingAttackActions: u.TotalAttackActionsPerTurn(),
				ActionUses:             map[string]int{},
			}
		}
	}
}

package screens

import (
	"fmt"
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

	u := src.Unit
	moveRange := u.TotalMoveRange()
	if moveRange <= 0 {
		return false
	}

	moved := false
	curX, curY := ex, ey

	for i := 0; i < moveRange; i++ {
		if !ps.enemyTrySingleStepTowardPlayer(g, curX, curY) {
			break
		}
		moved = true

		// unit has moved, update current position
		curX, curY = ps.findUnitOnBoard(g, u)
		if curX == -1 {
			break
		}

		// stop early if now adjacent to a player
		for _, p := range ps.playerUnitsOnBoard(g) {
			if manhattan(curX, curY, p.X, p.Y) == 1 {
				return true
			}
		}
	}

	if moved {
		ps.addBattleLog("Enemy moved")
	}
	return moved
}
func (ps *PlayScreen) enemyTrySingleStepTowardPlayer(g core.Game, ex, ey int) bool {
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

		for _, p := range players {
			if manhattan(nx, ny, p.X, p.Y) == 1 {
				score += 10
			}
		}

		if score > bestScore {
			bestScore = score
			bestX = nx
			bestY = ny
		}
	}

	if bestX == -1 {
		return false
	}

	board.TilePtr(bestX, bestY).Unit = src.Unit
	src.Unit = nil
	return true
}
func (ps *PlayScreen) findUnitOnBoard(g core.Game, u *core.Unit) (int, int) {
	board := g.Board()

	for y := range board.Location {
		for x := range board.Location[y] {
			if board.Location[y][x].Unit == u {
				return x, y
			}
		}
	}

	return -1, -1
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

/*
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
*/
func (ps *PlayScreen) runEnemyTurn(g core.Game) {
	enemies := ps.enemyUnitsOnBoard(g)

	for _, e := range enemies {
		board := g.Board()
		tile := board.TilePtr(e.X, e.Y)
		if tile == nil || tile.Unit == nil || tile.Unit.Playerid != -1 {
			continue
		}

		unit := tile.Unit

		ps.runEnemyUnitTurn(g, unit)

		if ps.checkBattleResult(g) != BattleOngoing {
			return
		}
	}

	ps.battle.Turn.Side = TurnPlayer
	ps.battle.Turn.Round++
	ps.resetUnitsForTurn(g, 0)
}
func (ps *PlayScreen) runEnemyUnitTurn(g core.Game, unit *core.Unit) {
	if unit == nil {
		return
	}

	switch unit.AIKind {
	case core.AIKindBroodLord:
		ps.runBroodLordTurn(g, unit)
	default:
		ps.runDefaultEnemyTurn(g, unit)
	}
}
func (ps *PlayScreen) runDefaultEnemyTurn(g core.Game, unit *core.Unit) {
	x, y, ok := ps.findUnitPosition(g, unit)
	if !ok {
		return
	}

	if ps.enemyTryBestAttack(g, x, y) {
		return
	}

	if ps.enemyTryMoveTowardPlayer(g, x, y) {
		newX, newY, ok := ps.findUnitPosition(g, unit)
		if ok {
			ps.enemyTryBestAttack(g, newX, newY)
		}
	}
}
func (ps *PlayScreen) runBroodLordTurn(g core.Game, unit *core.Unit) {
	x, y, ok := ps.findUnitPosition(g, unit)
	if !ok {
		return
	}

	// 1. Prefer spawning
	if ps.enemyTryUseSpawnRats(g, unit, x, y) {
		return
	}

	// 2. If adjacent to player, try to back away
	if ps.isEnemyAdjacentToAnyPlayer(g, x, y) {
		if ps.enemyTryStepAwayFromPlayer(g, x, y) {
			return
		}
	}

	// 3. Optional fallback: attack if it has something usable
	if ps.enemyTryBestAttack(g, x, y) {
		return
	}

	// 4. Otherwise reposition to safer distance
	ps.enemyTryKeepDistance(g, x, y, 3)
}
func (ps *PlayScreen) enemyTryUseSpawnRats(g core.Game, user *core.Unit, ux, uy int) bool {
	fmt.Println("tried to spawn rats")
	if user == nil {
		return false
	}

	for _, action := range user.AllActions() {
		if action.ID != "spawn_rats" {
			continue
		}

		if ps.useSpawnRats(g, user, &action, ux, uy) {
			ps.addBattleLog(user.Name + " used " + action.Name)
			return true
		}
	}

	return false
}
func (ps *PlayScreen) isEnemyAdjacentToAnyPlayer(g core.Game, ex, ey int) bool {
	for _, p := range ps.playerUnitsOnBoard(g) {
		if manhattan(ex, ey, p.X, p.Y) == 1 {
			return true
		}
	}
	return false
}
func (ps *PlayScreen) enemyTryStepAwayFromPlayer(g core.Game, ex, ey int) bool {
	board := g.Board()
	src := board.TilePtr(ex, ey)
	if src == nil || src.Unit == nil {
		return false
	}

	players := ps.playerUnitsOnBoard(g)
	if len(players) == 0 {
		return false
	}

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

		score := 0
		for _, p := range players {
			score += manhattan(nx, ny, p.X, p.Y)
		}

		if score > bestScore {
			bestScore = score
			bestX = nx
			bestY = ny
		}
	}

	if bestX == -1 {
		return false
	}

	board.TilePtr(bestX, bestY).Unit = src.Unit
	src.Unit = nil
	ps.addBattleLog("Enemy repositioned")
	return true
}
func (ps *PlayScreen) enemyTryKeepDistance(g core.Game, ex, ey int, preferredDist int) bool {
	board := g.Board()
	src := board.TilePtr(ex, ey)
	if src == nil || src.Unit == nil {
		return false
	}

	players := ps.playerUnitsOnBoard(g)
	if len(players) == 0 {
		return false
	}

	u := src.Unit
	moveRange := u.TotalMoveRange()
	if moveRange <= 0 {
		return false
	}

	moved := false
	curX, curY := ex, ey

	for i := 0; i < moveRange; i++ {
		nx, ny, ok := ps.bestStepForDistance(g, curX, curY, preferredDist)
		if !ok {
			break
		}

		board.TilePtr(nx, ny).Unit = board.TilePtr(curX, curY).Unit
		board.TilePtr(curX, curY).Unit = nil
		curX, curY = nx, ny
		moved = true
	}

	if moved {
		ps.addBattleLog(u.Name + " kept its distance")
	}
	return moved
}
func (ps *PlayScreen) bestStepForDistance(g core.Game, ex, ey int, preferredDist int) (int, int, bool) {
	board := g.Board()
	players := ps.playerUnitsOnBoard(g)
	if len(players) == 0 {
		return 0, 0, false
	}

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

		closest := 999999
		for _, p := range players {
			dist := manhattan(nx, ny, p.X, p.Y)
			if dist < closest {
				closest = dist
			}
		}

		score := -abs(closest - preferredDist)

		if score > bestScore {
			bestScore = score
			bestX = nx
			bestY = ny
		}
	}

	if bestX == -1 {
		return 0, 0, false
	}

	return bestX, bestY, true
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
				RemainingMoveActions: u.TotalMoveActionsPerTurn(),
				RemainingActions:     u.TotalAttackActionsPerTurn(),
				ActionUses:           map[string]int{},
			}
		}
	}
}

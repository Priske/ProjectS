package screens

import "github.com/Priske/ProjectS/internal/core"

type BattleResult int

const (
	BattleOngoing BattleResult = iota
	BattleVictory
	BattleDefeat
)

func (ps *PlayScreen) checkBattleResult(g core.Game) BattleResult {
	board := g.Board()

	hasPlayer := false
	hasEnemy := false

	for y := range board.Location {
		for x := range board.Location[y] {
			u := board.Location[y][x].Unit
			if u == nil {
				continue
			}

			if u.Playerid == -1 {
				hasEnemy = true
			} else {
				hasPlayer = true
			}
		}
	}

	if !hasEnemy {
		return BattleVictory
	}
	if !hasPlayer {
		return BattleDefeat
	}

	return BattleOngoing
}

func (ps *PlayScreen) resolveBattleResult(g core.Game) bool {
	switch ps.checkBattleResult(g) {
	case BattleVictory:
		ps.addBattleLog("Victory!")
		ps.battle.Active = false
		return true

	case BattleDefeat:
		ps.addBattleLog("Defeat!")
		ps.battle.Active = false
		return true
	}

	return false
}

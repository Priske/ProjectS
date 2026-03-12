package screens

import "github.com/Priske/ProjectS/internal/core"

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

		dst.Unit.Health -= src.Unit.AttackPower
		ps.addBattleLog("Enemy attacked")

		if dst.Unit.Health <= 0 {
			dst.Unit = nil
			ps.addBattleLog("Player unit defeated")
		}

		return true
	}

	return false
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

	nearest := players[0]
	bestDist := manhattan(ex, ey, nearest.X, nearest.Y)

	for _, p := range players[1:] {
		d := manhattan(ex, ey, p.X, p.Y)
		if d < bestDist {
			bestDist = d
			nearest = p
		}
	}

	dx := 0
	if nearest.X > ex {
		dx = 1
	} else if nearest.X < ex {
		dx = -1
	}

	dy := 0
	if nearest.Y > ey {
		dy = 1
	} else if nearest.Y < ey {
		dy = -1
	}

	candidates := [][2]int{}
	if dx != 0 {
		candidates = append(candidates, [2]int{ex + dx, ey})
	}
	if dy != 0 {
		candidates = append(candidates, [2]int{ex, ey + dy})
	}

	for _, c := range candidates {
		tx, ty := c[0], c[1]
		dst := board.TilePtr(tx, ty)
		if dst == nil || dst.Unit != nil {
			continue
		}

		dst.Unit = src.Unit
		src.Unit = nil
		ps.addBattleLog("Enemy moved")
		return true
	}

	return false
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

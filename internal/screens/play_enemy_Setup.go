package screens

import (
	"github.com/Priske/ProjectS/internal/core"
)

type EnemyEncounter struct {
	Name          string
	MinDifficulty int
	MaxDifficulty int
	Units         []core.UnitType
	Style         EncounterStyle
}

type EncounterStyle int

const (
	StyleCluster EncounterStyle = iota
	StyleLine
	StyleWedge
	StyleBox
	StyleSwarm
)

func makeEnemyUnit(unitType core.UnitType, playerId, unitId int) *core.Unit {
	switch unitType {
	case core.Enemy_cultist_knife:
		return core.MakeNewEnemyCultistKnife(playerId, unitId)
	case core.Enemy_cultist_lord:
		return core.MakeNewEnemyCultistLord(playerId, unitId)
	case core.Enemy_rat_knife:
		return core.MakeNewEnemyRatKnife(playerId, unitId)
	case core.Enemy_rat_brood_lord:
		return core.MakeNewEnemyRatBroodLord(playerId, unitId)
	}
	return nil
}

func encounterPositions(style EncounterStyle) []core.Pos {
	switch style {
	case StyleCluster:
		return []core.Pos{
			{X: 8, Y: 3},
			{X: 9, Y: 3},
			{X: 8, Y: 4},
			{X: 9, Y: 4},
			{X: 8, Y: 5},
			{X: 9, Y: 5},
		}
	case StyleLine:
		return []core.Pos{
			{X: 8, Y: 2},
			{X: 8, Y: 3},
			{X: 8, Y: 4},
			{X: 8, Y: 5},
			{X: 8, Y: 6},
		}
	case StyleBox:
		return []core.Pos{
			{X: 8, Y: 2},
			{X: 9, Y: 2},
			{X: 8, Y: 3},
			{X: 9, Y: 3},
			{X: 8, Y: 4},
			{X: 9, Y: 4},
		}
	case StyleWedge:
		return []core.Pos{
			{X: 9, Y: 3},
			{X: 8, Y: 4},
			{X: 9, Y: 4},
			{X: 10, Y: 4},
			{X: 9, Y: 5},
		}
	}
	return nil
}

func (ps *PlayScreen) spawnEnemySetup(g core.Game, encounter EnemyEncounter) {
	positions := encounterPositions(encounter.Style)
	board := g.Board()

	for i, ut := range encounter.Units {
		if i >= len(positions) {
			break
		}

		u := makeEnemyUnit(ut, -1, g.NewUnitID())
		if u == nil {
			continue
		}

		pos := positions[i]
		board.Location[pos.Y][pos.X].Unit = u
	}
}

func (ps *PlayScreen) deployEnemiesForEncounter(g core.Game, enemies []*core.Unit) {
	board := g.Board()

	positions := []core.Pos{
		{X: 8, Y: 2},
		{X: 9, Y: 2},
		{X: 8, Y: 3},
		{X: 9, Y: 3},
		{X: 8, Y: 4},
		{X: 9, Y: 4},
		{X: 8, Y: 5},
		{X: 9, Y: 5},
	}

	for i, u := range enemies {
		if u == nil || i >= len(positions) {
			break
		}

		pos := positions[i]
		board.Location[pos.Y][pos.X].Unit = u
	}
}

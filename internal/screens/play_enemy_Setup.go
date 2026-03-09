package screens

import (
	"math/rand"

	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) spawnEnemySetup(g core.Game) {
	enemyTopRight := core.Formation{
		Name: "Enemy Top Right",
		W:    2,
		H:    2,
		Wants: map[core.Pos]core.UnitType{
			{X: 0, Y: 0}: core.Enemy_cultist_knife,
			{X: 1, Y: 0}: core.Enemy_cultist_knife,
			{X: 0, Y: 1}: core.Enemy_cultist_knife,
			{X: 1, Y: 1}: core.Enemy_cultist_lord,
		},
	}

	enemyCenter := core.Formation{
		Name: "Enemy Center",
		W:    2,
		H:    2,
		Wants: map[core.Pos]core.UnitType{
			{X: 0, Y: 0}: core.Enemy_cultist_knife,
			{X: 1, Y: 0}: core.Enemy_cultist_lord,
			{X: 0, Y: 1}: core.Enemy_cultist_knife,
			{X: 1, Y: 1}: core.Enemy_cultist_knife,
		},
	}
	enemyBottomRight := core.Formation{
		Name: "Enemy Bottom Right",
		W:    2,
		H:    2,
		Wants: map[core.Pos]core.UnitType{
			{X: 0, Y: 0}: core.Enemy_cultist_lord,
			{X: 1, Y: 0}: core.Enemy_cultist_knife,
			{X: 0, Y: 1}: core.Enemy_cultist_knife,
			{X: 1, Y: 1}: core.Enemy_cultist_knife,
		},
	}
	units := g.InitializeStartingUnitsEnemy(-1)

	formations := []struct {
		F core.Formation
		X int
		Y int
	}{
		{enemyTopRight, 8, 0},
		{enemyCenter, 8, 4},
		{enemyBottomRight, 8, 7},
	}

	f := formations[rand.Intn(len(formations))]

	ps.deployFormation(g, &f.F, f.X, f.Y, units)
}

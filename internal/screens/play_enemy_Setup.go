package screens

import (
	"math/rand"

	"github.com/Priske/ProjectS/internal/core"
)

type EnemyEncounter struct {
	Name          string
	MinDifficulty int
	MaxDifficulty int
	Formation     core.Formation
	X             int
	Y             int
}

var enemyEncounters = []EnemyEncounter{
	{
		Name: "Enemy Top Right",
		Formation: core.Formation{
			Name: "Enemy Top Right",
			W:    2,
			H:    2,
			Wants: map[core.Pos]core.UnitType{
				{X: 0, Y: 0}: core.Enemy_cultist_knife,
				{X: 1, Y: 0}: core.Enemy_cultist_knife,
				{X: 0, Y: 1}: core.Enemy_cultist_knife,
				{X: 1, Y: 1}: core.Enemy_cultist_lord,
			},
		},
		X: 8,
		Y: 0,
	},
	{
		Name: "Enemy Center",
		Formation: core.Formation{
			Name: "Enemy Center",
			W:    2,
			H:    2,
			Wants: map[core.Pos]core.UnitType{
				{X: 0, Y: 0}: core.Enemy_cultist_knife,
				{X: 1, Y: 0}: core.Enemy_cultist_lord,
				{X: 0, Y: 1}: core.Enemy_cultist_knife,
				{X: 1, Y: 1}: core.Enemy_cultist_knife,
			},
		},
		X: 8,
		Y: 4,
	},
	{
		Name: "Enemy Bottom Right",
		Formation: core.Formation{
			Name: "Enemy Bottom Right",
			W:    2,
			H:    2,
			Wants: map[core.Pos]core.UnitType{
				{X: 0, Y: 0}: core.Enemy_cultist_lord,
				{X: 1, Y: 0}: core.Enemy_cultist_knife,
				{X: 0, Y: 1}: core.Enemy_cultist_knife,
				{X: 1, Y: 1}: core.Enemy_cultist_knife,
			},
		},
		X: 8,
		Y: 7,
	},
}

func (ps *PlayScreen) spawnEnemySetup(g core.Game) {
	e := enemyEncounters[rand.Intn(len(enemyEncounters))]
	units := g.InitializeStartingUnitsEnemy(-1)
	units = ps.deployFormation(g, &e.Formation, e.X, e.Y, units)
	_ = units
}

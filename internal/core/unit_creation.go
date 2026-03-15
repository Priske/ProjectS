package core

import "math/rand"

func MakeNewSoldier(playerId, unitId int) *Unit {
	return &Unit{
		Type:                 Soldier,
		UnitId:               unitId,
		MaxHealth:            2,
		CurrentHealth:        2,
		AttackPower:          2,
		Experience:           0,
		Playerid:             playerId,
		MoveRange:            2,
		MoveActionsPerTurn:   1,
		AttackActionsPerTurn: 1,
		Actions:              defaultSoldierActions(2, 1),
		UnitCategory:         Attack,
		Name:                 GenerateUnitName(Soldier),
	}
}
func MakeNewCommander(playerId, unitId int) *Unit {
	return &Unit{
		Type:                 Commander,
		UnitId:               unitId,
		MaxHealth:            5,
		CurrentHealth:        5,
		AttackPower:          1,
		Experience:           0,
		Playerid:             playerId,
		MoveRange:            1,
		MoveActionsPerTurn:   1,
		AttackActionsPerTurn: 1,
		Actions:              defaultMeleeActions(1, 1),
		UnitCategory:         Flag,
		Name:                 GenerateUnitName(Commander),
	}
}

func MakeNewEnemyCultistKnife(playerId, unitId int) *Unit {
	return &Unit{
		Type:                 Enemy_cultist_knife,
		UnitId:               unitId,
		MaxHealth:            4,
		CurrentHealth:        4,
		AttackPower:          2,
		Experience:           0,
		Playerid:             playerId,
		MoveRange:            3,
		MoveActionsPerTurn:   1,
		AttackActionsPerTurn: 1,
		Actions:              defaultCultistMeleeActions(3, 1),
		UnitCategory:         Attack,
	}
}
func MakeNewEnemyCultistLord(playerId, unitId int) *Unit {
	return &Unit{
		Type:                 Enemy_cultist_lord,
		UnitId:               unitId,
		MaxHealth:            4,
		CurrentHealth:        4,
		AttackPower:          1,
		Experience:           0,
		Playerid:             playerId,
		MoveRange:            1,
		MoveActionsPerTurn:   1,
		AttackActionsPerTurn: 1,
		Actions:              defaultCultistLordActions(1, 1),
		UnitCategory:         Flag,
	}
}

var soldierNames = []string{
	"Hale", "Brooks", "Keller", "Vargas", "Rowan",
	"Hayes", "Carter", "Maddox", "Pierce", "Griffin",
}

var supportNames = []string{
	"Lumen", "Solace", "Mercy", "Auriel", "Nova",
	"Helia", "Seren", "Lyra", "Elion", "Vale",
}

var defenseNames = []string{
	"Ironwall", "Bulwark", "Stone", "Aegis", "Bastion",
	"Atlas", "Garrick", "Tarkus", "Duran", "Bragg",
}

var commanderNames = []string{
	"Valerius", "Drake", "Arcturus", "Magnus", "Rhea",
	"Leonis", "Cassian", "Severin", "Octavia", "Tiber",
}

func GenerateUnitName(t UnitType) string {
	switch t {

	case Soldier:
		return soldierNames[rand.Intn(len(soldierNames))]

	case Medic:
		return supportNames[rand.Intn(len(supportNames))]

	case Shield:
		return defenseNames[rand.Intn(len(defenseNames))]

	case Commander:
		return commanderNames[rand.Intn(len(commanderNames))]
	}

	return "Unknown"
}

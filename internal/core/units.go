package core

import "math/rand"

type Unit struct {
	Type        UnitType
	UnitId      int
	Health      int
	AttackPower int
	Experience  int
	Playerid    int

	MoveRange            int
	MoveActionsPerTurn   int
	AttackActionsPerTurn int

	Name         string
	BattleStats  BattleStats
	UnitCategory UnitCategory
	Actions      []UnitAction
}
type BattleStats struct {
	Kills       int
	DamageTaken int
	DamageDealt int
}

type UnitTurnState struct {
	RemainingMoveActions   int
	RemainingAttackActions int

	ActionUses map[string]int

	TempMoveRangeBonus   int
	TempAttackRangeBonus int
}

type UnitType int

const (
	UnitNone UnitType = iota
	Soldier
	Commander
	Medic
	Shield
	Sniper
	Razor
	Enemy_cultist_knife
	Enemy_cultist_lord
)

type UnitCategory int

const (
	Attack UnitCategory = iota
	Defense
	Support
	Flag
)

func MakeNewSoldier(playerId, unitId int) *Unit {
	return &Unit{
		Type:                 Soldier,
		UnitId:               unitId,
		Health:               2,
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
		Health:               5,
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
		Health:               4,
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
		Health:               4,
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

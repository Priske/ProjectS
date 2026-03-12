package core

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

	Actions []UnitAction
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
	}
}

func MakeNewEnemyCultistKnife(playerId, unitId int) *Unit {
	return &Unit{
		Type:                 Enemy_cultist_knife,
		UnitId:               unitId,
		Health:               1,
		AttackPower:          1,
		Experience:           0,
		Playerid:             playerId,
		MoveRange:            3,
		MoveActionsPerTurn:   1,
		AttackActionsPerTurn: 1,
		Actions:              defaultMeleeActions(3, 1),
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
		Actions:              defaultMeleeActions(1, 1),
	}
}

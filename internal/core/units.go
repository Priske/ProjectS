package core

type Unit struct {
	Type       UnitType
	UnitId     int
	Health     int
	Attack     int
	Experience int
	Playerid   int
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
		Type:       Soldier,
		UnitId:     unitId,
		Health:     2,
		Attack:     2,
		Experience: 0,
		Playerid:   playerId,
	}
}
func MakeNewCommander(playerId, unitId int) *Unit {
	return &Unit{
		Type:       Commander,
		UnitId:     unitId,
		Health:     5,
		Attack:     1,
		Experience: 0,
		Playerid:   playerId,
	}
}

func MakeNewEnemyCultistKnife(playerId, unitId int) *Unit {
	return &Unit{
		Type:       Enemy_cultist_knife,
		UnitId:     unitId,
		Health:     1,
		Attack:     1,
		Experience: 0,
		Playerid:   playerId,
	}
}
func MakeNewEnemyCultistLord(playerId, unitId int) *Unit {
	return &Unit{
		Type:       Enemy_cultist_lord,
		UnitId:     unitId,
		Health:     4,
		Attack:     1,
		Experience: 0,
		Playerid:   playerId,
	}
}

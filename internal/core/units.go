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
	Soldier UnitType = iota
	Commander
)

func MakeNewSoldier(playerId, unitId int) *Unit {
	return &Unit{
		Type:       Soldier,
		UnitId:     unitId,
		Health:     1,
		Attack:     1,
		Experience: 0,
		Playerid:   playerId,
	}
}
func MakeNewCommander(playerId, unitId int) *Unit {
	return &Unit{
		Type:       Commander,
		UnitId:     unitId,
		Health:     1,
		Attack:     1,
		Experience: 0,
		Playerid:   playerId,
	}
}

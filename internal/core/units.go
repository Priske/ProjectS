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

package core

type Unit struct {
	Type       UnitType
	health     int
	attack     int
	experience int
	playerid   int
}

type UnitType int

const (
	Soldier UnitType = iota
	Commander
)

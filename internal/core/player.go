package core

type Player struct {
	Playerid   int
	Name       string
	Units      []*Unit
	Formations []Formation
}

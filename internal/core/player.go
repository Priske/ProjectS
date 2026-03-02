package core

type Player struct {
	Playerid   int
	Name       string
	Units      []*Unit
	Formations []Formation
}

func NewPlayer(id int, name string) *Player {
	return &Player{
		Playerid:   id,
		Name:       name,
		Units:      make([]*Unit, 0, 40),
		Formations: []Formation{},
	}
}

package core

type GameBoard struct {
	LocationXY [][]Tile
}
type Pos struct {
	X, Y int
}

type Tile struct {
	Unit         *Unit
	LocationType LocationType
}

type Formation struct {
	Name         string
	GridW, GridH int              // e.g. 5x5 editor (optional)
	Wants        map[Pos]UnitType // only cells the player reserved
}
type Deployment struct {
	At map[*Unit]Pos // board coords, or relative coords if you anchor later
}

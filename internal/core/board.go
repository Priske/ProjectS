package core

type GameBoard struct {
	LocationXY [][]Tile
}

type Tile struct {
	Unit         *Unit
	LocationType LocationType
}

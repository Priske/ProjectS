package core

type GameBoard struct {
	Location [][]Tile
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

func (b *GameBoard) Height() int {
	if b == nil {
		return 0
	}
	return len(b.Location)
}

func (b *GameBoard) Width() int {
	if b == nil || len(b.Location) == 0 {
		return 0
	}
	return len(b.Location[0])
}

func (b *GameBoard) InBounds(x, y int) bool {
	return y >= 0 && y < b.Height() && x >= 0 && x < b.Width()
}

func (b *GameBoard) TilePtr(x, y int) *Tile {
	if b == nil || !b.InBounds(x, y) {
		return nil
	}
	// IMPORTANT: LocationXY is [y][x]
	return &b.Location[y][x]
}

func (b *GameBoard) UnitAt(x, y int) (*Unit, bool) {
	t := b.TilePtr(x, y)
	if t == nil || t.Unit == nil {
		return nil, false
	}
	return t.Unit, true
}

func (b *GameBoard) SetUnit(x, y int, u *Unit) bool {
	t := b.TilePtr(x, y)
	if t == nil {
		return false
	}
	t.Unit = u
	return true
}

func (b *GameBoard) MoveUnit(fx, fy, tx, ty int) bool {
	src := b.TilePtr(fx, fy)
	dst := b.TilePtr(tx, ty)
	if src == nil || dst == nil {
		return false
	}
	if src.Unit == nil {
		return false
	}
	dst.Unit = src.Unit
	src.Unit = nil
	return true
}

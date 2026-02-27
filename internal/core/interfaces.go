package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game interface {
	SetScreen(Screen)
	Settings() Settings
	Input() Input
	Board() GameBoard
	Assets() Assets

	//MakeBoard() //
}

type Screen interface {
	Update(g Game) error
	Draw(g Game, dst *ebiten.Image)
}
type Widget interface {
	Update(in Input)          // handle clicks/typing/focus
	Draw(dst *ebiten.Image)   // render
	Bounds() (x, y, w, h int) // for layout/debug/hover
}
type Positionable interface {
	SetPos(x, y int)
}

func PointInBounds(mx, my int, w Widget) bool {
	x, y, ww, hh := w.Bounds()
	return mx >= x && mx < x+ww && my >= y && my < y+hh
}

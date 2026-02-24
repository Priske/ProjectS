package core

import "github.com/hajimehoshi/ebiten/v2"

// Keep this small: only what screens truly need.
// Add methods only when a screen *requires* them.
type Game interface {
	SetScreen(Screen)
	JustClickedLeft() bool
}

type Screen interface {
	Update(g Game) error
	Draw(g Game, dst *ebiten.Image)
}

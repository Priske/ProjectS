package screens

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SetupScreen struct {
	// Later you'll store placement state here.
}

func NewSetupScreen() *SetupScreen {
	ss := &SetupScreen{}
	return ss
}

func (ss *SetupScreen) Update(g core.Game) error {
	// ESC to go back to menu
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.SetScreen(NewMenuScreen())
	}

	// Example: click anywhere to return (just to prove flow)
	if g.JustClickedLeft() {
		g.SetScreen(NewMenuScreen())
	}

	return nil
}

func (ss *SetupScreen) Draw(g core.Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 24, 18, 255})
	ebitenutil.DebugPrintAt(screen, "Setup Screen", 270, 120)
	ebitenutil.DebugPrintAt(screen, "Press ESC or click to return to Menu", 170, 160)
}

package screens

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SetupScreen struct {
	widgets []core.Widget
}

func NewSetupScreen() *SetupScreen {
	ss := &SetupScreen{}
	w := core.MenuButtonW // 200
	h := core.MenuButtonH // 50
	spacing := 20

	totalH := 3*h + 2*spacing
	startY := (core.VirtualH - totalH) / 2
	x := (core.VirtualW - w) / 2
	ss.widgets = []core.Widget{
		GUI.MakeTextField(x, startY, w, h, w, "hello"),
	}
	return ss
}

func (ss *SetupScreen) Update(g core.Game) error {
	input := g.Input()
	for _, w := range ss.widgets {
		w.Update(input)
	}
	if input.Escape {
		g.SetScreen(NewMenuScreen(g))
	}
	return nil
}

func (ss *SetupScreen) Draw(g core.Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{12, 24, 18, 255})
	for _, w := range ss.widgets {
		w.Draw(screen)
	}

	ebitenutil.DebugPrintAt(screen, "Setup Screen", 270, 120)
	ebitenutil.DebugPrintAt(screen, "Press ESC or click to return to Menu", 170, 160)

}

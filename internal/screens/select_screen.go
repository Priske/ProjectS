package screens

import (
	"image/color"
	"os"

	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

type SelecScreen struct {
	widgets []core.Widget
}

func NewSelectScreen(g core.Game) SelecScreen {

	ss := SelecScreen{}

	w := core.MenuButtonW // 200
	h := core.MenuButtonH // 50
	spacing := 20

	totalH := 3*h + 2*spacing
	startY := (core.VirtualH - totalH) / 2
	x := (core.VirtualW - w) / 2

	ss.widgets = []core.Widget{
		GUI.MakeButton(x, startY, w, h, "New Game", func() {
			g.InitializeNewGame(1)
			g.SetScreen(NewPlayScreen(g))

		}),
		GUI.MakeButton(x, startY+(h+spacing), w, h, "Load", func() { g.SetScreen(NewSetupScreen()) }),
		GUI.MakeButton(x, startY+2*(h+spacing), w, h, "return to main menu", func() { os.Exit(0) }),
	}

	return ss
}

func (ms SelecScreen) Update(g core.Game) error {
	input := g.Input()
	for _, w := range ms.widgets {
		w.Update(input)
	}

	return nil
}

func (ms SelecScreen) Draw(g core.Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 18, 22, 255})

	for _, b := range ms.widgets {
		b.Draw(screen)

	}
}

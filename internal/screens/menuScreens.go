package screens

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

type MenuScreen struct {
	widgets []core.Widget
}

func NewMenuScreen(g core.Game) MenuScreen {

	ms := MenuScreen{}

	w := core.MenuButtonW // 200
	h := core.MenuButtonH // 50
	spacing := 20

	totalH := 3*h + 2*spacing
	startY := (core.VirtualH - totalH) / 2
	x := (core.VirtualW - w) / 2
	ms.widgets = []core.Widget{
		GUI.MakeButton(x, startY, w, h, "Play", func() { g.SetScreen(NewPlayScreen(g)) }),
		GUI.MakeButton(x, startY+(h+spacing), w, h, "Settings", func() { g.SetScreen(NewSetupScreen()) }),
		GUI.MakeButton(x, startY+2*(h+spacing), w, h, "Exit", func() {}),
	}

	return ms
}

func (ms MenuScreen) Update(g core.Game) error {
	input := g.Input()
	for _, w := range ms.widgets {
		w.Update(input)
	}

	return nil
}

func (ms MenuScreen) Draw(g core.Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 18, 22, 255})

	for _, b := range ms.widgets {
		b.Draw(screen)

	}
}

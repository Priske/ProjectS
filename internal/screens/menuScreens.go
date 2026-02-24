package screens

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type MenuScreen struct {
	buttons []GUI.Button
}

func NewMenuScreen() MenuScreen {
	ms := MenuScreen{}

	ms.buttons = []GUI.Button{
		{
			X: 220, Y: 160, W: 200, H: 50,
			Text: "Play",
			OnClick: func() {
				// We'll set this later via g.SetScreen(...)
				// can't access g here yet, so we handle in Update with closures that capture g
			},
		},
	}

	return ms
}

func (ms MenuScreen) Update(g core.Game) error {
	if g.JustClickedLeft() {
		mx, my := ebiten.CursorPosition()
		for _, b := range ms.buttons {
			if b.Contains(mx, my) {
				// act based on which button
				switch b.Text {
				case "Play":
					g.SetScreen(NewSetupScreen()) // you create this next
				}
				break
			}
		}
	}
	return nil
}

func (ms MenuScreen) Draw(g core.Game, screen *ebiten.Image) {
	screen.Fill(color.RGBA{18, 18, 22, 255})

	ebitenutil.DebugPrintAt(screen, "ProjectS", 280, 80)

	// draw buttons
	mx, my := ebiten.CursorPosition()
	for _, b := range ms.buttons {
		hover := b.Contains(mx, my)

		col := color.RGBA{80, 80, 95, 255}
		if hover {
			col = color.RGBA{110, 110, 130, 255}
		}

		ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y), float64(b.W), float64(b.H), col)
		ebitenutil.DebugPrintAt(screen, b.Text, b.X+20, b.Y+18)
	}
}

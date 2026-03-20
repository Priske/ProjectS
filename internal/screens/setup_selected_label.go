package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SetupSelectedLabelWidget struct {
	ps *PlayScreen
	X  int
	Y  int
}

func (w *SetupSelectedLabelWidget) Update(in core.Input) {}

func (w *SetupSelectedLabelWidget) Draw(dst *ebiten.Image) {
	text := "Selected: None"
	if w.ps.setup.Selected != nil {
		text = "Selected: " + w.ps.setup.Selected.Name
	}
	ebitenutil.DebugPrintAt(dst, text, w.X, w.Y)
}

func (w *SetupSelectedLabelWidget) Bounds() (x, y, width, height int) {
	return w.X, w.Y, 140, 16
}

func (w *SetupSelectedLabelWidget) SetPos(x, y int) {
	w.X = x
	w.Y = y
}

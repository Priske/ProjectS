package guiassets

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Label struct {
	X, Y  int
	Text  string
	Color color.Color
}

func MakeLabel(x, y int, text string) *Label {
	return &Label{
		X:     x,
		Y:     y,
		Text:  text,
		Color: color.White,
	}
}

func (l *Label) Bounds() (x, y, w, h int) {
	return l.X, l.Y, 0, 0
}

func (l *Label) Update(in core.Input) {}

func (l *Label) Draw(dst *ebiten.Image) {
	text.Draw(dst, l.Text, basicfont.Face7x13, l.X, l.Y, l.Color)
}

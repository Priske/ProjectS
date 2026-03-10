package guiassets

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Panel struct {
	X, Y, W, H int

	Title          string
	ShowBorder     bool
	ShowBackground bool
	Background     color.Color
	Border         color.Color

	Children []core.Widget
	Padding  int
	Gap      int
}

func MakePanel(x, y, w, h int, title string, children []core.Widget) *Panel {
	return &Panel{
		X:              x,
		Y:              y,
		W:              w,
		H:              h,
		Title:          title,
		ShowBorder:     true,
		ShowBackground: true,
		Background:     color.RGBA{35, 35, 45, 255},
		Border:         color.RGBA{90, 90, 110, 255},
		Children:       children,
		Padding:        10,
		Gap:            8,
	}
}

func (p *Panel) Bounds() (x, y, w, h int) {
	return p.X, p.Y, p.W, p.H
}

func (p *Panel) SetPos(x, y int) {
	p.X = x
	p.Y = y
}

func (p *Panel) contentStartY() int {
	y := p.Y + p.Padding
	if p.Title != "" {
		y += 20
	}
	return y
}

func (p *Panel) layoutChildren() {
	x := p.X + p.Padding
	y := p.contentStartY()

	for _, child := range p.Children {
		if pos, ok := child.(core.Positionable); ok {
			pos.SetPos(x, y)
		}

		_, _, _, h := child.Bounds()
		y += h + p.Gap
	}
}

func (p *Panel) Update(in core.Input) {
	p.layoutChildren()

	for _, child := range p.Children {
		child.Update(in)
	}
}

func (p *Panel) Draw(dst *ebiten.Image) {
	if p.ShowBackground {
		ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y), float64(p.W), float64(p.H), p.Background)
	}

	if p.ShowBorder {
		ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y), float64(p.W), 1, p.Border)
		ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y+p.H-1), float64(p.W), 1, p.Border)
		ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y), 1, float64(p.H), p.Border)
		ebitenutil.DrawRect(dst, float64(p.X+p.W-1), float64(p.Y), 1, float64(p.H), p.Border)
	}

	if p.Title != "" {
		ebitenutil.DebugPrintAt(dst, p.Title, p.X+p.Padding, p.Y+6)
	}

	p.layoutChildren()

	for _, child := range p.Children {
		child.Draw(dst)
	}
}

package guiassets

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Button struct {
	X, Y, W, H int
	Text       string
	OnClick    func()
	hovered    bool
}

func (b Button) Contains(mx, my int) bool {
	return mx >= b.X && mx < b.X+b.W && my >= b.Y && my < b.Y+b.H
}
func MakeButton(x, y, w, h int, txt string, Onclick func()) core.Widget {
	return &Button{
		X:       x,
		Y:       y,
		H:       h,
		W:       w,
		Text:    txt,
		OnClick: Onclick,
		hovered: false,
	}

}

func (b *Button) Draw(dst *ebiten.Image) {

	ebitenutil.DrawRect(dst,
		float64(b.X), float64(b.Y),
		float64(b.W), float64(b.H),
		color.RGBA{80, 80, 90, 255},
	)

	textW := len(b.Text) * 7
	textH := 12

	textX := b.X + (b.W-textW)/2
	textY := b.Y + (b.H-textH)/2

	ebitenutil.DebugPrintAt(dst, b.Text, textX, textY)
}
func (b *Button) Update(input core.Input) {
	b.hovered = core.PointInBounds(input.MX, input.MY, b)
	if b.hovered && input.LeftClicked && b.OnClick != nil {
		b.OnClick()
	}

}

func (b *Button) Bounds() (x, y, w, h int) {
	return b.X, b.Y, b.W, b.H
}
func (b *Button) SetPos(x, y int) {
	b.X = x
	b.Y = y
}

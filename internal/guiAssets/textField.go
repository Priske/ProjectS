package guiassets

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type TextField struct {
	X, Y, W, H int

	Text        string
	Placeholder string

	Focused bool
	MaxLen  int
}

func (t *TextField) Bounds() (x, y, w, h int) {
	return t.X, t.Y, t.W, t.H
}

func MakeTextField(X, Y, W, H, maxlen int, placeholder string) core.Widget {
	return &TextField{
		X: X,
		Y: Y,
		W: W,
		H: H,

		Placeholder: placeholder,

		MaxLen: maxlen,
	}

}
func (t *TextField) Update(input core.Input) {
	if input.LeftClicked {
		t.Focused = core.PointInBounds(input.MX, input.MY, t)
	}
	if !t.Focused {
		return
	}
	if input.Backspace && len(t.Text) > 0 {
		// Remove last rune safely (not last byte)
		r := []rune(t.Text)
		t.Text = string(r[:len(r)-1])
	}
	if len(input.RuneBuffer) > 0 {
		r := []rune(t.Text)
		for _, ch := range input.RuneBuffer {
			// optional: ignore control chars like \n
			if ch == '\n' || ch == '\r' || ch == '\t' {
				continue
			}
			if t.MaxLen > 0 && len(r) >= t.MaxLen {
				break
			}
			r = append(r, ch)
		}
		t.Text = string(r)
	}
}

func (t *TextField) Draw(screen *ebiten.Image) {
	col := color.RGBA{60, 60, 70, 255}
	if t.Focused {
		col = color.RGBA{90, 90, 110, 255}
	}

	ebitenutil.DrawRect(screen, float64(t.X), float64(t.Y), float64(t.W), float64(t.H), col)

	displayText := t.Text
	if displayText == "" && !t.Focused {
		displayText = t.Placeholder
	}

	ebitenutil.DebugPrintAt(screen, displayText, t.X+10, t.Y+12)

	// Simple cursor
	if t.Focused {
		cursorX := t.X + 10 + len(t.Text)*7 // approx width per char
		ebitenutil.DrawRect(screen, float64(cursorX), float64(t.Y+8), 2, float64(t.H-16), color.White)
	}

}
func (t *TextField) SetPos(x, y int) {
	t.X = x
	t.Y = y
}

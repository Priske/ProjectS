package guiassets

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/Priske/ProjectS/internal/core"
)

type Modal struct {
	X, Y, W, H int

	Open bool

	// behavior
	CloseOnEsc     bool
	CloseOnOutside bool

	// optional callback
	OnClose func()

	Children []core.Widget
}

func MakeModal(x, y, w, h int, children []core.Widget) *Modal {
	return &Modal{
		X:              x,
		Y:              y,
		W:              w,
		H:              h,
		Open:           true,
		CloseOnEsc:     true,
		CloseOnOutside: true,
		Children:       children,
	}
}

func (m *Modal) Bounds() (x, y, w, h int) { return m.X, m.Y, m.W, m.H }

func (m *Modal) Close() {
	if !m.Open {
		return
	}
	m.Open = false
	if m.OnClose != nil {
		m.OnClose()
	}
}

func (m *Modal) Update(in core.Input) {
	if !m.Open {
		return
	}

	if m.CloseOnEsc && in.Escape {
		m.Close()
		return
	}

	// click outside closes (optional)
	if m.CloseOnOutside && in.LeftClicked {
		if !pointInRect(in.MX, in.MY, m.X, m.Y, m.W, m.H) {
			m.Close()
			return
		}
	}

	// Update children (they receive input normally)
	for _, c := range m.Children {
		c.Update(in)
	}
}

func (m *Modal) Draw(dst *ebiten.Image) {
	if !m.Open {
		return
	}

	// 1) dim backdrop (covers full screen)
	sw, sh := dst.Bounds().Dx(), dst.Bounds().Dy()
	ebitenutil.DrawRect(dst, 0, 0, float64(sw), float64(sh), color.RGBA{0, 0, 0, 160})

	// 2) panel
	ebitenutil.DrawRect(dst, float64(m.X), float64(m.Y), float64(m.W), float64(m.H), color.RGBA{30, 30, 38, 255})
	// border
	ebitenutil.DrawRect(dst, float64(m.X), float64(m.Y), float64(m.W), 2, color.RGBA{90, 90, 110, 255})
	ebitenutil.DrawRect(dst, float64(m.X), float64(m.Y+m.H-2), float64(m.W), 2, color.RGBA{90, 90, 110, 255})
	ebitenutil.DrawRect(dst, float64(m.X), float64(m.Y), 2, float64(m.H), color.RGBA{90, 90, 110, 255})
	ebitenutil.DrawRect(dst, float64(m.X+m.W-2), float64(m.Y), 2, float64(m.H), color.RGBA{90, 90, 110, 255})

	// 3) children
	for _, c := range m.Children {
		c.Draw(dst)
	}
}

func pointInRect(mx, my, x, y, w, h int) bool {
	return mx >= x && mx < x+w && my >= y && my < y+h
}

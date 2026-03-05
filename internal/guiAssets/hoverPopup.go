package guiassets

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type HoverPopup struct {
	Target core.Widget
	// The hover "hot zone"
	X, Y, W, H int

	// popup placement offset relative to mouse
	OffsetX, OffsetY int

	// Optional: clamp popup to screen
	ClampToScreen bool

	// Popup size (for clamping and background)
	PopupW, PopupH int

	// Callbacks

	OnHover      func()           // optional (fire once when entering if you want)
	OnRightClick func(mx, my int) // optional
	DrawPopup    func(dst *ebiten.Image, x, y int)

	hovering bool
}

func (hp *HoverPopup) Bounds() (x, y, w, h int) {
	return hp.X, hp.Y, hp.W, hp.H
}
func (p *HoverPopup) Update(in core.Input) {
	bx, by, bw, bh := p.Target.Bounds()

	inside := pointInRect(in.MX, in.MY, bx, by, bw, bh)
	// enter
	if inside && !p.hovering {
		p.hovering = true
		if p.OnHover != nil {
			p.OnHover()
		}
	}

	// exit
	if !inside {
		p.hovering = false
		return
	}

	// right click while hovered
	if in.RightClicked && p.OnRightClick != nil {
		p.OnRightClick(in.MX, in.MY)
	}
}

func (p *HoverPopup) Draw(dst *ebiten.Image) {
	if !p.hovering || p.DrawPopup == nil {
		return
	}

	bx, by, bw, _ := p.Target.Bounds()

	px := bx + bw + p.OffsetX
	py := by + p.OffsetY

	// optionally clamp to screen bounds
	if p.ClampToScreen {
		if px+p.PopupW > core.VirtualW {
			px = core.VirtualW - p.PopupW
		}
		if py+p.PopupH > core.VirtualH {
			py = core.VirtualH - p.PopupH
		}
		if px < 0 {
			px = 0
		}
		if py < 0 {
			py = 0
		}
	}

	// simple background so it’s readable
	bg := color.RGBA{30, 30, 36, 235}
	border := color.RGBA{120, 120, 140, 255}
	ebitenutil.DrawRect(dst, float64(px), float64(py), float64(p.PopupW), float64(p.PopupH), bg)
	ebitenutil.DrawRect(dst, float64(px), float64(py), float64(p.PopupW), 1, border)
	ebitenutil.DrawRect(dst, float64(px), float64(py+p.PopupH-1), float64(p.PopupW), 1, border)
	ebitenutil.DrawRect(dst, float64(px), float64(py), 1, float64(p.PopupH), border)
	ebitenutil.DrawRect(dst, float64(px+p.PopupW-1), float64(py), 1, float64(p.PopupH), border)

	p.DrawPopup(dst, px, py)
}

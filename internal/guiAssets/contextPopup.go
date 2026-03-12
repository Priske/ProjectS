package guiassets

import (
	"fmt"
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ContextPopup struct {
	X, Y int
	W, H int

	Open bool

	OnClick func(mx, my int)
	DrawFn  func(dst *ebiten.Image, x, y int)
}

func (p *ContextPopup) Update(in core.Input) {
	if !p.Open {
		return
	}

	if in.LeftClicked {
		if pointInRect(in.MX, in.MY, p.X, p.Y, p.W, p.H) {
			if p.OnClick != nil {
				p.OnClick(in.MX, in.MY)
			}
		} else {
			p.Open = false
		}
	}
}

func (p *ContextPopup) Draw(dst *ebiten.Image) {
	fmt.Printf("draw open=%v x=%d y=%d w=%d h=%d drawfn_nil=%v\n",
		p.Open, p.X, p.Y, p.W, p.H, p.DrawFn == nil,
	)

	if !p.Open || p.DrawFn == nil {
		return
	}

	bg := color.RGBA{30, 30, 36, 235}
	border := color.RGBA{120, 120, 140, 255}

	ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y), float64(p.W), float64(p.H), bg)

	ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y), float64(p.W), 1, border)
	ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y+p.H-1), float64(p.W), 1, border)
	ebitenutil.DrawRect(dst, float64(p.X), float64(p.Y), 1, float64(p.H), border)
	ebitenutil.DrawRect(dst, float64(p.X+p.W-1), float64(p.Y), 1, float64(p.H), border)

	p.DrawFn(dst, p.X, p.Y)
}

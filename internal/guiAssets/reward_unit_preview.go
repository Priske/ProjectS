package guiassets

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

type RewardUnitPreviewWidget struct {
	X, Y int
	W, H int

	UnitType core.UnitType
	Game     core.Game
}

func (w *RewardUnitPreviewWidget) Update(in core.Input) {}

func (w *RewardUnitPreviewWidget) Draw(dst *ebiten.Image) {
	assets := w.Game.Assets()
	img := assets.UnitImages[w.UnitType]
	if img == nil {
		return
	}

	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	if sw == 0 || sh == 0 {
		return
	}

	scaleX := float64(w.W) / float64(sw)
	scaleY := float64(w.H) / float64(sh)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	drawW := float64(sw) * scale
	drawH := float64(sh) * scale

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		float64(w.X)+(float64(w.W)-drawW)/2,
		float64(w.Y)+(float64(w.H)-drawH)/2,
	)

	dst.DrawImage(img, op)
}

func (w *RewardUnitPreviewWidget) Bounds() (x, y, width, height int) {
	return w.X, w.Y, w.W, w.H
}

func (w *RewardUnitPreviewWidget) SetPos(x, y int) {
	w.X = x
	w.Y = y
}

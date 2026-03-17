package screens

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func drawImageCenteredColored(dst *ebiten.Image, img *ebiten.Image, px, py, w, h int, r, g, b, a float64) {
	if img == nil {
		return
	}

	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	if sw == 0 || sh == 0 {
		return
	}

	scaleX := float64(w) / float64(sw)
	scaleY := float64(h) / float64(sh)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	drawW := float64(sw) * scale
	drawH := float64(sh) * scale

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		float64(px)+(float64(w)-drawW)/2,
		float64(py)+(float64(h)-drawH)/2,
	)
	op.ColorScale.Scale(float32(r), float32(g), float32(b), float32(a))

	dst.DrawImage(img, op)
}
func opaqueBounds(img *ebiten.Image, alphaThreshold uint8) (minX, minY, maxX, maxY int, ok bool) {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()

	minX, minY = w, h
	maxX, maxY = -1, -1

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if uint8(a>>8) <= alphaThreshold {
				continue
			}
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}
	}

	if maxX < minX || maxY < minY {
		return 0, 0, 0, 0, false
	}
	return minX, minY, maxX + 1, maxY + 1, true
}
func drawImageCenteredTinted(dst *ebiten.Image, img *ebiten.Image, px, py, w, h int, brightness float64) {
	if img == nil {
		return
	}

	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	if sw == 0 || sh == 0 {
		return
	}

	scaleX := float64(w) / float64(sw)
	scaleY := float64(h) / float64(sh)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	drawW := float64(sw) * scale
	drawH := float64(sh) * scale

	op := &ebiten.DrawImageOptions{}
	scale = math.Floor(scale) // or math.Round(scale)
	if scale < 1 {
		scale = 1
	}
	op.GeoM.Translate(
		float64(px)+(float64(w)-drawW)/2,
		float64(py)+(float64(h)-drawH)/2,
	)

	op.ColorScale.Scale(float32(brightness), float32(brightness), float32(brightness), 1)

	dst.DrawImage(img, op)
}

func drawImageCenteredTrimmed(dst *ebiten.Image, img *ebiten.Image, px, py, w, h int) {
	if img == nil {
		return
	}

	minX, minY, maxX, maxY, ok := opaqueBounds(img, 8)
	if !ok {
		return
	}

	sub := img.SubImage(image.Rect(minX, minY, maxX, maxY)).(*ebiten.Image)

	sw, sh := sub.Bounds().Dx(), sub.Bounds().Dy()
	if sw == 0 || sh == 0 {
		return
	}

	scaleX := float64(w) / float64(sw)
	scaleY := float64(h) / float64(sh)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	drawW := float64(sw) * scale
	drawH := float64(sh) * scale

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		float64(px)+(float64(w)-drawW)/2,
		float64(py)+(float64(h)-drawH)/2,
	)

	dst.DrawImage(sub, op)
}

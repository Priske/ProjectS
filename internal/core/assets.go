package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Assets struct {
	UnitImages map[UnitType]*ebiten.Image
}

func (a *Assets) UnitImage(t UnitType) *ebiten.Image {
	if a == nil || a.UnitImages == nil {
		return nil
	}
	return a.UnitImages[t]
}

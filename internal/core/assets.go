package core

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Assets struct {
	UnitImages      map[UnitType]*ebiten.Image
	CategoryImages  map[UnitCategory]*ebiten.Image
	LocationImages  map[LocationType]*ebiten.Image
	ShopImages      map[ItemCategory]*ebiten.Image
	SlotIcons       map[ItemCategory]*ebiten.Image
	ItemIcons       map[ItemID]*ebiten.Image
	FrameTemplate   *ebiten.Image
	ChestButtonIcon *ebiten.Image
	ShopButtonIcon  *ebiten.Image
	Crosshair       *ebiten.Image
	HealCrosshair   *ebiten.Image
	DefendCrosshair *ebiten.Image
}

func (a *Assets) UnitImage(t UnitType) *ebiten.Image {
	if a == nil || a.UnitImages == nil {
		return nil
	}
	return a.UnitImages[t]
}
func (a *Assets) CategoryImage(t UnitCategory) *ebiten.Image {
	if a == nil || a.CategoryImages == nil {
		return nil
	}
	return a.CategoryImages[t]
}

func (a *Assets) LocationImage(l LocationType) *ebiten.Image {
	if a == nil || a.LocationImages == nil {
		return nil
	}
	return a.LocationImages[l]
}

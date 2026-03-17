package assets

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.png
var imageFiles embed.FS

func LoadImage(name string) (*ebiten.Image, error) {
	img, err := loadDecodedImage(name)
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}
func loadDecodedImage(name string) (image.Image, error) {
	data, err := imageFiles.ReadFile(name)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return img, nil
}
func MustLoadAll() core.Assets {
	return core.Assets{
		UnitImages: map[core.UnitType]*ebiten.Image{
			core.Soldier:             mustImage("soldier.png"),
			core.Commander:           mustImage("commander.png"),
			core.Medic:               mustImage("medic.png"),
			core.Shield:              mustImage("shield.png"),
			core.Sniper:              mustImage("sniper.png"),
			core.Razor:               mustImage("razor.png"),
			core.Enemy_cultist_knife: mustTrimmedImage("enemy_cultist_knife.png"),
			core.Enemy_cultist_lord:  mustTrimmedImage("enemy_cultist_lord.png"),
		},
		CategoryImages: map[core.UnitCategory]*ebiten.Image{
			core.Attack:  mustTrimmedImage("attack.png"),
			core.Defense: mustTrimmedImage("defense.png"),
			core.Support: mustTrimmedImage("support.png"),
		},
		LocationImages: map[core.LocationType]*ebiten.Image{
			core.Tile_01: mustImage("tile_01.png"),
			core.Tile_02: mustImage("tile_02.png"),
			core.Tile_03: mustImage("tile_03.png"),
			core.Tile_04: mustImage("tile_04.png"),
			core.Tile_05: mustImage("tile_05.png"),
			core.Tile_06: mustImage("tile_06.png"),
			core.Tile_07: mustImage("tile_07.png"),
			core.Tile_08: mustImage("tile_08.png"),
			core.Tile_09: mustImage("tile_09.png"),
			core.Tile_10: mustImage("tile_10.png"),
			core.Tile_11: mustImage("tile_11.png"),
			core.Tile_12: mustImage("tile_12.png"),
			core.Tile_13: mustImage("tile_13.png"),
			core.Tile_14: mustImage("tile_14.png"),
			core.Tile_15: mustImage("tile_15.png"),
			core.Tile_16: mustImage("tile_16.png"),
		},
		SlotIcons: map[core.ItemCategory]*ebiten.Image{
			core.CategoryAmmo:   mustTrimmedImage("ammo_gray.png"),
			core.CategoryArmor:  mustTrimmedImage("armor_gray.png"),
			core.CategoryCharm:  mustTrimmedImage("charm_gray.png"),
			core.CategoryWeapon: mustTrimmedImage("weapon_gray.png"),
			core.CategoryPotion: mustTrimmedImage("potion_gray.png"),
		},
		ShopImages: map[core.ItemCategory]*ebiten.Image{
			core.CategoryAmmo:   mustTrimmedImage("ammo_shop.png"),
			core.CategoryArmor:  mustTrimmedImage("armor_shop.png"),
			core.CategoryCharm:  mustTrimmedImage("charm_shop.png"),
			core.CategoryWeapon: mustTrimmedImage("weapon_shop.png"),
			core.CategoryPotion: mustTrimmedImage("potion_shop.png"),
		},
		FrameTemplate:   mustTrimmedImage("frame_template.png"),
		ChestButtonIcon: mustTrimmedImage("chest_shop.png"),
		ShopButtonIcon:  mustTrimmedImage("shop_shop.png"),
	}
}
func mustTrimmedImage(name string) *ebiten.Image {
	img, err := loadDecodedImage(name)
	if err != nil {
		log.Fatalf("failed to load asset %q: %v", name, err)
	}

	trimmed := trimTransparent(img)
	return ebiten.NewImageFromImage(trimmed)
}
func trimTransparent(img image.Image) image.Image {
	minX, minY, maxX, maxY, ok := opaqueBounds(img, 8)
	if !ok {
		return img
	}
	return cropImage(img, image.Rect(minX, minY, maxX, maxY))
}

func mustImage(name string) *ebiten.Image {
	img, err := LoadImage(name)
	if err != nil {
		log.Fatalf("failed to load asset %q: %v", name, err)
	}
	return img
}

func opaqueBounds(img image.Image, alphaThreshold uint8) (minX, minY, maxX, maxY int, ok bool) {
	b := img.Bounds()

	minX, minY = b.Max.X, b.Max.Y
	maxX, maxY = b.Min.X-1, b.Min.Y-1

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
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

func cropImage(src image.Image, r image.Rectangle) image.Image {
	dst := image.NewNRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	for y := 0; y < r.Dy(); y++ {
		for x := 0; x < r.Dx(); x++ {
			dst.Set(x, y, src.At(r.Min.X+x, r.Min.Y+y))
		}
	}
	return dst
}

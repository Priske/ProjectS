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
	data, err := imageFiles.ReadFile(name)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
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
			core.Enemy_cultist_knife: mustImage("enemy_cultist_knife.png"),
			core.Enemy_cultist_lord:  mustImage("enemy_cultist_lord.png"),
		},
		CategoryImages: map[core.UnitCategory]*ebiten.Image{
			core.Attack:  mustImage("attack.png"),
			core.Defense: mustImage("defense.png"),
			core.Support: mustImage("support.png"),
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
	}
}

func mustImage(name string) *ebiten.Image {
	img, err := LoadImage(name)
	if err != nil {
		log.Fatalf("failed to load asset %q: %v", name, err)
	}
	return img
}

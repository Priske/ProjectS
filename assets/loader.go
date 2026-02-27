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
			core.Soldier:   mustImage("soldier.png"),
			core.Commander: mustImage("commander.png"),
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

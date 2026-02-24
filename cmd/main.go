package main

import (
	"github.com/Priske/ProjectS/internal/game"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	cellSize = 64
	boardW   = 10
	boardH   = 10
	screenW  = boardW * cellSize
	screenH  = boardH * cellSize
)

func main() {
	ebiten.SetWindowSize(screenW, screenH)
	ebiten.SetWindowTitle("ProjectS")

	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}

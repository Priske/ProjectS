package main

import (
	"github.com/Priske/ProjectS/internal/game"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	g := game.NewGame()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}

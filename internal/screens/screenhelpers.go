package screens

import "github.com/Priske/ProjectS/internal/core"

func getOffXY(g core.Game) (int, int) {
	s := g.Settings()
	boardW := s.BoardW * s.CellSize
	boardH := s.BoardH * s.CellSize

	offX := (core.VirtualW - boardW) / 2
	offY := (core.VirtualH - boardH) / 2
	return offX, offY
}

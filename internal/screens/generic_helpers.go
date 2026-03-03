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
func removeByID(units []*core.Unit, id int) []*core.Unit {
	for i, u := range units {
		if u != nil && u.UnitId == id {
			copy(units[i:], units[i+1:])
			return units[:len(units)-1]
		}
	}
	return units
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

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
func unitTypesFor(cat core.UnitCategory) []core.UnitType {
	switch cat {
	case core.Attack:
		return []core.UnitType{core.Soldier, core.Commander}
	case core.Defense:
		return []core.UnitType{ /* later */ }
	default:
		return nil
	}
}

func pointInRect(mx, my, x, y, w, h int) bool {
	return mx >= x && mx < x+w && my >= y && my < y+h
}

func removeWidget(widgets []core.Widget, target core.Widget) []core.Widget {
	for i := range widgets {
		if widgets[i] == target {
			return append(widgets[:i], widgets[i+1:]...)
		}
	}
	return widgets
}

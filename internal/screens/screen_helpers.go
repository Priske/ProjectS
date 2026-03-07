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

func mouseToCell(g core.Game, mx, my int) (cx, cy int, ok bool) {
	s := g.Settings()
	offX, offY := getOffXY(g)

	mx -= offX
	my -= offY
	if mx < 0 || my < 0 {
		return 0, 0, false
	}

	cx = mx / s.CellSize
	cy = my / s.CellSize
	if cx < 0 || cx >= s.BoardW || cy < 0 || cy >= s.BoardH {
		return 0, 0, false
	}
	return cx, cy, true
}

func cellTopLeft(g core.Game, cx, cy int) (px, py int) {
	s := g.Settings()
	offX, offY := getOffXY(g)
	return offX + cx*s.CellSize, offY + cy*s.CellSize
}

func clickHitsWidget(mx, my int, widgets []core.Widget) bool {
	for _, w := range widgets {
		x, y, ww, hh := w.Bounds()
		if mx >= x && mx < x+ww && my >= y && my < y+hh {
			return true
		}
	}
	return false
}

// Named returns
func boardGeom(g core.Game) (offX, offY, w, h int, s core.Settings) {
	s = g.Settings()
	offX, offY = getOffXY(g)
	w = s.BoardW * s.CellSize
	h = s.BoardH * s.CellSize
	return
}

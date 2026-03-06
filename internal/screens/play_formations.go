package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

func copyFormationWants(src map[core.Pos]core.UnitType) map[core.Pos]core.UnitType {
	dst := make(map[core.Pos]core.UnitType, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (ps *PlayScreen) drawFormationPreview(dst *ebiten.Image, g core.Game, formationIndex int, x, y, cell int) {
	f := g.LocalPlayer().Formations[formationIndex]

	// draw grid lines (optional)
	// then draw units
	for pos, ut := range f.Wants {
		px := x + pos.X*cell
		py := y + pos.Y*cell
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, cell)
	}
}

func (ps *PlayScreen) formationFits(cx, cy int) bool {

	if cx < 0 || cx+3 > 3 { // formation width
		return false
	}

	if cy < 0 || cy+5 > 10 { // formation height
		return false
	}

	return true
}
func (ps *PlayScreen) deployFormation(g core.Game, f *core.Formation, cx, cy int) {

	board := g.Board()

	available := append([]*core.Unit{}, ps.unPlacedUnits...)

	for pos, ut := range f.Wants {

		unitIndex := -1

		for i, u := range available {
			if u.Type == ut {
				unitIndex = i
				break
			}
		}

		if unitIndex == -1 {
			continue // no unit available of that type
		}

		unit := available[unitIndex]
		available = append(available[:unitIndex], available[unitIndex+1:]...)

		bx := cx + pos.X
		by := cy + pos.Y

		board.Location[by][bx].Unit = unit
	}

	ps.unPlacedUnits = available
}

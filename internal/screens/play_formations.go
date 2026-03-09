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

func drawFormationPreview(dst *ebiten.Image, g core.Game, formationIndex int, x, y, cell int) {
	f := g.LocalPlayer().Formations[formationIndex]

	// draw grid lines (optional)
	// then draw units
	for pos, ut := range f.Wants {
		px := x + pos.X*cell
		py := y + pos.Y*cell
		drawUnitImage(dst, g.Assets(), ut, px, py, cell)
	}
}
func (ps *PlayScreen) formationFits(g core.Game, f *core.Formation, cx, cy int) bool {
	board := g.Board()

	for pos := range f.Wants {
		bx, by := formationBoardCell(f, cx, cy, pos)

		if by < 0 || by >= len(board.Location) {
			return false
		}
		if bx < 0 || bx >= len(board.Location[by]) {
			return false
		}

		// player deployment zone: first 3 columns only
		if bx < 0 || bx >= 3 {
			return false
		}

		if board.Location[by][bx].Unit != nil {
			return false
		}
	}

	return true
}

func (ps *PlayScreen) deployFormation(g core.Game, f *core.Formation, cx, cy int, available []*core.Unit) []*core.Unit {

	board := g.Board()

	available = append([]*core.Unit{}, available...)

	for pos, ut := range f.Wants {
		bx, by := formationBoardCell(f, cx, cy, pos)

		if by < 0 || by >= len(board.Location) {
			continue
		}
		if bx < 0 || bx >= len(board.Location[by]) {
			continue
		}
		if board.Location[by][bx].Unit != nil {
			continue
		}

		unitIndex := -1
		for i, u := range available {
			if u.Type == ut {
				unitIndex = i
				break
			}
		}

		if unitIndex == -1 {
			continue
		}

		unit := available[unitIndex]
		available = append(available[:unitIndex], available[unitIndex+1:]...)
		board.Location[by][bx].Unit = unit
	}

	return available
}

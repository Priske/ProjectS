package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) handleDrop(g core.Game, mx, my int) (bool, string) {
	defer func() { ps.drag.Active = false }()

	if f, ok := ps.drag.Payload.(*core.Formation); ok {

		cx, cy, ok := ps.mouseToCell(g, mx, my)
		if !ok {
			ps.drag.Active = false
			return false, "drop off board"
		}

		if ps.formationFits(cx, cy) {
			ps.deployFormation(g, f, cx, cy)
			ps.drag.Active = false
			return true, "formation deployed"
		}

		ps.drag.Active = false
		return false, "invalid formation placement"
	}
	// Return-to-reserve: dragging a unit from board onto reserve grid
	if ps.setupMode && ps.drag.Source == interaction.DragFromBoard && ps.mouseOverReserve(mx, my) {
		board := g.Board()
		src := board.TilePtr(ps.drag.FromX, ps.drag.FromY)
		if src == nil || src.Unit == nil {
			return false, "no src unit"
		}

		u := src.Unit
		src.Unit = nil
		ps.unPlacedUnits = append(ps.unPlacedUnits, u)
		return true, "returned"
	}
	toX, toY, ok := ps.mouseToCell(g, mx, my)
	if !ok {
		return false, "drop off board"
	}

	board := g.Board()
	dst := board.TilePtr(toX, toY) // <- prefer your new accessor
	if dst == nil {
		return false, "drop off board"
	}
	if dst.Unit != nil {
		return false, "dst occupied"
	}

	// ---- DRAG FROM GRID (placement)
	if ps.drag.Source == interaction.DragFromGrid {
		u, ok := ps.drag.Payload.(*core.Unit)
		if !ok || u == nil {
			return false, "invalid payload"
		}
		//Restrict players drop to first 3 columns
		if toX >= 3 {
			return false, "place only in first 3 columns"
		}
		dst.Unit = u
		ps.unPlacedUnits = removeByID(ps.unPlacedUnits, u.UnitId)
		return true, "placed"
	}

	// ---- DRAG FROM BOARD (movement)
	fromX, fromY := ps.drag.FromX, ps.drag.FromY
	if toX == fromX && toY == fromY {
		return false, "same cell"
	}
	if ps.setupMode {
		if fromX >= 3 || toX >= 3 {
			return false, "can't move outside placement zone"
		}
	}

	dx := abs(toX - fromX)
	dy := abs(toY - fromY)

	if dx+dy != 1 && !ps.setupMode {
		return false, "illegal move"
	}

	// mutate via pointers to avoid tile-copy bugs
	src := board.TilePtr(fromX, fromY)
	if src == nil || src.Unit == nil {
		return false, "no src unit"
	}

	dst.Unit = src.Unit
	src.Unit = nil
	return true, "moved"
}

func (ps *PlayScreen) handleFormationDrop(g core.Game, mx, my int) (bool, string) {
	defer func() { ps.drag.Active = false }()

	if ps.drag.Source != interaction.DragFromFormationPalette {
		return false, "not formation drag"
	}

	ut, ok := ps.drag.Payload.(core.UnitType)
	if !ok {
		return false, "bad payload"
	}

	// Convert mouse to formation cell
	// You already know formation grid position: gridX/gridY/cell, 3x5.
	// BEST: store a pointer on ps when you create it:
	// ps.formationGrid = formationGrid
	gf := ps.formationGrid
	if gf == nil {
		return false, "no formation grid"
	}

	cx, cy, ok := gf.MouseToCell(mx, my) // currently unexported
	if !ok {
		return false, "drop outside formation"
	}

	ps.formationWants[core.Pos{X: cx, Y: cy}] = ut
	return true, "placed in formation"
}

func (ps *PlayScreen) tryDropIntoFormation(mx, my int) {
	defer func() { ps.drag.Active = false }()

	if ps.drag.Source != interaction.DragFromFormationPalette {
		return
	}
	ut, ok := ps.drag.Payload.(core.UnitType)
	if !ok {
		return
	}
	if ps.formationGrid == nil {
		return
	}

	cx, cy, ok := ps.formationGrid.MouseToCell(mx, my) // add exported method
	if !ok {
		return
	}
	ps.formationWants[core.Pos{X: cx, Y: cy}] = ut
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

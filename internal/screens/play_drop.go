package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) handleDrop(g core.Game, mx, my int) (bool, string) {
	defer func() { ps.drag.Active = false }()

	if f, ok := ps.drag.Payload.(*core.Formation); ok {
		cx, cy, ok := mouseToCell(g, mx, my)
		if !ok {
			return false, "drop off board"
		}

		if ps.formationFits(cx, cy) {
			ps.deployFormation(g, f, cx, cy)
			return true, "formation deployed"
		}

		return false, "invalid formation placement"
	}

	if ps.setup.setupMode {
		returned, reason := ps.tryReturnUnitToReserve(g, mx, my)
		if returned {
			return true, reason
		}
	}

	toX, toY, ok := mouseToCell(g, mx, my)
	if !ok {
		return false, "drop off board"
	}

	board := g.Board()
	dst := board.TilePtr(toX, toY)
	if dst == nil {
		return false, "drop off board"
	}
	if dst.Unit != nil {
		return false, "dst occupied"
	}

	if ps.drag.Source == interaction.DragFromGrid {
		u, ok := ps.drag.Payload.(*core.Unit)
		if !ok || u == nil {
			return false, "invalid payload"
		}

		if ps.setup.setupMode {
			ok, reason := ps.validateSetupPlacement(toX)
			if !ok {
				return false, reason
			}
		}

		dst.Unit = u
		ps.setup.unPlacedUnits = removeByID(ps.setup.unPlacedUnits, u.UnitId)
		return true, "placed"
	}

	fromX, fromY := ps.drag.FromX, ps.drag.FromY
	if toX == fromX && toY == fromY {
		return false, "same cell"
	}

	if ps.setup.setupMode {
		ok, reason := ps.validateSetupMove(fromX, toX)
		if !ok {
			return false, reason
		}
	}

	dx := abs(toX - fromX)
	dy := abs(toY - fromY)
	if dx+dy != 1 && !ps.setup.setupMode {
		return false, "illegal move"
	}

	src := board.TilePtr(fromX, fromY)
	if src == nil || src.Unit == nil {
		return false, "no src unit"
	}

	dst.Unit = src.Unit
	src.Unit = nil
	return true, "moved"
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
	if ps.formation.formationGrid == nil {
		return
	}

	cx, cy, ok := ps.formation.formationGrid.MouseToCell(mx, my) // add exported method
	if !ok {
		return
	}
	ps.formation.formationWants[core.Pos{X: cx, Y: cy}] = ut
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
	gf := ps.formation.formationGrid
	if gf == nil {
		return false, "no formation grid"
	}

	cx, cy, ok := gf.MouseToCell(mx, my) // currently unexported
	if !ok {
		return false, "drop outside formation"
	}

	ps.formation.formationWants[core.Pos{X: cx, Y: cy}] = ut
	return true, "placed in formation"
}

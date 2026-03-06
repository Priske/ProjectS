package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

func (ps *PlayScreen) Update(g core.Game) error {
	input := g.Input()
	if ps.drag.Active {
		ps.drag.MX = input.MX
		ps.drag.MY = input.MY
	}

	if len(ps.unPlacedUnits) == 0 && !ps.readyAdded {
		ps.readyWidget = ps.makeReadyButton(g)
		ps.widgets = append(ps.widgets, ps.readyWidget)
		ps.readyAdded = true
	}

	if len(ps.unPlacedUnits) != 0 && ps.readyAdded {
		ps.widgets = removeWidget(ps.widgets, ps.readyWidget)
		ps.readyWidget = nil
		ps.readyAdded = false
	}
	// If modal open: update only it (and return)
	//Testing Modal

	if ps.modal != nil && ps.modal.Open {
		ps.modal.Update(input)

		if ps.drag.Active && !input.LeftPressed {
			ps.tryDropIntoFormation(input.MX, input.MY)
		}
		return nil
	}

	for _, w := range ps.widgets {
		w.Update(input)
	}

	// keep drag cursor fresh

	// Start drag on click
	if !ps.drag.Active && input.LeftClicked && !ps.clickHitsWidget(input.MX, input.MY) {
		cx, cy, ok := ps.mouseToCell(g, input.MX, input.MY)
		if ok {
			tile := g.Board().Location[cy][cx]
			if tile.Unit != nil {
				px, py := ps.cellTopLeft(g, cx, cy)

				ps.drag = interaction.DragState{
					Active:   true,
					FromX:    cx,
					FromY:    cy,
					Payload:  tile.Unit,
					GrabOffX: input.MX - px,
					GrabOffY: input.MY - py,
					MX:       input.MX,
					MY:       input.MY,
				}
			}
		}
	}

	// Drop when button released
	if ps.drag.Active && !input.LeftPressed {
		ok, reason := ps.handleDrop(g, input.MX, input.MY)
		ps.lastDrop = reason
		_ = ok
	}
	return nil
}

func (ps *PlayScreen) Draw(g core.Game, screen *ebiten.Image) {
	ps.drawBackground(screen)

	s := g.Settings()
	offX, offY := getOffXY(g)
	ps.makeRightSidebar(g)
	if ps.setupMode {
		drawPlacementZone(screen, offX, offY, s.BoardH, s.CellSize, 3)
	}

	drawGrid(screen, offX, offY, s.BoardW, s.BoardH, s.CellSize)
	ps.drawUnits(g, screen)

	ps.drawUI(screen)

	if ps.modal != nil && ps.modal.Open {
		ps.modal.Draw(screen)
	}

	// MUST be last so it shows above modal
	ps.drawDraggedUnit(g, screen)

	ps.drawDebug(screen)
}

func (ps *PlayScreen) mouseToCell(g core.Game, mx, my int) (cx, cy int, ok bool) {
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

func (ps *PlayScreen) cellTopLeft(g core.Game, cx, cy int) (px, py int) {
	s := g.Settings()
	offX, offY := getOffXY(g)
	return offX + cx*s.CellSize, offY + cy*s.CellSize
}

func (ps *PlayScreen) clickHitsWidget(mx, my int) bool {
	for _, w := range ps.widgets {
		x, y, ww, hh := w.Bounds()
		if mx >= x && mx < x+ww && my >= y && my < y+hh {
			return true
		}
	}
	return false
}

// Named returns
func (ps *PlayScreen) boardGeom(g core.Game) (offX, offY, w, h int, s core.Settings) {
	s = g.Settings()
	offX, offY = getOffXY(g)
	w = s.BoardW * s.CellSize
	h = s.BoardH * s.CellSize
	return
}
func pointInRect(mx, my, x, y, w, h int) bool {
	return mx >= x && mx < x+w && my >= y && my < y+h
}

func (ps *PlayScreen) mouseOverReserve(mx, my int) bool {
	if ps.reserveGrid == nil {
		return false
	}
	x, y, w, h := ps.reserveGrid.Bounds()
	return pointInRect(mx, my, x, y, w, h)
}
func removeWidget(widgets []core.Widget, target core.Widget) []core.Widget {
	for i := range widgets {
		if widgets[i] == target {
			return append(widgets[:i], widgets[i+1:]...)
		}
	}
	return widgets
}

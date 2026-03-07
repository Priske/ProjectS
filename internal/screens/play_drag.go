package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
)

func (ps *PlayScreen) tryStartBoardDrag(g core.Game, input core.Input) {
	if !ps.drag.Active && input.LeftClicked && !clickHitsWidget(input.MX, input.MY, ps.ui.widgets) {
		cx, cy, ok := mouseToCell(g, input.MX, input.MY)
		if ok {
			tile := g.Board().Location[cy][cx]
			if tile.Unit != nil {
				px, py := cellTopLeft(g, cx, cy)

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

}

func (ps *PlayScreen) updateDragCursor(input core.Input) {
	if ps.drag.Active {
		ps.drag.MX = input.MX
		ps.drag.MY = input.MY
	}
}

func (ps *PlayScreen) tryFinishDrag(g core.Game, input core.Input) {
	if !ps.drag.Active {
		return
	}
	if input.LeftPressed {
		return
	}

	ok, reason := ps.handleDrop(g, input.MX, input.MY)
	ps.ui.lastDrop = reason
	_ = ok
}

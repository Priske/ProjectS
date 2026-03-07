package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

func (ps *PlayScreen) Update(g core.Game) error {
	input := g.Input()
	ps.updateDragCursor(input)
	ps.updateSetupWidgets(g)

	if ps.handleModalUpdate(g, input) {
		return nil
	}

	for _, w := range ps.ui.widgets {
		w.Update(input)
	}
	ps.tryStartBoardDrag(g, input)
	ps.tryFinishDrag(g, input)
	return nil
}

func (ps *PlayScreen) Draw(g core.Game, screen *ebiten.Image) {
	ps.drawBackground(screen)
	ps.drawBoard(g, screen)
	ps.drawUI(screen)
	ps.drawModal(screen)
	ps.drawDraggedUnit(g, screen)
	ps.drawDebug(screen)
}

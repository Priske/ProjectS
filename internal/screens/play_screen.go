package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

func (ps *PlayScreen) Update(g core.Game) error {
	input := g.Input()

	ps.updateDragCursor(input)

	if ps.handleModalUpdate(g, input) {
		return nil
	}

	for _, w := range ps.ui.widgets {
		w.Update(input)
	}

	if ps.battle.Active {
		ps.updateBattle(g)
		ps.tryStartBoardDrag(g, input)
		ps.tryFinishDrag(g, input)
		return nil

	}

	ps.updateSetupWidgets(g)
	ps.tryStartBoardDrag(g, input)
	ps.tryFinishDrag(g, input)

	return nil
}
func (ps *PlayScreen) Draw(g core.Game, screen *ebiten.Image) {
	ps.drawBackground(screen)
	ps.drawBoard(g, screen)
	ps.drawSelectedUnitHighlight(g, screen)
	ps.drawMoveRange(g, screen)
	ps.drawUI(screen)
	ps.drawModal(screen)
	ps.drawDraggedUnit(g, screen)
	ps.drawDebug(screen)
	ps.drawHoveredUnitInfo(g, screen)
}

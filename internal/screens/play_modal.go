package screens

import "github.com/Priske/ProjectS/internal/core"

func (ps *PlayScreen) handleModalUpdate(g core.Game, input core.Input) bool {
	if ps.ui.modal == nil || !ps.ui.modal.Open {
		return false
	}

	// If overlay exists, it gets exclusive input
	if ps.ui.overlay != nil {
		ps.ui.overlay.Update(input)
		return true
	}

	ps.ui.modal.Update(input)

	if ps.drag.Active && !input.LeftPressed {
		ps.tryDropIntoFormation(input.MX, input.MY)
	}

	return true
}

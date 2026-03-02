package screens

import (
	"github.com/Priske/ProjectS/internal/core"
)

func NewPlayScreen(g core.Game) *PlayScreen {
	ps := &PlayScreen{}

	options := ps.makeOptionsSidebar(g)

	ps.widgets = []core.Widget{options}
	return ps
}

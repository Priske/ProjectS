package screens

import (
	"github.com/Priske/ProjectS/internal/core"
)

type List struct {
	g  core.Game
	ps PlayScreen

	x, y int
	w, h int

	rowHeight  int
	hoverIndex int
}

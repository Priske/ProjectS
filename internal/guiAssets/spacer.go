package guiassets

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

type Spacer struct {
	H int
}

func (s *Spacer) Bounds() (int, int, int, int) { return 0, 0, 0, s.H }
func (s *Spacer) Update(in core.Input)         {}
func (s *Spacer) Draw(dst *ebiten.Image)       {}

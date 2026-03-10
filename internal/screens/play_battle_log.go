package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type BattleLogWidget struct {
	ps   *PlayScreen
	X, Y int
}

func (w *BattleLogWidget) Bounds() (int, int, int, int) {
	return w.X, w.Y, 200, 180
}
func (w *BattleLogWidget) SetPos(x, y int) {
	w.X = x
	w.Y = y
}

func (w *BattleLogWidget) Update(in core.Input) {}

func (w *BattleLogWidget) Draw(dst *ebiten.Image) {
	maxLines := 10
	start := 0
	if len(w.ps.battle.Log) > maxLines {
		start = len(w.ps.battle.Log) - maxLines
	}

	for i, line := range w.ps.battle.Log[start:] {
		ebitenutil.DebugPrintAt(dst, line, w.X+10, w.Y+28+i*16)
	}
}

package screens

import (
	"fmt"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SelectedUnitInfoWidget struct {
	ps *PlayScreen
	X  int
	Y  int
}

func (w *SelectedUnitInfoWidget) SetPos(x, y int) {
	w.X = x
	w.Y = y
}

func (w *SelectedUnitInfoWidget) Bounds() (int, int, int, int) {
	return w.X, w.Y, 200, 180
}

func (w *SelectedUnitInfoWidget) Update(in core.Input) {}

func (w *SelectedUnitInfoWidget) Draw(dst *ebiten.Image) {
	if w.ps == nil {
		return
	}

	u := w.ps.battle.Selected
	if u == nil {
		ebitenutil.DebugPrintAt(dst, "No unit selected", w.X+10, w.Y+28)
		return
	}

	ebitenutil.DebugPrintAt(dst, "Selected Unit", w.X+10, w.Y+28)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Name: %v", u.Name), w.X+10, w.Y+48)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("HP: %d", u.Health), w.X+10, w.Y+64)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("ATK: %d", u.AttackPower), w.X+10, w.Y+80)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("XP: %d", u.Experience), w.X+10, w.Y+96)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Pos: %d,%d", w.ps.battle.SelectedX, w.ps.battle.SelectedY), w.X+10, w.Y+112)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("statistics:"), w.X+10, w.Y+128)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("- Kills: %v", u.BattleStats.Kills), w.X+30, w.Y+144)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("- Damage dealt: %v", u.BattleStats.DamageDealt), w.X+30, w.Y+160)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("- Damage Taken: %v", u.BattleStats.DamageTaken), w.X+30, w.Y+176)
}

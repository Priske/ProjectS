package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

type BattleResult int

const (
	BattleOngoing BattleResult = iota
	BattleVictory
	BattleDefeat
)

/*
DRAWS CURRENTLY NOT IMPLEMENTED MIGTH CONCIDER LATER IF I ADD ON DEATH EFFECT OR STATUS EFFECT LIKE POISON OR BLEED THAT COULD LEAD INTO DRAWS
*/

func (ps *PlayScreen) checkBattleResult(g core.Game) BattleResult {
	board := g.Board()

	hasFlagEnemy := false
	hasFlagPlayer := false

	for y := range board.Location {
		for x := range board.Location[y] {
			u := board.Location[y][x].Unit
			if u == nil {
				continue
			}

			if isEnemyUnit(u) {
				if isFlagUnit(u) {
					hasFlagEnemy = true
				}
			} else {
				if isFlagUnit(u) {
					hasFlagPlayer = true
				}
			}
			if hasFlagEnemy && hasFlagPlayer {
				return BattleOngoing
			}
		}
	}

	if !hasFlagPlayer {
		return BattleDefeat
	}
	if !hasFlagEnemy {
		return BattleVictory
	}

	return BattleOngoing
}

func (ps *PlayScreen) resolveBattleResult(g core.Game) bool {
	if ps.ui.modal != nil && ps.ui.modal.Open {
		return true
	}

	switch ps.checkBattleResult(g) {
	case BattleVictory:
		ps.addBattleLog("Victory!")
		ps.openBattleResultModal(BattleVictory)
		return true

	case BattleDefeat:
		ps.addBattleLog("Defeat!")
		ps.openBattleResultModal(BattleDefeat)
		return true
	}

	return false
}

func isEnemyUnit(u *core.Unit) bool {
	return u.Playerid == -1
}

func isFlagUnit(u *core.Unit) bool {
	return u.UnitCategory == core.Flag
}

func (ps *PlayScreen) openBattleResultModal(result BattleResult) {

	title := "Victory!"
	message := "Enemy commander defeated."

	if result == BattleDefeat {
		title = "Defeat"
		message = "Your commander has been defeated."
	}

	modalW := 280
	modalH := 140

	sw, sh := ebiten.WindowSize()

	x := sw/2 - modalW/2
	y := sh/2 - modalH/2

	modal := GUI.MakeModal(x, y, modalW, modalH, []core.Widget{
		GUI.MakeLabel(x+80, y+40, title),
		GUI.MakeLabel(x+40, y+70, message),
		GUI.MakeButton(x+80, y+100, 120, 30, "Continue", func() {
			if ps.ui.modal != nil {
				ps.ui.modal.Close()
			}
		}),
	})

	modal.CloseOnEsc = false
	modal.CloseOnOutside = false

	modal.OnClose = func() {
		ps.ui.modal = nil
		//ps.exitBattle()
		//ps.enterSetupAfterBattle()
	}

	ps.ui.modal = modal
}

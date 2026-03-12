package screens

import (
	"fmt"
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (ps *PlayScreen) enterBattle(g core.Game) {
	ps.drag.Active = false
	ps.setup.setupMode = false

	ps.battle.Active = true
	ps.battle.Turn.Side = TurnPlayer
	ps.battle.Turn.Round = 1
	ps.battle.Selected = nil
	ps.battle.Log = nil

	ps.swapAndResetUI(ps.buildBattleUI, g)
	ps.spawnEnemySetup(g)
	ps.addBattleLog("Battle started")
	ps.addBattleLog("Player turn")
}

func (ps *PlayScreen) endPlayerTurn(g core.Game) {
	ps.battle.Selected = nil
	ps.battle.Turn = TurnState{
		Side:  TurnEnemy,
		Round: ps.battle.Turn.Round,
		Units: make(map[int]UnitTurnState),
	}
	ps.battle.SelectedAction = nil
	ps.addBattleLog("Enemy turn")
}

func (ps *PlayScreen) endEnemyTurn(g core.Game) {
	ps.battle.Selected = nil
	ps.battle.Turn = TurnState{
		Side:  TurnPlayer,
		Round: ps.battle.Turn.Round + 1,
		Units: make(map[int]UnitTurnState),
	}
	ps.battle.SelectedAction = nil
	ps.addBattleLog("Player turn")
}

func (ps *PlayScreen) updateBattle(g core.Game) {
	in := g.Input()

	// Enemy turn runs automatically
	if ps.battle.Turn.Side == TurnEnemy {
		ps.runEnemyTurn(g)
		return
	}

	// Right click opens action popup
	if in.RightClicked {
		mx, my := in.MX, in.MY
		cx, cy, ok := mouseToCell(g, mx, my)
		if !ok {
			if ps.actionPopup != nil {
				ps.actionPopup.Open = false
			}
			return
		}

		board := g.Board()
		u := board.Location[cy][cx].Unit
		if u == nil || u.Playerid != 1 {
			if ps.actionPopup != nil {
				ps.actionPopup.Open = false
			}
			return
		}

		ps.battle.Selected = u
		ps.battle.SelectedX = cx
		ps.battle.SelectedY = cy

		px, py := cellTopLeft(g, cx, cy)

		if ps.actionPopup != nil {
			ps.actionPopup.X = px + g.Settings().CellSize + 6
			ps.actionPopup.Y = py
			ps.actionPopup.W = 140
			ps.actionPopup.H = 8 + len(u.Actions)*22
			ps.actionPopup.Open = true

			ps.actionPopup.DrawFn = func(dst *ebiten.Image, x, y int) {
				rowH := 22

				for i := range u.Actions {
					action := &u.Actions[i]
					rowY := y + 4 + i*rowH

					if ps.battle.SelectedAction != nil && ps.battle.SelectedAction.ID == action.ID {
						ebitenutil.DrawRect(
							dst,
							float64(x+2),
							float64(rowY-1),
							float64(ps.actionPopup.W-4),
							float64(rowH-2),
							color.RGBA{70, 70, 100, 255},
						)
					}

					label := action.Name
					if !ps.canUseAction(u, action) {
						label = "X " + label
					}

					ebitenutil.DebugPrintAt(dst, label, x+8, rowY+4)
				}
			}

			ps.actionPopup.OnClick = func(mx, my int) {
				rowH := 22
				index := (my - (ps.actionPopup.Y + 4)) / rowH
				if index < 0 || index >= len(u.Actions) {
					return
				}

				action := &u.Actions[index]
				if !ps.canUseAction(u, action) {
					return
				}

				ps.battle.SelectedAction = action
				ps.actionPopup.Open = false
				ps.addBattleLog(fmt.Sprintf("Selected action: %s", action.Name))
			}
		}

		return
	}

	// Popup consumes left click before board logic
	if in.LeftClicked && ps.actionPopup != nil && ps.actionPopup.Open {
		ps.actionPopup.Update(in)
		return
	}

	// Ignore frames without left click
	if !in.LeftClicked {
		return
	}

	mx, my := in.MX, in.MY
	cx, cy, ok := mouseToCell(g, mx, my)
	if !ok {
		return
	}

	board := g.Board()
	clickedUnit := board.Location[cy][cx].Unit

	// Select first unit
	if ps.battle.Selected == nil {
		ps.trySelectUnit(g, mx, my)
		return
	}

	// Reselect friendly unit
	if clickedUnit != nil && clickedUnit.Playerid == ps.battle.Selected.Playerid {
		ps.trySelectUnit(g, mx, my)
		return
	}

	// Attack enemy using selected action
	if clickedUnit != nil && clickedUnit.Playerid != ps.battle.Selected.Playerid {
		if ps.battle.SelectedAction != nil {
			ps.tryAttackWithSelectedAction(g, ps.battle.SelectedAction, cx, cy)
		}
		return
	}
}

func (ps *PlayScreen) runEnemyTurn(g core.Game) {
	enemies := ps.enemyUnitsOnBoard(g)

	for _, e := range enemies {
		if ps.enemyTryAttackAdjacent(g, e.X, e.Y) {
			continue
		}
		ps.enemyTryMoveTowardPlayer(g, e.X, e.Y)
	}

	ps.endEnemyTurn(g)
}

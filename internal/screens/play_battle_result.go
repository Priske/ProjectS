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
		ps.openBattleResultModal(BattleVictory, g)
		return true

	case BattleDefeat:
		ps.addBattleLog("Defeat!")
		ps.openBattleResultModal(BattleDefeat, g)
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

func (ps *PlayScreen) openBattleResultModal(result BattleResult, g core.Game) {
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
		ps.finishBattleResult(result, g)
	}

	ps.ui.modal = modal
}

func (ps *PlayScreen) finishBattleResult(result BattleResult, g core.Game) {
	switch result {
	case BattleVictory:
		ps.handleBattleVictory(g)
	case BattleDefeat:
		ps.handleBattleDefeat(g)
	}
}

func (ps *PlayScreen) handleBattleVictory(g core.Game) {
	ps.exitBattle()
	ps.openRewardModal(g)
	ps.setup.setupMode = true
	ps.healAllPlayerUnits(g)
}

func (ps *PlayScreen) handleBattleDefeat(g core.Game) {
	ps.exitBattle()
	g.SetScreen(NewMenuScreen(g))
}

func (ps *PlayScreen) exitBattle() {
	ps.battle.Active = false
	ps.battle.Selected = nil
	ps.battle.SelectedAction = nil
	ps.battle.ActionMenuOpen = false
	ps.battle.Log = nil
	ps.resetBattleTurnState()
}

func (ps *PlayScreen) openRewardModal(g core.Game) {
	choices := ps.generateRewardChoices(g)

	modalW := 520
	modalH := 260

	sw, sh := ebiten.WindowSize()
	x := sw/2 - modalW/2
	y := sh/2 - modalH/2

	cardW := 150
	cardH := 150
	cardGap := 15
	startX := x + (modalW-(3*cardW+2*cardGap))/2
	cardY := y + 55

	children := []core.Widget{
		GUI.MakeLabel(x+185, y+20, "Choose a Reward"),
	}

	for i, choice := range choices {
		cardX := startX + i*(cardW+cardGap)

		title := choice.Title
		if title == "" && choice.Unit != nil {
			title = choice.Unit.Name
		}

		card := GUI.MakePanel(cardX, cardY, cardW, cardH, "", []core.Widget{})
		card.AutoLayout = false

		card.Children = append(card.Children,
			GUI.MakeLabel(cardX+12, cardY+12, title),
		)

		if choice.Unit != nil {
			card.Children = append(card.Children,
				&GUI.RewardUnitPreviewWidget{
					X:        cardX + 35,
					Y:        cardY + 34,
					W:        80,
					H:        58,
					UnitType: choice.Unit.Type,
					Game:     g,
				},
			)
		}

		c := choice
		selectBtn := GUI.MakeButton(cardX+25, cardY+108, 100, 28, "Choose", func() {
			ps.applyReward(c, g)
			if ps.ui.modal != nil {
				ps.ui.modal.Close()
			}
			ps.resetSetupState(g)
			ps.swapAndResetUI(ps.buildSetupUI, g)
		})

		children = append(children, card, selectBtn)
	}

	modal := GUI.MakeModal(x, y, modalW, modalH, children)
	modal.CloseOnEsc = false
	modal.CloseOnOutside = false
	modal.OnClose = func() {
		ps.ui.modal = nil
	}

	ps.ui.modal = modal
}
func (ps *PlayScreen) generateRewardChoices(g core.Game) []RewardChoice {
	player := g.LocalPlayer()
	playerID := player.Playerid
	medic := core.MakeNewMedic(playerID, PreviewUnitID)
	shield := core.MakeNewShield(playerID, PreviewUnitID)
	sniper := core.MakeNewSniper(playerID, PreviewUnitID)
	return []RewardChoice{
		{Kind: RewardUnit, Title: "Medic", Description: "Add a Medic to your roster.", Unit: medic},
		{Kind: RewardUnit, Title: "Shield", Description: "Add a Shield unit to your roster.", Unit: shield},
		{Kind: RewardUnit, Title: "Sniper", Description: "Add a Sniper to your roster.", Unit: sniper},
	}
}

func (ps *PlayScreen) applyReward(choice RewardChoice, g core.Game) {
	player := g.LocalPlayer()

	switch choice.Kind {
	case RewardUnit:
		if choice.Unit == nil {
			return
		}

		choice.Unit.UnitId = g.NewUnitID()
		player.Units = append(player.Units, choice.Unit)
	}
}

func (ps *PlayScreen) healAllPlayerUnits(g core.Game) {
	player := g.LocalPlayer()
	if player == nil {
		return
	}

	for _, u := range player.Units {
		if u == nil {
			continue
		}
		u.CurrentHealth = u.MaxHealth
	}
}

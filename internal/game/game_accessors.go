package game

import "github.com/Priske/ProjectS/internal/core"

func (g *Game) LocalPlayer() *core.Player {
	if g.localSlot < 0 || g.localSlot >= len(g.players) {
		return nil
	}
	return g.players[g.localSlot]
}

func (g *Game) TurnPlayer() *core.Player {
	if g.turnSlot < 0 || g.turnSlot >= len(g.players) {
		return nil
	}
	return g.players[g.turnSlot]
}

func (g *Game) NextTurn() {
	if len(g.players) == 0 {
		return
	}
	g.turnSlot = (g.turnSlot + 1) % len(g.players)
}

func (g *Game) Settings() core.Settings {
	return g.settings
}
func (g *Game) Board() core.GameBoard {
	return g.board
}
func (g *Game) Assets() core.Assets {
	return g.assets
}
func (g *Game) Players() []*core.Player {
	return g.players
}
func (g *Game) Input() core.Input {
	return g.input
}

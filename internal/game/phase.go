package game

import "github.com/Priske/ProjectS/internal/screens"

type Phase int

const (
	PhaseMenu Phase = iota
	PhaseSetup
	PhasePlaying
	PhaseGameOver
	PhaseModeSelect
)

func (g *Game) SetPhase(p Phase) {
	g.phase = p

	switch p {
	case PhaseMenu:
		g.SetScreen(screens.NewMenuScreen(g))
	case PhaseModeSelect:
		g.SetScreen(screens.NewSelectScreen(g))
	case PhaseSetup:
		// Setup screen is where reserve grid + placement happens
		g.SetScreen(screens.NewPlayScreen(g)) // or NewSetupScreen(g) if you want it separate

	case PhasePlaying:
		// Can reuse same PlayScreen; rules differ by g.phase
		g.SetScreen(screens.NewPlayScreen(g))

	case PhaseGameOver:
		//	g.SetScreen(screens.NewGameOverScreen(g))
	}
}

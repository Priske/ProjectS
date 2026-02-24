package game

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/Priske/ProjectS/internal/screens"

	"github.com/hajimehoshi/ebiten/v2"
)

type Phase int

const (
	PhaseMenu Phase = iota
	PhaseSetup
	PhasePlaying
	PhaseGameOver
)

type Game struct {
	settings  Settings
	phase     Phase
	screen    core.Screen
	selectedX int
	selectedY int
	hasSel    bool

	mouseWasDown bool // simple click edge-detect
}

// innit currently in update Phases
func NewGame() *Game {
	g := &Game{
		phase: PhaseMenu,

		// load assets, init board, etc
	}
	g.screen = screens.NewMenuScreen()
	g.initBoard()
	return g
}
func (g *Game) SetScreen(s core.Screen) {
	g.screen = s
}

func (g *Game) Update() error {
	switch g.phase {
	case PhaseMenu:
		return g.updateMenu()
	case PhaseSetup:
		return g.updateSetup()
	case PhasePlaying:
		return g.updatePlaying()
	case PhaseGameOver:
		return g.updateGameOver()
	default:
		return nil
	}
}
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.phase {
	case PhaseMenu:
		g.drawMenu(screen)
	case PhaseSetup:
		g.drawSetup(screen)
	case PhasePlaying:
		g.drawPlaying(screen)
	case PhaseGameOver:
		g.drawGameOver(screen)
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 640, 480
}

func (g *Game) JustClickedLeft() bool {
	down := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	just := down && !g.mouseWasDown
	g.mouseWasDown = down
	return just
}

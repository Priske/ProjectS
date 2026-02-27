package game

import (
	"github.com/Priske/ProjectS/assets"
	"github.com/Priske/ProjectS/internal/core"
	"github.com/Priske/ProjectS/internal/screens"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Phase int

const (
	PhaseMenu Phase = iota
	PhaseSetup
	PhasePlaying
	PhaseGameOver
)

type Game struct {
	settings core.Settings
	phase    Phase
	screen   core.Screen
	board    core.GameBoard
	input    core.Input
	assets   core.Assets

	mouseWasDown bool // simple click edge-detect
}

func GetScreenWidthAndHeight(g *Game) (int, int) {
	s := g.Settings()
	w := s.BoardW * s.CellSize
	h := s.BoardH * s.CellSize
	return w, h
}

/*
using data as a 1 dimensional representation to enhance performance,
incase performance is important later on. tldr; data holds all tiles,
forloop assigns matrices to the data[0:5], data[5:10],data[10:15],..
*/
func MakeBoard(boardH, boardW int) core.GameBoard {
	data := make([]core.Tile, boardH*boardW)
	board := make([][]core.Tile, boardH)

	for i := 0; i < boardH; i++ {
		board[i] = data[i*boardW : (i+1)*boardW]
	}

	return core.GameBoard{
		LocationXY: board,
	}
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

// innit currently in update Phases
func NewGame() *Game {
	g := &Game{}
	g.SetScreen(screens.NewMenuScreen(g))
	g.settings = core.DefaultSettings()
	g.board = MakeBoard(g.settings.BoardH, g.settings.BoardW)
	g.assets = assets.MustLoadAll()
	g.board.LocationXY[0][0].Unit = &core.Unit{Type: core.Soldier}
	g.board.LocationXY[1][2].Unit = &core.Unit{Type: core.Commander}
	ebiten.SetWindowSize(core.VirtualW, core.VirtualH)
	ebiten.SetWindowTitle("ProjectS")
	ebiten.SetWindowResizable(true)
	return g
}
func (g *Game) SetScreen(s core.Screen) {
	g.screen = s
}

func (g *Game) Update() error {
	g.input = g.pollInput()
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

	return core.VirtualW, core.VirtualH
}

func (g *Game) Input() core.Input {
	return g.input
}
func (g *Game) pollInput() core.Input {
	mx, my := ebiten.CursorPosition()

	in := core.Input{
		MX:          mx,
		MY:          my,
		LeftPressed: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		LeftClicked: inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
		RuneBuffer:  ebiten.AppendInputChars(nil),
		Backspace:   inpututil.IsKeyJustPressed(ebiten.KeyBackspace),
		Escape:      inpututil.IsKeyJustPressed(ebiten.KeyEscape),
	}

	return in
}

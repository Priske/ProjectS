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
	players  []*core.Player

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
	//g.board.LocationXY[0][0].Unit = &core.Unit{Type: core.Soldier}
	//g.board.LocationXY[1][2].Unit = &core.Unit{Type: core.Commander}
	p := makeTestPlayer()
	g.players = append(g.players, &p)
	ebiten.SetWindowSize(core.VirtualW, core.VirtualH)
	ebiten.SetWindowTitle("ProjectS")
	ebiten.SetWindowResizable(true)

	return g
}
func (g *Game) SetScreen(s core.Screen) {
	g.screen = s
}
func (g *Game) Players() []*core.Player {
	return g.players
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

func makeTestFormation() core.Formation {
	p1 := core.Pos{X: 2, Y: 2}
	p2 := core.Pos{X: 3, Y: 2}
	p3 := core.Pos{X: 1, Y: 2}
	wants := make(map[core.Pos]core.UnitType)
	wants[p1] = core.Soldier
	wants[p2] = core.Commander
	wants[p3] = core.Soldier

	formationA := core.Formation{
		Name:  "Test1",
		GridW: 5,
		GridH: 5,
		Wants: wants,
	}
	return formationA
}

func makeTestPlayer() core.Player {
	u1 := core.Unit{
		Type:       core.Soldier,
		UnitId:     1,
		Health:     1,
		Attack:     1,
		Experience: 0,
		Playerid:   007,
	}
	u2 := core.Unit{
		Type:       core.Soldier,
		UnitId:     2,
		Health:     1,
		Attack:     1,
		Experience: 0,
		Playerid:   007,
	}
	u3 := core.Unit{
		Type:       core.Soldier,
		UnitId:     3,
		Health:     1,
		Attack:     1,
		Experience: 0,
		Playerid:   007,
	}
	units := make([]*core.Unit, 3)
	units[0] = &u1
	units[1] = &u2
	units[2] = &u3
	formations := make([]core.Formation, 1)
	formations[0] = makeTestFormation()
	return core.Player{
		Playerid:   007,
		Name:       "Bond James",
		Units:      units,
		Formations: formations,
	}
}

/*


 */

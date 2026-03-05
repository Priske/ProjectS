package game

import (
	"fmt"

	"github.com/Priske/ProjectS/assets"
	"github.com/Priske/ProjectS/internal/core"
	"github.com/Priske/ProjectS/internal/screens"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	nextUnitID int //nextUnitID Added, To keep track of adding unit id's remove and replace when DB is added with UUID's
	settings   core.Settings
	phase      Phase
	screen     core.Screen
	board      core.GameBoard
	input      core.Input
	assets     core.Assets
	players    []*core.Player //all players in a given game
	localSlot  int            // who is logged in on this machine
	turnSlot   int            // whose turn it is in the match

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
		Location: board,
	}
}
func PlacePlayersFormation(g *Game) {
	player := g.LocalPlayer()
	used := make(map[int]bool)

	units := player.Units
	form := player.Formations[0]

	for pos, wantType := range form.Wants {
		// IMPORTANT: board is [y][x]
		tile := g.board.TilePtr(pos.X, pos.Y)

		// place only if empty
		if tile.Unit != nil {
			continue
		}

		for _, u := range units {
			if u.Type == wantType && !used[u.UnitId] {
				tile.Unit = u
				used[u.UnitId] = true
				break
			}
		}
	}
}

// innit currently in update Phases
func NewGame() *Game {
	g := &Game{
		players: []*core.Player{
			core.NewPlayer(0, "P1"),
			core.NewPlayer(1, "P2"),
		},
		localSlot: 0,
		turnSlot:  0,
	}

	g.SetScreen(screens.NewMenuScreen(g))
	g.settings = core.DefaultSettings()

	g.assets = assets.MustLoadAll()

	ebiten.SetWindowSize(core.VirtualW, core.VirtualH)
	ebiten.SetWindowTitle("ProjectS")
	ebiten.SetWindowResizable(true)

	return g
}
func (g *Game) SetScreen(s core.Screen) {
	g.screen = s
}
func (g *Game) SetLocalPlayer(*core.Player) {
	g.LocalPlayer()

}

func (g *Game) Update() error {
	g.input = g.pollInput()
	if g.screen == nil {
		return nil
	}
	return g.screen.Update(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.screen == nil {
		return
	}
	g.screen.Draw(g, screen)
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {

	return core.VirtualW, core.VirtualH
}

func (g *Game) pollInput() core.Input {
	mx, my := ebiten.CursorPosition()

	in := core.Input{
		MX:           mx,
		MY:           my,
		LeftPressed:  ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		LeftClicked:  inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
		RightPressed: ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight),
		RightClicked: inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight),
		RuneBuffer:   ebiten.AppendInputChars(nil),
		Backspace:    inpututil.IsKeyJustPressed(ebiten.KeyBackspace),
		Escape:       inpututil.IsKeyJustPressed(ebiten.KeyEscape),
	}

	return in
}

// / Remove When UUID is added Temp Solution
func (g *Game) NewUnitID() int {
	g.nextUnitID++
	return g.nextUnitID
}
func (g *Game) InitializeNewGame(localPlayerID int) {
	g.settings = core.DefaultSettings()
	// Fresh match state
	g.players = make([]*core.Player, 0, 2)
	g.board = MakeBoard(g.settings.BoardH, g.settings.BoardW)
	// Local player
	local := &core.Player{Playerid: localPlayerID}
	local.Units = g.InitializeStartingUnits(localPlayerID)

	// Opponent (bot/placeholder for now)
	opponentID := -1
	opponent := &core.Player{Playerid: opponentID}
	opponent.Units = g.InitializeStartingUnits(opponentID)

	g.players = append(g.players, local, opponent)

	// Session state
	g.localSlot = 0
	g.turnSlot = 0
	g.SetPhase(PhaseMenu)
	fmt.Printf("players made: local=%d opponent=%d\n", local.Playerid, opponent.Playerid)
}

func (g *Game) InitializeStartingUnits(playerId int) []*core.Unit {
	return []*core.Unit{
		core.MakeNewSoldier(playerId, g.NewUnitID()),
		core.MakeNewSoldier(playerId, g.NewUnitID()),
		core.MakeNewCommander(playerId, g.NewUnitID()),
	}
}

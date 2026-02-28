package screens

import (
	"image/color"

	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type PlayScreen struct {
	widgets  []core.Widget
	drag     interaction.DragState
	lastDrop string
}

func NewPlayScreen(g core.Game) *PlayScreen {
	ps := &PlayScreen{}
	panelW := 260
	headerH := 44
	x := 20
	y := 40

	options := GUI.MakeCollapsible(x, y, panelW, headerH, "Options . . .", []core.Widget{
		GUI.MakeButton(0, 0, 240, 50, "Place Unit", func() {
			players := g.Players()
			for _, p := range players {

			}
		}),
	})

	ps.widgets = []core.Widget{options}
	return ps

}
func (ps *PlayScreen) Update(g core.Game) error {
	input := g.Input()

	for _, w := range ps.widgets {
		w.Update(input)
	}

	// keep drag cursor fresh
	if ps.drag.Active {
		ps.drag.MX = input.MX
		ps.drag.MY = input.MY
	}

	// Start drag on click
	if !ps.drag.Active && input.LeftClicked && !ps.clickHitsWidget(input.MX, input.MY) {
		cx, cy, ok := ps.mouseToCell(g, input.MX, input.MY)
		if ok {
			tile := g.Board().LocationXY[cy][cx]
			if tile.Unit != nil {
				px, py := ps.cellTopLeft(g, cx, cy)

				ps.drag = interaction.DragState{
					Active:   true,
					FromX:    cx,
					FromY:    cy,
					Payload:  tile.Unit,
					GrabOffX: input.MX - px,
					GrabOffY: input.MY - py,
					MX:       input.MX,
					MY:       input.MY,
				}
			}
		}
	}

	// Drop when button released
	if ps.drag.Active && !input.LeftPressed {
		ok, reason := ps.handleDrop(g, input.MX, input.MY)
		ps.lastDrop = reason
		_ = ok
	}
	return nil
}

func (ps *PlayScreen) Draw(g core.Game, screen *ebiten.Image) {
	ps.drawBackground(screen)
	s := g.Settings()
	offX, offY := getOffXY(g)
	drawGrid(screen, offX, offY, s.BoardW, s.BoardH, s.CellSize)
	ps.drawUnits(g, screen)       // skips dragged source tile
	ps.drawDraggedUnit(g, screen) // draws on top
	ps.drawUI(screen)
	ps.drawDebug(screen)

}

func (ps *PlayScreen) mouseToCell(g core.Game, mx, my int) (cx, cy int, ok bool) {
	s := g.Settings()
	offX, offY := getOffXY(g)

	mx -= offX
	my -= offY
	if mx < 0 || my < 0 {
		return 0, 0, false
	}

	cx = mx / s.CellSize
	cy = my / s.CellSize
	if cx < 0 || cx >= s.BoardW || cy < 0 || cy >= s.BoardH {
		return 0, 0, false
	}
	return cx, cy, true
}

func (ps *PlayScreen) cellTopLeft(g core.Game, cx, cy int) (px, py int) {
	s := g.Settings()
	offX, offY := getOffXY(g)
	return offX + cx*s.CellSize, offY + cy*s.CellSize
}

func (ps *PlayScreen) handleDrop(g core.Game, mx, my int) (committed bool, reason string) {
	toX, toY, ok := ps.mouseToCell(g, mx, my)
	if !ok {
		ps.drag.Active = false
		return false, "drop off board"
	}

	fromX, fromY := ps.drag.FromX, ps.drag.FromY
	if toX == fromX && toY == fromY {
		ps.drag.Active = false
		return false, "same cell"
	}

	dst := g.Board().LocationXY[toY][toX]
	if dst.Unit != nil {
		ps.drag.Active = false
		return false, "dst occupied"
	}

	dx := toX - fromX
	if dx < 0 {
		dx = -dx
	}
	dy := toY - fromY
	if dy < 0 {
		dy = -dy
	}
	if dx+dy != 1 {
		ps.drag.Active = false
		return false, "illegal move (dx+dy != 1)"
	}

	board := g.Board()

	board.LocationXY[toY][toX].Unit = board.LocationXY[fromY][fromX].Unit
	board.LocationXY[fromY][fromX].Unit = nil
	ps.drag.Active = false
	return true, "moved"
}

func (ps *PlayScreen) clickHitsWidget(mx, my int) bool {
	for _, w := range ps.widgets {
		x, y, ww, hh := w.Bounds()
		if mx >= x && mx < x+ww && my >= y && my < y+hh {
			return true
		}
	}
	return false
}

func (ps *PlayScreen) drawBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 24, 255})
}

// Named returns
func (ps *PlayScreen) boardGeom(g core.Game) (offX, offY, w, h int, s core.Settings) {
	s = g.Settings()
	offX, offY = getOffXY(g)
	w = s.BoardW * s.CellSize
	h = s.BoardH * s.CellSize
	return
}
func drawGrid(screen *ebiten.Image, offX, offY, boardW, boardH, cellSize int) {
	w := boardW * cellSize
	h := boardH * cellSize
	line := color.RGBA{60, 60, 70, 255}

	for x := 0; x <= boardW; x++ {
		px := float64(offX + x*cellSize)
		ebitenutil.DrawRect(screen, px, float64(offY), 1, float64(h), line)
	}
	for y := 0; y <= boardH; y++ {
		py := float64(offY + y*cellSize)
		ebitenutil.DrawRect(screen, float64(offX), py, float64(w), 1, line)
	}
}
func (ps *PlayScreen) drawUnits(g core.Game, screen *ebiten.Image) {
	offX, offY, _, _, s := ps.boardGeom(g)
	assets := g.Assets()

	for y := 0; y < s.BoardH; y++ {
		for x := 0; x < s.BoardW; x++ {
			if ps.drag.Active && x == ps.drag.FromX && y == ps.drag.FromY {
				continue
			}
			tile := g.Board().LocationXY[y][x]
			if tile.Unit == nil {
				continue
			}
			ps.drawUnitImage(screen, assets, tile.Unit.Type, offX+x*s.CellSize, offY+y*s.CellSize, s.CellSize)
		}
	}
}

func (ps *PlayScreen) drawUnitImage(
	screen *ebiten.Image,
	assets core.Assets,
	unitType core.UnitType,
	px, py, cellSize int,
) {
	img := assets.UnitImages[unitType]
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	op.GeoM.Scale(float64(cellSize)/float64(sw), float64(cellSize)/float64(sh))
	op.GeoM.Translate(float64(px), float64(py))
	screen.DrawImage(img, op)
}

func (ps *PlayScreen) drawDraggedUnit(g core.Game, screen *ebiten.Image) {
	if !ps.drag.Active || ps.drag.Payload == nil {
		return
	}
	u, ok := ps.drag.Payload.(*core.Unit)
	if !ok || u == nil {
		return
	}
	s := g.Settings()
	assets := g.Assets()
	img := assets.UnitImages[u.Type]
	if img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	op.GeoM.Scale(float64(s.CellSize)/float64(sw), float64(s.CellSize)/float64(sh))

	drawX := ps.drag.MX - ps.drag.GrabOffX
	drawY := ps.drag.MY - ps.drag.GrabOffY
	op.GeoM.Translate(float64(drawX), float64(drawY))
	screen.DrawImage(img, op)
}

func (ps *PlayScreen) drawUI(screen *ebiten.Image) {
	for _, b := range ps.widgets {
		b.Draw(screen)
	}
}

func (ps *PlayScreen) drawDebug(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Drop: "+ps.lastDrop)
}

package screens

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

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

func drawPlacementZone(screen *ebiten.Image, offX, offY, boardH, cellSize int, cols int) {
	// translucent green fill
	fill := color.RGBA{0, 200, 0, 70}

	for y := 0; y < boardH; y++ {
		for x := 0; x < cols; x++ {
			px := float64(offX + x*cellSize)
			py := float64(offY + y*cellSize)
			ebitenutil.DrawRect(screen, px, py, float64(cellSize), float64(cellSize), fill)
		}
	}
}

func (ps *PlayScreen) drawUnits(g core.Game, screen *ebiten.Image) {
	offX, offY, _, _, s := boardGeom(g)
	assets := g.Assets()

	for y := 0; y < s.BoardH; y++ {
		for x := 0; x < s.BoardW; x++ {
			if ps.drag.Active && x == ps.drag.FromX && y == ps.drag.FromY {
				continue
			}
			tile := g.Board().Location[y][x]
			if tile.Unit == nil {
				continue
			}
			drawUnitImage(screen, assets, tile.Unit.Type, offX+x*s.CellSize, offY+y*s.CellSize, s.CellSize)
		}
	}
}

func drawUnitImage(screen *ebiten.Image, assets core.Assets, unitType core.UnitType, px, py, cellSize int) {
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

	s := g.Settings()
	assets := g.Assets()

	// 1) UnitType drag
	if ut, ok := ps.drag.Payload.(core.UnitType); ok {
		img := assets.UnitImages[ut]
		if img == nil {
			return
		}
		op := &ebiten.DrawImageOptions{}
		sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
		op.GeoM.Scale(float64(s.CellSize)/float64(sw), float64(s.CellSize)/float64(sh))
		op.GeoM.Translate(float64(ps.drag.MX-ps.drag.GrabOffX), float64(ps.drag.MY-ps.drag.GrabOffY))
		screen.DrawImage(img, op)
		return
	}
	//Draw formation drag
	if f, ok := ps.drag.Payload.(*core.Formation); ok {

		cx, cy, ok := mouseToCell(g, ps.drag.MX, ps.drag.MY)
		if !ok {
			return
		}

		offX, offY := getOffXY(g)
		s := g.Settings()

		for pos, ut := range f.Wants {

			px := offX + (cx+pos.X)*s.CellSize
			py := offY + (cy+pos.Y)*s.CellSize

			drawUnitImage(screen, g.Assets(), ut, px, py, s.CellSize)
		}

		return
	}
	// 2) *Unit drag
	u, ok := ps.drag.Payload.(*core.Unit)
	if !ok || u == nil {
		return
	}
	img := assets.UnitImages[u.Type]
	if img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	op.GeoM.Scale(float64(s.CellSize)/float64(sw), float64(s.CellSize)/float64(sh))
	op.GeoM.Translate(float64(ps.drag.MX-ps.drag.GrabOffX), float64(ps.drag.MY-ps.drag.GrabOffY))
	screen.DrawImage(img, op)
}

func (ps *PlayScreen) drawUI(screen *ebiten.Image) {
	for _, b := range ps.ui.widgets {
		b.Draw(screen)
	}
}

func (ps *PlayScreen) drawDebug(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Drop: "+ps.ui.lastDrop)
}

func (ps *PlayScreen) drawBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 24, 255})
}

func (ps *PlayScreen) drawBoard(g core.Game, screen *ebiten.Image) {
	s := g.Settings()
	offX, offY := getOffXY(g)

	if ps.setup.setupMode {
		drawPlacementZone(screen, offX, offY, s.BoardH, s.CellSize, 3)
	}

	drawGrid(screen, offX, offY, s.BoardW, s.BoardH, s.CellSize)
	ps.drawUnits(g, screen)
}
func (ps *PlayScreen) drawModal(screen *ebiten.Image) {
	if ps.ui.modal != nil && ps.ui.modal.Open {
		ps.ui.modal.Draw(screen)
	}
	if ps.ui.overlay != nil {
		ps.ui.overlay.Draw(screen)
	}
}

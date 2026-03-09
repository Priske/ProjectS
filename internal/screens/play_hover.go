package screens

import (
	"fmt"
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type HoveredUnitInfo struct {
	Unit *core.Unit
	X    int
	Y    int
	W    int
	H    int
}

func drawUnitInfoCard(dst *ebiten.Image, g core.Game, info HoveredUnitInfo) {
	if info.Unit == nil {
		return
	}

	u := info.Unit
	panelX := info.X + info.W + 8
	panelY := info.Y
	panelW := 140
	panelH := 84

	ebitenutil.DrawRect(dst, float64(panelX), float64(panelY), float64(panelW), float64(panelH), color.RGBA{25, 25, 35, 230})

	assets := g.Assets()
	img := assets.UnitImage(u.Type)
	if img != nil {
		op := &ebiten.DrawImageOptions{}
		sw, sh := img.Bounds().Dx(), img.Bounds().Dy()

		scaleX := 32.0 / float64(sw)
		scaleY := 32.0 / float64(sh)
		scale := scaleX
		if scaleY < scale {
			scale = scaleY
		}

		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(float64(panelX+8), float64(panelY+8))
		dst.DrawImage(img, op)
	}

	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Type: %v", u.Type), panelX+48, panelY+8)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("HP: %d", u.Health), panelX+48, panelY+24)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("ATK: %d", u.Attack), panelX+48, panelY+38)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("XP: %d", u.Experience), panelX+48, panelY+52)
}

func (ps *PlayScreen) hoveredReserveUnitInfo(g core.Game) (HoveredUnitInfo, bool) {
	grid := ps.reserve.grid
	if grid == nil || grid.Get == nil {
		return HoveredUnitInfo{}, false
	}

	in := g.Input()
	mx, my := in.MX, in.MY

	x, y, w, h := grid.Bounds()
	if !pointInRect(mx, my, x, y, w, h) {
		return HoveredUnitInfo{}, false
	}

	cx := (mx - x) / grid.Cell
	cy := (my - y) / grid.Cell
	if cx < 0 || cx >= grid.Cols || cy < 0 || cy >= grid.Rows {
		return HoveredUnitInfo{}, false
	}

	payload := grid.Get(cx, cy)
	u, ok := payload.(*core.Unit)
	if !ok || u == nil {
		return HoveredUnitInfo{}, false
	}

	return HoveredUnitInfo{
		Unit: u,
		X:    x + cx*grid.Cell,
		Y:    y + cy*grid.Cell,
		W:    grid.Cell,
		H:    grid.Cell,
	}, true
}
func (ps *PlayScreen) hoveredBoardUnitInfo(g core.Game) (HoveredUnitInfo, bool) {
	in := g.Input()

	cx, cy, ok := mouseToCell(g, in.MX, in.MY)
	if !ok {
		return HoveredUnitInfo{}, false
	}

	tile := g.Board().Location[cy][cx]
	if tile.Unit == nil {
		return HoveredUnitInfo{}, false
	}

	px, py := cellTopLeft(g, cx, cy)
	cell := g.Settings().CellSize

	return HoveredUnitInfo{
		Unit: tile.Unit,
		X:    px,
		Y:    py,
		W:    cell,
		H:    cell,
	}, true
}

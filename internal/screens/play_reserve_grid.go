package screens

import (
	"math"

	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

func (ps *PlayScreen) makeUnitsGrid(g core.Game) *GUI.GridField {
	unitCount := len(ps.unPlacedUnits)

	// fixed 5 columns, variable rows
	cols := 5
	if unitCount < cols {
		cols = unitCount
		if cols == 0 {
			cols = 1
		}
	}
	rows := int(math.Ceil(float64(unitCount) / 5.0))
	if rows == 0 {
		rows = 1
	}

	grid := GUI.MakeGridField(0, 0, cols, rows, 48)
	grid.ShowGrid = true

	// (cx,cy) -> unit pointer (from unPlacedUnits)
	grid.Get = func(cx, cy int) any {
		i := cy*grid.Cols + cx
		if i < 0 || i >= len(ps.unPlacedUnits) {
			return nil
		}
		return ps.unPlacedUnits[i]
	}
	grid.OnBeginDrag = func(cx, cy int, payload any) {
		u, ok := payload.(*core.Unit)
		if !ok || u == nil {
			return
		}

		in := g.Input()

		// cell top-left in screen coords
		px := grid.X + cx*grid.Cell
		py := grid.Y + cy*grid.Cell

		ps.drag = interaction.DragState{
			Active:   true,
			Source:   interaction.DragFromGrid,
			FromX:    cx,
			FromY:    cy,
			Payload:  u,
			GrabOffX: in.MX - px,
			GrabOffY: in.MY - py,
			MX:       in.MX,
			MY:       in.MY,
		}
	}

	// draw the unit image
	grid.DrawCell = func(dst *ebiten.Image, cx, cy int, px, py, size int, payload any) {
		u, ok := payload.(*core.Unit)
		if !ok || u == nil {
			return
		}

		img := g.Assets().UnitImages[u.Type]
		if img == nil {
			return
		}

		op := &ebiten.DrawImageOptions{}
		sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
		op.GeoM.Scale(float64(size)/float64(sw), float64(size)/float64(sh))
		op.GeoM.Translate(float64(px), float64(py))
		dst.DrawImage(img, op)
	}

	return grid
}

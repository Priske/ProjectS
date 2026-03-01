package guiassets

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GridField struct {
	X, Y       int
	Cols, Rows int
	Cell       int

	ShowGrid bool

	// Data source: caller provides getters/setters so GridField doesn’t own game state
	Get func(cx, cy int) any
	Set func(cx, cy int, v any) // optional (for internal rearranging)

	// Optional callbacks
	OnCellClick func(cx, cy int)
	OnBeginDrag func(cx, cy int, payload any)
	DrawCell    func(dst *ebiten.Image, cx, cy int, x, y, size int, payload any)
}

func MakeGridField(x, y, cols, rows, cell int) GridField {
	return GridField{
		X:    x,
		Y:    y,
		Rows: rows,
		Cols: cols,
		Cell: cell,
	}
}

func (gf *GridField) Bounds() (x, y, w, h int) {
	return gf.X, gf.Y, gf.Cols * gf.Cell, gf.Rows * gf.Cell
}

func (gf *GridField) mouseToCell(mx, my int) (cx, cy int, ok bool) {
	mx -= gf.X
	my -= gf.Y
	if mx < 0 || my < 0 {
		return 0, 0, false
	}
	cx = mx / gf.Cell
	cy = my / gf.Cell
	if cx < 0 || cx >= gf.Cols || cy < 0 || cy >= gf.Rows {
		return 0, 0, false
	}
	return cx, cy, true
}
func (gf *GridField) Update(in core.Input) {
	if !in.LeftClicked {
		return
	}

	cx, cy, ok := gf.mouseToCell(in.MX, in.MY)
	if !ok {
		return
	}

	if gf.OnCellClick != nil {
		gf.OnCellClick(cx, cy)
	}

	if gf.OnBeginDrag != nil && gf.Get != nil {
		payload := gf.Get(cx, cy)
		if payload != nil {
			gf.OnBeginDrag(cx, cy, payload)
		}
	}
}

func (gf *GridField) Draw(dst *ebiten.Image) {
	// optional background
	//ebitenutil.DrawRect(dst, float64(gf.X), float64(gf.Y), float64(gf.Cols*gf.Cell), float64(gf.Rows*gf.Cell), gf.BG)

	if gf.ShowGrid {
		gf.drawGrid(dst)
	}

	if gf.Get == nil || gf.DrawCell == nil {
		return
	}

	for cy := 0; cy < gf.Rows; cy++ {
		for cx := 0; cx < gf.Cols; cx++ {
			payload := gf.Get(cx, cy)
			if payload == nil {
				continue
			}
			px := gf.X + cx*gf.Cell
			py := gf.Y + cy*gf.Cell
			gf.DrawCell(dst, cx, cy, px, py, gf.Cell, payload)
		}
	}
}

func (gf *GridField) drawGrid(dst *ebiten.Image) {
	if gf.ShowGrid {
		line := color.RGBA{60, 60, 70, 255}
		w := gf.Cols * gf.Cell
		h := gf.Rows * gf.Cell
		for x := 0; x <= gf.Cols; x++ {
			px := float64(gf.X + x*gf.Cell)
			ebitenutil.DrawRect(dst, px, float64(gf.Y), 1, float64(h), line)
		}
		for y := 0; y <= gf.Rows; y++ {
			py := float64(gf.Y + y*gf.Cell)
			ebitenutil.DrawRect(dst, float64(gf.X), py, float64(w), 1, line)
		}
	}
}

func (b *GridField) SetPos(x, y int) {
	b.X = x
	b.Y = y
}

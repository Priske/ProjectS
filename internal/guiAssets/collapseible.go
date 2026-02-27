package guiassets

import (
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Collapsible struct {
	X, Y, W, H int
	Title      string

	Open    bool
	Hovered bool

	Padding int // inner padding left/right
	Gap     int // vertical gap between children

	Children []core.Widget
}

func MakeCollapsible(x, y, w, h int, title string, children []core.Widget) core.Widget {
	return &Collapsible{
		X: x, Y: y, W: w, H: h,
		Title:    title,
		Padding:  10,
		Gap:      8,
		Children: children,
	}
}

func (c *Collapsible) Bounds() (int, int, int, int) {
	// Bounds represents the HEADER only for hit-testing.
	return c.X, c.Y, c.W, c.H
}

func (c *Collapsible) SetPos(x, y int) { // optional, but useful if parent wants to move it
	c.X = x
	c.Y = y
}

func (c *Collapsible) Update(in core.Input) {
	c.Hovered = core.PointInBounds(in.MX, in.MY, c)

	if in.LeftClicked && c.Hovered {
		c.Open = !c.Open
		return
	}

	if !c.Open {
		return
	}

	c.layoutChildren()

	for _, child := range c.Children {
		child.Update(in)
	}
}

func (c *Collapsible) Draw(dst *ebiten.Image) {
	// Header
	headerCol := color.RGBA{60, 60, 70, 255}
	if c.Hovered {
		headerCol = color.RGBA{90, 90, 110, 255}
	}
	ebitenutil.DrawRect(dst, float64(c.X), float64(c.Y), float64(c.W), float64(c.H), headerCol)

	// Title + indicator
	indicator := "▸"
	if c.Open {
		indicator = "▾"
	}
	ebitenutil.DebugPrintAt(dst, indicator+" "+c.Title, c.X+10, c.Y+18)

	if !c.Open {
		return
	}

	// Children
	c.layoutChildren()
	for _, child := range c.Children {
		child.Draw(dst)
	}
}

func (c *Collapsible) layoutChildren() {
	// Stack children below header.
	y := c.Y + c.H + c.Gap
	x := c.X + c.Padding

	for _, child := range c.Children {
		// Set child position if supported
		if p, ok := child.(core.Positionable); ok {
			p.SetPos(x, y)
		}

		_, _, _, h := child.Bounds()
		y += h + c.Gap
	}
}

// Optional helper if you want to know total height when open
func (c *Collapsible) TotalHeight() int {
	if !c.Open {
		return c.H
	}
	total := c.H + c.Gap
	for _, ch := range c.Children {
		_, _, _, h := ch.Bounds()
		total += h + c.Gap
	}
	return total
}

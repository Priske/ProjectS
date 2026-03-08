package screens

import (
	"image/color"

	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type DynamicFormationSection struct {
	X, Y, W, H int

	Title   string
	Open    bool
	Hovered bool

	Padding int
	Gap     int

	ps *PlayScreen
	g  core.Game

	cachedChildren []core.Widget
	lastCount      int
}

func (s *DynamicFormationSection) Bounds() (int, int, int, int) {
	return s.X, s.Y, s.W, s.TotalHeight()
}

func (s *DynamicFormationSection) SetPos(x, y int) {
	s.X = x
	s.Y = y
}

func (s *DynamicFormationSection) headerHit(mx, my int) bool {
	return mx >= s.X && mx < s.X+s.W && my >= s.Y && my < s.Y+s.H
}

func (s *DynamicFormationSection) refreshChildrenIfNeeded() {
	count := len(s.g.LocalPlayer().Formations)

	if s.cachedChildren != nil && s.lastCount == count {
		return
	}

	widgets := []core.Widget{}

	createFormationBtn := GUI.MakeButton(0, 0, 240, 50, "Create formation", func() {
		s.ps.resetFormationDraft()
		s.ps.openFormationEditorModal(s.g)
	})
	widgets = append(widgets, createFormationBtn)

	widgets = append(widgets, s.ps.makeFormationListWidget(s.g)...)

	s.cachedChildren = widgets
	s.lastCount = count
}

func (s *DynamicFormationSection) layoutChildren(children []core.Widget) {
	y := s.Y + s.H + s.Gap
	x := s.X + s.Padding

	for _, child := range children {
		if p, ok := child.(core.Positionable); ok {
			p.SetPos(x, y)
		}
		_, _, _, h := child.Bounds()
		y += h + s.Gap
	}
}
func (s *DynamicFormationSection) Update(in core.Input) {
	s.Hovered = s.headerHit(in.MX, in.MY)
	if in.LeftClicked && s.Hovered {
		s.Open = !s.Open
		return
	}
	if !s.Open {
		return
	}

	s.refreshChildrenIfNeeded()
	s.layoutChildren(s.cachedChildren)

	for _, child := range s.cachedChildren {
		child.Update(in)
	}
}
func (s *DynamicFormationSection) Draw(dst *ebiten.Image) {
	headerCol := color.RGBA{60, 60, 70, 255}
	if s.Hovered {
		headerCol = color.RGBA{90, 90, 110, 255}
	}

	ebitenutil.DrawRect(dst, float64(s.X), float64(s.Y), float64(s.W), float64(s.H), headerCol)

	indicator := "▸"
	if s.Open {
		indicator = "▾"
	}
	ebitenutil.DebugPrintAt(dst, indicator+" "+s.Title, s.X+10, s.Y+18)

	if !s.Open {
		return
	}

	s.refreshChildrenIfNeeded()
	s.layoutChildren(s.cachedChildren)

	for _, child := range s.cachedChildren {
		child.Draw(dst)
	}
}
func (s *DynamicFormationSection) TotalHeight() int {
	if !s.Open {
		return s.H
	}

	s.refreshChildrenIfNeeded()

	total := s.H + s.Gap
	for _, ch := range s.cachedChildren {
		_, _, _, h := ch.Bounds()
		total += h + s.Gap
	}
	return total
}
func (ps *PlayScreen) makeFormationSection(g core.Game) core.Widget {
	return &DynamicFormationSection{
		X:       0,
		Y:       0,
		W:       240,
		H:       50,
		Title:   "Formations",
		Padding: 10,
		Gap:     8,
		ps:      ps,
		g:       g,
	}
}

func (ps *PlayScreen) makeFormationListWidget(g core.Game) []core.Widget {
	widgets := []core.Widget{}

	for i := range g.LocalPlayer().Formations {
		index := i
		f := g.LocalPlayer().Formations[index]

		name := f.Name
		if name == "" {
			name = "Unnamed"
		}

		btn := GUI.MakeButton(0, 0, 240, 44, name, func() {
			ps.drag.Active = true
			ps.drag.Source = interaction.DragFromFormation
			ps.drag.Payload = &g.LocalPlayer().Formations[index]
			ps.drag.GrabOffX = 0
			ps.drag.GrabOffY = 0
		})

		hover := &GUI.HoverPopup{
			Target:        btn,
			OffsetX:       8,
			OffsetY:       0,
			PopupW:        3*32 + 16,
			PopupH:        5*32 + 16,
			ClampToScreen: true,
			DrawPopup: func(dst *ebiten.Image, px, py int) {
				drawFormationPreview(dst, g, index, px+8, py+8, 32)
			},
		}

		widgets = append(widgets, btn, hover)
	}

	return widgets
}

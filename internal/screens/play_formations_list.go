package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

type List struct {
	g  core.Game
	ps PlayScreen

	x, y int
	w, h int

	rowHeight  int
	hoverIndex int
}

func (ps *PlayScreen) makeFormationSection(g core.Game) core.Widget {
	widgets := []core.Widget{}

	createFormationBtn := GUI.MakeButton(0, 0, 240, 50, "Create formation", func() {
		ps.resetFormationDraft()
		ps.openFormationEditorModal(g)
	})
	widgets = append(widgets, createFormationBtn)
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

		const iconZone = 40

		hover := &GUI.HoverPopup{
			Target:        btn,
			OffsetX:       8,
			OffsetY:       0,
			PopupW:        3*32 + 16,
			PopupH:        5*32 + 16,
			ClampToScreen: true,
			DrawPopup: func(dst *ebiten.Image, px, py int) {
				ps.drawFormationPreview(dst, g, index, px+8, py+8, 32)
			},
		}
		widgets = append(widgets, btn, hover)
	}

	return GUI.MakeCollapsible(0, 0, 240, 50, "Formations . . . ", widgets)
}

package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

func (ps *PlayScreen) makeOptionsSidebar(g core.Game) core.Widget {
	panelW := 260
	headerH := 44
	x := 20
	y := 40

	placeUnit := ps.makePlaceUnitSection(g)
	formation := ps.makeFormationSection(g)
	exit := GUI.MakeButton(0, 0, 240, 50, "Save & Quit", func() {
		g.SetScreen(NewMenuScreen(g))
	})

	widgetsSidebar := []core.Widget{placeUnit, formation, exit}
	return GUI.MakeCollapsible(x, y, panelW, headerH, "Options . . .", widgetsSidebar)
}

func (ps *PlayScreen) makePlaceUnitSection(g core.Game) core.Widget {
	grid := ps.makeUnitsGrid(g) // returns core.Widget
	ps.reserveGrid = grid
	return GUI.MakeCollapsible(0, 0, 240, 50, "Place Unit . . .", []core.Widget{
		grid,
	})
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

func (ps *PlayScreen) rebuildLeftSidebar(g core.Game) {
	ps.widgets = nil
	ps.widgets = append(ps.widgets,
		ps.makeOptionsSidebar(g), // or whatever you call it
	)
}

func (ps *PlayScreen) makeUnitOptionsGrid(g core.Game, x, y int) *GUI.GridField {
	const columns = 3
	const rows = 4
	const cellSize = 48

	grid := GUI.MakeGridField(x, y, columns, rows, cellSize)
	grid.ShowGrid = true

	grid.Get = func(cx, cy int) any {
		index := cy*columns + cx
		if index < 0 || index >= len(ps.availableUnitTypesForCategory) {
			return nil
		}
		return ps.availableUnitTypesForCategory[index]
	}

	grid.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut, ok := payload.(core.UnitType)
		if !ok {
			return
		}
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, size)
	}

	// IMPORTANT: do NOT reference ps.modal here (it might still be nil while building)
	grid.OnBeginDrag = func(cx, cy int, payload any) {
		ut, ok := payload.(core.UnitType)
		if !ok {
			return
		}

		ps.drag.Active = true
		ps.drag.Source = interaction.DragFromFormationPalette
		ps.drag.Payload = ut

		// center sprite under cursor
		ps.drag.GrabOffX = grid.Cell / 2
		ps.drag.GrabOffY = grid.Cell / 2

		// IMPORTANT: set mouse immediately so it draws same frame

		ps.drag.MX = g.Input().MX
		ps.drag.MY = g.Input().MY
	}
	return grid
}

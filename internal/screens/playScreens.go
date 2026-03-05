package screens

import (
	"github.com/Priske/ProjectS/interaction"
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayScreen struct {
	widgets                []core.Widget
	drag                   interaction.DragState
	lastDrop               string
	modal                  *GUI.Modal
	unPlacedUnits          []*core.Unit
	readyAdded             bool
	setupMode              bool
	reserveGrid            *GUI.GridField
	readyWidget            core.Widget
	formationGrid          *GUI.GridField
	unitOptionsGrid        *GUI.GridField
	nameFormationTextField *GUI.TextField

	formationWants                map[core.Pos]core.UnitType
	selectedUnitCategory          core.UnitCategory
	availableUnitTypesForCategory []core.UnitType
	formationBrushUnitType        core.UnitType
}

func (ps *PlayScreen) Update(g core.Game) error {
	input := g.Input()
	if ps.drag.Active {
		ps.drag.MX = input.MX
		ps.drag.MY = input.MY
	}

	if len(ps.unPlacedUnits) == 0 && !ps.readyAdded {
		ps.readyWidget = ps.makeReadyButton(g)
		ps.widgets = append(ps.widgets, ps.readyWidget)
		ps.readyAdded = true
	}

	if len(ps.unPlacedUnits) != 0 && ps.readyAdded {
		ps.widgets = removeWidget(ps.widgets, ps.readyWidget)
		ps.readyWidget = nil
		ps.readyAdded = false
	}
	// If modal open: update only it (and return)
	//Testing Modal

	if ps.modal != nil && ps.modal.Open {
		ps.modal.Update(input)

		if ps.drag.Active && !input.LeftPressed {
			ps.tryDropIntoFormation(input.MX, input.MY)
		}
		return nil
	}

	for _, w := range ps.widgets {
		w.Update(input)
	}

	// keep drag cursor fresh

	// Start drag on click
	if !ps.drag.Active && input.LeftClicked && !ps.clickHitsWidget(input.MX, input.MY) {
		cx, cy, ok := ps.mouseToCell(g, input.MX, input.MY)
		if ok {
			tile := g.Board().Location[cy][cx]
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

func (ps *PlayScreen) clickHitsWidget(mx, my int) bool {
	for _, w := range ps.widgets {
		x, y, ww, hh := w.Bounds()
		if mx >= x && mx < x+ww && my >= y && my < y+hh {
			return true
		}
	}
	return false
}

// Named returns
func (ps *PlayScreen) boardGeom(g core.Game) (offX, offY, w, h int, s core.Settings) {
	s = g.Settings()
	offX, offY = getOffXY(g)
	w = s.BoardW * s.CellSize
	h = s.BoardH * s.CellSize
	return
}
func pointInRect(mx, my, x, y, w, h int) bool {
	return mx >= x && mx < x+w && my >= y && my < y+h
}

func (ps *PlayScreen) mouseOverReserve(mx, my int) bool {
	if ps.reserveGrid == nil {
		return false
	}
	x, y, w, h := ps.reserveGrid.Bounds()
	return pointInRect(mx, my, x, y, w, h)
}
func removeWidget(widgets []core.Widget, target core.Widget) []core.Widget {
	for i := range widgets {
		if widgets[i] == target {
			return append(widgets[:i], widgets[i+1:]...)
		}
	}
	return widgets
}

func (ps *PlayScreen) makeCategoryBar(g core.Game, x, y int) *GUI.GridField {
	categories := []core.UnitCategory{core.Attack, core.Defense, core.Support}

	gf := GUI.MakeGridField(x, y, len(categories), 1, 48)
	gf.ShowGrid = true

	gf.Get = func(cx, cy int) any { return categories[cx] }

	gf.OnCellClick = func(cx, cy int) {
		cat := categories[cx]
		ps.openUnitPickerModal(g, cat) // opens popup of unit types for that category
	}

	gf.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		// draw icon/text for category here
	}

	return gf
}
func (ps *PlayScreen) openUnitPickerModal(g core.Game, category core.UnitCategory) {
	types := unitTypesFor(category)
	cols := 4
	rows := (len(types) + cols - 1) / cols

	pw, ph := 320, 240
	px, py := (core.VirtualW-pw)/2, (core.VirtualH-ph)/2

	picker := GUI.MakeGridField(px+16, py+16, cols, rows, 48)
	picker.ShowGrid = true

	picker.Get = func(cx, cy int) any {
		i := cy*cols + cx
		if i < 0 || i >= len(types) {
			return nil
		}
		return types[i]
	}

	picker.OnCellClick = func(cx, cy int) {
		i := cy*cols + cx
		if i < 0 || i >= len(types) {
			return
		}
		ps.formationBrushUnitType = types[i] // set brush
		if ps.modal != nil {
			ps.modal.Close()
		}
	}

	picker.DrawCell = func(dst *ebiten.Image, cx, cy, px, py, size int, payload any) {
		ut := payload.(core.UnitType)
		ps.drawUnitImage(dst, g.Assets(), ut, px, py, 60)
	}

	closeBtn := GUI.MakeButton(px+pw-110, py+ph-60, 90, 40, "Close", func() {
		if ps.modal != nil {
			ps.modal.Close()
		}
	})

	ps.modal = GUI.MakeModal(px, py, pw, ph, []core.Widget{picker, closeBtn})
}

func (ps *PlayScreen) openNameFormationPopup(mx, my int) {
	ps.nameFormationPopupOpen = true

	if ps.nameFormationTextField == nil {
		tf := GUI.MakeTextField(mx+20, my+40, 260, 36) // adjust ctor to your API
		ps.nameFormationTextField = tf
	} else {
		ps.nameFormationTextField.SetPos(mx+20, my+40) // if you have SetPos
		ps.nameFormationTextField.SetText("")          // clear if you have it
	}

	ps.nameFormationConfirmBtn = GUI.MakeButton(mx+20, my+90, 120, 36, "Confirm", func() {
		ps.nameFormationPopupOpen = false
	})

	ps.nameFormationCancelBtn = GUI.MakeButton(mx+160, my+90, 120, 36, "Cancel", func() {
		ps.nameFormationPopupOpen = false
	})
}

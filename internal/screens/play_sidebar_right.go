package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
	"github.com/hajimehoshi/ebiten/v2"
)

func (ps *PlayScreen) makeRightSidebarSetup(g core.Game, width int) core.Widget {
	return GUI.MakeButton(0, 0, width, 50, "reset", func() {
		ps.resetSetupState(g)
	})
}
func (ps *PlayScreen) makeReadyButton(g core.Game, width int) core.Widget {
	return GUI.MakeButton(0, 0, width, 50, "Ready", func() {
		ps.confirmSetup(g)
	})
}
func (ps *PlayScreen) makeBattleRightSidebar(g core.Game) core.Widget {
	offX, offY, boardW, boardH, _ := boardGeom(g)

	sidebarGap := 20
	sidebarX := offX + boardW + sidebarGap
	sidebarY := offY
	sidebarW := 240
	sidebarH := boardH

	padding := 10
	gap := 10

	innerW := sidebarW - padding*2
	innerH := sidebarH - padding*2 - gap
	halfH := innerH / 2

	logPanel := GUI.MakePanel(0, 0, innerW, halfH, "Battle Log", []core.Widget{
		&BattleLogWidget{ps: ps},
	})

	infoPanel := GUI.MakePanel(0, 0, innerW, innerH-halfH, "Info", []core.Widget{})
	infoPanel.AutoLayout = false

	infoWidget := &SelectedUnitInfoWidget{ps: ps}
	endTurnButton := GUI.MakeButton(0, 0, 120, 32, "End Turn", func() {
		ps.endPlayerTurn(g)
	})

	infoPanel.Children = []core.Widget{
		infoWidget,
		endTurnButton,
	}

	sidebar := GUI.MakePanel(sidebarX, sidebarY, sidebarW, sidebarH, "", []core.Widget{
		logPanel,
		infoPanel,
	})

	// position children inside infoPanel manually
	infoPanel.SetPos(sidebarX+padding, sidebarY+padding+halfH+gap)

	infoWidget.SetPos(infoPanel.X+10, infoPanel.Y+20)

	buttonX := infoPanel.X + (infoPanel.W-120)/2
	buttonY := infoPanel.Y + infoPanel.H - 32 - 12
	endTurnButton.(core.Positionable).SetPos(buttonX, buttonY)

	return sidebar
}

func (ps *PlayScreen) makeSetupRightSidebar(g core.Game) core.Widget {
	offX, _, boardW, _, _ := boardGeom(g)

	sidebarGap := 20
	sidebarX := offX + boardW + sidebarGap
	sidebarY := 40
	sidebarW := 260
	sidebarH := 640

	padding := 10
	gap := 10

	innerW := sidebarW - padding*2
	if innerW < 0 {
		innerW = 0
	}

	actionsH := 180
	inventoryH := sidebarH - padding*2 - gap - actionsH
	if inventoryH < 0 {
		inventoryH = 0
	}

	actionsPanel := GUI.MakePanel(0, 0, innerW, actionsH, "Actions", []core.Widget{
		ps.makeRightSidebarSetup(g, innerW-20),
	})
	actionsPanel.AutoLayout = true

	if len(ps.setup.unPlacedUnits) == 0 {
		actionsPanel.Children = append(actionsPanel.Children, ps.makeReadyButton(g, innerW-20))
	}

	inventoryPanel := GUI.MakePanel(0, 0, innerW, inventoryH, "Inventory", []core.Widget{})
	inventoryPanel.AutoLayout = false

	actionsPanel.SetPos(sidebarX+padding, sidebarY+padding)
	inventoryPanel.SetPos(sidebarX+padding, sidebarY+padding+actionsH+gap)

	contentX := inventoryPanel.X + 14
	contentY := inventoryPanel.Y + 40
	contentW := inventoryPanel.W - 28
	if contentW < 0 {
		contentW = 0
	}

	selectedLabel := &SetupSelectedLabelWidget{
		ps: ps,
		X:  contentX,
		Y:  contentY,
	}

	slotGrid := ps.makeSetupInventoryGrid(g, 0, 0)
	actionGrid := ps.makeSetupInventoryActionGrid(g, 0, 0)

	slotCell := 64
	slotGridW := 3 * slotCell

	actionCell := 112
	actionGridW := 2 * actionCell
	actionGridH := actionCell

	slotGridX := contentX + (contentW-slotGridW)/2
	slotGridY := contentY + 26

	actionGridX := contentX + (contentW-actionGridW)/2
	actionGridY := inventoryPanel.Y + inventoryPanel.H - actionGridH - 28

	if p, ok := slotGrid.(core.Positionable); ok {
		p.SetPos(slotGridX, slotGridY)
	}
	if p, ok := actionGrid.(core.Positionable); ok {
		p.SetPos(actionGridX, actionGridY)
	}

	inventoryPanel.Children = []core.Widget{
		selectedLabel,
		slotGrid,
		actionGrid,
	}

	sidebar := GUI.MakePanel(sidebarX, sidebarY, sidebarW, sidebarH, "", []core.Widget{
		actionsPanel,
		inventoryPanel,
	})
	sidebar.AutoLayout = false

	return sidebar
}
func (ps *PlayScreen) makeSetupInventoryActionGrid(g core.Game, x, y int) core.Widget {
	grid := GUI.MakeGridField(x, y, 2, 1, 112)
	actions := []string{"chest", "shop"}

	grid.ShowGrid = false

	grid.Get = func(cx, cy int) any {
		i := cy*2 + cx
		if i < 0 || i >= len(actions) {
			return nil
		}
		return actions[i]
	}

	grid.DrawCell = func(dst *ebiten.Image, cx, cy int, px, py, size int, payload any) {
		action, ok := payload.(string)
		if !ok {
			return
		}

		switch action {
		case "chest":
			drawInventoryActionButton(dst, g.Assets().ChestButtonIcon, px, py, size)
		case "shop":
			drawInventoryActionButton(dst, g.Assets().ShopButtonIcon, px, py, size)
		}
	}

	grid.OnCellClick = func(cx, cy int) {
		switch cy*2 + cx {
		case 0:
			// chest later
		case 1:
			// shop later
		}
	}

	return grid
}

type invSlot struct {
	Category core.ItemCategory
}

func (ps *PlayScreen) makeSetupInventoryGrid(g core.Game, x, y int) core.Widget {
	cell := 68

	equipmentLayout := [][]core.EquipmentSlot{
		{core.SlotCharm, core.SlotHead, core.SlotBag},
		{core.SlotWeapon1, core.SlotArmor, core.SlotWeapon2},
		{core.SlotAmmo1, core.SlotLegs, core.SlotAmmo2},
	}

	cols := 3
	rows := len(equipmentLayout)

	grid := GUI.MakeGridField(x, y, cols, rows, cell)
	grid.ShowGrid = false

	grid.Get = func(cx, cy int) any {
		if cy < 0 || cy >= len(equipmentLayout) {
			return nil
		}

		slot := equipmentLayout[cy][cx]

		u := ps.setup.Selected
		if u == nil {
			return slot
		}

		if slot == core.SlotWeapon1 || slot == core.SlotWeapon2 {
			item1 := u.Equipped[core.SlotWeapon1]
			item2 := u.Equipped[core.SlotWeapon2]

			if item1 != nil && item1 == item2 {
				return item1
			}

			if item, ok := u.Equipped[slot]; ok {
				return item
			}

			return slot
		}

		if item, ok := u.Equipped[slot]; ok {
			return item
		}

		return slot
	}

	grid.DrawCell = func(dst *ebiten.Image, cx, cy int, px, py, size int, payload any) {
		assets := g.Assets()

		if assets.FrameTemplate != nil {
			sw, sh := assets.FrameTemplate.Bounds().Dx(), assets.FrameTemplate.Bounds().Dy()
			if sw > 0 && sh > 0 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(float64(size)/float64(sw), float64(size)/float64(sh))
				op.GeoM.Translate(float64(px), float64(py))
				dst.DrawImage(assets.FrameTemplate, op)
			}
		}

		if item, ok := payload.(core.Item); ok && item != nil {
			if icon := assets.ItemIcons[item.Base().ID]; icon != nil {
				padding := size / 10
				drawImageCentered(
					dst,
					icon,
					px+padding,
					py+padding,
					size-(padding*2),
					size-(padding*2),
				)
				return
			}
		}

		if slot, ok := payload.(core.EquipmentSlot); ok {
			cat := slotToCategory(slot)

			icon := assets.SlotIcons[cat]
			if icon == nil {
				return
			}

			padding := size / 12
			drawImageCentered(
				dst,
				icon,
				px+padding,
				py+padding,
				size-(padding*2),
				size-(padding*2),
			)
		}
	}

	return grid
}
func slotToCategory(slot core.EquipmentSlot) core.ItemCategory {
	switch slot {
	case core.SlotWeapon1, core.SlotWeapon2:
		return core.CategoryWeapon
	case core.SlotArmor:
		return core.CategoryArmor
	case core.SlotHead:
		return core.CategoryHead
	case core.SlotLegs:
		return core.CategoryLegs
	case core.SlotCharm:
		return core.CategoryCharm
	case core.SlotBag:
		return core.CategoryBag
	case core.SlotAmmo1, core.SlotAmmo2, core.SlotAmmo3:
		return core.CategoryAmmo
	}
	return core.CategoryWeapon
}
func (ps *PlayScreen) drawInventoryCell(dst *ebiten.Image, g core.Game, cx, cy int, px, py, size int, payload any) {
	assets := g.Assets()

	if assets.FrameTemplate != nil {
		sw, sh := assets.FrameTemplate.Bounds().Dx(), assets.FrameTemplate.Bounds().Dy()
		if sw > 0 && sh > 0 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(size)/float64(sw), float64(size)/float64(sh))
			op.GeoM.Translate(float64(px), float64(py))
			dst.DrawImage(assets.FrameTemplate, op)
		}
	}

	// Real item icon
	if item, ok := payload.(core.Item); ok && item != nil {
		if icon := assets.ItemIcons[item.Base().ID]; icon != nil {
			padding := size / 10
			drawImageCentered(
				dst,
				icon,
				px+padding,
				py+padding,
				size-(padding*2),
				size-(padding*2),
			)
			return
		}
	}

	// Empty-slot placeholder by layout position
	var cat core.ItemCategory

	if cy == 0 {
		switch cx {
		case 0:
			cat = core.CategoryWeapon
		case 1:
			cat = core.CategoryArmor
		case 2:
			cat = core.CategoryWeapon // change later if you add a true 2nd weapon / carry / accessory slot
		}
	} else if cy == 1 {
		cat = core.CategoryAccessory
	} else {
		cat = core.CategoryAmmo
	}

	icon := assets.SlotIcons[cat]
	if icon == nil {
		return
	}

	padding := size / 12
	drawImageCentered(
		dst,
		icon,
		px+padding,
		py+padding,
		size-(padding*2),
		size-(padding*2),
	)
}

func (ps *PlayScreen) makeSlot(cat core.ItemCategory) core.Widget {
	return GUI.MakeButton(0, 0, 48, 48, "", nil) // temp

	// later:
	// custom widget that draws:
	// - frame_template
	// - icon
}

func (ps *PlayScreen) makeSetupInventoryPanel(g core.Game) core.Widget {
	x := 999
	y := 230
	w := 240
	h := 480

	children := []core.Widget{
		GUI.MakeLabel(x+12, y+24, "Selected: None"),
	}

	return GUI.MakePanel(x, y, w, h, "Inventory", children)
}
func (ps *PlayScreen) makeSetupActionsPanel(g core.Game) core.Widget {
	x := 999
	y := 40
	w := 240
	h := 180

	padding := 10
	contentW := w - padding*2
	if contentW < 0 {
		contentW = 0
	}

	children := []core.Widget{
		ps.makeRightSidebarSetup(g, contentW),
	}

	if len(ps.setup.unPlacedUnits) == 0 {
		children = append(children, ps.makeReadyButton(g, contentW))
	}

	return GUI.MakePanel(x, y, w, h, "Actions", children)
}

func (ps *PlayScreen) makeSetupRightPanel(g core.Game) core.Widget {
	x := 999
	y := 40
	w := 240
	h := 410

	padding := 10
	contentW := w - padding*2
	if contentW < 0 {
		contentW = 0
	}

	children := []core.Widget{
		ps.makeRightSidebarSetup(g, contentW),
	}

	if len(ps.setup.unPlacedUnits) == 0 {
		ready := ps.makeReadyButton(g, contentW)
		ps.setup.readyWidget = ready
		ps.setup.readyAdded = true
		children = append(children, ready)
	} else {
		ps.setup.readyWidget = nil
		ps.setup.readyAdded = false
	}

	return GUI.MakePanel(x, y, w, h, "Actions", children)
}

package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
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
	inventoryPanel.AutoLayout = true

	// IMPORTANT: absolute screen coordinates, not local coordinates
	actionsPanel.SetPos(sidebarX+padding, sidebarY+padding)
	inventoryPanel.SetPos(sidebarX+padding, sidebarY+padding+actionsH+gap)

	sidebar := GUI.MakePanel(sidebarX, sidebarY, sidebarW, sidebarH, "", []core.Widget{
		actionsPanel,
		inventoryPanel,
	})
	sidebar.AutoLayout = false

	return sidebar
}
func (ps *PlayScreen) makeSetupInventoryPanel(g core.Game) core.Widget {
	x := 999
	y := 230
	w := 240
	h := 480

	children := []core.Widget{
		GUI.MakeLabel(0, 0, "Placed units: "+itoa(countPlacedPlayerUnits(g))),
		GUI.MakeLabel(0, 0, "Unplaced units: "+itoa(len(ps.setup.unPlacedUnits))),
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

package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeRightSidebarSetup(g core.Game) core.Widget {

	x := 999
	y := 40
	reset := GUI.MakeButton(x, y, 240, 50, "reset", func() {
		ps.resetSetupState(g)
	})
	return reset

}

func (ps *PlayScreen) makeReadyButton(g core.Game) core.Widget {
	x := 999
	y := 400
	ready := GUI.MakeButton(x, y, 240, 50, "Ready", func() {
		ps.confirmSetup(g)
	})
	return ready
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
	infoPanel := GUI.MakePanel(0, 0, innerW, innerH-halfH, "Info", []core.Widget{
		&SelectedUnitInfoWidget{ps: ps},
	})

	sidebar := GUI.MakePanel(sidebarX, sidebarY, sidebarW, sidebarH, "", []core.Widget{
		logPanel,
		infoPanel,
	})

	return sidebar
}

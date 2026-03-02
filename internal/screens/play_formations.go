package screens

import (
	"github.com/Priske/ProjectS/internal/core"
	GUI "github.com/Priske/ProjectS/internal/guiAssets"
)

func (ps *PlayScreen) makeFormationSection(g core.Game) core.Widget {
	modalTest := GUI.MakeButton(0, 0, 240, 50, "FormationA", func() {
		mx, my := 200, 120
		mw, mh := 420, 360

		closeBtn := GUI.MakeButton(mx+mw-110, my+mh-60, 90, 40, "Close", func() {
			if ps.modal != nil {
				ps.modal.Close()
			}
		})

		ps.modal = GUI.MakeModal(mx, my, mw, mh, []core.Widget{closeBtn})
	})

	return GUI.MakeCollapsible(0, 0, 240, 50, "Formations . . . ", []core.Widget{modalTest})
}

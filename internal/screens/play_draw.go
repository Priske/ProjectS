package screens

import (
	"fmt"
	"image/color"

	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func drawGrid(screen *ebiten.Image, offX, offY, boardW, boardH, cellSize int) {
	w := boardW * cellSize
	h := boardH * cellSize
	line := color.RGBA{60, 60, 70, 255}

	for x := 0; x <= boardW; x++ {
		px := float64(offX + x*cellSize)
		ebitenutil.DrawRect(screen, px, float64(offY), 1, float64(h), line)
	}
	for y := 0; y <= boardH; y++ {
		py := float64(offY + y*cellSize)
		ebitenutil.DrawRect(screen, float64(offX), py, float64(w), 1, line)
	}
}

func drawPlacementZone(screen *ebiten.Image, offX, offY, boardH, cellSize int, cols int) {
	// translucent green fill
	fill := color.RGBA{0, 200, 0, 70}

	for y := 0; y < boardH; y++ {
		for x := 0; x < cols; x++ {
			px := float64(offX + x*cellSize)
			py := float64(offY + y*cellSize)
			ebitenutil.DrawRect(screen, px, py, float64(cellSize), float64(cellSize), fill)
		}
	}
}

func (ps *PlayScreen) drawUnits(g core.Game, screen *ebiten.Image) {
	offX, offY, _, _, s := boardGeom(g)
	assets := g.Assets()

	for y := 0; y < s.BoardH; y++ {
		for x := 0; x < s.BoardW; x++ {
			if ps.drag.Active && x == ps.drag.FromX && y == ps.drag.FromY {
				continue
			}
			tile := g.Board().Location[y][x]
			if tile.Unit == nil {
				continue
			}
			drawUnitImage(screen, assets, tile.Unit.Type, offX+x*s.CellSize, offY+y*s.CellSize, s.CellSize)
		}
	}
}

func drawUnitImage(screen *ebiten.Image, assets core.Assets, unitType core.UnitType, px, py, cellSize int) {
	img := assets.UnitImages[unitType]
	if img == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	op.GeoM.Scale(float64(cellSize)/float64(sw), float64(cellSize)/float64(sh))
	op.GeoM.Translate(float64(px), float64(py))
	screen.DrawImage(img, op)
}
func drawCategoryImage(screen *ebiten.Image, assets core.Assets, categoryType core.UnitCategory, px, py, cellSize int) {
	img := assets.CategoryImage(categoryType)
	if img == nil {
		return
	}
	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(cellSize)/float64(sw), float64(cellSize)/float64(sh))
	op.GeoM.Translate(float64(px), float64(py))
	screen.DrawImage(img, op)
}

func (ps *PlayScreen) drawDraggedUnit(g core.Game, screen *ebiten.Image) {
	if !ps.drag.Active || ps.drag.Payload == nil {
		return
	}

	s := g.Settings()
	assets := g.Assets()

	// 1) UnitType drag
	if ut, ok := ps.drag.Payload.(core.UnitType); ok {
		img := assets.UnitImages[ut]
		if img == nil {
			return
		}
		op := &ebiten.DrawImageOptions{}
		sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
		op.GeoM.Scale(float64(s.CellSize)/float64(sw), float64(s.CellSize)/float64(sh))
		op.GeoM.Translate(float64(ps.drag.MX-ps.drag.GrabOffX), float64(ps.drag.MY-ps.drag.GrabOffY))
		screen.DrawImage(img, op)
		return
	}
	//Draw formation drag
	if f, ok := ps.drag.Payload.(*core.Formation); ok {

		cx, cy, ok := mouseToCell(g, ps.drag.MX, ps.drag.MY)
		if !ok {
			return
		}

		offX, offY := getOffXY(g)
		s := g.Settings()

		for pos, ut := range f.Wants {

			px := offX + (cx+pos.X)*s.CellSize
			py := offY + (cy+pos.Y)*s.CellSize

			drawUnitImage(screen, g.Assets(), ut, px, py, s.CellSize)
		}

		return
	}
	// 2) *Unit drag
	u, ok := ps.drag.Payload.(*core.Unit)
	if !ok || u == nil {
		return
	}
	img := assets.UnitImages[u.Type]
	if img == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}
	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	op.GeoM.Scale(float64(s.CellSize)/float64(sw), float64(s.CellSize)/float64(sh))
	op.GeoM.Translate(float64(ps.drag.MX-ps.drag.GrabOffX), float64(ps.drag.MY-ps.drag.GrabOffY))
	screen.DrawImage(img, op)
}

func (ps *PlayScreen) drawUI(screen *ebiten.Image) {
	for _, b := range ps.ui.widgets {
		b.Draw(screen)
	}
}

func (ps *PlayScreen) drawDebug(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Drop: "+ps.ui.lastDrop)
}

func (ps *PlayScreen) drawBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 24, 255})
}

func (ps *PlayScreen) drawBoard(g core.Game, screen *ebiten.Image) {
	s := g.Settings()
	offX, offY := getOffXY(g)
	ps.drawBoardTiles(g, screen, offX, offY, s.CellSize)
	if ps.setup.setupMode {
		drawPlacementZone(screen, offX, offY, s.BoardH, s.CellSize, 3)
	}

	drawGrid(screen, offX, offY, s.BoardW, s.BoardH, s.CellSize)
	ps.drawUnits(g, screen)
}
func (ps *PlayScreen) drawModal(screen *ebiten.Image) {
	if ps.ui.modal != nil && ps.ui.modal.Open {
		ps.ui.modal.Draw(screen)
	}
	if ps.ui.overlay != nil {
		ps.ui.overlay.Draw(screen)
	}
}
func (ps *PlayScreen) drawHoveredUnitInfo(g core.Game, screen *ebiten.Image) {
	if !ebiten.IsKeyPressed(ebiten.KeyControl) {
		return
	}
	if info, ok := ps.hoveredBoardUnitInfo(g); ok {
		drawUnitInfoCard(screen, g, info)
		return
	}

	if info, ok := ps.hoveredReserveUnitInfo(g); ok {
		drawUnitInfoCard(screen, g, info)
	}
}
func (ps *PlayScreen) drawBoardTiles(g core.Game, screen *ebiten.Image, offX, offY, cellSize int) {
	s := g.Settings()
	board := g.Board()
	assets := g.Assets()

	for y := 0; y < s.BoardH; y++ {
		for x := 0; x < s.BoardW; x++ {
			tile := board.Location[y][x]

			img := assets.LocationImage(tile.LocationType)
			if img == nil {
				continue
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(
				float64(cellSize)/float64(img.Bounds().Dx()),
				float64(cellSize)/float64(img.Bounds().Dy()),
			)
			op.GeoM.Translate(
				float64(offX+x*cellSize),
				float64(offY+y*cellSize),
			)

			screen.DrawImage(img, op)
		}
	}
}
func (ps *PlayScreen) drawBattleLog(screen *ebiten.Image, g core.Game) {
	if !ps.battle.Active {
		return
	}

	// location of log panel
	x := 1000 + 20
	y := 40

	for i, line := range ps.battle.Log {
		ebitenutil.DebugPrintAt(screen, line, x, y+i*14)
	}
}

func (ps *PlayScreen) drawSelectedUnitHighlight(g core.Game, screen *ebiten.Image) {
	if !ps.battle.Active || ps.battle.Selected == nil {
		return
	}

	px, py := cellTopLeft(g, ps.battle.SelectedX, ps.battle.SelectedY)
	size := g.Settings().CellSize

	ebitenutil.DrawRect(screen, float64(px), float64(py), float64(size), 2, color.RGBA{255, 255, 0, 255})
	ebitenutil.DrawRect(screen, float64(px), float64(py+size-2), float64(size), 2, color.RGBA{255, 255, 0, 255})
	ebitenutil.DrawRect(screen, float64(px), float64(py), 2, float64(size), color.RGBA{255, 255, 0, 255})
	ebitenutil.DrawRect(screen, float64(px+size-2), float64(py), 2, float64(size), color.RGBA{255, 255, 0, 255})
}
func (ps *PlayScreen) drawMoveRange(g core.Game, screen *ebiten.Image) {
	if !ps.battle.Active || ps.battle.Selected == nil {
		return
	}
	if ps.battle.Turn.Side != TurnPlayer {
		return
	}
	if !ps.battle.Turn.CanMove(ps.battle.Selected) {
		return
	}

	board := g.Board()
	size := g.Settings().CellSize

	x := ps.battle.SelectedX
	y := ps.battle.SelectedY
	r := ps.battle.Selected.MoveRange

	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			dist := abs(dx) + abs(dy)
			if dist == 0 || dist > r {
				continue
			}

			cx := x + dx
			cy := y + dy

			if cy < 0 || cy >= len(board.Location) {
				continue
			}
			if cx < 0 || cx >= len(board.Location[cy]) {
				continue
			}
			if board.Location[cy][cx].Unit != nil {
				continue
			}

			px, py := cellTopLeft(g, cx, cy)
			col := color.RGBA{0, 255, 0, 255}

			ebitenutil.DrawRect(screen, float64(px), float64(py), float64(size), 2, col)
			ebitenutil.DrawRect(screen, float64(px), float64(py+size-2), float64(size), 2, col)
			ebitenutil.DrawRect(screen, float64(px), float64(py), 2, float64(size), col)
			ebitenutil.DrawRect(screen, float64(px+size-2), float64(py), 2, float64(size), col)
		}
	}
}

func (ps *PlayScreen) drawActionMenu(g core.Game, screen *ebiten.Image) {
	fmt.Println("draw Action menu")
	if !ps.battle.Active || !ps.battle.ActionMenuOpen || ps.battle.Selected == nil {
		return
	}

	u := ps.battle.Selected
	x := ps.battle.ActionMenuX
	y := ps.battle.ActionMenuY
	w := 140
	rowH := 22
	h := len(u.Actions)*rowH + 8

	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), color.RGBA{35, 35, 45, 255})

	border := color.RGBA{90, 90, 110, 255}
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), 1, border)
	ebitenutil.DrawRect(screen, float64(x), float64(y+h-1), float64(w), 1, border)
	ebitenutil.DrawRect(screen, float64(x), float64(y), 1, float64(h), border)
	ebitenutil.DrawRect(screen, float64(x+w-1), float64(y), 1, float64(h), border)

	for i := range u.Actions {
		a := &u.Actions[i]
		rowY := y + 4 + i*rowH

		if ps.battle.SelectedAction != nil && ps.battle.SelectedAction.ID == a.ID {
			ebitenutil.DrawRect(screen, float64(x+2), float64(rowY-1), float64(w-4), float64(rowH-2), color.RGBA{70, 70, 100, 255})
		}

		label := a.Name
		if !ps.canUseAction(u, a) {
			label = "X " + label
		}

		ebitenutil.DebugPrintAt(screen, label, x+8, rowY+4)
	}
}

func (ps *PlayScreen) drawBattleActionCursor(g core.Game, screen *ebiten.Image) {
	if !ps.battle.Active || ps.battle.SelectedAction == nil {
		return
	}

	switch ps.battle.SelectedAction.Kind {
	case core.ActionAttack:
		ps.drawAttackCursor(g, screen)
	}
}

func (ps *PlayScreen) drawAttackCursor(g core.Game, screen *ebiten.Image) {
	if !ps.canAttackHoveredCell(g) {
		return
	}

	in := g.Input()
	cx, cy, ok := mouseToCell(g, in.MX, in.MY)
	if !ok {
		return
	}

	px, py := cellTopLeft(g, cx, cy)
	size := g.Settings().CellSize
	col := color.RGBA{255, 60, 60, 255}

	cornerLen := size / 4
	thick := 2

	// top-left corner
	ebitenutil.DrawRect(screen, float64(px), float64(py), float64(cornerLen), float64(thick), col)
	ebitenutil.DrawRect(screen, float64(px), float64(py), float64(thick), float64(cornerLen), col)

	// top-right corner
	ebitenutil.DrawRect(screen, float64(px+size-cornerLen), float64(py), float64(cornerLen), float64(thick), col)
	ebitenutil.DrawRect(screen, float64(px+size-thick), float64(py), float64(thick), float64(cornerLen), col)

	// bottom-left corner
	ebitenutil.DrawRect(screen, float64(px), float64(py+size-thick), float64(cornerLen), float64(thick), col)
	ebitenutil.DrawRect(screen, float64(px), float64(py+size-cornerLen), float64(thick), float64(cornerLen), col)

	// bottom-right corner
	ebitenutil.DrawRect(screen, float64(px+size-cornerLen), float64(py+size-thick), float64(cornerLen), float64(thick), col)
	ebitenutil.DrawRect(screen, float64(px+size-thick), float64(py+size-cornerLen), float64(thick), float64(cornerLen), col)

	// small center dot
	ebitenutil.DrawRect(
		screen,
		float64(px+size/2-1),
		float64(py+size/2-1),
		2,
		2,
		col,
	)
}

func (ps *PlayScreen) canAttackHoveredCell(g core.Game) bool {
	if ps.battle.Selected == nil || ps.battle.SelectedAction == nil {
		return false
	}

	in := g.Input()
	cx, cy, ok := mouseToCell(g, in.MX, in.MY)
	if !ok {
		return false
	}

	board := g.Board()
	if cy < 0 || cy >= len(board.Location) || cx < 0 || cx >= len(board.Location[cy]) {
		return false
	}

	u := board.Location[cy][cx].Unit
	if u == nil {
		return false
	}
	if u.Playerid == ps.battle.Selected.Playerid {
		return false
	}

	fromX := ps.battle.SelectedX
	fromY := ps.battle.SelectedY
	dist := abs(cx-fromX) + abs(cy-fromY)

	return dist > 0 && dist <= ps.battle.SelectedAction.Range
}
func (ps *PlayScreen) drawSelectedActionOverlay(g core.Game, screen *ebiten.Image) {
	if !ps.battle.Active || ps.battle.SelectedAction == nil {
		return
	}

	switch ps.battle.SelectedAction.Kind {
	case core.ActionMove:
		ps.drawMoveRange(g, screen)
	case core.ActionAttack:
		ps.drawAttackRange(g, screen)
		ps.drawAttackTargetLine(g, screen)
		ps.drawAttackCursor(g, screen)
	}
}

func (ps *PlayScreen) drawAttackRange(g core.Game, screen *ebiten.Image) {
	if !ps.battle.Active || ps.battle.Selected == nil || ps.battle.SelectedAction == nil {
		return
	}
	if ps.battle.SelectedAction.Kind != core.ActionAttack {
		return
	}

	board := g.Board()
	size := g.Settings().CellSize

	x := ps.battle.SelectedX
	y := ps.battle.SelectedY
	r := ps.battle.SelectedAction.Range

	for dx := -r; dx <= r; dx++ {
		for dy := -r; dy <= r; dy++ {
			dist := abs(dx) + abs(dy)
			if dist == 0 || dist > r {
				continue
			}

			cx := x + dx
			cy := y + dy

			if cy < 0 || cy >= len(board.Location) {
				continue
			}
			if cx < 0 || cx >= len(board.Location[cy]) {
				continue
			}

			u := board.Location[cy][cx].Unit
			if u == nil {
				continue
			}
			if u.Playerid == ps.battle.Selected.Playerid {
				continue
			}

			px, py := cellTopLeft(g, cx, cy)
			col := color.RGBA{255, 60, 60, 255}

			ebitenutil.DrawRect(screen, float64(px), float64(py), float64(size), 2, col)
			ebitenutil.DrawRect(screen, float64(px), float64(py+size-2), float64(size), 2, col)
			ebitenutil.DrawRect(screen, float64(px), float64(py), 2, float64(size), col)
			ebitenutil.DrawRect(screen, float64(px+size-2), float64(py), 2, float64(size), col)
		}
	}
}

func (ps *PlayScreen) drawAttackTargetLine(g core.Game, screen *ebiten.Image) {
	if !ps.canAttackHoveredCell(g) || ps.battle.Selected == nil {
		return
	}

	in := g.Input()
	cx, cy, ok := mouseToCell(g, in.MX, in.MY)
	if !ok {
		return
	}

	sx, sy := cellTopLeft(g, ps.battle.SelectedX, ps.battle.SelectedY)
	tx, ty := cellTopLeft(g, cx, cy)
	size := g.Settings().CellSize

	x1 := sx + size/2
	y1 := sy + size/2
	x2 := tx + size/2
	y2 := ty + size/2

	ebitenutil.DrawLine(screen, float64(x1), float64(y1), float64(x2), float64(y2), color.RGBA{255, 60, 60, 180})
}

func drawImageCentered(dst *ebiten.Image, img *ebiten.Image, px, py, w, h int) {
	if img == nil {
		return
	}

	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	if sw == 0 || sh == 0 {
		return
	}

	scaleX := float64(w) / float64(sw)
	scaleY := float64(h) / float64(sh)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	drawW := float64(sw) * scale
	drawH := float64(sh) * scale

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(
		float64(px)+(float64(w)-drawW)/2,
		float64(py)+(float64(h)-drawH)/2,
	)
	dst.DrawImage(img, op)
}
func drawSlotFrame(screen *ebiten.Image, assets core.Assets, px, py, size int) {
	img := assets.FrameTemplate
	if img == nil {
		return
	}

	sw, sh := img.Bounds().Dx(), img.Bounds().Dy()
	if sw == 0 || sh == 0 {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(size)/float64(sw), float64(size)/float64(sh))
	op.GeoM.Translate(float64(px), float64(py))
	screen.DrawImage(img, op)
}

func drawEquipmentSlot(screen *ebiten.Image, assets core.Assets, cat core.ItemCategory, px, py, size int) {
	drawSlotFrame(screen, assets, px, py, size)

	icon := assets.SlotIcons[cat]
	if icon == nil {
		return
	}

	padding := size / 7 // ← smaller padding = bigger icon

	drawImageCenteredTinted(
		screen,
		icon,
		px+padding,
		py+padding,
		size-(padding*2),
		size-(padding*2),
		1.8, // ← brighten
	)
}

func drawInventoryActionButton(screen *ebiten.Image, icon *ebiten.Image, px, py, size int) {
	if icon == nil {
		return
	}

	padding := size / 18

	drawImageCentered(
		screen,
		icon,
		px+padding,
		py+padding,
		size-(padding*2),
		size-(padding*2),
	)
}

func drawInventoryPanelLayout(screen *ebiten.Image, assets core.Assets, px, py, w, h int) {
	text.Draw(screen, "Inventory", basicfont.Face7x13, px+12, py+18, color.White)
	text.Draw(screen, "Selected: None", basicfont.Face7x13, px+12, py+38, color.White)

	slotSize := 52
	gap := 10

	startX := px + 16
	startY := py + 56

	// Row 1: weapon, weapon
	drawEquipmentSlot(screen, assets, core.CategoryWeapon, startX, startY, slotSize)
	drawEquipmentSlot(screen, assets, core.CategoryWeapon, startX+slotSize+gap, startY, slotSize)

	// Row 2: armor, carry (temporary stand-in still needed if no carry icon exists yet)
	row2Y := startY + slotSize + gap
	drawEquipmentSlot(screen, assets, core.CategoryArmor, startX, row2Y, slotSize)
	drawEquipmentSlot(screen, assets, core.CategoryChest, startX+slotSize+gap, row2Y, slotSize)

	// Row 3: accessories
	row3Y := row2Y + slotSize + gap + 8
	drawEquipmentSlot(screen, assets, core.CategoryCharm, startX, row3Y, slotSize)
	drawEquipmentSlot(screen, assets, core.CategoryCharm, startX+slotSize+gap, row3Y, slotSize)
	drawEquipmentSlot(screen, assets, core.CategoryCharm, startX+2*(slotSize+gap), row3Y, slotSize)

	// Row 4: ammo

	row4Y := row3Y + slotSize + gap
	drawEquipmentSlot(screen, assets, core.CategoryAmmo, startX, row4Y, slotSize)
	drawEquipmentSlot(screen, assets, core.CategoryAmmo, startX+slotSize+gap, row4Y, slotSize)
	drawEquipmentSlot(screen, assets, core.CategoryAmmo, startX+2*(slotSize+gap), row4Y, slotSize)

	// Bottom action row
	buttonSize := 56
	buttonGap := 16
	buttonY := py + h - buttonSize - 16

	drawInventoryActionButton(screen, assets.ChestButtonIcon, startX, buttonY, buttonSize)
	drawInventoryActionButton(screen, assets.ShopButtonIcon, startX+buttonSize+buttonGap, buttonY, buttonSize)
}

func (ps *PlayScreen) drawInventorySlot(dst *ebiten.Image, g core.Game, cat core.ItemCategory, px, py, size int) {
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

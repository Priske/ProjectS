package screens

import (
	"fmt"

	"github.com/Priske/ProjectS/internal/core"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SelectedUnitInfoWidget struct {
	ps *PlayScreen
	X  int
	Y  int
}

func (w *SelectedUnitInfoWidget) SetPos(x, y int) {
	w.X = x
	w.Y = y
}

func (w *SelectedUnitInfoWidget) Bounds() (int, int, int, int) {
	return w.X, w.Y, 200, 180
}

func (w *SelectedUnitInfoWidget) Update(in core.Input) {}

func (w *SelectedUnitInfoWidget) Draw(dst *ebiten.Image) {
	if w.ps == nil {
		return
	}

	u := w.ps.battle.Selected
	if u == nil {
		ebitenutil.DebugPrintAt(dst, "No unit selected", w.X+10, w.Y+28)
		return
	}

	ebitenutil.DebugPrintAt(dst, "Selected Unit", w.X+10, w.Y+28)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Type: %v", u.Type), w.X+10, w.Y+48)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("HP: %d", u.Health), w.X+10, w.Y+64)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("ATK: %d", u.AttackPower), w.X+10, w.Y+80)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("XP: %d", u.Experience), w.X+10, w.Y+96)
	ebitenutil.DebugPrintAt(dst, fmt.Sprintf("Pos: %d,%d", w.ps.battle.SelectedX, w.ps.battle.SelectedY), w.X+10, w.Y+112)
}

type BattleState struct {
	Active bool

	Turn      TurnState
	Selected  *core.Unit
	SelectedX int
	SelectedY int

	Log []string
}

type boardUnitRef struct {
	U *core.Unit
	X int
	Y int
}

func (ps *PlayScreen) enterBattle(g core.Game) {
	ps.drag.Active = false
	ps.setup.setupMode = false

	ps.battle.Active = true
	ps.battle.Turn.Side = TurnPlayer
	ps.battle.Turn.Round = 1
	ps.battle.Selected = nil
	ps.battle.Log = nil

	ps.swapAndResetUI(ps.buildBattleUI, g)
	ps.spawnEnemySetup(g)
	ps.addBattleLog("Battle started")
	ps.addBattleLog("Player turn")
}

func (ps *PlayScreen) addBattleLog(line string) {
	ps.battle.Log = append(ps.battle.Log, line)

	const maxLog = 30
	if len(ps.battle.Log) > maxLog {
		ps.battle.Log = ps.battle.Log[len(ps.battle.Log)-maxLog:]
	}
}

func (ps *PlayScreen) trySelectUnit(g core.Game, mx, my int) bool {
	if !ps.battle.Active || ps.battle.Turn.Side != TurnPlayer {
		return false
	}

	cx, cy, ok := mouseToCell(g, mx, my)
	if !ok {
		return false
	}

	u := g.Board().Location[cy][cx].Unit
	if u == nil {
		return false
	}

	if u.Playerid != 1 {
		return false
	}

	ps.battle.Selected = u
	ps.battle.SelectedX = cx
	ps.battle.SelectedY = cy

	ps.addBattleLog(fmt.Sprintf("Selected %v (%d,%d)", u.Type, cx, cy))
	return true
}

func (ps *PlayScreen) tryMoveUnit(g core.Game, u *core.Unit, fromX, fromY, toX, toY int) bool {
	if u == nil {
		ps.addBattleLog("Move blocked: no unit")
		return false
	}
	if u.Playerid != 1 {
		ps.addBattleLog("Move blocked: not your unit")
		return false
	}
	if ps.battle.Turn.Side != TurnPlayer {
		ps.addBattleLog("Move blocked: not player turn")
		return false
	}
	if !ps.battle.Turn.CanMove(u) {
		ps.addBattleLog("Move blocked: unit already moved")
		return false
	}

	dx := abs(toX - fromX)
	dy := abs(toY - fromY)
	dist := dx + dy

	if dist == 0 {
		ps.addBattleLog("Move blocked: same tile")
		return false
	}
	if dist > u.MoveRange {
		ps.addBattleLog("Move blocked: out of range")
		return false
	}

	board := g.Board()

	if fromY < 0 || fromY >= len(board.Location) || fromX < 0 || fromX >= len(board.Location[fromY]) {
		ps.addBattleLog("Move blocked: invalid source")
		return false
	}
	if toY < 0 || toY >= len(board.Location) || toX < 0 || toX >= len(board.Location[toY]) {
		ps.addBattleLog("Move blocked: invalid destination")
		return false
	}

	src := &board.Location[fromY][fromX]
	dst := &board.Location[toY][toX]

	if src.Unit == nil {
		ps.addBattleLog("Move blocked: source empty")
		return false
	}
	if src.Unit != u {
		ps.addBattleLog("Move blocked: dragged unit mismatch")
		return false
	}
	if dst.Unit != nil {
		ps.addBattleLog("Move blocked: destination occupied")
		return false
	}

	dst.Unit = src.Unit
	src.Unit = nil

	ps.battle.Selected = dst.Unit
	ps.battle.SelectedX = toX
	ps.battle.SelectedY = toY
	ps.battle.Turn.MarkMoved(dst.Unit)

	ps.addBattleLog(fmt.Sprintf("Unit moved to (%d,%d)", toX, toY))
	return true
}
func (ps *PlayScreen) tryAttackWithSelectedUnit(g core.Game, targetX, targetY int) bool {
	if ps.battle.Selected == nil {
		return false
	}
	if ps.battle.Turn.Side != TurnPlayer {
		return false
	}
	if !ps.battle.Turn.CanAttack(ps.battle.Selected) {
		return false
	}

	fromX := ps.battle.SelectedX
	fromY := ps.battle.SelectedY

	dx := abs(targetX - fromX)
	dy := abs(targetY - fromY)
	if dx+dy != 1 {
		return false
	}

	board := g.Board()
	src := board.TilePtr(fromX, fromY)
	dst := board.TilePtr(targetX, targetY)

	if src == nil || src.Unit == nil || dst == nil || dst.Unit == nil {
		return false
	}

	attacker := src.Unit
	defender := dst.Unit

	if attacker.Playerid == defender.Playerid {
		return false
	}

	defender.Health -= attacker.AttackPower
	ps.battle.Turn.MarkAttacked(ps.battle.Selected)

	ps.addBattleLog("Unit attacked")

	if defender.Health <= 0 {
		dst.Unit = nil
		ps.addBattleLog("Enemy defeated")
	}

	return true
}

func (ps *PlayScreen) endPlayerTurn(g core.Game) {
	ps.battle.Selected = nil
	ps.battle.Turn = TurnState{
		Side:  TurnEnemy,
		Round: ps.battle.Turn.Round,
		Units: make(map[int]UnitTurnState),
	}
	ps.addBattleLog("Enemy turn")
}

func (ps *PlayScreen) endEnemyTurn(g core.Game) {
	ps.battle.Selected = nil
	ps.battle.Turn = TurnState{
		Side:  TurnPlayer,
		Round: ps.battle.Turn.Round + 1,
		Units: make(map[int]UnitTurnState),
	}
	ps.addBattleLog("Player turn")
}
func (ps *PlayScreen) enemyUnitsOnBoard(g core.Game) []boardUnitRef {
	board := g.Board()
	out := []boardUnitRef{}

	for y := range board.Location {
		for x := range board.Location[y] {
			u := board.Location[y][x].Unit
			if u == nil {
				continue
			}
			if u.Playerid != -1 {
				continue
			}
			out = append(out, boardUnitRef{
				U: u,
				X: x,
				Y: y,
			})
		}
	}

	return out
}

func (ps *PlayScreen) playerUnitsOnBoard(g core.Game) []boardUnitRef {
	board := g.Board()
	out := []boardUnitRef{}

	for y := range board.Location {
		for x := range board.Location[y] {
			u := board.Location[y][x].Unit
			if u == nil {
				continue
			}
			if u.Playerid == -1 {
				continue
			}
			out = append(out, boardUnitRef{
				U: u,
				X: x,
				Y: y,
			})
		}
	}

	return out
}

func (ps *PlayScreen) enemyTryAttackAdjacent(g core.Game, ex, ey int) bool {
	board := g.Board()
	src := board.TilePtr(ex, ey)
	if src == nil || src.Unit == nil {
		return false
	}

	dirs := [][2]int{
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, 0},
	}

	for _, d := range dirs {
		tx := ex + d[0]
		ty := ey + d[1]

		dst := board.TilePtr(tx, ty)
		if dst == nil || dst.Unit == nil {
			continue
		}
		if dst.Unit.Playerid == src.Unit.Playerid {
			continue
		}

		dst.Unit.Health -= src.Unit.AttackPower
		ps.addBattleLog("Enemy attacked")

		if dst.Unit.Health <= 0 {
			dst.Unit = nil
			ps.addBattleLog("Player unit defeated")
		}

		return true
	}

	return false
}

func (ps *PlayScreen) enemyTryMoveTowardPlayer(g core.Game, ex, ey int) bool {
	board := g.Board()
	src := board.TilePtr(ex, ey)
	if src == nil || src.Unit == nil {
		return false
	}

	players := ps.playerUnitsOnBoard(g)
	if len(players) == 0 {
		return false
	}

	nearest := players[0]
	bestDist := manhattan(ex, ey, nearest.X, nearest.Y)

	for _, p := range players[1:] {
		d := manhattan(ex, ey, p.X, p.Y)
		if d < bestDist {
			bestDist = d
			nearest = p
		}
	}

	dx := 0
	if nearest.X > ex {
		dx = 1
	} else if nearest.X < ex {
		dx = -1
	}

	dy := 0
	if nearest.Y > ey {
		dy = 1
	} else if nearest.Y < ey {
		dy = -1
	}

	candidates := [][2]int{}
	if dx != 0 {
		candidates = append(candidates, [2]int{ex + dx, ey})
	}
	if dy != 0 {
		candidates = append(candidates, [2]int{ex, ey + dy})
	}

	for _, c := range candidates {
		tx, ty := c[0], c[1]
		dst := board.TilePtr(tx, ty)
		if dst == nil || dst.Unit != nil {
			continue
		}

		dst.Unit = src.Unit
		src.Unit = nil
		ps.addBattleLog("Enemy moved")
		return true
	}

	return false
}
func (ps *PlayScreen) updateBattle(g core.Game) {
	in := g.Input()

	if ps.battle.Turn.Side == TurnEnemy {
		ps.runEnemyTurn(g)
		return
	}

	if !in.LeftClicked {
		return
	}

	mx, my := in.MX, in.MY
	cx, cy, ok := mouseToCell(g, mx, my)
	if !ok {
		return
	}

	board := g.Board()
	clickedUnit := board.Location[cy][cx].Unit

	if ps.battle.Selected == nil {
		ps.trySelectUnit(g, mx, my)
		return
	}

	if clickedUnit != nil && clickedUnit.Playerid == ps.battle.Selected.Playerid {
		ps.trySelectUnit(g, mx, my)
		return
	}

	if clickedUnit != nil && clickedUnit.Playerid != ps.battle.Selected.Playerid {
		ps.tryAttackWithSelectedUnit(g, cx, cy)
		return
	}
}

func (ps *PlayScreen) runEnemyTurn(g core.Game) {
	enemies := ps.enemyUnitsOnBoard(g)

	for _, e := range enemies {
		if ps.enemyTryAttackAdjacent(g, e.X, e.Y) {
			continue
		}
		ps.enemyTryMoveTowardPlayer(g, e.X, e.Y)
	}

	ps.endEnemyTurn(g)
}

func manhattan(x1, y1, x2, y2 int) int {
	dx := x1 - x2
	if dx < 0 {
		dx = -dx
	}

	dy := y1 - y2
	if dy < 0 {
		dy = -dy
	}

	return dx + dy
}

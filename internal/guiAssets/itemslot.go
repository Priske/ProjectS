package guiassets

import (
	"github.com/Priske/ProjectS/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

type ItemSlotKind int

const (
	SlotWeaponPrimary ItemSlotKind = iota
	SlotWeaponSecondary
	SlotArmor
	SlotCarry
	SlotAccessory
	SlotAmmo
)

type ItemSlotWidget struct {
	X, Y, W, H int

	Kind  ItemSlotKind
	Label string

	Item core.Item

	Hovered  bool
	Selected bool

	OnClick func()
}

func (s *ItemSlotWidget) Bounds() (int, int, int, int) {
	return s.X, s.Y, s.W, s.H
}

func (s *ItemSlotWidget) SetPos(x, y int) {
	s.X = x
	s.Y = y
}

func (s *ItemSlotWidget) Update(in core.Input) {
	s.Hovered = in.MX >= s.X && in.MX < s.X+s.W &&
		in.MY >= s.Y && in.MY < s.Y+s.H

	if s.Hovered && in.LeftClicked && s.OnClick != nil {
		s.OnClick()
	}
}

func (s *ItemSlotWidget) Draw(dst *ebiten.Image) {
	// draw slot background by Kind
	// draw hover/selected border
	// if s.Item != nil draw item icon
	// else draw faint slot hint/icon
}

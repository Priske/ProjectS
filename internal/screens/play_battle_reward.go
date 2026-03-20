package screens

import "github.com/Priske/ProjectS/internal/core"

type RewardKind int

const (
	RewardUnit RewardKind = iota
	RewardItem
	RewardBuff
	RewardXP
)

type RewardChoice struct {
	Kind        RewardKind
	Title       string
	Description string

	Unit   *core.Unit
	Item   core.Item
	Amount int
}

const PreviewUnitID = -1

package core

type MetaProgression struct {
	UnlockedUnitTypes map[UnitType]bool
	UnlockedItemIDs   map[ItemID]bool
	MetaCurrency      int
}

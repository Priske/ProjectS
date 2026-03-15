package core

type Armor struct {
	ItemBase

	Modifiers StatModifiers
	Capacity  map[ItemCategory]int
	Stored    []Item
	Actions   []UnitAction
}

func (a *Armor) Base() *ItemBase                     { return &a.ItemBase }
func (a *Armor) Category() ItemCategory              { return CategoryArmor }
func (a *Armor) EquipSlot() EquipmentSlot            { return SlotArmor }
func (a *Armor) StoredItems() []Item                 { return a.Stored }
func (a *Armor) CarryCapacity() map[ItemCategory]int { return a.Capacity }
func (a *Armor) StatModifiers() StatModifiers        { return a.Modifiers }
func (a *Armor) GrantedActions() []UnitAction        { return a.Actions }

type StatModifiers struct {
	HealthBonus        int
	AttackBonus        int
	MoveRangeBonus     int
	MoveActionsBonus   int
	AttackActionsBonus int
	CarryLimitBonus    int
}

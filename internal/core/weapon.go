package core

type Weapon struct {
	ItemBase

	AttackBonus int
	AmmoCurrent int
	AmmoMax     int

	Actions []UnitAction
	Mods    []Item
}

func (w *Weapon) Base() *ItemBase              { return &w.ItemBase }
func (w *Weapon) Category() ItemCategory       { return CategoryWeapon }
func (w *Weapon) EquipSlot() EquipmentSlot     { return SlotWeapon }
func (w *Weapon) CurrentAmmo() int             { return w.AmmoCurrent }
func (w *Weapon) MaxAmmo() int                 { return w.AmmoMax }
func (w *Weapon) GrantedActions() []UnitAction { return w.Actions }
func (w *Weapon) ConsumeAmmo(amount int) bool {
	if w.AmmoCurrent < amount {
		return false
	}
	w.AmmoCurrent -= amount
	return true
}

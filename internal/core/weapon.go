package core

type Weapon struct {
	ItemBase
	slot        EquipmentSlot
	statMods    StatModifiers
	actions     []UnitAction
	ammoCurrent int
	ammoMax     int
	TwoHanded   bool
}

func (w *Weapon) Base() *ItemBase              { return &w.ItemBase }
func (w *Weapon) Category() ItemCategory       { return CategoryWeapon }
func (w *Weapon) EquipSlot() EquipmentSlot     { return w.slot }
func (w *Weapon) StatModifiers() StatModifiers { return w.statMods }
func (w *Weapon) GrantedActions() []UnitAction { return w.actions }
func (w *Weapon) CurrentAmmo() int             { return w.ammoCurrent }
func (w *Weapon) MaxAmmo() int                 { return w.ammoMax }

func (w *Weapon) ConsumeAmmo(amount int) bool {
	if amount <= 0 {
		return true
	}
	if w.ammoCurrent < amount {
		return false
	}
	w.ammoCurrent -= amount
	return true
}

func (w *Weapon) Reload(ammo Item) bool {
	// stub for now
	return false
}
func (u *Unit) EquipWeapon(slot EquipmentSlot, w Equippable) {
	weapon, ok := w.(*Weapon)
	if !ok {
		return
	}

	// clear both slots first (prevents weird overlaps)
	delete(u.Equipped, SlotWeapon1)
	delete(u.Equipped, SlotWeapon2)

	if weapon.TwoHanded {
		// occupy BOTH slots with same pointer
		u.Equipped[SlotWeapon1] = w
		u.Equipped[SlotWeapon2] = w
		return
	}

	// one-handed
	u.Equipped[slot] = w
}

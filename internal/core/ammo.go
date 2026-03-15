package core

type Ammo struct {
	ItemBase

	Amount int
}

func (a *Ammo) Base() *ItemBase        { return &a.ItemBase }
func (a *Ammo) Category() ItemCategory { return CategoryAmmo }

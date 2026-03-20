package core

func MakeCommandersPistol() Equippable {
	return &Weapon{
		ItemBase: ItemBase{
			ID:          "commanders_pistol",
			Name:        "Commander's Pistol",
			Description: "A reliable sidearm.",
			Weight:      1,
		},
		slot:        SlotWeapon1,
		ammoCurrent: 6,
		ammoMax:     6,
	}
}

func MakeCommandersCoat() Equippable {
	return &Armor{
		ItemBase: ItemBase{
			ID:          "commanders_coat",
			Name:        "Commander's Coat",
			Description: "A distinguished officer's coat.",
			Weight:      1,
		},
		slot: SlotArmor,
	}
}

package core

import "math/rand"

func newBaseUnit(playerId, unitId int, unitType UnitType) Unit {
	return Unit{
		Type:         unitType,
		UnitId:       unitId,
		Playerid:     playerId,
		CarryLimit:   0,
		Equipped:     map[EquipmentSlot]Equippable{},
		CarriedItems: []Item{},
		BattleStats:  BattleStats{},
		Name:         GenerateUnitName(unitType),
	}
}

func MakeNewSoldier(playerId, unitId int) *Unit {
	u := newBaseUnit(playerId, unitId, Soldier)
	u.MaxHealth = 2
	u.CurrentHealth = 2
	u.AttackPower = 2
	u.Experience = 0
	u.MoveRange = 2
	u.MoveActionsPerTurn = 1
	u.AttackActionsPerTurn = 1
	u.Actions = defaultSoldierActions(2, 1)
	u.UnitCategory = Attack
	return &u
}

func MakeNewCommander(playerId, unitId int) *Unit {
	u := newBaseUnit(playerId, unitId, Commander)
	u.MaxHealth = 5
	u.CurrentHealth = 5
	u.AttackPower = 1
	u.Experience = 6
	u.MoveRange = 1
	u.MoveActionsPerTurn = 1
	u.AttackActionsPerTurn = 1
	u.Actions = defaultMeleeActions(1, 1)
	u.UnitCategory = Flag
	u.Equipped[SlotWeapon1] = MakeCommandersPistol()
	u.Equipped[SlotArmor] = MakeCommandersCoat()
	return &u
}

func MakeNewMedic(playerId, unitId int) *Unit {
	u := newBaseUnit(playerId, unitId, Medic)
	u.MaxHealth = 5
	u.CurrentHealth = 5
	u.AttackPower = 1
	u.Experience = 0
	u.MoveRange = 1
	u.MoveActionsPerTurn = 1
	u.AttackActionsPerTurn = 1
	u.Actions = defaultMeleeActions(1, 1)
	u.UnitCategory = Support
	return &u
}

func MakeNewShield(playerId, unitId int) *Unit {
	u := newBaseUnit(playerId, unitId, Shield)
	u.MaxHealth = 6
	u.CurrentHealth = 6
	u.AttackPower = 1
	u.Experience = 0
	u.MoveRange = 1
	u.MoveActionsPerTurn = 1
	u.AttackActionsPerTurn = 1
	u.Actions = defaultMeleeActions(1, 1)
	u.UnitCategory = Defense
	return &u
}

func MakeNewSniper(playerId, unitId int) *Unit {
	u := newBaseUnit(playerId, unitId, Sniper)
	u.MaxHealth = 1
	u.CurrentHealth = 1
	u.AttackPower = 3
	u.Experience = 0
	u.MoveRange = 1
	u.MoveActionsPerTurn = 1
	u.AttackActionsPerTurn = 1
	u.Actions = defaulSniperActions(1, 1)
	u.UnitCategory = Attack
	return &u
}

func MakeNewEnemyCultistKnife(playerId, unitId int) *Unit {
	u := newBaseUnit(playerId, unitId, Enemy_cultist_knife)
	u.MaxHealth = 4
	u.CurrentHealth = 4
	u.AttackPower = 2
	u.Experience = 0
	u.MoveRange = 3
	u.MoveActionsPerTurn = 1
	u.AttackActionsPerTurn = 1
	u.Actions = defaultCultistMeleeActions(3, 1)
	u.UnitCategory = Attack
	return &u
}

func MakeNewEnemyCultistLord(playerId, unitId int) *Unit {
	u := newBaseUnit(playerId, unitId, Enemy_cultist_lord)
	u.MaxHealth = 4
	u.CurrentHealth = 4
	u.AttackPower = 1
	u.Experience = 0
	u.MoveRange = 1
	u.MoveActionsPerTurn = 1
	u.AttackActionsPerTurn = 1
	u.Actions = defaultCultistLordActions(1, 1)
	u.UnitCategory = Flag
	return &u
}

var soldierNames = []string{
	"Hale", "Brooks", "Keller", "Vargas", "Rowan",
	"Hayes", "Carter", "Maddox", "Pierce", "Griffin",
}

var supportNames = []string{
	"Lumen", "Solace", "Mercy", "Auriel", "Nova",
	"Helia", "Seren", "Lyra", "Elion", "Vale",
}

var defenseNames = []string{
	"Ironwall", "Bulwark", "Stone", "Aegis", "Bastion",
	"Atlas", "Garrick", "Tarkus", "Duran", "Bragg",
}

var commanderNames = []string{
	"Valerius", "Drake", "Arcturus", "Magnus", "Rhea",
	"Leonis", "Cassian", "Severin", "Octavia", "Tiber",
}

func GenerateUnitName(t UnitType) string {
	switch t {

	case Soldier, Sniper:
		return soldierNames[rand.Intn(len(soldierNames))]

	case Medic:
		return supportNames[rand.Intn(len(supportNames))]

	case Shield:
		return defenseNames[rand.Intn(len(defenseNames))]

	case Commander:
		return commanderNames[rand.Intn(len(commanderNames))]
	}

	return "Unknown"
}

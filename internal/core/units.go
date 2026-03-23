package core

type Unit struct {
	Type          UnitType
	UnitId        int
	MaxHealth     int
	CurrentHealth int
	AttackPower   int
	Experience    int
	Playerid      int
	CarryLimit    int
	Defending     bool

	MoveRange            int
	MoveActionsPerTurn   int
	AttackActionsPerTurn int

	Equipped     map[EquipmentSlot]Equippable
	CarriedItems []Item

	Name         string
	BattleStats  BattleStats
	UnitCategory UnitCategory
	Actions      []UnitAction
	AIKind       UnitAIKind
}
type BattleStats struct {
	Kills       int
	DamageTaken int
	DamageDealt int
}

type UnitTurnState struct {
	RemainingMoveActions   int
	RemainingAttackActions int

	ActionUses map[string]int

	TempMoveRangeBonus   int
	TempAttackRangeBonus int
}

type UnitType int

const (
	UnitNone UnitType = iota
	Soldier
	Commander
	Medic
	Shield
	Sniper
	Razor
	Enemy_cultist_knife
	Enemy_cultist_lord
	Enemy_cultist_shield
	Enemy_rat_brood_lord
	Enemy_rat_lord
	Enemy_rat_knife
	Enemy_rat_axes
)

type UnitCategory int

const (
	Attack UnitCategory = iota
	Defense
	Support
	Flag
)

type UnitAIKind int

const (
	AIKindDefaultAggro UnitAIKind = iota
	AIKindBroodLord
)

func (u *Unit) TotalWeight() int {
	total := 0

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}
		total += eq.Base().Weight
	}

	for _, item := range u.CarriedItems {
		if item == nil {
			continue
		}
		total += item.Base().Weight
	}

	return total
}

func (u *Unit) TotalCarryLimit() int {
	total := u.CarryLimit

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}
		if cm, ok := eq.(CarryModifier); ok {
			total += cm.CarryBonus()
		}
	}

	if total < 0 {
		total = 0
	}

	return total
}

func (u *Unit) CanCarryByWeight(item Item) bool {
	if item == nil {
		return false
	}
	return u.TotalWeight()+item.Base().Weight <= u.TotalCarryLimit()
}

func (u *Unit) CanCarry(item Item) bool {
	if !u.CanCarryByCategory(item) {
		return false
	}
	if !u.CanCarryByWeight(item) {
		return false
	}
	return true
}

func (u *Unit) Equip(item Equippable) Equippable {
	if u.Equipped == nil {
		u.Equipped = map[EquipmentSlot]Equippable{}
	}

	slot := item.EquipSlot()
	old := u.Equipped[slot]
	u.Equipped[slot] = item
	u.ClampCurrentHealth()
	return old
}
func (u *Unit) Unequip(slot EquipmentSlot) Equippable {
	if u.Equipped == nil {
		return nil
	}

	old := u.Equipped[slot]
	delete(u.Equipped, slot)
	u.ClampCurrentHealth()
	return old
}

func (u *Unit) EquippedItem(slot EquipmentSlot) Equippable {
	if u.Equipped == nil {
		return nil
	}
	return u.Equipped[slot]
}

func (u *Unit) TotalCategoryCapacity() map[ItemCategory]int {
	out := map[ItemCategory]int{}

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}

		cp, ok := eq.(CapacityProvider)
		if !ok {
			continue
		}

		for cat, amount := range cp.CarryCapacity() {
			out[cat] += amount
		}
	}

	return out
}

func (u *Unit) usedCategoryCount(cat ItemCategory) int {
	count := 0

	for _, item := range u.CarriedItems {
		if item.Category() == cat {
			count++
		}
	}

	return count
}
func (u *Unit) CanCarryByCategory(item Item) bool {
	if item == nil {
		return false
	}

	cat := item.Category()

	capacity := u.TotalCategoryCapacity()
	used := u.usedCategoryCount(cat)

	return used < capacity[cat]
}

func (u *Unit) TotalMoveRange() int {
	total := u.MoveRange

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}

		sm, ok := eq.(StatModifierProvider)
		if !ok {
			continue
		}

		total += sm.StatModifiers().MoveRangeBonus
	}

	if total < 0 {
		total = 0
	}

	return total
}

func (u *Unit) AllActions() []UnitAction {
	out := make([]UnitAction, 0, len(u.Actions))
	out = append(out, u.Actions...)

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}

		ap, ok := eq.(ActionProvider)
		if !ok {
			continue
		}

		out = append(out, ap.GrantedActions()...)
	}

	return out
}

func (u *Unit) TotalAttackPower() int {
	total := u.AttackPower

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}

		sm, ok := eq.(StatModifierProvider)
		if !ok {
			continue
		}

		total += sm.StatModifiers().AttackBonus
	}

	if total < 0 {
		total = 0
	}

	return total
}

func (u *Unit) TotalMaxHealth() int {
	total := u.MaxHealth

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}

		sm, ok := eq.(StatModifierProvider)
		if !ok {
			continue
		}

		total += sm.StatModifiers().HealthBonus
	}

	if total < 0 {
		total = 0
	}

	return total
}
func (u *Unit) TotalMoveActionsPerTurn() int {
	total := u.MoveActionsPerTurn

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}

		sm, ok := eq.(StatModifierProvider)
		if !ok {
			continue
		}

		total += sm.StatModifiers().MoveActionsBonus
	}

	if total < 0 {
		total = 0
	}

	return total
}

func (u *Unit) TotalAttackActionsPerTurn() int {
	total := u.AttackActionsPerTurn

	for _, eq := range u.Equipped {
		if eq == nil {
			continue
		}

		sm, ok := eq.(StatModifierProvider)
		if !ok {
			continue
		}

		total += sm.StatModifiers().AttackActionsBonus
	}

	if total < 0 {
		total = 0
	}

	return total
}

func (u *Unit) AddCarriedItem(item Item) bool {
	if !u.CanCarry(item) {
		return false
	}
	u.CarriedItems = append(u.CarriedItems, item)
	return true
}

func (u *Unit) RemoveCarriedItem(item Item) bool {
	for i, it := range u.CarriedItems {
		if it == item {
			u.CarriedItems = append(u.CarriedItems[:i], u.CarriedItems[i+1:]...)
			return true
		}
	}
	return false
}

func (u *Unit) ClampCurrentHealth() {
	maxHP := u.TotalMaxHealth()

	if u.CurrentHealth > maxHP {
		u.CurrentHealth = maxHP
	}
	if u.CurrentHealth < 0 {
		u.CurrentHealth = 0
	}
}

func (u *Unit) TakeDamage(amount int) {
	if amount < 0 {
		amount = 0
	}

	u.CurrentHealth -= amount
	if u.CurrentHealth < 0 {
		u.CurrentHealth = 0
	}
}

func (u *Unit) Heal(amount int) {
	if amount < 0 {
		amount = 0
	}

	u.CurrentHealth += amount
	if u.CurrentHealth > u.TotalMaxHealth() {
		u.CurrentHealth = u.TotalMaxHealth()
	}
}

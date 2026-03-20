package core

type ItemCategory int

type EquipmentSlot int

const (
	SlotWeapon1 EquipmentSlot = iota
	SlotWeapon2
	SlotArmor
	SlotHead
	SlotLegs
	SlotCharm
	SlotBag
	SlotAmmo1
	SlotAmmo2
	SlotAmmo3
	SlotAccessory
	SlotCarry
)
const (
	CategoryAmmo ItemCategory = iota
	CategoryPotion
	CategoryWeapon
	CategoryArmor
	CategoryLegs
	CategoryHead
	CategoryAccessory
	CategoryGrenade
	CategoryTool
	CategoryCharm
	CategoryBag
)

type ItemID string
type ItemBase struct {
	ID          ItemID
	Name        string
	Description string
	Weight      int
}

type Item interface {
	Base() *ItemBase
	Category() ItemCategory
}

type StatModifiers struct {
	HealthBonus        int
	AttackBonus        int
	MoveRangeBonus     int
	MoveActionsBonus   int
	AttackActionsBonus int
	CarryLimitBonus    int
}

/////
/////
/*	Interfaces to modify items*/
/////
/////
type Equippable interface {
	Item
	EquipSlot() EquipmentSlot
}

type ContainerItem interface {
	Item
	CanStore(item Item) bool
	Store(item Item) bool
	Remove(item Item) bool
	StoredItems() []Item
}

type Modifiable interface {
	Item
	CanAttach(mod Item) bool
	Attach(mod Item) bool
	Detach(mod Item) bool
	AttachedItems() []Item
}

type AmmoUser interface {
	Item
	CurrentAmmo() int
	MaxAmmo() int
	ConsumeAmmo(amount int) bool
	Reload(ammo Item) bool
}
type Usable interface {
	Item
	CanUse(user *Unit) bool
	Use(user *Unit, target *Unit) bool
}

type CarryModifier interface {
	Item
	CarryBonus() int
}
type CapacityProvider interface {
	Item
	CarryCapacity() map[ItemCategory]int
}
type StatModifierProvider interface {
	Item
	StatModifiers() StatModifiers
}
type ActionProvider interface {
	Item
	GrantedActions() []UnitAction
}

////
////
/* ITEM HELPERS*/
////
////

func ItemWeight(item Item) int {
	if item == nil {
		return 0
	}
	return item.Base().Weight
}

package core

type ItemCategory int

type EquipmentSlot int

const (
	SlotWeapon EquipmentSlot = iota
	SlotArmor
	SlotAccessory
)
const (
	CategoryAmmo ItemCategory = iota
	CategoryTile
	CategoryPotion
	CategoryWeapon
	CategoryArmor
	CategoryAccessory
	CategoryGrenade
	CategoryTool
	CategoryChest
	CategoryCharm
	CategoryBag
)

type ItemBase struct {
	ID          string
	Name        string
	Description string
	Weight      int
}

type Item interface {
	Base() *ItemBase
	Category() ItemCategory
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

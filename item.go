package main

// ItemType
const (
	ARMOR          ItemType = iota
	HEALTH         ItemType = iota
	INSTANT_DAMAGE ItemType = iota
)

type ItemType int8

type Item struct {
	id     int64
	iType  ItemType
	effect float64
}

var itemIDCounter int64

func NewItem(iType ItemType, effect float64) *Item {
	item := new(Item)
	item.id = itemIDCounter
	itemIDCounter++
	item.iType = iType
	item.effect = effect
	return item
}

func getStringFromItemType(iType ItemType) string {
	switch iType {
	case ARMOR:
		return "Armor"
	case HEALTH:
		return "Health"
	case INSTANT_DAMAGE:
		return "Instant Damage"
	default:
		return "INVALID"
	}
}

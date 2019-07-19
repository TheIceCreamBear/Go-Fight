package main

import "fmt"

// ItemType
const (
	ARMOR          ItemType = iota // 0
	HEALTH         ItemType = iota // 1
	INSTANT_DAMAGE ItemType = iota // 2
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
		return "Damage" // TODO determine if name shoule be instant damage or just damage
	default:
		return "INVALID"
	}
}

func (item *Item) print() {
	fmt.Printf("Item: Type=%-7s Effect=%7.3f\n", getStringFromItemType(item.iType), item.effect)
}

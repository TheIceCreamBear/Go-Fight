package main

import "fmt"

// ItemType
const (
	KEY            ItemType = iota // 0 // NOTE: the effect for a key is how many chests it can open
	ARMOR          ItemType = iota // 1
	HEALTH         ItemType = iota // 2
	INSTANT_DAMAGE ItemType = iota // 3
)

type ItemType int8

type Item struct {
	id     int64
	iType  ItemType
	effect float64
}

var itemIDCounter int64

func getGenetateableItemsWithChance() (chances map[ItemType]float64) {
	chances = make(map[ItemType]float64, INSTANT_DAMAGE+1)
	chances[KEY] = 0.4
	chances[ARMOR] = 0.1
	chances[HEALTH] = 0.2
	chances[INSTANT_DAMAGE] = 0.1
	// chances[No Item Chance] = 0.2
	// NOTE: TOTAL CHANCES DO NOT HAVE TO EQUAL 1
	return
}

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
	case KEY:
		return "Key"
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

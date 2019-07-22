package main

import (
	"fmt"
	"math/rand"
)

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
	chances[KEY] = 0.45 // THIS MUST REMAIN ABOVE OR EQUAL TO chestLockedChance SO ALL CHESTS CAN BE OPENED
	chances[ARMOR] = 0.15
	chances[HEALTH] = 0.25
	chances[INSTANT_DAMAGE] = 0.15
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

func createItemWithType(iType ItemType) *Item {
	var effect float64
	switch iType {
	case KEY:
		chanceNeeded := rand.Float64()
		switch {
		case .6 > chanceNeeded:
			effect = 1
		case .6+.3 > chanceNeeded:
			effect = 2
		case .6+.3+.1 > chanceNeeded:
			effect = 3
		default:
			effect = 1
		}
	case ARMOR:
		chanceNeeded := rand.Float64()
		switch {
		case .4 > chanceNeeded:
			effect = 1
		case .4+.3 > chanceNeeded:
			effect = 2
		case .4+.3+.2 > chanceNeeded:
			effect = 3
		case .4+.3+.2+.1 > chanceNeeded:
			effect = 4
		default:
			effect = 1
		}
	case HEALTH:
		chanceNeeded := rand.Float64()
		switch {
		case .5 > chanceNeeded:
			effect = 20
		case .5+.35 > chanceNeeded:
			effect = 50
		case .5+.35+.1 > chanceNeeded:
			effect = 100
		case .5+.35+.1+.05 > chanceNeeded:
			effect = 200
		default:
			effect = 1
		}
	case INSTANT_DAMAGE:
		chanceNeeded := rand.Float64()
		switch {
		case .525 > chanceNeeded:
			effect = 20
		case .525+.375 > chanceNeeded:
			effect = 50
		case .525+.375+.075 > chanceNeeded:
			effect = 100
		case .525+.375+.075+.025 > chanceNeeded:
			effect = 200
		default:
			effect = 1
		}
	default:
		fmt.Println("IMPOSSIBLE CASE: default case in generated item type", iType)
	}

	item := NewItem(iType, effect)

	return item
}

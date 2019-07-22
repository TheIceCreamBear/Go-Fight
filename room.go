package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// room types
const (
	START      RoomType = iota
	HALLWAY    RoomType = iota
	GREAT_HALL RoomType = iota
	DUNGEON    RoomType = iota
	CHEST      RoomType = iota
	MYSTIC     RoomType = iota
)

const chestLockedChance float64 = 0.4
const (
	noChest    float64 = 0.4
	oneChest   float64 = 0.4
	twoChest   float64 = 0.15
	threeChest float64 = 0.05
)

type Chest struct {
	locked bool
	item   *Item
}

type RoomType int8

type Room struct {
	rType  RoomType
	id     int64
	loc    Location
	chests []*Chest
	dUp    Door
	dDown  Door
	dLeft  Door
	dRight Door
}

func getGenetateableTypes() [6]RoomType {
	return [6]RoomType{START, HALLWAY, GREAT_HALL, DUNGEON, CHEST, MYSTIC}
}

func (r *Room) canLeaveFrom(direction Direction) bool {
	switch direction {
	case UP:
		return r.dUp.exists /* && !r.dUp.locked */
	case DOWN:
		return r.dDown.exists /* && !r.dDown.locked */
	case LEFT:
		return r.dLeft.exists /* && !r.dLeft.locked */
	case RIGHT:
		return r.dRight.exists /* && !r.dRight.locked */
	default:
		fmt.Println("D E F A U L T  C A S E ")
		return false
	}
}

func (r *Room) initChests() {
	r.chests = make([]*Chest, 3)
	var numChests int
	chanceNeeded := rand.Float64()
	switch {
	case noChest > chanceNeeded:
		return // numChests is 0
	case noChest+oneChest > chanceNeeded:
		numChests = 1
	case noChest+oneChest+twoChest > chanceNeeded:
		numChests = 2
	case noChest+oneChest+twoChest+threeChest > chanceNeeded:
		numChests = 3
	default:
		fmt.Println("Bad chance, check room.initChests")
	}

	for i := 0; i < numChests; i++ {
		chest := new(Chest)
		chanceNeeded = rand.Float64()
		if chestLockedChance > chanceNeeded {
			chest.locked = true
		}

		itemChances := getGenetateableItemsWithChance()
		var generatedType ItemType
		chance := 0.0
		chanceNeeded = rand.Float64()
		keys := make([]ItemType, len(itemChances))
		ind := 0
		for key := range itemChances {
			keys[ind] = key
			ind++
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] > keys[j]
		})

		for _, val := range keys {
			chance += itemChances[val]
			if chance > chanceNeeded {
				generatedType = val
				break
			}
		}

		chest.item = createItemWithType(generatedType)
		r.chests[i] = chest
	}
}

func (r *Room) initDoors(up Door, do Door, le Door, ri Door) {
	if DEBUG_MODE {
		fmt.Println("init doors ", r.loc.x, r.loc.y, up, do, le, ri)
	}
	r.dUp = Door{up.exists, up.locked}
	r.dDown = Door{do.exists, do.locked}
	r.dLeft = Door{le.exists, le.locked}
	r.dRight = Door{ri.exists, ri.locked}
	if DEBUG_MODE {
		fmt.Println("post init doors", r.dUp, r.dDown, r.dLeft, r.dRight)
	}
}

func (r *Room) getNumChests() int {
	numChests := 0
	for _, val := range r.chests {
		if val != nil {
			numChests++
		}
	}
	return numChests
}

func (r *Room) getNumChestsWithItem() int {
	numChests := 0
	for _, val := range r.chests {
		if val != nil {
			if val.item != nil {
				numChests++
			}
		}
	}
	return numChests
}

func getPrintStringFromRoomType(rType RoomType) string {
	switch rType {
	case START:
		return "Start Room"
	case HALLWAY:
		return "Hallway"
	case GREAT_HALL:
		return "Great Hall"
	case DUNGEON:
		return "Dungeon"
	case CHEST:
		return "Chest Room"
	case MYSTIC:
		return "Mystical Room"
	default:
		return "_"
	}
}

func getPrintCharFromRoomType(rType RoomType) string {
	switch rType {
	case START:
		return "S"
	case HALLWAY:
		return "H"
	case GREAT_HALL:
		return "G"
	case DUNGEON:
		return "D"
	case CHEST:
		return "C"
	case MYSTIC:
		return "M"
	default:
		return "_"
	}
}

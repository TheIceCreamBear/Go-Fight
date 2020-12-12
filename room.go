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

type Door struct {
	exists bool
	locked bool
}

type Chest struct {
	locked bool
	item   *Item
}

type RoomType int8

type Room struct {
	rType   RoomType
	id      int64
	loc     Location
	chests  []*Chest
	enemies []*Enemy
	dUp     Door
	dDown   Door
	dLeft   Door
	dRight  Door
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

func (r *Room) canRunFrom(chance float64) bool {
	switch r.rType {
	case START:
		// should not have to call this but what ever
		return true
	case HALLWAY:
		return chance < .7
	case GREAT_HALL:
		return chance < .6
	case DUNGEON:
		return chance < .2
	case CHEST:
		return chance < .4
	case MYSTIC:
		return chance < .3
	default:
		fmt.Println("D E F A U L T  C A S E ")
		return false
	}
}

func (r *Room) canRunTo(chance float64) bool {
	switch r.rType {
	case START:
		// should not have to call this but what ever
		return true
	case HALLWAY:
		return chance < 1
	case GREAT_HALL:
		return chance < .9
	case DUNGEON:
		return chance < .3
	case CHEST:
		return chance < .7
	case MYSTIC:
		return chance < .6
	default:
		fmt.Println("D E F A U L T  C A S E ")
		return false
	}
}

func (r *Room) initChests() {
	r.chests = make([]*Chest, 3)
	var numChests int
	chanceNeeded := rand.Float64()
	switch r.rType {
	case START:
		return // numChests is 0
	case HALLWAY:
		switch {
		case .85 > chanceNeeded:
			return // numChests is 0
		case .85+.15 > chanceNeeded:
			numChests = 1
		}
	case GREAT_HALL:
		switch {
		case .65 > chanceNeeded:
			return // numChests is 0
		case .65+.3 > chanceNeeded:
			numChests = 1
		case .65+.3+.05 > chanceNeeded:
			numChests = 2
		}
	case DUNGEON:
		switch {
		case .5 > chanceNeeded:
			return // numChests is 0
		case .5+.3 > chanceNeeded:
			numChests = 1
		case .5+.3+.075 > chanceNeeded:
			numChests = 2
		}
	case CHEST:
		switch {
		case .0 > chanceNeeded:
			return // numChests is 0
		case .0+.25 > chanceNeeded:
			numChests = 1
		case .0+.25+.5 > chanceNeeded:
			numChests = 2
		case .0+.25+.5+.25 > chanceNeeded:
			numChests = 3
		}
	case MYSTIC:
		switch {
		case .05 > chanceNeeded:
			return // numChests is 0
		case .05+.35 > chanceNeeded:
			numChests = 1
		case .05+.35+.45 > chanceNeeded:
			numChests = 2
		case .05+.35+.45+.15 > chanceNeeded:
			numChests = 3
		}
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

func (r *Room) initEnemies(x, y, raid int64) {
	if r.rType == DUNGEON || r.rType == MYSTIC {
		r.enemies = make([]*Enemy, 2)
	} else {
		r.enemies = make([]*Enemy, 1)
	}
	chanceNeeded := rand.Float64()

	// NOTE: Any chance that is greater than one is assuming that the radius is large enough to
	switch r.rType {
	case START:
		// impossible case but...
		return
	case HALLWAY:
		switch {
		case .7 > chanceNeeded:
			return
		case .7+.25 > chanceNeeded:
			r.enemies[0] = NewEnemy(PEON)
		case .7+.25+.05 > chanceNeeded:
			r.enemies[0] = NewEnemy(WARRIOR)
		}
	case GREAT_HALL:
		switch {
		case .6 > chanceNeeded:
			return
		case .6+.1 > chanceNeeded:
			r.enemies[0] = NewEnemy(PEON)
		case .6+.1+.3 > chanceNeeded:
			r.enemies[0] = NewEnemy(WARRIOR)
		}
	case DUNGEON:
		switch {
		case .025 > chanceNeeded:
			r.enemies[0] = NewEnemy(PEON)
		case .025+.025 > chanceNeeded:
			r.enemies[0] = NewEnemy(PEON)
			r.enemies[1] = NewEnemy(PEON)
		case .025+.025+.05 > chanceNeeded:
			r.enemies[0] = NewEnemy(WARRIOR)
		case .025+.025+.05+.15 > chanceNeeded:
			r.enemies[0] = NewEnemy(WARRIOR)
			r.enemies[1] = NewEnemy(WARRIOR)
		case .025+.025+.05+.15+.65 > chanceNeeded:
			r.enemies[0] = NewEnemy(BRUTE)
		case .025+.025+.05+.15+.65+.1 > chanceNeeded:
			r.enemies[0] = NewEnemy(BRUTE)
			r.enemies[1] = NewEnemy(BRUTE)
		}
	case CHEST:
		switch {
		case .4 > chanceNeeded:
			return
		case .4+.4 > chanceNeeded:
			r.enemies[0] = NewEnemy(WARRIOR)
		case .4+.4+.2 > chanceNeeded:
			r.enemies[0] = NewEnemy(BRUTE)
		}
	case MYSTIC:
		switch {
		case .3 > chanceNeeded:
			return
		case .3+.2 > chanceNeeded:
			r.enemies[0] = NewEnemy(WARRIOR)
		case .3+.2+.25 > chanceNeeded:
			r.enemies[0] = NewEnemy(E_MYSTIC)
		case .3+.2+.25+.25 > chanceNeeded:
			r.enemies[0] = NewEnemy(E_MYSTIC)
			r.enemies[1] = NewEnemy(E_MYSTIC)
		}
	default:
		// also impossible but...
		return

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

func (r *Room) getCurrentEnemy() *Enemy {
	for _, enemy := range r.enemies {
		if enemy != nil && enemy.health >= 0 {
			return enemy
		}
	}
	return nil
}

func (r *Room) getNumEnemies() int {
	num := 0
	for _, val := range r.enemies {
		if val != nil {
			num++
		}
	}
	return num
}

func (r *Room) getNumEnemiesAlive() int {
	num := 0
	for _, val := range r.enemies {
		if val != nil && val.health >= 0 {
			num++
		}
	}
	return num
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

func (r *Room) getNumLockedChests() int {
	numChests := 0
	for _, val := range r.chests {
		if val != nil {
			if val.locked {
				numChests++
			}
		}
	}
	return numChests
}

func (r *Room) getNumLootableChests() int {
	numChests := 0
	for _, val := range r.chests {
		if val != nil {
			if !val.locked {
				if val.item != nil {
					numChests++
				}
			}
		}
	}
	return numChests
}

func (r *Room) unlockChests(amount int) {
	if amount > r.getNumLockedChests() {
		amount = r.getNumLockedChests()
	}
	for _, val := range r.chests {
		if val != nil {
			if val.locked {
				val.locked = false
				amount--
			}
			if amount == 0 {
				break
			}
		}
	}
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

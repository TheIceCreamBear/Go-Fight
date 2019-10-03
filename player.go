package main

import (
	"fmt"
)

const (
	BasePlayerHealth   = 100.0
	BasePlayerStrength = 1.0
	BasePlayerDefense  = 1.0
)

type PlayerState int8

const (
	Exploring PlayerState = iota
	Fighting  PlayerState = iota
)

type Move struct {
	id        uint8 // For future features, including saving.
	minDamage float64
	maxDamage float64
	name      string

	// TODO Cooldowns
}

var moveIdCounter uint8

func NewMove(min, max float64, name string) *Move {
	m := new(Move)
	m.id = moveIdCounter
	moveIdCounter++
	m.minDamage = min
	m.maxDamage = max
	m.name = name
	return m
}

type Player struct {
	state       PlayerState
	movedLast   bool
	loc         *Location
	currentRoom *Room
	inventory   *Inventory
	moves       []*Move
	health      float64
	defense     float64 // enemyDamage = 1 / defense
	strength    float64 // playerDamage = (min + rand.Int64N(max - min)) * strength
}

func NewPlayer(current *Room, loc *Location, moves []*Move) *Player {
	p := new(Player)
	p.state = Exploring
	p.loc = loc
	p.currentRoom = current
	p.inventory = NewInventory()
	p.moves = moves
	p.health = BasePlayerHealth
	p.defense = 1.0
	p.strength = 1.0
	return p
}

func (p *Player) update() bool {
	// var reset
	p.movedLast = false
	if p.currentRoom.getNumEnemies() > 0 {
		p.state = Fighting
	}

	// Player turn
	var run bool
	if p.state == Exploring {
		run = p.printChoices()
	} else if p.state == Fighting {
		run = p.printFightingChoices()

		// Enemy turn
		// TODO
	} else {
		// IMPOSSIBLE CASE - Try and reset and recover
		p.state = Exploring
		run = true
	}

	return run
}

func (p *Player) printChoices() bool {
	var choice int8

	for {
		fmt.Println("\nWhat would you like to do?")
		fmt.Println("1. Explore current room")
		fmt.Println("2. Move to another room")
		fmt.Println("3. View Inventory Options")
		fmt.Println("4. View Player Stats")
		fmt.Println("5. Exit")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("An error occured while reading your choice in, please try again: ", err)
			continue
		}
		switch choice {
		case cheatInputNumber:
			p.doCheatLoop()
			continue
		case 1:
			p.printRoomOptions()
			return true
		case 2:
			p.printMoveChoices()
			return true
		case 3:
			turnConsumed := p.printInventoryChoices()
			if turnConsumed {
				return true
			}
			continue
		case 4:
			p.printPlayerStats()
		case 5:
			return false
		default:
			fmt.Println("Invalid Input, try again")
		}
	}
}

func (p *Player) printRoomOptions() {
	if DEBUG_MODE {
		fmt.Println("\nDebug prints.")
		fmt.Printf("Player Location: %+v", *p.loc)
		fmt.Printf("Room Location: %+v", p.currentRoom.loc)
		fmt.Printf("Room Type=%s\n", getPrintStringFromRoomType(p.currentRoom.rType))
		fmt.Println("Room ID", p.currentRoom.id)
	}

	fmt.Printf("\nYou are in a %s, located at %+v\n", getPrintStringFromRoomType(p.currentRoom.rType), *p.loc)
	totalChests := p.currentRoom.getNumChests()
	numUnlockedChest := p.currentRoom.getNumLootableChests()
	numLockedChests := p.currentRoom.getNumLockedChests()

	if numLockedChests+numUnlockedChest == 0 && totalChests != 0 {
		fmt.Println("All chests in this room have been looted.")
		return
	}
	if totalChests == 0 {
		fmt.Println("There are no chests in this room.")
		return
	} else if totalChests == 1 {
		if numLockedChests == totalChests {
			fmt.Println("There is 1 locked chest and no unlocked chests in the room")
			fmt.Println("To unlock the chest, use a key from the inventory menu")
			return
		} else if numUnlockedChest == totalChests {
			fmt.Println("There are no locked chests and 1 unlocked chest in the room")
			fmt.Println("Would you like to loot it?")
		}
	} else { // totalChests > 1
		if numLockedChests == totalChests {
			fmt.Println("There are", numLockedChests, "locked chests and no unlocked chests in the room")
			fmt.Println("To unlock the chests, use a key or keys from the inventory menu")
			return
		} else if numLockedChests == 0 {
			fmt.Println("There are no locked chests and", numUnlockedChest, "unlocked chests in the room")
			fmt.Println("Would you like to loot them all?")
		} else { // one or more for both
			if numLockedChests == 1 {
				fmt.Println("There is 1 locked chest and", numUnlockedChest, "unlocked chests in the room")
				fmt.Println("To unlock the locked chest, use a key from the inventory menu")
				fmt.Println("Would you like to loot all the unlocked chests?")
			} else if numUnlockedChest == 1 {
				fmt.Println("There are", numLockedChests, "locked chests and 1 unlocked chest in the room")
				fmt.Println("To unlock the locked chests, use a key from the inventory menu")
				fmt.Println("Would you like to loot the unlocked chest?")
			} else {
				fmt.Println("There are", numLockedChests, "locked chests and", numUnlockedChest, "unlocked chests in the room")
				fmt.Println("To unlock the locked chests, use a key from the inventory menu")
				fmt.Println("Would you like to loot all the unlocked chests?")
			}
		}
	}
	var choice int8
	fmt.Println("  1: Yes")
	fmt.Println("Any: No")
	_, err := fmt.Scanln(&choice)
	if err != nil {
		fmt.Println("An error occured while reading your choice in, please try again: ", err)
	}
	if choice == 1 {
		if p.inventory.isFull() {
			if numUnlockedChest == 1 {
				fmt.Println("Your inventory is full.\nTo loot the chest in this room, discard an item to free up space")
			} else {
				fmt.Println("Your inventory is full.\nTo loot the chests in this room, discard an item to free up space")
			}
			return
		} else if numUnlockedChest > p.inventory.slotsNotUsed() {
			fmt.Println("There are more chests than your inventoy has space for.")
			fmt.Println("Only some chests will be looted. To loot them all, discard items to free inventory space")
		}
		count := 0
		for _, chest := range p.currentRoom.chests {
			if chest == nil {
				continue
			}
			if chest.locked {
				continue
			}
			if p.inventory.addItem(chest.item) {
				chest.item = nil
				count++
			}
			if p.inventory.isFull() {
				break
			}
		}
		if count == 1 {
			fmt.Println("Looted 1 chest")
		} else {
			fmt.Println("Looted", count, "chests")
		}
	} else {
		fmt.Println("You can come back to loot the chests at any time")
	}
}

func (p *Player) printFightingChoices() bool {
	// temp
	return true
}

func (p *Player) printPlayerStats() {
	fmt.Println("This feature is not currently implemented", "Plyr stats")
	fmt.Println("\nPlayer Stats:")
	fmt.Println("Health   =", p.health)
	fmt.Println("Defense  =", p.defense)
	fmt.Println("Strength =", p.strength)
}

func (p *Player) printMoveChoices() {
	var choice int8
	valid := false
	for !valid {
		fmt.Println("\nWhere would you like to go?")
		if p.currentRoom.canLeaveFrom(UP) {
			fmt.Println("1. UP")
		}
		if p.currentRoom.canLeaveFrom(DOWN) {
			fmt.Println("2. DOWN")
		}
		if p.currentRoom.canLeaveFrom(LEFT) {
			fmt.Println("3. LEFT")
		}
		if p.currentRoom.canLeaveFrom(RIGHT) {
			fmt.Println("4. RIGHT")
		}

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("An error occured while reading your choice in, please try again: ", err)
			continue
		}

		choice-- // due to directions being index 0 based and prints being index 1 based
		dir := Direction(choice)
		if p.currentRoom.canLeaveFrom(dir) {
			switch dir {
			case UP:
				p.loc.add(&Location{0, -1})
				if DEBUG_MODE {
					fmt.Println("UP")
				}
				valid = true
			case DOWN:
				p.loc.add(&Location{0, 1})
				if DEBUG_MODE {
					fmt.Println("DOWN")
				}
				valid = true
			case LEFT:
				p.loc.add(&Location{-1, 0})
				if DEBUG_MODE {
					fmt.Println("LEFT")
				}
				valid = true
			case RIGHT:
				p.loc.add(&Location{1, 0})
				if DEBUG_MODE {
					fmt.Println("RIGHT")
				}
				valid = true
			default:
				fmt.Println("Invalid Input, try again")
			}
		} else {
			fmt.Println("Invalid Input, try again")
		}
	}
	fmt.Println("You have entered a new room")
	if DEBUG_MODE {
		p.debugPrintLoc()
	}
	p.movedLast = true
}

func (p *Player) printInventoryChoices() (turnConsumed bool) {
	var choice int8
	done := false
	for !done {
		fmt.Println("\nWhat inventory action would you like to do?")
		fmt.Println("1. View Inventory")
		fmt.Println("2. Use Item")
		fmt.Println("3. Equip Item")
		fmt.Println("4. Discard Item")
		fmt.Println("5. Leave Inventory")

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("An error occured while reading your choice in, please try again: ", err)
			continue
		}

		switch choice {
		case cheatInputNumber:
			p.doCheatLoop()
		case 1:
			p.inventory.printFullInventory()
		case 2:
			validIn := false
			if p.inventory.slotsUsed() == 0 {
				fmt.Println("There are no items in your inventory")
				break
			}

			if p.inventory.numUseables() <= 0 {
				fmt.Println("There are no useable items in your inventory")
				break
			}

			for !validIn {
				fmt.Println("\nWhich item would you like to use? (Select by number):")
				p.inventory.printItemInventory()
				_, err := fmt.Scanln(&choice)
				if err != nil {
					fmt.Println("An error occured while reading your choice in, please try again: ", err)
					continue
				}

				if choice >= 0 && choice < inventorySize {
					if item, ok := p.inventory.isUseable(int(choice)); ok {
						switch item.iType {
						case KEY:
							numLocked := p.currentRoom.getNumLockedChests()
							if numLocked > 0 {
								if item.effect > float64(numLocked) {
									fmt.Println("Unlocking all chests")
									item.effect -= float64(numLocked)
									p.currentRoom.unlockChests(numLocked)
									fmt.Printf("This key can unlcok %3.1f more locked chests", item.effect)
									if item.effect > 1 {
										fmt.Print("s")
									}
									fmt.Println()
								} else {

								}
							} else {
								fmt.Println("There are no locked chests in this room, this item cannot be used.")
								validIn = true
								// valid input, but kick them back to the inventory choices list
							}
						case HEALTH: // TODO
							fallthrough
						case INSTANT_DAMAGE:
							fmt.Println("This feature is not implemented for this item type.")
							fmt.Println("This will be impelented when enemies are implemented.")
							continue
						default:
							fmt.Println("Impossible case: Default case from inv.isUseable")
							continue
						}
						validIn = true
						turnConsumed = true
						done = true
					} else {
						fmt.Println("The selected item is not a useable item")
						continue
					}
				} else {
					fmt.Println("Selected index does not exist.")
					fmt.Printf("Please pick from the range 0-%-2d\n", inventorySize)
				}
			}
		case 3:
			validIn := false
			if p.inventory.slotsUsed() == 0 {
				fmt.Println("There are no items in your inventory")
				break
			}

			if p.inventory.numEquipables() <= 0 {
				fmt.Println("There are no equipable items in your inventory")
				break
			}

			var tempItem *Item

			// TODO: when there are more than just armor equips, this will have to change
			if p.inventory.armorSlot != nil {
				fmt.Println("There is already an equiped ARMOR item.")
				fmt.Println("It must be unequipped before a new ARMOR item can be equiped.")
				fmt.Println("Would you like to unequip it?")
				fmt.Println("  1: Yes")
				fmt.Println("Any: No")
				_, err := fmt.Scanln(&choice)
				if err != nil {
					fmt.Println("An error occured while reading your choice in, please try again: ", err)
					break
				}

				if choice == 1 {
					if p.inventory.isFull() {
						tempItem = p.inventory.armorSlot
					} else {
						bo := p.inventory.addItem(p.inventory.armorSlot)
						if !bo {
							// ERROR
							fmt.Println("THIS ALSO SHOULDNT BE POSSIBLE BUT IM CATCHING IT ANYWAY")
							fmt.Println("inventoryChoices() case 3: unequipIetm")
						} else {
							p.inventory.armorSlot = nil
						}
					}
				} else {
					fmt.Println("Canceling equip process")
					break
				}
			}

			for !validIn {
				fmt.Println("\nWhich item would you like to equip? (Select by number):")
				p.inventory.printItemInventory()
				_, err := fmt.Scanln(&choice)
				if err != nil {
					fmt.Println("An error occured while reading your choice in, please try again: ", err)
					continue
				}

				if choice >= 0 && choice < inventorySize {
					if item, ok := p.inventory.isEquipable(int(choice)); ok {
						// not currently necessary
						switch item.iType {
						case ARMOR:
							fmt.Print("Equipped item: ")
							item.print()
							p.inventory.armorSlot = item
							p.inventory.itemSlots[choice] = nil
						default:
							fmt.Println("Impossible case: Default case from inv.isEquipable")
							continue
						}
						if tempItem != nil {
							p.inventory.itemSlots[choice] = tempItem
						}
						validIn = true
						turnConsumed = true
						done = true
					} else {
						fmt.Println("The selected item is not an equipable item")
						continue
					}
				} else {
					fmt.Println("Selected index does not exist.")
					fmt.Printf("Please pick from the range 0-%-2d\n", inventorySize)
				}
			}
		case 4:
			validIn := false
			if p.inventory.slotsUsed() == 0 {
				fmt.Println("There are no items in your inventory")
				break
			}

			for !validIn {
				fmt.Println("\nWhich item would you like to discard? (Select by number):")
				fmt.Println("Enter -1 to cancel")
				p.inventory.printItemInventory()
				_, err := fmt.Scanln(&choice)
				if err != nil {
					fmt.Println("An error occured while reading your choice in, please try again: ", err)
					continue
				}
				if choice == -1 {
					validIn = true
					fmt.Println("Canceling Discard Process")
					continue
				}

				index := int(choice)

				if index >= 0 && index < inventorySize {
					if p.inventory.itemSlots[index] == nil {
						fmt.Println("There is no item in that slot")
						continue
					}
					fmt.Println("You are about to discard the following item:")
					p.inventory.printItemAt(index)
					fmt.Println("\nDo you wish to continue?")
					fmt.Println("  1: Yes, discard the item")
					fmt.Println("Any: No, keep the item")

					_, err := fmt.Scanln(&choice)
					if err != nil {
						fmt.Println("An error occured while reading your choice in, please try again: ", err)
						continue
					}

					if choice == 1 {
						fmt.Println("Discarded item")
						p.inventory.itemSlots[index] = nil
						validIn = true
						turnConsumed = true
						done = true
					} else {
						fmt.Println("Item will not be discarded")
						validIn = true
					}
				} else {
					fmt.Println("Selected index does not exist.")
					fmt.Printf("Please pick from the range 0-%-2d\n", inventorySize)
				}
			}
		case 5:
			done = true
		default:
			fmt.Println("Invalid choice")
		}
	}
	return
}

func (p *Player) doCheatLoop() {
	var choice int8
	valid := false
	for !valid {
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("An error occured while reading your choice in, please try again: ", err)
			continue
		}
		switch choice {
		case -1: // leave cheat loop
			valid = true
		case 1: // give item
			effect := 0.0
			fmt.Scanln(&choice, &effect)
			item := NewItem(ItemType(choice), effect)
			success := p.inventory.addItem(item)
			if success {
				fmt.Printf("Given item %+v\n", item)
			} else {
				fmt.Println("failed to give item")
			}
		// todo more cheat options
		default:
			// do nothing
		}
	}
}

func (p *Player) debugPrintLoc() {
	fmt.Println("Player Loc:", p.loc.x, p.loc.y)
}

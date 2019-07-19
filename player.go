package main

import "fmt"

const (
	BASE_PLAYER_HEALTH = 100
)

type Player struct {
	movedLast   bool
	loc         *Location
	currentRoom *Room
	inventory   *Inventory
}

func NewPlayer(current *Room, loc *Location) *Player {
	p := new(Player)
	p.currentRoom = current
	p.loc = loc
	p.inventory = NewInventory()
	return p
}

func (p *Player) update() bool {
	// var reset
	p.movedLast = false

	// player turn
	run := p.printChoices()

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
	fmt.Println("\nDebug prints. TODO")
	if DEBUG_MODE {
		fmt.Println("Player Location: x=", p.loc.x, " y=", p.loc.y)
		fmt.Println("Room Location: x=", p.currentRoom.loc.x, " y=", p.currentRoom.loc.y)
	} else {
		fmt.Printf("Location: %+v\n", *p.loc)
	}
	fmt.Printf("Room Type=%s\n", getPrintStringFromRoomType(p.currentRoom.rType))
	fmt.Println("Room ID", p.currentRoom.id)
}

func (p *Player) printPlayerStats() {
	fmt.Println("This feature is not currently implemented")
}

func (p *Player) printMoveChoices() {
	var choice int8
	valid := false
	for !valid {
		fmt.Println("\nWhere would you like to go?")
		if p.currentRoom.CanLeaveFrom(UP) {
			fmt.Println("1. UP")
		}
		if p.currentRoom.CanLeaveFrom(DOWN) {
			fmt.Println("2. DOWN")
		}
		if p.currentRoom.CanLeaveFrom(LEFT) {
			fmt.Println("3. LEFT")
		}
		if p.currentRoom.CanLeaveFrom(RIGHT) {
			fmt.Println("4. RIGHT")
		}

		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("An error occured while reading your choice in, please try again: ", err)
			continue
		}

		choice-- // due to directions being index 0 based and prints being index 1 based
		dir := Direction(choice)
		if p.currentRoom.CanLeaveFrom(dir) {
			switch dir {
			case UP:
				p.loc.Add(&Location{0, -1})
				if DEBUG_MODE {
					fmt.Println("UP")
				}
				valid = true
			case DOWN:
				p.loc.Add(&Location{0, 1})
				if DEBUG_MODE {
					fmt.Println("DOWN")
				}
				valid = true
			case LEFT:
				p.loc.Add(&Location{-1, 0})
				if DEBUG_MODE {
					fmt.Println("LEFT")
				}
				valid = true
			case RIGHT:
				p.loc.Add(&Location{1, 0})
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

			/*
				fmt.Println("This feature is currently not implemented.\nThis option is a place holder for furture features.")
				fmt.Println("This will be impelented when enemies are implemented.")
				validIn = true
			*/
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
							// TODO relies on rooms having chests
							fmt.Println("Feautre will be implemented soon")
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

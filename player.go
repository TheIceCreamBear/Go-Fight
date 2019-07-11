package main

import "fmt"

const (
	BASE_PLAYER_HEALTH = 100
)

type Player struct {
	movedLast   bool
	loc         *Location
	currentRoom *Room
}

func NewPlayer(current *Room, loc *Location) *Player {
	p := new(Player)
	p.currentRoom = current
	p.loc = loc
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
	choice := 0

	for {
		fmt.Println("\nWhat would you like to do?")
		fmt.Println("1. Explore current room")
		fmt.Println("2. Move to another room")
		fmt.Println("3. Exit")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("An error occured while reading your choice in, please try again: ", err)
			continue
		}
		/*
			_, err := fmt.Scanf("%d", &choice)
			if err != nil {
				fmt.Println("An error occured while reading your choice in, please try again: ", err)
				continue
			}
		*/
		switch choice {
		case 1:
			fmt.Println("\nDebug prints. TODO")
			fmt.Println("Player Location: x=", p.loc.x, " y=", p.loc.y)
			fmt.Println("Room Location: x=", p.currentRoom.loc.x, " y=", p.currentRoom.loc.y)
			fmt.Println("Room ID", p.currentRoom.id)
			return true
		case 2:
			p.printMoveChoices()
			return true
		case 3:
			return false
		default:
			fmt.Println("Invalid Input, try again")
		}
	}
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
		/*
			_, err := fmt.Scanf("%d", &choice)
			if err != nil {
				fmt.Println("An error occured while reading your choice in, please try again: ", err)
				continue
			}
		*/

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

func (p *Player) debugPrintLoc() {
	fmt.Println("Player Loc:", p.loc.x, p.loc.y)
}

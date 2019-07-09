package main

import "fmt"

type Game struct {
	rooms [GameHeight][GameWidth]Room
}

const (
	UP    int8 = iota
	DOWN  int8 = iota
	LEFT  int8 = iota
	RIGHT int8 = iota
)

const GameWidth int64 = 11
const GameHeight int64 = 11
const DEBUG_MODE bool = true

func main() {
	game := new(Game)
	initRooms(game)
	if DEBUG_MODE {
		printRooms(game)
	}

	// Temp for 11x11
	var playerStartX int64 = 5
	var playerStartY int64 = 5

	plyr := NewPlayer(&game.rooms[playerStartY][playerStartX], &Location{playerStartX, playerStartY})

	for {
		run := plyr.update()
		if !run {
			break
		}
		if plyr.movedLast {
			plyr.currentRoom = &game.rooms[plyr.loc.y][plyr.loc.x]
		}
	}
}

func initRooms(game *Game) {
	roomID := int64(0)
	// TODO rng room types

	// end room type loops

	// Init doors
	// TODO dynamic doors
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]
			current.loc = Location{x, y}
			current.id = roomID
			roomID++

			up := Door{false, false}
			down := Door{false, false}
			left := Door{false, false}
			right := Door{false, false}

			if y == 0 {
				down.exists = true
			} else if y == GameHeight-1 {
				up.exists = true
			} else {
				up.exists = true
				down.exists = true
			}

			if x == 0 {
				right.exists = true
			} else if x == GameWidth-1 {
				left.exists = true
			} else {
				left.exists = true
				right.exists = true
			}
			if DEBUG_MODE {
				fmt.Println("Door init  ", x, y, up, down, left, right)
			}
			current.InitDoors(up, down, left, right)
		}
	}
}

func printRooms(game *Game) {
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - -")
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]
			fmt.Printf("%T\n", current)
			fmt.Println("id", current.id)
			// prints whole struct
			// fmt.Println(&current)
			fmt.Println("Location", "x", current.loc.x, "y", current.loc.y)
			fmt.Println("Vars", "x", x, "y", y)
			fmt.Println("Doors-Method:", current.CanLeaveFrom(UP), current.CanLeaveFrom(LEFT), current.CanLeaveFrom(RIGHT), current.CanLeaveFrom(DOWN))
			fmt.Println("Doors:", current.dUP, current.dLEFT, current.dRIGHT, current.dDOWN)
			fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - -")
		}
	}
}

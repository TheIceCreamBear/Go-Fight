package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

type Game struct {
	rooms   [GameRaidus*2 + 1][GameRaidus*2 + 1]Room
	chances map[RoomType]float64
}

// directions
const (
	UP    Direction = iota
	DOWN  Direction = iota
	LEFT  Direction = iota
	RIGHT Direction = iota
)

type Direction int8

const GameRaidus int64 = 30 // 30 tiles on each side

// left included for legacy use
const (
	GameWidth  int64 = GameRaidus*2 + 1
	GameHeight int64 = GameRaidus*2 + 1
)

var DEBUG_MODE bool = true

const HowSticky float64 = 0.25

func main() {
	args := os.Args[1:]
	debugStr := args[0]
	debugBool, err := strconv.ParseBool(debugStr)
	if err == nil {
		DEBUG_MODE = debugBool
	}

	game := new(Game)
	game.initRoomTypeChances()
	game.initDefaultRoomType()
	initRooms(game)
	game.calcStats()
	if DEBUG_MODE {
		printRooms(game)
	}

	// Temp for 11x11
	var playerStartX = GameRaidus
	var playerStartY = GameRaidus

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
	// rng room generation spiraling out from the center
	game.rooms[GameRaidus][GameRaidus].rType = START
	for r := int64(1); r <= GameRaidus; r++ {
		for t := int64(0); t < r*8; t++ {
			var x int64
			var y int64
			if t < 2*r {
				x = GameRaidus - r + t
				y = GameRaidus - r
			} else if t < 4*r {
				x = GameRaidus + r
				y = GameRaidus - (3 * r) + t
			} else if t < 6*r {
				x = GameRaidus + (5 * r) - t
				y = GameRaidus + r
			} else {
				x = GameRaidus - r
				y = GameRaidus + (7 * r) - t
			}

			game.rooms[y][x].rType = initRoomType(game, x, y)
		}
	}
	// end room type loops

	fmt.Println("=====================END TYPE=====================")
	pause()
	fmt.Println("=====================DOORS=====================")
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

func initRoomType(game *Game, x int64, y int64) RoomType {
	var adjecents [4]RoomType
	if y > 0 {
		adjecents[0] = game.rooms[y-1][x].rType
	} else {
		adjecents[0] = -1
	}
	if y < GameRaidus*2 {
		adjecents[1] = game.rooms[y+1][x].rType
	} else {
		adjecents[1] = -1
	}
	if x < GameRaidus*2 {
		adjecents[2] = game.rooms[y][x+1].rType
	} else {
		adjecents[2] = -1
	}
	if x > 0 {
		adjecents[3] = game.rooms[y][x-1].rType
	} else {
		adjecents[3] = -1
	}

	chances := make(map[RoomType]float64, MYSTIC+1)
	stickyLeft := 1.0

	for i := 0; i < 4; i++ {
		if adjecents[i] == -1 {
			continue
		}
		if adjecents[i] == START {
			if DEBUG_MODE {
				fmt.Println("\nYeet, next to start")
				fmt.Println("Location", Location{x, y})
				fmt.Println("Adjacents", adjecents)
			}
			return HALLWAY
		} else {
			if val, ok := chances[adjecents[i]]; ok {
				chances[adjecents[i]] = val + HowSticky
			} else {
				chances[adjecents[i]] = HowSticky
			}
			stickyLeft -= HowSticky
		}
	}

	rTypes := GetGenetateableTypes()
	for _, rType := range rTypes {
		chance := 0.0
		if val, ok := chances[rType]; ok {
			chance += val
		}
		chance += game.chances[rType] * stickyLeft
		chances[rType] = chance
	}

	keys := make([]RoomType, len(chances))
	i := 0
	for key := range chances {
		keys[i] = key
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	chance := 0.0
	chanceNeeded := rand.Float64()
	if DEBUG_MODE {
		defer func() {
			fmt.Println("Location", Location{x, y})
			fmt.Println("Chance : ChanceNeeded ->", chance, ":", chanceNeeded)
			fmt.Println("Chances map: Type map[RoomType]float64", chances)
			//pause()
		}()
	}
	for _, val := range keys {
		chance += chances[val]
		if chance > chanceNeeded {
			if DEBUG_MODE {
				fmt.Println("\nYate, returning from good spot", val)
			}
			return val
		}
	}
	if DEBUG_MODE {
		fmt.Println("\nYoot, chance wasnt high enough")
	}
	return HALLWAY
}

func (game *Game) initRoomTypeChances() {
	game.chances = make(map[RoomType]float64, MYSTIC+1)
	game.chances[START] = 0
	game.chances[HALLWAY] = 0.5
	game.chances[GREAT_HALL] = 0.2
	game.chances[DUNGEON] = 0.15
	game.chances[CHEST] = 0.1
	game.chances[MYSTIC] = 0.05
}

func (game *Game) initDefaultRoomType() {
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]
			current.rType = -1
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
			fmt.Println("Doors:", current.dUp, current.dLeft, current.dRight, current.dDown)
			fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - -")
		}
	}
	fmt.Println("===============================================================")
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]
			fmt.Print(getPrintStringFromRoomType(current.rType))
		}
		fmt.Println("")
	}
	fmt.Println("===============================================================")
}

func (game *Game) calcStats() {
	total := GameWidth * GameHeight
	numS, numH, numG, numD, numC, numM := 0, 0, 0, 0, 0, 0
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]
			switch current.rType {
			case START:
				numS++
			case HALLWAY:
				numH++
			case GREAT_HALL:
				numG++
			case DUNGEON:
				numD++
			case CHEST:
				numC++
			case MYSTIC:
				numM++
			default:
			}
		}
	}
	fmt.Println("=====================Stats=====================")
	fmt.Printf("START        %6d/%-7d = %9.6f%%\n", numS, total, (float64(numS) / float64(total) * 100.0))
	fmt.Printf("HALLWAY      %6d/%-7d = %9.6f%%\n", numH, total, (float64(numH) / float64(total) * 100.0))
	fmt.Printf("GREAT_HALL   %6d/%-7d = %9.6f%%\n", numG, total, (float64(numG) / float64(total) * 100.0))
	fmt.Printf("DUNGEON      %6d/%-7d = %9.6f%%\n", numD, total, (float64(numD) / float64(total) * 100.0))
	fmt.Printf("CHEST        %6d/%-7d = %9.6f%%\n", numC, total, (float64(numC) / float64(total) * 100.0))
	fmt.Printf("MYSTIC       %6d/%-7d = %9.6f%%\n", numM, total, (float64(numM) / float64(total) * 100.0))
	pause()
}

func pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

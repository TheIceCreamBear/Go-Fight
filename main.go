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
	moves   []*Move
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

var DEBUG_MODE = false

const cheatInputNumber int8 = -111
const HowSticky float64 = 0.25

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		debugStr := args[0]
		debugBool, err := strconv.ParseBool(debugStr)
		if err == nil {
			DEBUG_MODE = debugBool
		}
	}

	game := new(Game)
	game.initRoomTypeChances()
	game.initDefaultRoomType()
	game.initRooms()
	game.initRoomChests()
	game.initEnemies()
	game.initMoves()
	game.calcStats()
	if DEBUG_MODE {
		printRooms(game)
	}

	var playerStartX = GameRaidus
	var playerStartY = GameRaidus

	plyr := NewPlayer(&game.rooms[playerStartY][playerStartX], &Location{playerStartX, playerStartY}, game.moves[:3])

	for {
		run := plyr.update()
		if !run {
			break
		}
		if plyr.movedLast {
			plyr.currentRoom = &game.rooms[plyr.loc.y][plyr.loc.x]
			if plyr.currentRoom.getNumEnemies() > 0 {
				fmt.Println("\n\nYou have Encountered an Enemy!\nPrepare to Fight!")
			}
		}
	}
}

func (game *Game) initRooms() {
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
			current.initDoors(up, down, left, right)
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

	rTypes := getGenetateableTypes()
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

func (game *Game) initRoomChests() {
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]
			current.initChests()
		}
	}
}

func (game *Game) initEnemies() {
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

			game.rooms[y][x].initEnemies(x, y, r)
		}
	}
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

func (game *Game) initMoves() {
	game.moves = make([]*Move, 3)
	game.moves[0] = NewMove(3.0, 6.0, "Punch")
	game.moves[1] = NewMove(5.0, 10.0, "Kick")
	game.moves[2] = NewMove(9.0, 15.0, "Special")
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
			fmt.Println("Doors-Method:", current.canLeaveFrom(UP), current.canLeaveFrom(LEFT), current.canLeaveFrom(RIGHT), current.canLeaveFrom(DOWN))
			fmt.Println("Doors:", current.dUp, current.dLeft, current.dRight, current.dDown)
			fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - -")
		}
	}
	fmt.Println("===============================================================")
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]
			fmt.Print(getPrintCharFromRoomType(current.rType))
		}
		fmt.Println("")
	}
	fmt.Println("===============================================================")
}

func (game *Game) calcStats() {
	total := GameWidth * GameHeight
	s, h, g, d, c, m := 0, 0, 0, 0, 0, 0
	chests, lChests := 0, 0
	rWch, rWe, rTot := 0, 0, 0
	ep, en, eb, em, et := 0, 0, 0, 0, 0
	// items, each new time on new line, each diff effective type new var
	itemTotal := 0
	k1, k2, k3, kT := 0, 0, 0, 0
	a1, a2, a3, a4, aT := 0, 0, 0, 0, 0
	h1, h2, h3, h4, hT := 0, 0, 0, 0, 0
	d1, d2, d3, d4, dT := 0, 0, 0, 0, 0
	for y := int64(0); y < GameHeight; y++ {
		for x := int64(0); x < GameWidth; x++ {
			current := &game.rooms[y][x]

			rTot++
			currentNumChests := current.getNumChests()
			if currentNumChests > 0 {
				rWch++
				chests += currentNumChests
				for _, chest := range current.chests {
					if chest == nil {
						continue
					}
					if chest.locked {
						lChests++
					}
					item := chest.item
					if item != nil {
						itemTotal++
						switch item.iType {
						case KEY:
							kT++
							switch item.effect {
							case 1:
								k1++
							case 2:
								k2++
							case 3:
								k3++
							default:
							}
						case ARMOR:
							aT++
							switch item.effect {
							case 1:
								a1++
							case 2:
								a2++
							case 3:
								a3++
							case 4:
								a4++
							default:
							}
						case HEALTH:
							hT++
							switch item.effect {
							case 20:
								h1++
							case 50:
								h2++
							case 100:
								h3++
							case 200:
								h4++
							default:
							}
						case INSTANT_DAMAGE:
							dT++
							switch item.effect {
							case 20:
								d1++
							case 50:
								d2++
							case 100:
								d3++
							case 200:
								d4++
							default:
							}
						default:
						}
					}
				}
			}

			if current.getNumEnemies() > 0 {
				rWe++
				et += current.getNumEnemies()
				for _, val := range current.enemies {
					if val == nil {
						continue
					}
					switch val.eType {
					case PEON:
						ep++
					case NORMAL:
						en++
					case BRUTE:
						eb++
					case E_MYSTIC:
						em++
					}
				}
			}

			switch current.rType {
			case START:
				s++
			case HALLWAY:
				h++
			case GREAT_HALL:
				g++
			case DUNGEON:
				d++
			case CHEST:
				c++
			case MYSTIC:
				m++
			default:
			}
		}
	}
	fmt.Println("=====================Stats=====================")
	fmt.Println("---------------------RTYPE---------------------")
	fmt.Printf("START        %6d/%-7d = %9.6f%%\n", s, total, (float64(s) / float64(total) * 100.0))
	fmt.Printf("HALLWAY      %6d/%-7d = %9.6f%%\n", h, total, (float64(h) / float64(total) * 100.0))
	fmt.Printf("GREAT_HALL   %6d/%-7d = %9.6f%%\n", g, total, (float64(g) / float64(total) * 100.0))
	fmt.Printf("DUNGEON      %6d/%-7d = %9.6f%%\n", d, total, (float64(d) / float64(total) * 100.0))
	fmt.Printf("CHEST        %6d/%-7d = %9.6f%%\n", c, total, (float64(c) / float64(total) * 100.0))
	fmt.Printf("MYSTIC       %6d/%-7d = %9.6f%%\n", m, total, (float64(m) / float64(total) * 100.0))
	fmt.Println("---------------------ENEMY---------------------")
	fmt.Printf("RWE/RTot     %6d/%-7d = %9.6f%%\n", rWe, rTot, (float64(rWe) / float64(rTot) * 100.0))
	fmt.Printf("PEON         %6d/%-7d = %9.6f%%\n", ep, et, (float64(ep) / float64(et) * 100.0))
	fmt.Printf("NORMAL       %6d/%-7d = %9.6f%%\n", en, et, (float64(en) / float64(et) * 100.0))
	fmt.Printf("BRUTE        %6d/%-7d = %9.6f%%\n", eb, et, (float64(eb) / float64(et) * 100.0))
	fmt.Printf("MYSTIC       %6d/%-7d = %9.6f%%\n", em, et, (float64(em) / float64(et) * 100.0))
	fmt.Println("---------------------CHEST---------------------")
	fmt.Printf("RWC/RTot     %6d/%-7d = %9.6f%%\n", rWch, rTot, (float64(rWch) / float64(rTot) * 100.0))
	fmt.Printf("Locked/Total %6d/%-7d = %9.6f%%\n", lChests, chests, (float64(lChests) / float64(chests) * 100.0))
	fmt.Println("---------------------ITEMS---------------------")
	fmt.Printf("KEY          %6d/%-7d = %9.6f%%\n", kT, itemTotal, (float64(kT) / float64(itemTotal) * 100.0))
	fmt.Printf("ARMOR        %6d/%-7d = %9.6f%%\n", aT, itemTotal, (float64(aT) / float64(itemTotal) * 100.0))
	fmt.Printf("HEALTH       %6d/%-7d = %9.6f%%\n", hT, itemTotal, (float64(hT) / float64(itemTotal) * 100.0))
	fmt.Printf("DAMAGE       %6d/%-7d = %9.6f%%\n", dT, itemTotal, (float64(dT) / float64(itemTotal) * 100.0))
	fmt.Println("----------------------KEY----------------------")
	fmt.Printf("KEY 1        %6d/%-7d = %9.6f%%\n", k1, kT, (float64(k1) / float64(kT) * 100.0))
	fmt.Printf("KEY 2        %6d/%-7d = %9.6f%%\n", k2, kT, (float64(k2) / float64(kT) * 100.0))
	fmt.Printf("KEY 3        %6d/%-7d = %9.6f%%\n", k3, kT, (float64(k3) / float64(kT) * 100.0))
	fmt.Println("---------------------ARMOR---------------------")
	fmt.Printf("ARMOR 1      %6d/%-7d = %9.6f%%\n", a1, aT, (float64(a1) / float64(aT) * 100.0))
	fmt.Printf("ARMOR 2      %6d/%-7d = %9.6f%%\n", a2, aT, (float64(a2) / float64(aT) * 100.0))
	fmt.Printf("ARMOR 3      %6d/%-7d = %9.6f%%\n", a3, aT, (float64(a3) / float64(aT) * 100.0))
	fmt.Printf("ARMOR 4      %6d/%-7d = %9.6f%%\n", a4, aT, (float64(a4) / float64(aT) * 100.0))
	fmt.Println("---------------------HEALTH--------------------")
	fmt.Printf("HEALTH 1     %6d/%-7d = %9.6f%%\n", h1, hT, (float64(h1) / float64(hT) * 100.0))
	fmt.Printf("HEALTH 2     %6d/%-7d = %9.6f%%\n", h2, hT, (float64(h2) / float64(hT) * 100.0))
	fmt.Printf("HEALTH 3     %6d/%-7d = %9.6f%%\n", h3, hT, (float64(h3) / float64(hT) * 100.0))
	fmt.Printf("HEALTH 4     %6d/%-7d = %9.6f%%\n", h4, hT, (float64(h4) / float64(hT) * 100.0))
	fmt.Println("---------------------DAMAGE--------------------")
	fmt.Printf("DAMAGE 1     %6d/%-7d = %9.6f%%\n", d1, dT, (float64(d1) / float64(dT) * 100.0))
	fmt.Printf("DAMAGE 2     %6d/%-7d = %9.6f%%\n", d2, dT, (float64(d2) / float64(dT) * 100.0))
	fmt.Printf("DAMAGE 3     %6d/%-7d = %9.6f%%\n", d3, dT, (float64(d3) / float64(dT) * 100.0))
	fmt.Printf("DAMAGE 4     %6d/%-7d = %9.6f%%\n", d4, dT, (float64(d4) / float64(dT) * 100.0))
	fmt.Println("======================END======================")
	pause()
}

func pause() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

package main

import "fmt"

// room types
const (
	START      RoomType = iota
	HALLWAY    RoomType = iota
	GREAT_HALL RoomType = iota
	DUNGEON    RoomType = iota
	CHEST      RoomType = iota
	MYSTIC     RoomType = iota
)

type RoomType int8

type Room struct {
	rType  RoomType
	id     int64
	loc    Location
	dUp    Door
	dDown  Door
	dLeft  Door
	dRight Door
}

func GetGenetateableTypes() [6]RoomType {
	return [6]RoomType{START, HALLWAY, GREAT_HALL, DUNGEON, CHEST, MYSTIC}
}

func (r *Room) CanLeaveFrom(direction Direction) bool {
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

func (room *Room) InitDoors(u Door, d Door, l Door, r Door) {
	if DEBUG_MODE {
		fmt.Println("init doors ", room.loc.x, room.loc.y, u, d, l, r)
	}
	room.dUp = Door{u.exists, u.locked}
	room.dDown = Door{d.exists, d.locked}
	room.dLeft = Door{l.exists, l.locked}
	room.dRight = Door{r.exists, r.locked}
	if DEBUG_MODE {
		fmt.Println("post init doors", room.dUp, room.dDown, room.dLeft, room.dRight)
	}
}

func getPrintStringFromRoomType(rType RoomType) string {
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

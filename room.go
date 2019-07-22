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

func getPrintStringFromRoomType(rType RoomType) string {
	switch rType {
	case START:
		return "Start room"
	case HALLWAY:
		return "Hallway"
	case GREAT_HALL:
		return "Great Hall"
	case DUNGEON:
		return "Dungeon"
	case CHEST:
		return "Chest Room"
	case MYSTIC:
		return "Mystical"
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

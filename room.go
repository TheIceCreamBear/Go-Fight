package main

import "fmt"

type Room struct {
	id     int64
	loc    Location
	dUP    Door
	dDOWN  Door
	dLEFT  Door
	dRIGHT Door
}

func (r *Room) CanLeaveFrom(direction int8) bool {
	switch direction {
	case UP:
		return r.dUP.exists /* && !r.dUP.locked */
	case DOWN:
		return r.dDOWN.exists /* && !r.dDOWN.locked */
	case LEFT:
		return r.dLEFT.exists /* && !r.dLEFT.locked */
	case RIGHT:
		return r.dRIGHT.exists /* && !r.dRIGHT.locked */
	default:
		fmt.Println("D E F A U L T  C A S E ")
		return false
	}
}

func (room *Room) InitDoors(u Door, d Door, l Door, r Door) {
	if DEBUG_MODE {
		fmt.Println("init doors ", room.loc.x, room.loc.y, u, d, l, r)
	}
	room.dUP = Door{u.exists, u.locked}
	room.dDOWN = Door{d.exists, d.locked}
	room.dLEFT = Door{l.exists, l.locked}
	room.dRIGHT = Door{r.exists, r.locked}
	if DEBUG_MODE {
		fmt.Println("post init doors", room.dUP, room.dDOWN, room.dLEFT, room.dRIGHT)
	}
}

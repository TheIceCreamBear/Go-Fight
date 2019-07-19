package main

type Location struct {
	x int64
	y int64
}

func (l *Location) add(other *Location) {
	l.x += other.x
	l.y += other.y
}

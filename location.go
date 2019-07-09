package main

type Location struct {
	x int64
	y int64
}

func (l *Location) Add(other *Location) {
	l.x += other.x
	l.y += other.y
}

func (l *Location) Equal(other Location) bool {
	return l.x == other.x && l.y == other.y
}

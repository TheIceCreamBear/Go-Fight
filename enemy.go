package main

import "math/rand"

type EnemyType int8

const (
	PEON     EnemyType = iota
	WARRIOR  EnemyType = iota
	BRUTE    EnemyType = iota
	E_MYSTIC EnemyType = iota
)

const (
	BaseEnemyMinDamage = 5
	BaseEnemyMaxDamage = 10
)

type Enemy struct {
	eType       EnemyType
	health      float64
	strength    float64
	turnCounter int
}

func NewEnemy(eType EnemyType) *Enemy {
	e := new(Enemy)
	e.eType = eType
	switch eType {
	case PEON:
		e.health = 75
		e.strength = 0.75
	case WARRIOR:
		e.health = 100
		e.strength = 1.0
	case BRUTE:
		e.health = 150
		e.strength = 1.25
	case E_MYSTIC:
		e.health = 50
		e.strength = 1.5
	}
	e.turnCounter = 1
	return e
}

func (e *Enemy) getDamageFromAttack() float64 {
	min, max := BaseEnemyMinDamage*e.strength, BaseEnemyMaxDamage*e.strength
	return min + rand.Float64()*(max-min)
}

func getEnemyNameFromType(eType EnemyType) string {
	switch eType {
	case PEON:
		return "Peon"
	case WARRIOR:
		return "Warrior"
	case BRUTE:
		return "Brute"
	case E_MYSTIC:
		return "Mystic"
	default:
		return "Invalid"
	}
}

package main

import (
	"image/color"
)

type Entity struct {
	x              int
	y              int
	char           string
	color          color.RGBA
	Name           string
	BlocksMovement bool
}

func NewEntity(x, y int, char string, col color.RGBA, name string, blocksMovement bool) *Entity {
	return &Entity{
		x:              x,
		y:              y,
		char:           char,
		color:          col,
		Name:           name,
		BlocksMovement: blocksMovement,
	}
}

func (e Entity) Spawn(gameMap *GameMap, x, y int) Entity {
	clone := e
	clone.x = x
	clone.y = y
	gameMap.Entities = append(gameMap.Entities, &clone)
	return clone
}

func (e Entity) X() int {
	return e.x
}

func (e Entity) Y() int {
	return e.y
}

func (e *Entity) SetPostion(pos [2]int) {
	e.x = pos[0]
	e.y = pos[1]
}

func (e Entity) Char() string {
	return e.char
}

func (e Entity) Color() color.RGBA {
	return e.color
}

func (e *Entity) Move(dx, dy int) {
	e.x += dx
	e.y += dy
}

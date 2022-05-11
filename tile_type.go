package main

import (
	"image/color"
)

type Tile struct {
	walkable    bool
	transparent bool
	char        string
	color       color.RGBA
}

func NewFloor() *Tile {
	return &Tile{
		walkable:    true,
		transparent: true,
		char:        string(0xdb),
		color:       color.RGBA{R: 50, G: 50, B: 150, A: 255},
	}
}

func NewWall() *Tile {
	return &Tile{
		walkable:    false,
		transparent: false,
		char:        string(0xdb),
		color:       color.RGBA{R: 0, G: 0, B: 100, A: 255},
	}
}

func (t Tile) Walkable() bool {
	return t.walkable
}

func (t Tile) Transparent() bool {
	return t.transparent
}

func (t Tile) Char() string {
	return t.char
}

func (t Tile) Color() color.RGBA {
	return t.color
}

package main

import (
	"image/color"
)

type Tile struct {
	walkable    bool
	transparent bool
	char        string
	dark        color.RGBA
	light       color.RGBA
	shroud      color.RGBA
}

func NewFloor() *Tile {
	return &Tile{
		walkable:    true,
		transparent: true,
		char:        string(0xdb),
		dark:        color.RGBA{R: 50, G: 50, B: 150, A: 255},
		light:       color.RGBA{R: 200, G: 180, B: 50, A: 255},
		shroud:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

func NewWall() *Tile {
	return &Tile{
		walkable:    false,
		transparent: false,
		char:        string(0xdb),
		dark:        color.RGBA{R: 0, G: 0, B: 100, A: 255},
		light:       color.RGBA{R: 130, G: 110, B: 50, A: 255},
		shroud:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
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

func (t Tile) Dark() color.RGBA {
	return t.dark
}

func (t Tile) Light() color.RGBA {
	return t.light
}

func (t Tile) Shroud() color.RGBA {
	return t.shroud
}

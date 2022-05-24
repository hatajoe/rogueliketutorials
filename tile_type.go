package main

import (
	"image/color"
)

type tile struct {
	Walkable    bool
	Transparent bool
	Char        string
	Dark        color.RGBA
	Light       color.RGBA
	Shroud      color.RGBA
}

func newFloor() *tile {
	return &tile{
		Walkable:    true,
		Transparent: true,
		Char:        string([]rune{0xdb}),
		Dark:        color.RGBA{R: 50, G: 50, B: 150, A: 255},
		Light:       color.RGBA{R: 200, G: 180, B: 50, A: 255},
		Shroud:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

func newWall() *tile {
	return &tile{
		Walkable:    false,
		Transparent: false,
		Char:        string([]rune{0xdb}),
		Dark:        color.RGBA{R: 0, G: 0, B: 100, A: 255},
		Light:       color.RGBA{R: 130, G: 110, B: 50, A: 255},
		Shroud:      color.RGBA{R: 0, G: 0, B: 0, A: 255},
	}
}

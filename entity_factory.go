package main

import "image/color"

func NewPlayer() *Entity {
	return &Entity{
		char: "@",
		color: color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
		Name:           "Player",
		BlocksMovement: true,
	}
}

func NewOrc() *Entity {
	return &Entity{
		char: "o",
		color: color.RGBA{
			R: 63,
			G: 127,
			B: 63,
			A: 255,
		},
		Name:           "Orc",
		BlocksMovement: true,
	}
}

func NewTroll() *Entity {
	return &Entity{
		char: "T",
		color: color.RGBA{
			R: 0,
			G: 127,
			B: 0,
			A: 255,
		},
		Name:           "Troll",
		BlocksMovement: true,
	}
}

package main

import "image/color"

func newPlayer() *actor {
	return newActor(
		0,
		0,
		"@",
		color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
		"Player",
		&fighter{
			baseComponent: &baseComponent{
				Parent: nil,
			},
			MaxHP:   30,
			HP:      30,
			Defense: 2,
			Power:   5,
		},
		newInventory(26),
	)
}

func newOrc() *actor {
	return newActor(
		0,
		0,
		"o",
		color.RGBA{
			R: 63,
			G: 127,
			B: 63,
			A: 255,
		},
		"Orc",
		&fighter{
			baseComponent: &baseComponent{
				Parent: nil,
			},
			MaxHP:   10,
			HP:      10,
			Defense: 0,
			Power:   3,
		},
		newInventory(0),
	)
}

func newTroll() *actor {
	return newActor(
		0,
		0,
		"T",
		color.RGBA{
			R: 0,
			G: 127,
			B: 0,
			A: 255,
		},
		"Troll",
		&fighter{
			baseComponent: &baseComponent{
				Parent: nil,
			},
			MaxHP:   16,
			HP:      16,
			Defense: 1,
			Power:   4,
		},
		newInventory(0),
	)
}

func newHealthPortion() *item {
	return newItem(
		0,
		0,
		"!",
		color.RGBA{
			R: 127,
			G: 0,
			B: 255,
			A: 255,
		},
		"Health Portion",
		newHealingConsumable(4),
	)
}

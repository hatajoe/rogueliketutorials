package main

import "image/color"

func newPlayer() *actor {
	e := newActor(
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
		nil,
		nil,
	)
	e.AI = NewHostileEnemy(e)
	e.Fighter = &fighter{
		baseComponent: &baseComponent{
			Entity: e,
		},
		MaxHP:   30,
		HP:      30,
		Defense: 2,
		Power:   5,
	}
	return e
}

func newOrc() *actor {
	e := newActor(
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
		nil,
		nil,
	)
	e.AI = NewHostileEnemy(e)
	e.Fighter = &fighter{
		baseComponent: &baseComponent{
			Entity: e,
		},
		MaxHP:   10,
		HP:      10,
		Defense: 0,
		Power:   3,
	}
	return e
}

func newTroll() *actor {
	e := newActor(
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
		nil,
		nil,
	)
	e.AI = NewHostileEnemy(e)
	e.Fighter = &fighter{
		baseComponent: &baseComponent{
			Entity: e,
		},
		MaxHP:   16,
		HP:      16,
		Defense: 1,
		Power:   4,
	}
	return e
}

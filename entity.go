package main

import (
	"image/color"
)

type entity struct {
	GameMap        *gameMap
	X              int
	Y              int
	Char           string
	Color          color.RGBA
	Name           string
	BlocksMovement bool
	RenderOrder    RenderOrder
}

type actor struct {
	*entity
	AI      *hostileEnemy
	Fighter *fighter
}

func newActor(x, y int, char string, color color.RGBA, name string, ai *hostileEnemy, fighter *fighter) *actor {
	a := &actor{
		entity: &entity{
			X:              x,
			Y:              y,
			Char:           char,
			Color:          color,
			Name:           name,
			BlocksMovement: true,
			RenderOrder:    RenderOrder_Actor,
		},
		AI:      ai,
		Fighter: fighter,
	}
	return a
}

func (e actor) IsAlive() bool {
	return e.Fighter.Entity.AI != nil
}

func (e actor) Spawn(gm *gameMap, x, y int) actor {
	clone := e
	clone.X = x
	clone.Y = y
	clone.GameMap = gm
	gm.Entities = append(gm.Entities, &clone)
	return clone
}

func (e *actor) Place(x, y int, gm *gameMap) {
	e.X = x
	e.Y = y
	if gm != nil {
		if e.GameMap != nil {
			e.GameMap.Entities = []*actor{}
		}
		e.GameMap = gm
		gm.Entities = append(gm.Entities, e)
	}
}

func (e *actor) SetPostion(pos [2]int) {
	e.X = pos[0]
	e.Y = pos[1]
}

func (e *actor) Move(dx, dy int) {
	e.X += dx
	e.Y += dy
}

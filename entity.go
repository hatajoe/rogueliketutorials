package main

import (
	"fmt"
	"image/color"
)

type entityHolder interface {
	GetEntities() []entity
	SetEntities(entities []entity)
}

type entity interface {
	Entity() *baseEntity
	GameMap() *gameMap
	RenderOrder() RenderOrder
	Spawn(gm *gameMap, x, y int) entity
	Place(x, y int, gm *gameMap)
}

type baseEntity struct {
	Parent         entityHolder
	X              int
	Y              int
	Char           string
	Color          color.RGBA
	Name           string
	BlocksMovement bool
	RO             RenderOrder
}

func (e *baseEntity) Entity() *baseEntity {
	return e
}

func (e *baseEntity) GameMap() *gameMap {
	switch t := e.Parent.(type) {
	case *gameMap:
		return t
	case *inventory:
		return t.Parent.GameMap()
	default:
		panic(fmt.Sprintf("undefined type: %v", t))
	}
}

func (e *baseEntity) RenderOrder() RenderOrder {
	return e.RO
}

type actor struct {
	*baseEntity
	AI        *hostileEnemy
	Fighter   *fighter
	Inventory *inventory
}

func newActor(x, y int, char string, color color.RGBA, name string, fig *fighter, inv *inventory) *actor {
	a := &actor{
		baseEntity: &baseEntity{
			X:              x,
			Y:              y,
			Char:           char,
			Color:          color,
			Name:           name,
			BlocksMovement: true,
			RO:             RenderOrder_Actor,
		},
		Fighter:   fig,
		Inventory: inv,
	}
	a.Fighter.Parent = a
	a.Inventory.Parent = a
	a.AI = NewHostileEnemy(a)
	return a
}

func (e actor) IsAlive() bool {
	if p, ok := e.Fighter.Parent.(*actor); ok {
		return p.AI != nil
	}
	return false
}

func (e *actor) Spawn(gm *gameMap, x, y int) entity {
	clone := e
	clone.X = x
	clone.Y = y
	clone.Parent = gm
	gm.Entities = append(gm.Entities, clone)
	return clone
}

func (e *actor) Place(x, y int, gm *gameMap) {
	e.X = x
	e.Y = y
	if gm != nil {
		if e.Parent != nil {
			if e.Parent == e.GameMap() {
				e.Parent.SetEntities([]entity{})
			}
		}
		e.Parent = gm
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

type item struct {
	*baseEntity
	Consumable consumable
}

func newItem(x, y int, char string, color color.RGBA, name string, c consumable) *item {
	i := &item{
		baseEntity: &baseEntity{
			X:              x,
			Y:              y,
			Char:           char,
			Color:          color,
			Name:           name,
			BlocksMovement: true,
			RO:             RenderOrder_Actor,
		},
	}
	i.Consumable = c
	i.Consumable.SetParent(i)
	return i
}

func (e *item) Spawn(gm *gameMap, x, y int) entity {
	clone := e
	clone.X = x
	clone.Y = y
	clone.Parent = gm
	gm.Entities = append(gm.Entities, clone)
	return clone
}

func (e *item) Place(x, y int, gm *gameMap) {
	e.X = x
	e.Y = y
	if gm != nil {
		if e.Parent != nil {
			if e.Parent == e.GameMap() {
				entities := e.Parent.GetEntities()
				for i, entity := range entities {
					if e == entity {
						entities = append(entities[:i], entities[i+1:]...)
						break
					}
				}
				e.Parent.SetEntities(entities)
			}
		}
		e.Parent = gm
		gm.Entities = append(gm.Entities, e)
	}
}

package main

import (
	"errors"
)

var (
	regularTermination = errors.New("regular termination")
)

type Action interface {
	Perform(engine *Engine, entity *Entity) error
}

type NoneAction struct{}

func (a NoneAction) Perform(engine *Engine, entity *Entity) error {
	return nil
}

type EscapeAction struct{}

func (a EscapeAction) Perform(engine *Engine, entity *Entity) error {
	return regularTermination
}

type MovementAction struct {
	Dx int
	Dy int
}

func (a MovementAction) Perform(engine *Engine, entity *Entity) error {
	destX := entity.X() + a.Dx
	destY := entity.Y() + a.Dy
	if !engine.GameMap.InBounds(destX, destY) {
		return nil
	}
	if !engine.GameMap.Walkable(destX, destY) {
		return nil
	}
	entity.Move(a.Dx, a.Dy)
	return nil
}

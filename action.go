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
		destX := entity.X() + a.Dx * 10
		destY := entity.Y() + a.Dy * 10
		if !engine.GameMap.InBounds(destX, destY) {
			return nil
		}
		if !engine.GameMap.Walkable(destX/10, destY/10) {
			return nil
		}
		entity.Move(a.Dx * 10, a.Dy * 10)
		return nil
}

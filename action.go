package main

import (
	"errors"
	"fmt"
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

type MeleeAction struct {
	Dx int
	Dy int
}

func (a MeleeAction) Perform(engine *Engine, entity *Entity) error {
	destX := entity.X() + a.Dx
	destY := entity.Y() + a.Dy
	target := engine.GameMap.GetBlockingEntityAtLocation(destX, destY)
	if target == nil {
		return nil
	}
	fmt.Printf("You kick the %s, much to its annoynance!\n", target.Name)
	return nil
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
	if engine.GameMap.GetBlockingEntityAtLocation(destX, destY) != nil {
		return nil
	}

	entity.Move(a.Dx, a.Dy)
	return nil
}

type BumpAction struct {
	Dx int
	Dy int
}

func (a BumpAction) Perform(engine *Engine, entity *Entity) error {
	destX := entity.X() + a.Dx
	destY := entity.Y() + a.Dy

	if engine.GameMap.GetBlockingEntityAtLocation(destX, destY) != nil {
		return MeleeAction{Dx: a.Dx, Dy: a.Dy}.Perform(engine, entity)
	} else {
		return MovementAction{Dx: a.Dx, Dy: a.Dy}.Perform(engine, entity)
	}
}

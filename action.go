package main

import (
	"errors"
	"fmt"
)

var (
	regularTermination = errors.New("regular termination")
)

type action interface {
	Perform() error
}

type noneAction struct{}

func (a noneAction) Perform() error {
	return nil
}

type escapeAction struct{}

func (a escapeAction) Perform() error {
	return regularTermination
}

type waitAction struct{}

func (a waitAction) Perform() error {
	return nil
}

type baseAction struct {
	Entity *actor
}

func (a baseAction) Engine() *engine {
	return a.Entity.GameMap.Engine
}

type actionWithDirection struct {
	baseAction
	Dx int
	Dy int
}

func (a actionWithDirection) DestXY() (int, int) {
	return a.Entity.X + a.Dx, a.Entity.Y + a.Dy
}

func (a actionWithDirection) BlockingEntity() *actor {
	return a.Engine().GameMap.GetBlockingEntityAtLocation(a.DestXY())
}

func (a actionWithDirection) TargetActor() *actor {
	return a.Engine().GameMap.GetAcotrAtLocation(a.DestXY())
}

type meleeAction struct {
	actionWithDirection
}

func (a meleeAction) Perform() error {
	target := a.TargetActor()
	if target == nil {
		return nil
	}

	damage := a.Entity.Fighter.Power - target.Fighter.Defense

	attackDesc := fmt.Sprintf("%s attacks %s", a.Entity.Name, target.Name)
	if damage > 0 {
		fmt.Printf("%s for %d hit points.\n", attackDesc, damage)
		target.Fighter.SetHP(target.Fighter.HP - damage)
	} else {
		fmt.Printf("%s but does no damage.\n", attackDesc)
	}
	return nil
}

type movementAction struct {
	actionWithDirection
}

func (a movementAction) Perform() error {
	destX, destY := a.DestXY()
	if !a.Engine().GameMap.InBounds(destX, destY) {
		return nil
	}
	if !a.Engine().GameMap.Walkable(destX, destY) {
		return nil
	}
	if a.Engine().GameMap.GetBlockingEntityAtLocation(destX, destY) != nil {
		return nil
	}

	a.Entity.Move(a.Dx, a.Dy)
	return nil
}

type bumpAction struct {
	actionWithDirection
}

func (a bumpAction) Perform() error {
	if a.TargetActor() != nil {
		return meleeAction{
			actionWithDirection{
				baseAction: baseAction{
					Entity: a.Entity,
				},
				Dx: a.Dx,
				Dy: a.Dy,
			},
		}.Perform()
	} else {
		return movementAction{
			actionWithDirection{
				baseAction: baseAction{
					Entity: a.Entity,
				},
				Dx: a.Dx,
				Dy: a.Dy,
			},
		}.Perform()
	}
}

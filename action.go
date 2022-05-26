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

type dropItem struct {
	itemAction
}

func (a dropItem) Perform() error {
	a.Entity.Inventory.Drop(a.Item)
	return nil
}

type waitAction struct{}

func (a waitAction) Perform() error {
	return nil
}

type baseAction struct {
	Entity *actor
}

func (a *baseAction) Engine() *engine {
	return a.Entity.Parent.(*gameMap).Engine
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
		return impossible{"Nothing to attack."}
	}

	damage := a.Entity.Fighter.Power - target.Fighter.Defense

	attackDesc := fmt.Sprintf("%s attacks %s", a.Entity.Name, target.Name)
	attackColor := ColorEnemyAtk
	if a.Entity == a.Engine().Player {
		attackColor = ColorPlayerAtk
	}
	if damage > 0 {
		a.Engine().MessageLog.AddMessage(
			fmt.Sprintf("%s for %d hit points.", attackDesc, damage),
			attackColor,
			true,
		)
		target.Fighter.TakeDamage(damage)
	} else {
		a.Engine().MessageLog.AddMessage(
			fmt.Sprintf("%s but does no damage.", attackDesc),
			attackColor,
			true,
		)
	}
	return nil
}

type movementAction struct {
	actionWithDirection
}

func (a movementAction) Perform() error {
	destX, destY := a.DestXY()
	if !a.Engine().GameMap.InBounds(destX, destY) {
		return impossible{"That way is blocked."}
	}
	if !a.Engine().GameMap.Walkable(destX, destY) {
		return impossible{"That way is blocked."}
	}
	if a.Engine().GameMap.GetBlockingEntityAtLocation(destX, destY) != nil {
		return impossible{"That way is blocked."}
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

type pickupAction struct {
	baseAction
}

func newPickupAction(entity *actor) *pickupAction {
	return &pickupAction{
		baseAction: baseAction{
			Entity: entity,
		},
	}
}

func (a *pickupAction) Perform() error {
	actorLocationX := a.Entity.X
	actorLocationY := a.Entity.Y
	inv := a.Entity.Inventory

	for _, it := range a.Engine().GameMap.Items() {
		if actorLocationX == it.X && actorLocationY == it.Y {
			if len(inv.Items) >= inv.Capacity {
				return impossible{"Your inventory is full."}
			}

			entities := a.Engine().GameMap.Entities
			for i, e := range entities {
				if e == it {
					entities = append(entities[:i], entities[i+1:]...)
					break
				}
			}
			a.Engine().GameMap.Entities = entities

			it.Parent = a.Entity.Inventory
			inv.Items = append(inv.Items, it)

			a.Engine().MessageLog.AddMessage(fmt.Sprintf("You picked up the %s!", it.Entity().Name), ColorWhite, true)
			return nil
		}
	}
	return impossible{"There is nothing here to pick up."}
}

type itemAction struct {
	baseAction
	Item     *item
	TargetXY [2]int
}

func newItemAction(entity *actor, i *item, targetXY *[2]int) *itemAction {
	a := &itemAction{
		baseAction: baseAction{
			Entity: entity,
		},
		Item: i,
	}
	if targetXY == nil {
		targetXY = &[2]int{entity.X, entity.Y}
	}
	a.TargetXY = *targetXY

	return a
}

func (a *itemAction) TargetActor() *actor {
	return a.Entity.Parent.(*gameMap).GetAcotrAtLocation(a.TargetXY[0], a.TargetXY[1])
}

func (a *itemAction) Perform() error {
	return a.Item.Consumable.Activate(a)
}

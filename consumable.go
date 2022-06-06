package main

import (
	"errors"
	"fmt"
)

type consumable interface {
	SetParent(i *item)
	GetAction(consumer *actor) action
	Activate(action *itemAction) error
}

type baseConsumable struct {
	baseComponent
}

func (c *baseComponent) SetParent(i *item) {
	c.Parent = i
}

func (c *baseConsumable) GetAction(consumer *actor) action {
	return newItemAction(consumer, c.Parent.(*item), nil)
}

func (c *baseConsumable) Activate(action *itemAction) error {
	return errors.New("Not implemented error")
}

func (c *baseConsumable) Consume() {
	entity := c.Parent
	inv := entity.Entity().Parent
	if di, ok := inv.(*inventory); ok {
		items := di.Items
		for i, item := range items {
			if entity == item {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}
		di.Items = items
	}
}

type healingConsumable struct {
	baseConsumable
	Amount int
}

func newHealingConsumable(amount int) *healingConsumable {
	return &healingConsumable{
		baseConsumable: baseConsumable{
			baseComponent: baseComponent{
				Parent: nil,
			},
		},
		Amount: amount,
	}
}

func (c *healingConsumable) Activate(act *itemAction) error {
	consumer := act.Entity
	amountRecoverd := consumer.Fighter.Heal(c.Amount)

	if amountRecoverd > 0 {
		c.Engine().MessageLog.AddMessage(fmt.Sprintf("You consume the %s, and recover %d HP!", c.Parent.Entity().Name, amountRecoverd), ColorHealthRecovered, true)
		c.Consume()
	} else {
		return impossible{"Your health is already full."}
	}
	return nil
}

type lightningDamageConsumable struct {
	baseConsumable
	Damage       int
	MaximumRange int
}

func newLightningDamageConsumable(damage, maximumRange int) *lightningDamageConsumable {
	return &lightningDamageConsumable{
		baseConsumable: baseConsumable{
			baseComponent: baseComponent{
				Parent: nil,
			},
		},
		Damage:       damage,
		MaximumRange: maximumRange,
	}
}

func (c *lightningDamageConsumable) Activate(act *itemAction) error {
	consumer := act.Entity
	var target *actor

	closestDistance := float64(c.MaximumRange + 1.0)
	for _, a := range c.Engine().GameMap.Actors() {
		if a != consumer && c.Parent.GameMap().Visible[a.X][a.Y] {
			distance := consumer.Distance(a.X, a.Y)
			if distance < closestDistance {
				target = a
				closestDistance = distance
			}
		}
	}
	if target != nil {
		c.Engine().MessageLog.AddMessage(fmt.Sprintf("A lightning bolt strikes the %s with a loud thunder, for %d damage!", target.Name, c.Damage), ColorWhite, true)
		target.Fighter.TakeDamage(c.Damage)
		c.Consume()
	} else {
		return impossible{"No enemy is close enough to strike."}
	}
	return nil
}

type confusionConsumable struct {
	baseConsumable
	NumberOfTurns int
}

func newConfusionConsumable(numberOfTurns int) *confusionConsumable {
	return &confusionConsumable{
		baseConsumable: baseConsumable{
			baseComponent: baseComponent{
				Parent: nil,
			},
		},
		NumberOfTurns: numberOfTurns,
	}
}

func (c *confusionConsumable) GetAction(consumer *actor) action {
	c.Engine().MessageLog.AddMessage("Select a target location.", ColorNeedsTarget, true)
	c.Engine().EventHandler = &singleRangedAttackHandler{
		selectIndexHandler: newSelectIndexHandler(c.Engine()),
		Callback: func(x, y int) action {
			return newItemAction(consumer, c.Parent.(*item), &[2]int{x, y})
		},
	}
	return noneAction{}
}

func (c *confusionConsumable) Activate(action *itemAction) error {
	consumer := action.Entity
	target := action.TargetActor()

	if !c.Engine().GameMap.Visible[action.TargetXY[0]][action.TargetXY[1]] {
		return impossible{"You cannot target an area that you cannot see."}
	} else if target == nil {
		return impossible{"You must select an enemy to target."}
	} else if target == consumer {
		return impossible{"You cannot confuse yourself!"}
	}

	c.Engine().MessageLog.AddMessage(fmt.Sprintf("The eyes of the %s look vacant, as it starts to stumble around!", target.Name), ColorStatusEffectApplied, true)
	target.AI = newConfusedEnemy(target, target.AI, c.NumberOfTurns)

	c.Consume()
	return nil
}

type fireballDamageConsumable struct {
	baseConsumable
	Damage int
	Radius int
}

func newFireballDamageConsumable(damage, radius int) *fireballDamageConsumable {
	return &fireballDamageConsumable{
		baseConsumable: baseConsumable{
			baseComponent: baseComponent{
				Parent: nil,
			},
		},
		Damage: damage,
		Radius: radius,
	}
}

func (c *fireballDamageConsumable) GetAction(consumer *actor) action {
	c.Engine().MessageLog.AddMessage("Select a target location.", ColorNeedsTarget, true)
	c.Engine().EventHandler = &areaRangedAttackHandler{
		selectIndexHandler: newSelectIndexHandler(c.Engine()),
		Radius:             c.Radius,
		Callback: func(x, y int) action {
			return newItemAction(consumer, c.Parent.(*item), &[2]int{x, y})
		},
	}
	return noneAction{}
}

func (c *fireballDamageConsumable) Activate(action *itemAction) error {
	targetXY := action.TargetXY

	if !c.Engine().GameMap.Visible[targetXY[0]][targetXY[1]] {
		return impossible{"You cannot target an area that you cannot see."}
	}

	targetHit := false
	for _, a := range c.Engine().GameMap.Actors() {
		if a.Distance(targetXY[0], targetXY[1]) <= float64(c.Radius+1) {
			c.Engine().MessageLog.AddMessage(fmt.Sprintf("The %s is engulfed in a fiery explosion, taking %d damage!", a.Name, c.Damage), ColorWhite, true)
			a.Fighter.TakeDamage(c.Damage)
			targetHit = true
		}
	}
	if !targetHit {
		return impossible{"There are no targets in the radius."}
	}

	c.Consume()
	return nil
}

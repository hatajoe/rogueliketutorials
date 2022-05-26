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

func (c *healingConsumable) Activate(action *itemAction) error {
	consumer := action.Entity
	amountRecoverd := consumer.Fighter.Heal(c.Amount)

	if amountRecoverd > 0 {
		c.Engine().MessageLog.AddMessage(fmt.Sprintf("You consume the %s, and recover %d HP!", c.Parent.Entity().Name, amountRecoverd), ColorHealthRecovered, true)
		c.Consume()
	} else {
		return impossible{"Your health is already full."}
	}
	return nil
}

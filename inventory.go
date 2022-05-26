package main

import "fmt"

type inventory struct {
	*baseComponent
	Capacity int
	Items    []*item
}

func newInventory(capacity int) *inventory {
	return &inventory{
		baseComponent: &baseComponent{
			Parent: nil,
		},
		Capacity: capacity,
		Items:    []*item{},
	}
}

func (c *inventory) GetEntities() []entity {
	entities := make([]entity, 0, len(c.Items))
	for _, i := range c.Items {
		entities = append(entities, i)
	}
	return entities
}

func (c *inventory) SetEntities(entities []entity) {
	items := make([]*item, 0, len(entities))
	for _, e := range entities {
		if i, ok := e.(*item); ok {
			items = append(items, i)
		}
	}
	c.Items = items
}

func (c *inventory) Drop(drop *item) {
	for i, it := range c.Items {
		if it == drop {
			c.Items = append(c.Items[:i], c.Items[i+1:]...)
			break
		}
	}
	drop.Place(c.Parent.Entity().X, c.Parent.Entity().Y, c.GameMap())

	c.Engine().MessageLog.AddMessage(fmt.Sprintf("You dropped the %s.", drop.Entity().Name), ColorWhite, true)
}

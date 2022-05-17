package main

type baseComponent struct {
	Entity *actor
}

func (c baseComponent) Engine() *engine {
	return c.Entity.GameMap.Engine
}

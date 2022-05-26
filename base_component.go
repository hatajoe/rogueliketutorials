package main

type baseComponent struct {
	Parent entity
}

func (c *baseComponent) GameMap() *gameMap {
	return c.Parent.GameMap()
}

func (c baseComponent) Engine() *engine {
	return c.GameMap().Engine
}

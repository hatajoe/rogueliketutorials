package main

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type gameMap struct {
	Engine   *engine
	Width    int
	Height   int
	Tiles    [][]*tile
	Visible  [][]bool
	Explored [][]bool
	Entities []*actor
}

func newGameMap(en *engine, width, height int, entities []*actor) *gameMap {
	tiles := make([][]*tile, width)
	visible := make([][]bool, width)
	explored := make([][]bool, width)

	// fill map by floor
	for w := 0; w < width; w += 1 {
		tiles[w] = make([]*tile, height)
		visible[w] = make([]bool, height)
		explored[w] = make([]bool, height)
		for h := 0; h < height; h += 1 {
			tiles[w][h] = newWall()
			visible[w][h] = false
			explored[w][h] = false
		}
	}

	return &gameMap{
		Engine:   en,
		Width:    width,
		Height:   height,
		Tiles:    tiles,
		Visible:  visible,
		Explored: explored,
		Entities: entities,
	}
}

func (g gameMap) Actors() []*actor {
	es := []*actor{}
	for _, a := range g.Entities {
		if a.IsAlive() {
			es = append(es, a)
		}
	}
	return es
}

func (g gameMap) GetAcotrAtLocation(x, y int) *actor {
	for _, a := range g.Actors() {
		if a.X == x && a.Y == y {
			return a
		}
	}
	return nil
}

func (g gameMap) GetBlockingEntityAtLocation(x, y int) *actor {
	for _, e := range g.Entities {
		if e.BlocksMovement && e.X == x && e.Y == y {
			return e
		}
	}
	return nil
}

func (g gameMap) Render(screen *ebiten.Image, font font.Face) {
	for w, ts := range g.Tiles {
		for h, t := range ts {
			color := t.Shroud
			if g.IsVisible(w, h) {
				color = t.Light
			} else if g.IsExplored(w, h) {
				color = t.Dark
			}
			text.Draw(screen, t.Char, font, w*10, h*10, color)
		}
	}
	entities := g.Entities
	sort.Slice(entities, func(i, j int) bool { return entities[i].RenderOrder < entities[i].RenderOrder })
	for _, entity := range entities {
		if g.IsVisible(entity.X, entity.Y) {
			text.Draw(screen, entity.Char, font, entity.X*10, entity.Y*10, entity.Color)
		}
	}
}

func (g gameMap) InBounds(x, y int) bool {
	return 0 <= x && x < g.Width && 0 <= y && y < g.Height
}

func (g gameMap) Walkable(x, y int) bool {
	return g.Tiles[x][y].Walkable
}

func (g gameMap) IsVisible(x, y int) bool {
	return g.Visible[x][y]
}

func (g gameMap) IsExplored(x, y int) bool {
	return g.Explored[x][y]
}

package main

import (
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type GameMap struct {
	Width    int
	Height   int
	Tiles    [][]*Tile
	Visible  [][]bool
	Explored [][]bool
	Entities []*Entity
}

func NewGameMap(width, height int, entities []*Entity) *GameMap {
	tiles := make([][]*Tile, width)
	visible := make([][]bool, width)
	explored := make([][]bool, width)

	// fill map by floor
	for w := 0; w < width; w += 1 {
		tiles[w] = make([]*Tile, height)
		visible[w] = make([]bool, height)
		explored[w] = make([]bool, height)
		for h := 0; h < height; h += 1 {
			tiles[w][h] = NewWall()
			visible[w][h] = false
			explored[w][h] = false
		}
	}

	return &GameMap{
		Width:    width,
		Height:   height,
		Tiles:    tiles,
		Visible:  visible,
		Explored: explored,
		Entities: entities,
	}
}

func (g GameMap) GetBlockingEntityAtLocation(x, y int) *Entity {
	for _, e := range g.Entities {
		if e.BlocksMovement && e.X() == x && e.Y() == y {
			return e
		}
	}
	return nil
}

func (g GameMap) Render(screen *ebiten.Image, font font.Face) {
	for w, ts := range g.Tiles {
		for h, t := range ts {
			color := t.Shroud()
			if g.IsVisible(w, h) {
				color = t.Light()
			} else if g.IsExplored(w, h) {
				color = t.Dark()
			}
			text.Draw(screen, t.Char(), font, w*10, h*10, color)
		}
	}
	for _, entity := range g.Entities {
		if g.IsVisible(entity.X(), entity.Y()) {
			text.Draw(screen, entity.Char(), font, entity.X()*10, entity.Y()*10, entity.Color())
		}
	}
}

func (g GameMap) InBounds(x, y int) bool {
	return 0 <= x && x < g.Width && 0 <= y && y < g.Height
}

func (g GameMap) Walkable(x, y int) bool {
	return g.Tiles[x][y].Walkable()
}

func (g GameMap) IsVisible(x, y int) bool {
	return g.Visible[x][y]
}

func (g GameMap) IsExplored(x, y int) bool {
	return g.Explored[x][y]
}

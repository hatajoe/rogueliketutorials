package main

type GameMap struct {
	Width  int
	Height int
	Tiles  [][]*Tile
}

func NewGameMap(width, height int) *GameMap {
	tiles := make([][]*Tile, width)

	// fill map by floor
	for w := 0; w < width; w += 1 {
		tiles[w] = make([]*Tile, height)
		for h := 0; h < height; h += 1 {
			tiles[w][h] = NewWall()
		}
	}

	return &GameMap{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}
}

func (g GameMap) InBounds(x, y int) bool {
	return 0 <= x && x < g.Width && 0 <= y && y < g.Height
}

func (g GameMap) Walkable(x, y int) bool {
	return g.Tiles[x][y].Walkable()
}

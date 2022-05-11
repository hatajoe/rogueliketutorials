package main

type GameMap struct {
	width  int
	height int
	tiles  [][]*Tile
}

func NewGameMap(width, height int) *GameMap {
	tiles := make([][]*Tile, width/10)

	// fill map by floor
	for w := 0; w < width/10; w += 1 {
		tiles[w] = make([]*Tile, height/10)
		for h := 0; h < height/10; h += 1 {
			tiles[w][h] = NewFloor()
		}
	}

	// make some floors to wall
	for _, w := range []int{30, 31, 32, 33} {
		for _, h := range []int{22} {
			tiles[w][h] = NewWall()
		}
	}

	return &GameMap{
		width:  width,
		height: height,
		tiles:  tiles,
	}
}

func (g GameMap) Tiles() [][]*Tile {
	return g.tiles
}

func (g GameMap) InBounds(x, y int) bool {
	return 0 <= x && x < g.width && 0 <= y && y < g.height
}

func (g GameMap) Walkable(x, y int) bool {
	return g.tiles[x][y].Walkable()
}

package main

type GameMap struct {
	Width    int
	Height   int
	Tiles    [][]*Tile
	Visible  [][]bool
	Explored [][]bool
}

func NewGameMap(width, height int) *GameMap {
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

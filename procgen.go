package main

import (
	"math"
	"math/rand"
	"time"
)

type rectangularRoom struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

func newRectangularRoom(x, y, width, height int) rectangularRoom {
	return rectangularRoom{
		X1: x,
		Y1: y,
		X2: x + width,
		Y2: y + height,
	}
}

func (r rectangularRoom) Center() (int, int) {
	centerX := int((r.X1 + r.X2) / 2)
	centerY := int((r.Y1 + r.Y2) / 2)
	return centerX, centerY
}

func (r rectangularRoom) Inner() ([2]int, [2]int) {
	return [2]int{r.X1 + 1, r.X2}, [2]int{r.Y1 + 1, r.Y2}
}

func (r rectangularRoom) Intersects(other rectangularRoom) bool {
	return r.X1 <= other.X2 &&
		r.X2 >= other.X1 &&
		r.Y1 <= other.Y2 &&
		r.Y2 >= other.Y1
}

func generateDungeon(maxRooms, roomMinSize, roomMaxSize, mapWidth, mapHeight, maxMonsterPerRoom, maxItemsPerRoom int, en *engine) *gameMap {
	rand.Seed(time.Now().UnixNano())

	player := en.Player
	dungeon := newGameMap(en, mapWidth, mapHeight, []entity{player})

	rooms := []rectangularRoom{}
	for i := 0; i < maxRooms; i++ {
		roomWidth := rand.Intn(roomMaxSize-roomMinSize) + roomMinSize
		roomHeight := rand.Intn(roomMaxSize-roomMinSize) + roomMinSize

		x := rand.Intn(dungeon.Width - roomWidth - 1)
		y := rand.Intn(dungeon.Height - roomHeight - 1)

		newRoom := newRectangularRoom(x, y, roomWidth, roomHeight)
		intersects := false
		for _, otherRoom := range rooms {
			if newRoom.Intersects(otherRoom) {
				intersects = true
				break
			}
		}
		if intersects {
			continue
		}

		wtup, htup := newRoom.Inner()
		for w := wtup[0]; w < wtup[1]; w++ {
			for h := htup[0]; h < htup[1]; h++ {
				dungeon.Tiles[w][h] = newFloor()
			}
		}

		if len(rooms) == 0 {
			px, py := newRoom.Center()
			player.Place(px, py, dungeon)
		} else {
			x1, y1 := rooms[len(rooms)-1].Center()
			x2, y2 := newRoom.Center()
			ch := tunnelBetween(x1, y1, x2, y2)
			for tup := range ch {
				w, h := tup[0], tup[1]
				dungeon.Tiles[w][h] = newFloor()
			}
		}

		placeEntities(newRoom, dungeon, maxMonsterPerRoom, maxItemsPerRoom)

		rooms = append(rooms, newRoom)
	}

	return dungeon
}

func tunnelBetween(x1, y1, x2, y2 int) chan [2]int {
	cornerX, cornerY := x1, y2
	if rand.Float32() < 0.5 {
		cornerX, cornerY = x2, y1
	}

	ch := make(chan [2]int)
	go func() {
		defer close(ch)
		ctup := bresenham(x1, y1, cornerX, cornerY)
		for tup := range ctup {
			ch <- tup
		}
		ctup = bresenham(cornerX, cornerY, x2, y2)
		for tup := range ctup {
			ch <- tup
		}
	}()

	return ch
}

// https://41j.com/blog/2012/09/bresenhams-line-drawing-algorithm-implemetations-in-go-and-c/
func bresenham(x1, y1, x2, y2 int) chan [2]int {
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	e := dx - dy

	ch := make(chan [2]int)
	go func() {
		defer close(ch)
		for {
			ch <- [2]int{x1, y1}
			if x1 == x2 && y1 == y2 {
				break
			}

			e2 := 2 * e
			if e2 > -dy {
				e -= dy
				x1 += sx
			}
			if e2 < dx {
				e += dx
				y1 += sy
			}
		}
	}()
	return ch
}

func placeEntities(room rectangularRoom, dungeon *gameMap, maximumMonsters, maximumItems int) {
	numberOfMonsters := rand.Intn(maximumMonsters + 1)
	numberOfItems := rand.Intn(maximumItems)

	for i := 0; i < numberOfMonsters; i++ {
		x := rand.Intn((room.X2-1)-(room.X1+1)) + room.X1 + 1
		y := rand.Intn((room.Y2-1)-(room.Y1+1)) + room.Y1 + 1

		for _, entity := range dungeon.Entities {
			e := entity.Entity()
			if !(e.X == x && e.Y == y) {
				if rand.Float32() < 0.8 {
					newOrc().Spawn(dungeon, x, y)
				} else {
					newTroll().Spawn(dungeon, x, y)
				}
				break
			}
		}
	}

	for i := 0; i < numberOfItems; i++ {
		x := rand.Intn((room.X2-1)-(room.X1+1)) + room.X1 + 1
		y := rand.Intn((room.Y2-1)-(room.Y1+1)) + room.Y1 + 1

		for _, entity := range dungeon.Entities {
			e := entity.Entity()
			if !(e.X == x && e.Y == y) {
				newHealthPortion().Spawn(dungeon, x, y)
				break
			}
		}
	}
}

package main

import (
	"math"
	"math/rand"
	"time"
)

type rectangularRoom struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func newRectangularRoom(x, y, width, height int) rectangularRoom {
	return rectangularRoom{
		x1: x,
		y1: y,
		x2: x + width,
		y2: y + height,
	}
}

func (r rectangularRoom) Center() [2]int {
	centerX := int((r.x1 + r.x2) / 2)
	centerY := int((r.y1 + r.y2) / 2)
	return [2]int{centerX, centerY}
}

func (r rectangularRoom) Inner() ([2]int, [2]int) {
	return [2]int{r.x1 + 1, r.x2}, [2]int{r.y1 + 1, r.y2}
}

func (r rectangularRoom) Intersects(other rectangularRoom) bool {
	return r.x1 <= other.x2 &&
		r.x2 >= other.x1 &&
		r.y1 <= other.y2 &&
		r.y2 >= other.y1
}

func GenerateDungeon(maxRooms, roomMinSize, roomMaxSize, mapWidth, mapHeight int, player *Entity) *GameMap {
	rand.Seed(time.Now().UnixNano())

	dungeon := NewGameMap(mapWidth, mapHeight)

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
				dungeon.Tiles[w][h] = NewFloor()
			}
		}

		if len(rooms) == 0 {
			player.SetPostion(newRoom.Center())
		} else {
			ch := tunnelBetween(rooms[len(rooms)-1].Center(), newRoom.Center())
			for tup := range ch {
				w, h := tup[0], tup[1]
				dungeon.Tiles[w][h] = NewFloor()
			}
		}
		rooms = append(rooms, newRoom)
	}

	return dungeon
}

func tunnelBetween(start [2]int, end [2]int) chan [2]int {
	x1, y1 := start[0], start[1]
	x2, y2 := end[0], end[1]

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

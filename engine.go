package main

import (
	"math"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

type engine struct {
	GameMap       *gameMap
	MessageLog    *MessageLog
	MouseLocation [2]int
	Player        *actor
	Font          font.Face
}

func NewEngine(pl *actor, font font.Face) *engine {
	e := &engine{
		Player: pl,
		Font:   font,
	}
	e.MessageLog = NewMessageLog()
	e.MouseLocation = [2]int{0, 0}
	return e
}

func (e *engine) HandleEnemyTurns() error {
	for _, entity := range e.GameMap.Entities {
		if entity == e.Player {
			continue
		}
		if a, ok := entity.(*actor); ok {
			if a.IsAlive() {
				if err := a.AI.Perform(); err != nil {
					switch err.(type) {
					case impossible:
						// pass
					default:
						return err
					}
				}
			}
		}
	}
	return nil
}

func (e engine) Render(screen *ebiten.Image) {
	e.GameMap.Render(screen, e.Font)

	e.MessageLog.Render(screen, e.Font, 21, 45, 40, 5)

	RenderBar(screen, e.Font, e.Player.Fighter.HP, e.Player.Fighter.MaxHP, 200)

	RenderNamesAtMouseLocation(screen, e.Font, 21, 44, &e)
}

func (e *engine) UpdateFov() {
	for w, wv := range e.GameMap.Visible {
		for h := range wv {
			e.GameMap.Visible[w][h] = false
		}
	}
	e.computeFov(e.Player.X, e.Player.Y, 8)
	for w, v := range e.GameMap.Visible {
		for h := range v {
			if e.GameMap.Visible[w][h] {
				e.GameMap.Explored[w][h] = true
			}
		}
	}
}

// https://github.com/norendren/go-fov/blob/master/fov/fov.go
func (e *engine) computeFov(x, y, radius int) {
	e.GameMap.Visible[x][y] = true

	for i := 0; i < 8; i++ {
		e.fov(x, y, 1, 0, 1, i, radius)
	}
}

func (e *engine) fov(x, y, dist int, lowSlope, highSlope float64, oct, rad int) {
	if dist > rad {
		return
	}

	low := math.Floor(lowSlope*float64(dist) + 0.5)
	high := math.Floor(highSlope*float64(dist) + 0.5)

	inGap := false
	for height := low; height <= high; height++ {
		mapx, mapy := distHeightXY(x, y, dist, int(height), oct)
		if e.GameMap.InBounds(mapx, mapy) && distTo(x, y, mapx, mapy) < rad {
			e.GameMap.Visible[mapx][mapy] = true
		}
		if e.GameMap.InBounds(mapx, mapy) && !e.GameMap.Walkable(mapx, mapy) {
			if inGap {
				e.fov(x, y, dist+1, lowSlope, (height-0.5)/float64(dist), oct, rad)
			}
			lowSlope = (height + 0.5) / float64(dist)
			inGap = false
		} else {
			inGap = true
			if height == high {
				e.fov(x, y, dist+1, lowSlope, highSlope, oct, rad)
			}
		}
	}
}

func distHeightXY(x, y, d, h, oct int) (int, int) {
	if oct&0x1 > 0 {
		d = -d
	}
	if oct&0x2 > 0 {
		h = -h
	}
	if oct&0x4 > 0 {
		return x + h, y + d
	}
	return x + d, y + h
}

func distTo(x1, y1, x2, y2 int) int {
	vx := math.Pow(float64(x1-x2), 2)
	vy := math.Pow(float64(y1-y2), 2)
	return int(math.Sqrt(vx + vy))
}

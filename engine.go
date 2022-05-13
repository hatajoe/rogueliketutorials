package main

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type Engine struct {
	entities     []*Entity
	eventHandler *EventHandler
	GameMap      *GameMap
	player       *Entity
	font         font.Face
}

func NewEngine(es []*Entity, eh *EventHandler, gm *GameMap, pl *Entity, font font.Face) *Engine {
	engine := &Engine{
		entities:     es,
		eventHandler: eh,
		GameMap:      gm,
		player:       pl,
		font:         font,
	}
	engine.UpdateFov()

	return engine
}

func (e *Engine) HandleEvent(keys []ebiten.Key) error {
	action := e.eventHandler.KeyDown(keys)
	switch act := action.(type) {
	case NoneAction:
		return nil
	default:
		if err := act.Perform(e, e.player); err != nil {
			return err
		}
		e.UpdateFov()
	}
	return nil
}

func (e *Engine) Render(screen *ebiten.Image) {
	for w, ts := range e.GameMap.Tiles {
		for h, t := range ts {
			color := t.Shroud()
			if e.GameMap.IsVisible(w, h) {
				color = t.Light()
			} else if e.GameMap.IsExplored(w, h) {
				color = t.Dark()
			}
			text.Draw(screen, t.Char(), e.font, w*10, h*10, color)
		}
	}
	for _, entity := range e.entities {
		text.Draw(screen, entity.Char(), e.font, entity.X()*10, entity.Y()*10, entity.Color())
	}
}

func (e *Engine) UpdateFov() {
	for w, wv := range e.GameMap.Visible {
		for h := range wv {
			e.GameMap.Visible[w][h] = false
		}
	}
	e.computeFov(e.player.X(), e.player.Y(), 8)
	for w, v := range e.GameMap.Visible {
		for h := range v {
			if e.GameMap.Visible[w][h] {
				e.GameMap.Explored[w][h] = true
			}
		}
	}
}

// https://github.com/norendren/go-fov/blob/master/fov/fov.go
func (e *Engine) computeFov(x, y, radius int) {
	e.GameMap.Visible[x][y] = true

	for i := 0; i < 8; i++ {
		e.fov(x, y, 1, 0, 1, i, radius)
	}
}

func (e *Engine) fov(x, y, dist int, lowSlope, highSlope float64, oct, rad int) {
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
			lowSlope = height / float64(dist)
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

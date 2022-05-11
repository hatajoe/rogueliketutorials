package main

import (
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

func (e *Engine) HandleEvent(keys []ebiten.Key) error {
	action := e.eventHandler.KeyDown(keys)
	switch act := action.(type) {
	case NoneAction:
		return nil
	default:
		return act.Perform(e, e.player)
	}
}

func (e *Engine) Render(screen *ebiten.Image) {
	for w, ts := range e.GameMap.Tiles() {
		for h, t := range ts {
			text.Draw(screen, t.Char(), e.font, w*10, h*10, t.Color())
		}
	}
	for _, entity := range e.entities {
		text.Draw(screen, entity.Char(), e.font, entity.X(), entity.Y(), entity.Color())
	}
}

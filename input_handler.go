package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EventHandler struct{}

func (e *EventHandler) KeyDown(keys []ebiten.Key) Action {
	ma := MovementAction{}
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		switch p {
		case ebiten.KeyArrowUp:
			ma.Dy -= 1
		case ebiten.KeyArrowDown:
			ma.Dy += 1
		case ebiten.KeyArrowLeft:
			ma.Dx -= 1
		case ebiten.KeyArrowRight:
			ma.Dx += 1
		case ebiten.KeyEscape:
			return EscapeAction{}
		default:
			return NoneAction{}
		}
	}
	return ma
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)

	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

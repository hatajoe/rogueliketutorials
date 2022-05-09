package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type EventHandler struct{}

func (e *EventHandler) KeyDown(keys []ebiten.Key) Action {
	ma := MovementAction{}
	for _, p := range keys {
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

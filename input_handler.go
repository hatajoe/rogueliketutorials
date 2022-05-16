package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type EventHandler struct{}

func (e *EventHandler) KeyDown(keys []ebiten.Key) Action {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		switch p {
		case ebiten.KeyArrowUp:
			return BumpAction{Dx: 0, Dy: -1}
		case ebiten.KeyArrowDown:
			return BumpAction{Dx: 0, Dy: 1}
		case ebiten.KeyArrowLeft:
			return BumpAction{Dx: -1, Dy: 0}
		case ebiten.KeyArrowRight:
			return BumpAction{Dx: 1, Dy: 0}
		case ebiten.KeyEscape:
			return EscapeAction{}
		default:
			return NoneAction{}
		}
	}
	return NoneAction{}
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

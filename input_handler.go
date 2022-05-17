package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	moveKeys = map[ebiten.Key][2]int{
		// Arrow keys.
		ebiten.KeyArrowUp:    {0, -1},
		ebiten.KeyArrowDown:  {0, 1},
		ebiten.KeyArrowLeft:  {-1, 0},
		ebiten.KeyArrowRight: {1, 0},
		ebiten.KeyHome:       {-1, -1},
		ebiten.KeyEnd:        {-1, 1},
		ebiten.KeyPageUp:     {1, -1},
		ebiten.KeyPageDown:   {1, 1},
		// Numpad keys.
		ebiten.KeyNumpad1: {-1, 1},
		ebiten.KeyNumpad2: {0, 1},
		ebiten.KeyNumpad3: {1, 1},
		ebiten.KeyNumpad4: {-1, 0},
		ebiten.KeyNumpad6: {1, 0},
		ebiten.KeyNumpad7: {-1, -1},
		ebiten.KeyNumpad8: {0, -1},
		ebiten.KeyNumpad9: {1, -1},
		// Vi keys
		ebiten.KeyH: {-1, 0},
		ebiten.KeyJ: {0, 1},
		ebiten.KeyK: {0, -1},
		ebiten.KeyL: {1, 0},
		ebiten.KeyY: {-1, -1},
		ebiten.KeyU: {1, -1},
		ebiten.KeyB: {-1, 1},
		ebiten.KeyN: {1, 1},
	}
	waitKeys = map[ebiten.Key]struct{}{
		ebiten.KeyPeriod:  {},
		ebiten.KeyNumpad5: {},
	}
)

type eventHandler interface {
	HandleEvent(keys []ebiten.Key) error
}

type mainGameEventHandler struct {
	engine *engine
}

func (e *mainGameEventHandler) HandleEvent(keys []ebiten.Key) error {
	action := e.EvkeyDown(keys)
	switch act := action.(type) {
	case noneAction:
		return nil
	default:
		if err := act.Perform(); err != nil {
			return err
		}
		if err := e.engine.HandleEnemyTurns(); err != nil {
			return err
		}
		e.engine.UpdateFov()
	}
	return nil
}

func (e *mainGameEventHandler) EvkeyDown(keys []ebiten.Key) action {
	player := e.engine.Player
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		if d, ok := moveKeys[p]; ok {
			return bumpAction{
				actionWithDirection{
					baseAction: baseAction{
						Entity: player,
					},
					Dx: d[0],
					Dy: d[1],
				},
			}
		} else if _, ok := waitKeys[p]; ok {
			return waitAction{}
		}
	}
	return noneAction{}
}

type gameOverEventHandler struct {
	engine *engine
}

func (e gameOverEventHandler) HandleEvent(keys []ebiten.Key) error {
	action := e.EvkeyDown(keys)
	switch act := action.(type) {
	case noneAction:
		return nil
	default:
		if err := act.Perform(); err != nil {
			return err
		}
	}
	return nil
}

func (e *gameOverEventHandler) EvkeyDown(keys []ebiten.Key) action {
	for _, p := range keys {
		switch p {
		case ebiten.KeyEscape:
			return escapeAction{}
		default:
			return noneAction{}
		}
	}
	return noneAction{}
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

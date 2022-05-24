package main

import (
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mattn/go-runewidth"
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
	OnRender(screen *ebiten.Image)
}

type eventHandlerBase struct {
	engine *engine
}

func (e *eventHandlerBase) HandleEvent(keys []ebiten.Key) error {
	return nil
}

func (e *eventHandlerBase) OnRender(screen *ebiten.Image) {
	e.engine.Render(screen)
}

type mainGameEventHandler struct {
	eventHandlerBase
}

func (e *mainGameEventHandler) HandleEvent(keys []ebiten.Key) error {
	mx, my := ebiten.CursorPosition()
	if e.engine.GameMap.InBounds(int(mx/10), int(my/10)) {
		e.engine.MouseLocation = [2]int{int(mx / 10), int(my/10) + 1} // I don't know why +1 is needed, however this worked well
	}
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
		} else {
			switch p {
			case ebiten.KeyEscape:
				return escapeAction{}
			case ebiten.KeyV:
				e.engine.EventHandler = newHistoryViewer(e.engine)
			default:
			}
		}
	}
	return noneAction{}
}

type gameOverEventHandler struct {
	eventHandlerBase
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

var (
	cursorYKeys = map[ebiten.Key]int{
		ebiten.KeyArrowUp:   -1,
		ebiten.KeyArrowDown: 1,
		ebiten.KeyPageUp:    -10,
		ebiten.KeyPageDown:  10,
	}
)

type historyViewer struct {
	eventHandlerBase
	LogLength int
	Cursor    int
	Window    *ebiten.Image
}

func newHistoryViewer(e *engine) *historyViewer {
	logLength := len(e.MessageLog.Messages)
	return &historyViewer{
		eventHandlerBase: eventHandlerBase{
			engine: e,
		},
		LogLength: logLength,
		Cursor:    logLength - 1,
		Window:    nil,
	}
}

func (e *historyViewer) HandleEvent(keys []ebiten.Key) error {
	_ = e.EvkeyDown(keys)
	return nil
}

func (e *historyViewer) EvkeyDown(keys []ebiten.Key) action {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		if adjust, ok := cursorYKeys[p]; ok {
			if adjust < 0 && e.Cursor == 0 {
				e.Cursor = e.LogLength - 1
			} else if adjust > 0 && e.Cursor == e.LogLength-1 {
				e.Cursor = 0
			} else {
				e.Cursor = int(math.Max(0, math.Min(float64(e.Cursor+adjust), float64(e.LogLength-1))))
			}
		} else {
			switch p {
			case ebiten.KeyHome:
				e.Cursor = 0
			case ebiten.KeyEnd:
				e.Cursor = e.LogLength - 1
			default:
				e.engine.EventHandler = &mainGameEventHandler{eventHandlerBase{engine: e.engine}}
			}
		}
	}
	return noneAction{}
}

func (e *historyViewer) OnRender(screen *ebiten.Image) {
	e.eventHandlerBase.OnRender(screen)

	width, height := screen.Size()
	if e.Window != nil {
		ww, wh := e.Window.Size()
		if ww != width-60 || wh != height-60 {
			e.Window = ebiten.NewImage(width-60, height-60)
		}
	} else {
		e.Window = ebiten.NewImage(width-60, height-60)
	}
	e.Window.Fill(ColorBlack)

	windowTitle := "Message history"
	windowTitleLen := (len(windowTitle) + 2) * 10
	ww, wh := e.Window.Size()
	windowTitleStart := int(math.Floor(float64(ww-windowTitleLen)/20) * 10)
	for w := 0; w < ww; w += 10 {
		for h := 0; h < wh; h += 10 {
			if w == 0 && h == 0 {
				text.Draw(e.Window, string([]rune{0xda}), e.engine.Font, w, h+10, ColorWhite) // left upper
			} else if w == ww-10 && h == 0 {
				text.Draw(e.Window, string([]rune{0xbf}), e.engine.Font, w, h+10, ColorWhite) // right upper
			} else if w == 0 && h == wh-10 {
				text.Draw(e.Window, string([]rune{0xc0}), e.engine.Font, w, h+10, ColorWhite) // left lower
			} else if w == ww-10 && h == wh-10 {
				text.Draw(e.Window, string([]rune{0xd9}), e.engine.Font, w, h+10, ColorWhite) // right lower
			} else if w == 0 || w == ww-10 {
				text.Draw(e.Window, string([]rune{0xb3}), e.engine.Font, w, h+10, ColorWhite) // vertical line
			} else if h == 0 || h == wh-10 {
				if h == 0 {
					if w == windowTitleStart {
						text.Draw(e.Window, string([]rune{0xb4}), e.engine.Font, w, h+10, ColorWhite) // title left
					} else if w == windowTitleStart+windowTitleLen-10 {
						text.Draw(e.Window, string([]rune{0xc3}), e.engine.Font, w, h+10, ColorWhite) // title right
					} else if w > windowTitleStart && w < windowTitleStart+windowTitleLen-10 {
						text.Draw(e.Window, string(windowTitle[(w-windowTitleStart-10)/10]), e.engine.Font, w, h+10, ColorWhite)
					} else {
						text.Draw(e.Window, string([]rune{0xc4}), e.engine.Font, w, h+10, ColorWhite) // horizontal line
					}
				} else {
					text.Draw(e.Window, string([]rune{0xc4}), e.engine.Font, w, h+10, ColorWhite) // horizontal line
				}
			}
		}
	}

	yOffset := wh - 10
	reversedMsg := make([]Message, 0, len(e.engine.MessageLog.Messages[:e.Cursor+1]))
	for _, msg := range e.engine.MessageLog.Messages[:e.Cursor+1] {
		reversedMsg = append(reversedMsg, *msg)
	}
	for i := 0; i < len(reversedMsg)/2; i++ {
		reversedMsg[i], reversedMsg[len(reversedMsg)-i-1] = reversedMsg[len(reversedMsg)-i-1], reversedMsg[i]
	}
	endRender := false
	for _, msg := range reversedMsg {
		wrapped := strings.Split(runewidth.Wrap(msg.FullText(), width), "\n")
		for i := 0; i < len(wrapped)/2; i++ {
			wrapped[i], wrapped[len(wrapped)-i-1] = wrapped[len(wrapped)-i-1], wrapped[i]
		}
		for _, line := range wrapped {
			text.Draw(e.Window, line, e.engine.Font, 10, yOffset, msg.Fg)
			yOffset -= 10
			if yOffset <= 10 {
				endRender = true
				break
			}
		}
		if endRender {
			break
		}
	}

	op := &ebiten.DrawImageOptions{}
	x := float64((width - ww) / 2)
	y := float64((height - wh) / 2)
	op.GeoM.Translate(x, y)
	screen.DrawImage(e.Window, op)
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

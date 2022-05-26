package main

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/mattn/go-runewidth"
	"golang.org/x/image/font"
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
	EvKeyDown(keys []ebiten.Key) action
	HandleAction(act action) (bool, error)
	OnRender(screen *ebiten.Image)
}

type eventHandlerBase struct {
	engine *engine
}

func (e *eventHandlerBase) HandleEvent(keys []ebiten.Key) error {
	_, err := e.HandleAction(e.EvKeyDown(keys))
	return err
}

func (e *eventHandlerBase) EvKeyDown(keys []ebiten.Key) action {
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
			case ebiten.KeyG:
				return newPickupAction(e.engine.Player)
			case ebiten.KeyI:
				e.engine.EventHandler = newInventoryActivateHandler(e.engine)
			case ebiten.KeyD:
				e.engine.EventHandler = newInventoryDropHandler(e.engine)
			default:
			}
		}
	}
	return noneAction{}
}

func (e *eventHandlerBase) HandleAction(act action) (bool, error) {
	switch act := act.(type) {
	case noneAction:
		return false, nil
	default:
		if err := act.Perform(); err != nil {
			switch err.(type) {
			case impossible:
				e.engine.MessageLog.AddMessage(err.Error(), ColorImpossible, true)
				return false, nil
			default:
				return false, err
			}
		}
		if err := e.engine.HandleEnemyTurns(); err != nil {
			return false, err
		}
		e.engine.UpdateFov()
		return true, nil
	}
}

func (e *eventHandlerBase) OnRender(screen *ebiten.Image) {
	e.engine.Render(screen)
}

type askUserEventHandler struct {
	eventHandlerBase
}

func (e *askUserEventHandler) HandleEvent(keys []ebiten.Key) error {
	_, err := e.HandleAction(e.EvKeyDown(keys))
	return err
}

func (e *askUserEventHandler) EvKeyDown(keys []ebiten.Key) action {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		switch p {
		case ebiten.KeyShiftLeft, ebiten.KeyShiftRight, ebiten.KeyControlLeft, ebiten.KeyControlRight, ebiten.KeyAltLeft, ebiten.KeyAltRight:
			return noneAction{}
		default:
			return e.OnExit()
		}
	}
	return noneAction{}
}

func (e *askUserEventHandler) HandleAction(act action) (bool, error) {
	ok, err := e.eventHandlerBase.HandleAction(act)
	if err != nil {
		return false, err
	}
	if ok {
		e.engine.EventHandler = &mainGameEventHandler{eventHandlerBase{engine: e.engine}}
		return true, nil
	}
	return false, nil
}

func (e *askUserEventHandler) OnExit() action {
	e.engine.EventHandler = &mainGameEventHandler{eventHandlerBase{engine: e.engine}}
	return noneAction{}
}

type inventoryEventHandler struct {
	askUserEventHandler
	Title  string
	Window *ebiten.Image
}

func (e *inventoryEventHandler) OnRender(screen *ebiten.Image) {
	e.askUserEventHandler.OnRender(screen)

	numberOfItemsInInventory := len(e.engine.Player.Inventory.Items)
	height := numberOfItemsInInventory*10 + 20

	width := len(e.Title)*10 + 40
	if height <= 30 {
		height = 30
	}
	if e.Window == nil {
		e.Window = ebiten.NewImage(width, height)
	}
	fillWindow(e.Window, width, height, e.Title, e.engine.Font, ColorBlack, ColorWhite)

	if numberOfItemsInInventory > 0 {
		for i, it := range e.engine.Player.Inventory.Items {
			text.Draw(e.Window, fmt.Sprintf("(%s) %s", string(0x41+i), it.Name), e.engine.Font, 10, 20+i*10, ColorWhite)
		}
	} else {
		text.Draw(e.Window, "(Empty)", e.engine.Font, 10, 20, ColorWhite)
	}
	op := &ebiten.DrawImageOptions{}
	x := 0.0
	y := 0.0
	if e.engine.Player.X <= 30 {
		x = 400
	}
	op.GeoM.Translate(x, y)
	screen.DrawImage(e.Window, op)
}

type inventoryActivateHandler struct {
	inventoryEventHandler
}

func newInventoryActivateHandler(e *engine) *inventoryActivateHandler {
	return &inventoryActivateHandler{
		inventoryEventHandler: inventoryEventHandler{
			askUserEventHandler: askUserEventHandler{
				eventHandlerBase: eventHandlerBase{
					engine: e,
				},
			},
			Title: "Select an item to use",
		},
	}
}

func (e *inventoryActivateHandler) HandleEvent(keys []ebiten.Key) error {
	_, err := e.HandleAction(e.EvKeyDown(keys))
	return err
}

func (e *inventoryActivateHandler) EvKeyDown(keys []ebiten.Key) action {
	player := e.engine.Player
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		idx := p - ebiten.KeyA

		if 0 <= idx && idx <= 26 {
			if len(player.Inventory.Items) > int(idx) {
				return e.OnItemSelected(player.Inventory.Items[idx])
			} else {
				e.engine.MessageLog.AddMessage("Invalid entry.", ColorInvalid, true)
				return noneAction{}
			}
		}
	}
	return e.askUserEventHandler.EvKeyDown(keys)
}

func (e inventoryActivateHandler) OnItemSelected(it *item) action {
	return it.Consumable.GetAction(e.engine.Player)
}

type inventoryDropHandler struct {
	inventoryEventHandler
}

func newInventoryDropHandler(e *engine) *inventoryDropHandler {
	return &inventoryDropHandler{
		inventoryEventHandler: inventoryEventHandler{
			askUserEventHandler: askUserEventHandler{
				eventHandlerBase: eventHandlerBase{
					engine: e,
				},
			},
			Title: "Select an item to drop",
		},
	}
}

func (e *inventoryDropHandler) HandleEvent(keys []ebiten.Key) error {
	_, err := e.HandleAction(e.EvKeyDown(keys))
	return err
}

func (e *inventoryDropHandler) EvKeyDown(keys []ebiten.Key) action {
	player := e.engine.Player
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		idx := p - ebiten.KeyA

		if 0 <= idx && idx <= 26 {
			if len(player.Inventory.Items) > int(idx) {
				return e.OnItemSelected(player.Inventory.Items[idx])
			} else {
				e.engine.MessageLog.AddMessage("Invalid entry.", ColorInvalid, true)
				return noneAction{}
			}
		}
	}
	return e.askUserEventHandler.EvKeyDown(keys)
}

func (e inventoryDropHandler) OnItemSelected(it *item) action {
	return &dropItem{
		itemAction: *newItemAction(e.engine.Player, it, nil),
	}
}

type mainGameEventHandler struct {
	eventHandlerBase
}

type gameOverEventHandler struct {
	eventHandlerBase
}

func (e *gameOverEventHandler) EvKeyDown(keys []ebiten.Key) action {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
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
	_ = e.EvKeyDown(keys)
	return nil
}

func (e *historyViewer) EvKeyDown(keys []ebiten.Key) action {
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
	if e.Window == nil {
		e.Window = ebiten.NewImage(width-60, height-60)
	}
	fillWindow(e.Window, width-60, height-60, "Message history", e.engine.Font, ColorWhite, ColorBlack)

	ww, wh := e.Window.Size()
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

func fillWindow(window *ebiten.Image, width, height int, title string, f font.Face, fg, bg color.RGBA) {
	window.Fill(ColorBlack)

	windowTitleLen := (len(title) + 2) * 10
	ww, wh := window.Size()
	windowTitleStart := int(math.Floor(float64(ww-windowTitleLen)/20) * 10)
	for w := 0; w < ww; w += 10 {
		for h := 0; h < wh; h += 10 {
			if w == 0 && h == 0 {
				text.Draw(window, string([]rune{0xda}), f, w, h+10, ColorWhite) // left upper
			} else if w == ww-10 && h == 0 {
				text.Draw(window, string([]rune{0xbf}), f, w, h+10, ColorWhite) // right upper
			} else if w == 0 && h == wh-10 {
				text.Draw(window, string([]rune{0xc0}), f, w, h+10, ColorWhite) // left lower
			} else if w == ww-10 && h == wh-10 {
				text.Draw(window, string([]rune{0xd9}), f, w, h+10, ColorWhite) // right lower
			} else if w == 0 || w == ww-10 {
				text.Draw(window, string([]rune{0xb3}), f, w, h+10, ColorWhite) // vertical line
			} else if h == 0 || h == wh-10 {
				if h == 0 {
					if w == windowTitleStart {
						text.Draw(window, string([]rune{0xdb}), f, w, h+10, bg)
						text.Draw(window, string([]rune{0xb4}), f, w, h+10, ColorWhite) // title left
					} else if w == windowTitleStart+windowTitleLen-10 {
						text.Draw(window, string([]rune{0xdb}), f, w, h+10, bg)
						text.Draw(window, string([]rune{0xc3}), f, w, h+10, ColorWhite) // title right
					} else if w > windowTitleStart && w < windowTitleStart+windowTitleLen-10 {
						text.Draw(window, string([]rune{0xdb}), f, w, h+10, bg)
						text.Draw(window, string(title[(w-windowTitleStart-10)/10]), f, w, h+10, fg)
					} else {
						text.Draw(window, string([]rune{0xc4}), f, w, h+10, ColorWhite) // horizontal line
					}
				} else {
					text.Draw(window, string([]rune{0xc4}), f, w, h+10, ColorWhite) // horizontal line
				}
			}
		}
	}
}

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
	confirmKeys = map[ebiten.Key]struct{}{
		ebiten.KeyEnter: {},
	}
)

type eventHandler interface {
	HandleEvent(keys []ebiten.Key) (eventHandler, error)
	OnRender(screen *ebiten.Image)
}

type eventHandlerBase struct {
	engine *engine
}

func (e *eventHandlerBase) Perform() error {
	return nil
}

func (e *eventHandlerBase) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
		if !e.engine.Player.IsAlive() {
			return &gameOverEventHandler{
				eventHandlerBase: eventHandlerBase{
					engine: e.engine,
				},
			}, nil
		}
	}
	return e, nil
}

func (e *eventHandlerBase) EvKeyDown(keys []ebiten.Key) interface{} {
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
				return newHistoryViewer(e.engine)
			case ebiten.KeyG:
				return newPickupAction(e.engine.Player)
			case ebiten.KeyI:
				return newInventoryActivateHandler(e.engine)
			case ebiten.KeyD:
				return newInventoryDropHandler(e.engine)
			case ebiten.KeySlash:
				return &lookHandler{selectIndexHandler: newSelectIndexHandler(e.engine)}
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

func (e *askUserEventHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *askUserEventHandler) EvKeyDown(keys []ebiten.Key) interface{} {
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

func (e *askUserEventHandler) HandleAction(act action) (eventHandler, error) {
	ok, err := e.eventHandlerBase.HandleAction(act)
	if err != nil {
		return nil, err
	}
	if ok {
		return &mainGameEventHandler{eventHandlerBase{engine: e.engine}}, nil
	}
	return nil, nil
}

func (e *askUserEventHandler) OnExit() action {
	return &mainGameEventHandler{eventHandlerBase{engine: e.engine}}
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

func (e *inventoryActivateHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *inventoryActivateHandler) EvKeyDown(keys []ebiten.Key) interface{} {
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

func (e *inventoryDropHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *inventoryDropHandler) EvKeyDown(keys []ebiten.Key) interface{} {
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

type selectIndexHandler struct {
	askUserEventHandler
	Player *actor
}

func newSelectIndexHandler(e *engine) selectIndexHandler {
	e.MouseLocation = [2]int{e.Player.X, e.Player.Y}
	return selectIndexHandler{
		askUserEventHandler: askUserEventHandler{
			eventHandlerBase: eventHandlerBase{
				engine: e,
			},
		},
		Player: e.Player,
	}
}

func (e *selectIndexHandler) OnRender(screen *ebiten.Image) {
	e.askUserEventHandler.OnRender(screen)

	x, y := e.engine.MouseLocation[0], e.engine.MouseLocation[1]

	text.Draw(screen, string([]rune{0xdb}), e.engine.Font, x*10, y*10, ColorSelect)
}

func (e *selectIndexHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *selectIndexHandler) EvKeyDown(keys []ebiten.Key) interface{} {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		if _, ok := moveKeys[p]; ok {
			modifier := 1 // Holding modifier keys will speed up key movement
			for _, m := range keys {
				if m == ebiten.KeyShiftLeft || m == ebiten.KeyShiftRight {
					modifier *= 5
					break
				} else if m == ebiten.KeyControlLeft || m == ebiten.KeyControlRight {
					modifier *= 10
					break
				} else if m == ebiten.KeyAltLeft || m == ebiten.KeyAltRight {
					modifier *= 20
					break
				}
			}
			x, y := e.engine.MouseLocation[0], e.engine.MouseLocation[1]
			dx, dy := moveKeys[p][0], moveKeys[p][1]
			x += dx * modifier
			y += dy * modifier
			// Clamp the cursor index to the map size
			x = int(math.Max(0, math.Min(float64(x), float64(e.engine.GameMap.Width-1))))
			y = int(math.Max(0, math.Min(float64(y), float64(e.engine.GameMap.Height-1))))
			e.engine.MouseLocation = [2]int{x, y}
			return noneAction{}
		}
		e.askUserEventHandler.EvKeyDown(keys)
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return e.OnMouseButtonDown()
	}
	return noneAction{}
}

func (e *selectIndexHandler) OnMouseButtonDown() action {
	return e.OnIndexSelected(e.engine.MouseLocation[0], e.engine.MouseLocation[1])
}

func (e *selectIndexHandler) OnIndexSelected(x, y int) action {
	return noneAction{}
}

type lookHandler struct {
	selectIndexHandler
}

func (e *lookHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *lookHandler) EvKeyDown(keys []ebiten.Key) interface{} {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		if _, ok := confirmKeys[p]; ok {
			return e.OnIndexSelected(e.engine.MouseLocation[0], e.engine.MouseLocation[1])
		}
	}
	act := e.selectIndexHandler.EvKeyDown(keys)
	if _, ok := act.(noneAction); !ok {
		return act
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return e.OnMouseButtonDown()
	}
	return noneAction{}
}

func (e *lookHandler) OnMouseButtonDown() eventHandler {
	return e.OnIndexSelected(e.engine.MouseLocation[0], e.engine.MouseLocation[1])
}

func (e *lookHandler) OnIndexSelected(x, y int) eventHandler {
	return &mainGameEventHandler{eventHandlerBase{engine: e.engine}}
}

type singleRangedAttackHandler struct {
	selectIndexHandler
	Callback func(x, y int) action
}

func (e *singleRangedAttackHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *singleRangedAttackHandler) EvKeyDown(keys []ebiten.Key) interface{} {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		if _, ok := confirmKeys[p]; ok {
			return e.OnIndexSelected(e.engine.MouseLocation[0], e.engine.MouseLocation[1])
		}
	}
	act := e.selectIndexHandler.EvKeyDown(keys)
	if _, ok := act.(noneAction); !ok {
		return act
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return e.OnMouseButtonDown()
	}
	return noneAction{}
}

func (e *singleRangedAttackHandler) OnMouseButtonDown() action {
	return e.OnIndexSelected(e.engine.MouseLocation[0], e.engine.MouseLocation[1])
}

func (e *singleRangedAttackHandler) OnIndexSelected(x, y int) action {
	return e.Callback(x, y)
}

type areaRangedAttackHandler struct {
	selectIndexHandler
	Radius   int
	Callback func(x, y int) action
}

func (e *areaRangedAttackHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *areaRangedAttackHandler) OnRender(screen *ebiten.Image) {
	e.selectIndexHandler.OnRender(screen)

	x, y := e.engine.MouseLocation[0], e.engine.MouseLocation[1]

	drawRange(screen, x*10-e.Radius*10-10, y*10-e.Radius*10-20, int(math.Pow(float64(e.Radius), 2))*10, int(math.Pow(float64(e.Radius), 2))*10, e.engine.Font, ColorRed)
}

func (e *areaRangedAttackHandler) EvKeyDown(keys []ebiten.Key) interface{} {
	for _, p := range keys {
		if !repeatingKeyPressed(p) {
			continue
		}
		if _, ok := confirmKeys[p]; ok {
			return e.OnIndexSelected(e.engine.MouseLocation[0], e.engine.MouseLocation[1])
		}
	}
	act := e.selectIndexHandler.EvKeyDown(keys)
	if _, ok := act.(noneAction); !ok {
		return act
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return e.OnMouseButtonDown()
	}
	return noneAction{}
}

func (e *areaRangedAttackHandler) OnMouseButtonDown() action {
	return e.OnIndexSelected(e.engine.MouseLocation[0], e.engine.MouseLocation[1])
}

func (e *areaRangedAttackHandler) OnIndexSelected(x, y int) action {
	return e.Callback(x, y)
}

type mainGameEventHandler struct {
	eventHandlerBase
}

type gameOverEventHandler struct {
	eventHandlerBase
}

func (e *gameOverEventHandler) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *gameOverEventHandler) EvKeyDown(keys []ebiten.Key) interface{} {
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

func (e *historyViewer) HandleEvent(keys []ebiten.Key) (eventHandler, error) {
	state := e.EvKeyDown(keys)
	if h, ok := state.(eventHandler); ok {
		return h, nil
	}
	if a, ok := state.(action); ok {
		_, err := e.HandleAction(a)
		if err != nil {
			return nil, err
		}
	}
	return e, nil
}

func (e *historyViewer) EvKeyDown(keys []ebiten.Key) interface{} {
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
				return &mainGameEventHandler{eventHandlerBase{engine: e.engine}}
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

func drawRange(screen *ebiten.Image, x, y, width, height int, f font.Face, c color.RGBA) {
	for w := x; w < x+width; w += 10 {
		for h := y; h < y+height; h += 10 {
			if w == x && h == y {
				text.Draw(screen, string([]rune{0xda}), f, w, h+10, c) // left upper
			} else if w == x+width-10 && h == y {
				text.Draw(screen, string([]rune{0xbf}), f, w, h+10, c) // right upper
			} else if w == x && h == y+height-10 {
				text.Draw(screen, string([]rune{0xc0}), f, w, h+10, c) // left lower
			} else if w == x+width-10 && h == y+height-10 {
				text.Draw(screen, string([]rune{0xd9}), f, w, h+10, c) // right lower
			} else if w == x || w == x+width-10 {
				text.Draw(screen, string([]rune{0xb3}), f, w, h+10, c) // vertical line
			} else if h == y || h == y+height-10 {
				text.Draw(screen, string([]rune{0xc4}), f, w, h+10, c) // horizontal line
			}
		}
	}
}

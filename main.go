package main

import (
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hatajoe/rogueliketutorials/resources/fonts"
)

const (
	screenWidth    int = 800
	screenHeight   int = 500
	screenTileSize int = 10

	mapWidth           int = 80
	mapHeight          int = 43
	roomMaxSize        int = 10
	roomMinSize        int = 6
	maxRooms           int = 30
	maxMonstersPerRoom int = 2
	maxItemsPerRoom    int = 2
)

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	mx, my := ebiten.CursorPosition()
	if gameEngine.GameMap.InBounds(int(mx/10), int(my/10)) {
		gameEngine.MouseLocation = [2]int{int(mx / 10), int(my/10) + 1} // I don't know why +1 is needed, however this worked well
	}
	var err error
	if handler, err = handler.HandleEvent(g.keys); err != nil {
		switch err.(type) {
		case impossible:
			gameEngine.MessageLog.AddMessage(err.Error(), ColorError, false)
			return nil
		case QuitWithoutSaving:
			// TODO: save
			return err
		default:
			// TODO: save
			return err
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	handler.OnRender(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	qbicfeetFont font.Face
	gameEngine   *engine
	handler      eventHandler
)

func init() {
	tt, err := opentype.Parse(fonts.Qbicfeet_ttf)
	if err != nil {
		log.Fatal(err)
	}

	qbicfeetFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(screenTileSize),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	player := newPlayer()

	gameEngine = NewEngine(player, qbicfeetFont)
	gameEngine.GameMap = generateDungeon(
		maxRooms,
		roomMinSize,
		roomMaxSize,
		screenWidth/10,
		(screenHeight-50)/10,
		maxMonstersPerRoom,
		maxItemsPerRoom,
		gameEngine,
	)
	gameEngine.UpdateFov()
	gameEngine.MessageLog.AddMessage(
		"Hello and welcome, adventure, to yet another dungeon!",
		ColorWelcomText,
		true,
	)

	handler = &mainGameEventHandler{
		eventHandlerBase: eventHandlerBase{
			engine: gameEngine,
		},
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")

	if err := ebiten.RunGame(&Game{}); err != nil && err != regularTermination {
		log.Fatal(err)
	}
	os.Exit(0)
}

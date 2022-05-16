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

	roomMaxSize        int = 10
	roomMinSize        int = 6
	maxRooms           int = 30
	maxMonstersPerRoom int = 2
)

type Game struct {
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	return engine.HandleEvent(g.keys)
}

func (g *Game) Draw(screen *ebiten.Image) {
	engine.Render(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 500
}

var (
	qbicfeetFont font.Face
	engine       *Engine
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

	eventHandler := &EventHandler{}
	player := NewPlayer()

	gameMap := GenerateDungeon(
		maxRooms,
		roomMinSize,
		roomMaxSize,
		screenWidth/10,
		(screenHeight-50)/10,
		maxMonstersPerRoom,
		player,
	)

	engine = NewEngine(eventHandler, gameMap, player, qbicfeetFont)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Yet Another Roguelike Tutorial")

	if err := ebiten.RunGame(&Game{}); err != nil && err != regularTermination {
		log.Fatal(err)
	}
	os.Exit(0)
}

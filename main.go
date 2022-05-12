package main

import (
	"image/color"
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

	roomMaxSize int = 10
	roomMinSize int = 6
	maxRooms    int = 30
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
	player := &Entity{
		x:    int(screenWidth / 10 / 2),
		y:    int(screenHeight / 10 / 2),
		char: "@",
		color: color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		},
	}
	npc := &Entity{
		x:    int(screenWidth/10/2) - 5,
		y:    int(screenHeight / 10 / 2),
		char: "@",
		color: color.RGBA{
			R: 255,
			G: 255,
			B: 0,
			A: 255,
		},
	}
	entities := []*Entity{npc, player}

	gameMap := GenerateDungeon(
		maxRooms,
		roomMinSize,
		roomMaxSize,
		screenWidth/10,
		(screenHeight-50)/10,
		player,
	)

	engine = &Engine{
		entities:     entities,
		eventHandler: eventHandler,
		GameMap:      gameMap,
		player:       player,
		font:         qbicfeetFont,
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

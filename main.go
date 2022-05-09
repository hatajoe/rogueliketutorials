package main

import (
	"errors"
	"image/color"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hatajoe/rogueliketutorials/resources/fonts"
)

const (
	screenWidth  int = 640
	screenHeight int = 480
)

type Game struct{
	keys []ebiten.Key
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	action := eventHandler.KeyDown(g.keys)
	switch act := action.(type) {
	case MovementAction:
		player.X += act.Dx
		player.Y += act.Dy
	case EscapeAction:
		return regularTermination
	default:
		return nil
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, "@", qbicfeetFont, player.X, player.Y, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

type Player struct {
	X int
	Y int
}

var (
	qbicfeetFont font.Face
	eventHandler *EventHandler
	player       *Player

	regularTermination = errors.New("regular termination")
)

func init() {
	tt, err := opentype.Parse(fonts.Qbicfeet_ttf)
	if err != nil {
		log.Fatal(err)
	}

	qbicfeetFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	eventHandler = &EventHandler{}

	player = &Player{
		X: int(screenWidth / 4),
		Y: int(screenHeight / 4),
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(&Game{}); err != nil && err != regularTermination {
		log.Fatal(err)
	}
	os.Exit(0)
}

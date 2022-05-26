package main

import (
	"fmt"
	"strings"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func RenderBar(screen *ebiten.Image, font font.Face, currentValue, maximumValue, totalWidth int) {
	barWidth := float64(currentValue) / float64(maximumValue) * float64(totalWidth)

	ebitenutil.DrawRect(screen, 0, 450, float64(totalWidth), 10, ColorBarEmpty)

	if barWidth > 0 {
		ebitenutil.DrawRect(screen, 0, 450, barWidth, 10, ColorBarFilled)
	}
	text.Draw(screen, fmt.Sprintf("HP: %d/%d", currentValue, maximumValue), font, 1, 450, ColorBarText)
}

func RenderNamesAtMouseLocation(screen *ebiten.Image, font font.Face, x, y int, e *engine) {
	mx, my := e.MouseLocation[0], e.MouseLocation[1]

	namesAtMouseLocation := getNamesAtLocation(mx, my, e.GameMap)
	text.Draw(screen, namesAtMouseLocation, font, x*10, y*10, ColorWhite)
}

func getNamesAtLocation(x, y int, gm *gameMap) string {
	if !gm.InBounds(x, y) || !gm.Visible[x][y] {
		return ""
	}
	names := []string{}
	for _, entity := range gm.Entities {
		e := entity.Entity()
		if e.X == x && e.Y == y {
			names = append(names, e.Name)
		}
	}
	return strings.Join(names, ",")
}

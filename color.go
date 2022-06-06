package main

import "image/color"

var (
	ColorWhite = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	ColorBlack = color.RGBA{R: 0x0, G: 0x0, B: 0x0, A: 0xFF}
	ColorRed   = color.RGBA{R: 0xFF, G: 0x0, B: 0x0, A: 0xFF}

	ColorPlayerAtk           = color.RGBA{R: 0xE0, G: 0xE0, B: 0xE0, A: 0xFF}
	ColorEnemyAtk            = color.RGBA{R: 0xFF, G: 0xC0, B: 0xC0, A: 0xFF}
	ColorNeedsTarget         = color.RGBA{R: 0x3F, G: 0xFF, B: 0xFF, A: 0xFF}
	ColorStatusEffectApplied = color.RGBA{R: 0x3F, G: 0xFF, B: 0x3F, A: 0xFF}

	ColorPlayerDie = color.RGBA{R: 0xFF, G: 0x30, B: 0x30, A: 0xFF}
	ColorEnemyDie  = color.RGBA{R: 0xFF, G: 0xA0, B: 0x30, A: 0xFF}

	ColorInvalid    = color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}
	ColorImpossible = color.RGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xFF}
	ColorError      = color.RGBA{R: 0xFF, G: 0x40, B: 0x40, A: 0xFF}

	ColorWelcomText      = color.RGBA{R: 0x20, G: 0xA0, B: 0xFF, A: 0xFF}
	ColorHealthRecovered = color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}

	ColorBarText   = ColorWhite
	ColorBarFilled = color.RGBA{R: 0x00, G: 0x60, B: 0x00, A: 0xFF}
	ColorBarEmpty  = color.RGBA{R: 0x40, G: 0x10, B: 0x10, A: 0xFF}

	ColorSelect = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x88}
)

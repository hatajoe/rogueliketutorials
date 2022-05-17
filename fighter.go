package main

import (
	"fmt"
	"image/color"
	"math"
)

type fighter struct {
	*baseComponent
	MaxHP   int
	HP      int
	Defense int
	Power   int
}

func (c *fighter) SetHP(val int) {
	c.HP = int(math.Max(0, math.Min(float64(val), float64(c.MaxHP))))
	if c.HP == 0 && c.Entity.AI != nil {
		c.Die()
	}
}

func (c *fighter) Die() {
	deathMessage := ""
	if c.Entity == c.Engine().Player {
		deathMessage = "You died!"
		c.Engine().EventHandler = &gameOverEventHandler{engine: c.Engine()}
	} else {
		deathMessage = fmt.Sprintf("%s is dead!", c.Entity.Name)
	}

	c.Entity.Char = "%"
	c.Entity.Color = color.RGBA{R: 191, G: 0, B: 0, A: 255}
	c.Entity.BlocksMovement = false
	c.Entity.AI = nil
	c.Entity.Name = fmt.Sprintf("remains of %s", c.Entity.Name)
	c.Entity.RenderOrder = RenderOrder_Corpse

	fmt.Println(deathMessage)
}

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

func (c *fighter) setHP(val int) {
	c.HP = int(math.Max(0, math.Min(float64(val), float64(c.MaxHP))))
	if p, ok := c.Parent.(*actor); ok {
		if c.HP == 0 && p.AI != nil {
			c.Die()
		}
	}
}

func (c *fighter) Die() {
	p, ok := c.Parent.(*actor)
	if !ok {
		return
	}
	deathMessage := fmt.Sprintf("%s is dead!", p.Name)
	deathMessageColor := ColorEnemyDie
	if c.Parent == c.Engine().Player {
		deathMessage = "You died!"
		deathMessageColor = ColorPlayerDie
	}

	p.Char = "%"
	p.Color = color.RGBA{R: 191, G: 0, B: 0, A: 255}
	p.BlocksMovement = false
	p.AI = nil
	p.Name = fmt.Sprintf("remains of %s", p.Name)
	p.RO = RenderOrder_Corpse

	c.Engine().MessageLog.AddMessage(deathMessage, deathMessageColor, true)
}

func (c *fighter) Heal(amount int) int {
	if c.HP == c.MaxHP {
		return 0
	}

	newHpValue := c.HP + amount
	if newHpValue > c.MaxHP {
		newHpValue = c.MaxHP
	}

	amountRecoverd := newHpValue - c.HP

	c.HP = newHpValue

	return amountRecoverd
}

func (c *fighter) TakeDamage(amount int) {
	c.HP -= amount
	c.setHP(c.HP)
}

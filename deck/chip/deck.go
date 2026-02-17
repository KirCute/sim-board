package chip

import (
	"fmt"
	"math/rand"

	"github.com/KirCute/sim-board"
)

type Chip struct {
	Pool     []sim_board.Card
	shuffled bool
	*Params
}

func (c *Chip) Type() string {
	return Name
}

func (c *Chip) Name() string {
	if len(c.CustomName) == 0 {
		return Name
	}
	return c.CustomName
}

func (c *Chip) RestLen() int {
	return len(c.Pool)
}

func (c *Chip) MaxLen() int {
	return c.Count1 + c.Count5 + c.Count20 + c.Count100 + c.Count500 + c.Count2k + c.Count1w
}

func (c *Chip) Return(card sim_board.Card) {
	c.Pool = append(c.Pool, card)
	c.shuffled = false
}

func (c *Chip) Draw(count int) []sim_board.Card {
	if !c.shuffled && count < c.RestLen() {
		rand.Shuffle(len(c.Pool), func(i, j int) {
			c.Pool[i], c.Pool[j] = c.Pool[j], c.Pool[i]
		})
		c.shuffled = true
	}
	ret := c.Pool[:count]
	c.Pool = c.Pool[count:]
	return ret
}

func GetHTML(card sim_board.Card) (string, bool) {
	var color string
	var size, fontSize int
	switch card {
	case "1":
		color = "grey"
		size = 60
		fontSize = 30
	case "5":
		color = "limegreen"
		size = 65
		fontSize = 30
	case "20":
		color = "dodgerblue"
		size = 70
		fontSize = 30
	case "100":
		color = "rebeccapurple"
		size = 75
		fontSize = 25
	case "500":
		color = "gold"
		size = 80
		fontSize = 25
	case "2000":
		color = "firebrick"
		size = 85
		fontSize = 20
	case "10000":
		color = "black"
		size = 90
		fontSize = 20
	}
	return fmt.Sprintf(`<div style="width: %dpx; aspect-ratio: 1; border-radius: 50%%; background: radial-gradient(circle at center, white 50%%, %s 50.1%%); display: grid; place-items: center; font-size: %dpx; font-weight: bold; color: black; box-shadow: 0 2px 5px rgba(0,0,0,0.4)">%s</div>`, size, color, fontSize, card), true
}

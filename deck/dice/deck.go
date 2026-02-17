package dice

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/KirCute/sim-board"
)

type Dice struct {
	*Params
}

func (p *Dice) Type() string {
	return Name
}

func (p *Dice) Name() string {
	if len(p.CustomName) == 0 {
		if p.Face == 6 {
			return Name
		}
		return fmt.Sprintf("%d面%s", p.Face, Name)
	}
	return p.CustomName
}

func (p *Dice) RestLen() int {
	return -1
}

func (p *Dice) MaxLen() int {
	return -1
}

func (p *Dice) Return(_ sim_board.Card) {
}

func (p *Dice) Draw(count int) []sim_board.Card {
	ret := make([]sim_board.Card, 0, count)
	for i := 0; i < count; i++ {
		point := rand.Intn(p.Face) + 1
		ret = append(ret, sim_board.Card(strconv.Itoa(point)))
	}
	return ret
}

func GetHTML(card sim_board.Card) (string, bool) {
	var content string
	switch card {
	case "1":
		content = `<div style="display: block; width: 12px; height: 12px; background: red; border-radius: 50%; grid-area: 2/2/3/3;"/>`
	case "2":
		content = `
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 1/1/2/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 3/3/4/4;"></div>
`
	case "3":
		content = `
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 1/1/2/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 2/2/3/3;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 3/3/4/4;"></div>
`
	case "4":
		content = `
<div style="display: block; width: 8px; height: 8px; background: red; border-radius: 50%; grid-area: 1/1/2/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: red; border-radius: 50%; grid-area: 1/3/2/4;"></div>
<div style="display: block; width: 8px; height: 8px; background: red; border-radius: 50%; grid-area: 3/1/4/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: red; border-radius: 50%; grid-area: 3/3/4/4;"></div>`
	case "5":
		content = `
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 1/1/2/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 1/3/2/4;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 2/2/3/3;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 3/1/4/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 3/3/4/4;"></div>
`
	case "6":
		content = `
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 1/1/2/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 1/2/2/3;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 1/3/2/4;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 3/1/4/2;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 3/2/4/3;"></div>
<div style="display: block; width: 8px; height: 8px; background: #1e1e2f; border-radius: 50%; grid-area: 3/3/4/4;"></div>
`
	default:
		content = fmt.Sprintf(`<div style="font-size: 20px; font-weight: bold">%s</div>`, card)
	}
	return fmt.Sprintf(`
<div style="width: 40px; aspect-ratio: 1; background: white; border-radius: 15px; box-shadow: 0 4px 8px rgba(0,0,0,0.2), 0 2px 4px rgba(0,0,0,0.1); display: grid; grid-template-columns: repeat(3, 1fr); place-items: center; overflow: hidden; padding: 7px">
  %s
</div>
`, content), true
}

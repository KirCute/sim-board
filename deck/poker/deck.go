package poker

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/KirCute/sim-board"
)

type Poker struct {
	rest     []sim_board.Card
	shuffled bool
	*Params
}

func (p *Poker) Type() string {
	return Name
}

func (p *Poker) Name() string {
	if len(p.CustomName) == 0 {
		if p.Count == 1 && p.CountRank == 13 && p.CountSuit == 4 {
			return Name
		}
		return fmt.Sprintf("%s%d副%d色%d值", Name, p.Count, p.CountSuit, p.CountRank)
	}
	return p.CustomName
}

func (p *Poker) RestLen() int {
	return len(p.rest)
}

func (p *Poker) MaxLen() int {
	return p.CountRank*p.CountSuit*p.Count + p.CountRedJoker + p.CountBlackJoker
}

func (p *Poker) Return(card sim_board.Card) {
	p.rest = append(p.rest, card)
	p.shuffled = false
}

func (p *Poker) Draw(count int) []sim_board.Card {
	if !p.shuffled {
		rand.Shuffle(len(p.rest), func(i, j int) {
			p.rest[i], p.rest[j] = p.rest[j], p.rest[i]
		})
		p.shuffled = true
	}
	ret := p.rest[:count]
	p.rest = p.rest[count:]
	return ret
}

func GetHTML(card sim_board.Card) (string, bool) {
	if card == "rj" {
		return `<div style="color: red; width: 65px; aspect-ratio: 0.7222; border-radius: 10px; background-color: white; display: grid; place-items: center; box-shadow: 0 2px 5px rgba(0,0,0,0.4)">JOKER</div>`, true
	}
	if card == "bj" {
		return `<div style="color: black; width: 65px; aspect-ratio: 0.7222; border-radius: 10px; background-color: white; display: grid; place-items: center; box-shadow: 0 2px 5px rgba(0,0,0,0.4)">JOKER</div>`, true
	}
	var suit, rank int
	n, _ := fmt.Sscanf(string(card), "%d-%d", &suit, &rank)
	if n != 2 || suit < 0 || suit > 3 || rank < 0 || rank > 12 {
		return "", false
	}
	var color, suitStr, rankStr string
	switch suit {
	case 0:
		color = "black"
		suitStr = "♠"
	case 1:
		color = "red"
		suitStr = "♥"
	case 2:
		color = "green"
		suitStr = "♣"
	case 3:
		color = "blue"
		suitStr = "♦"
	}
	switch rank {
	case 12:
		rankStr = "A"
	case 11:
		rankStr = "K"
	case 10:
		rankStr = "Q"
	case 9:
		rankStr = "J"
	default:
		rankStr = strconv.Itoa(rank + 2)
	}
	return fmt.Sprintf(`<div style="color: %s; width: 65px; aspect-ratio: 0.7222; border-radius: 10px; background-color: white; display: grid; place-items: center; box-shadow: 0 2px 5px rgba(0,0,0,0.4)">%s<br/>%s</div>`, color, suitStr, rankStr), true
}

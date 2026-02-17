package sim_board

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Card string

type Deck interface {
	Type() string
	Name() string
	RestLen() int
	MaxLen() int
	Return(card Card)
	Draw(count int) []Card
}

type DeckCard struct {
	DeckId int
	Card
}

func (d *DeckCard) UnmarshalText(data []byte) error {
	s := string(data)
	deckStr, card, ok := strings.Cut(s, "@")
	if !ok {
		return errors.New("invalid deck card: no sep @")
	}
	deckId, err := strconv.Atoi(deckStr)
	if err != nil {
		return fmt.Errorf("failed to parse deck id: %+v", err)
	}
	d.DeckId = deckId
	d.Card = Card(card)
	return nil
}

func (d DeckCard) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%d@%s", d.DeckId, d.Card)
	return []byte(s), nil
}

type PublicCard struct {
	Card DeckCard `json:"card"`
	X    float32  `json:"x"`
	Y    float32  `json:"y"`
	OpID uint     `json:"op_id"`
	PlID uint     `json:"pl_id"`
}

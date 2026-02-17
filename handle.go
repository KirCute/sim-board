package sim_board

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type MarshaledDeck struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	MaxLen  int    `json:"max_len"`
	RestLen int    `json:"rest_len"`
}

func (r *Room) marshalDeck() []MarshaledDeck {
	ret := make([]MarshaledDeck, 0, len(r.decks))
	for _, d := range r.decks {
		ret = append(ret, MarshaledDeck{
			Type:    d.Type(),
			Name:    d.Name(),
			MaxLen:  d.MaxLen(),
			RestLen: d.RestLen(),
		})
	}
	return ret
}

type BroadcastResponse struct {
	Board   map[string]*PublicCard `json:"board"`
	Hole    map[DeckCard]int       `json:"hole"`
	Decks   []MarshaledDeck        `json:"decks"`
	Players []string               `json:"players"`
}

func (r *Room) makeBroadcastResp() *BroadcastResponse {
	players := make([]string, 0, len(r.conn))
	for player := range r.conn {
		players = append(players, player)
	}
	return &BroadcastResponse{
		Board:   r.board,
		Decks:   r.marshalDeck(),
		Players: players,
	}
}

func (r *Room) broadcast(except ...string) {
	ret := r.makeBroadcastResp()
	for player, hole := range r.hole {
		if SliceContains(except, player) {
			continue
		}
		ret.Hole = hole
		msg := &ServerMessage{Type: "broadcast", Data: ret}
		r.sendMsgTo(player, msg)
	}
}

type WelcomeResponse struct {
	Broadcast      *BroadcastResponse          `json:"broadcast"`
	AvailableDecks map[string][]map[string]any `json:"available_decks"`
}

func (r *Room) handleWelcome(player string) {
	b := r.makeBroadcastResp()
	b.Hole = r.hole[player]
	ret := &WelcomeResponse{
		Broadcast:      b,
		AvailableDecks: GetAllAvailableDecks(),
	}
	r.sendMsgTo(player, &ServerMessage{Type: "welcome", Data: ret})
	r.broadcast(player)
}

func getRandomPos() (float32, float32) {
	x := rand.Float32()/2.0 + 0.25
	y := rand.Float32()/2.0 + 0.25
	return x, y
}

type DrawArgs struct {
	Deck   int    `json:"deck"`
	Num    int    `json:"num"`
	Target string `json:"target"`
}

func (r *Room) handleDraw(player string, args DrawArgs) {
	if args.Deck >= len(r.decks) {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "牌堆不存在"})
		return
	}
	d := r.decks[args.Deck]
	if d.RestLen() >= 0 && args.Num > d.RestLen() {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "数量不足"})
		return
	}
	hole, ok := r.hole[args.Target]
	if !ok && args.Target != "" {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "目标不存在"})
		return
	}
	cards := d.Draw(args.Num)
	if args.Target == "" {
		for _, card := range cards {
			x, y := getRandomPos()
			r.board[uuid.NewString()] = &PublicCard{
				Card: DeckCard{
					DeckId: args.Deck,
					Card:   card,
				},
				X:    x,
				Y:    y,
				OpID: 0,
				PlID: r.placeCnter,
			}
			r.placeCnter++
		}
	} else {
		for _, card := range cards {
			dc := DeckCard{
				DeckId: args.Deck,
				Card:   card,
			}
			if _, ok = hole[dc]; !ok {
				hole[dc] = 0
			}
			hole[dc]++
		}
	}
	r.broadcast()
}

type AnnounceArgs struct {
	DeckCard DeckCard `json:"card"`
	X        float32  `json:"x"`
	Y        float32  `json:"y"`
}

func (r *Room) handleAnnounce(player string, args AnnounceArgs) {
	hole := r.hole[player]
	if i, ok := hole[args.DeckCard]; !ok || i <= 0 {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "手牌余量不足"})
		return
	}
	if hole[args.DeckCard] <= 1 {
		delete(hole, args.DeckCard)
	} else {
		hole[args.DeckCard]--
	}
	r.board[uuid.NewString()] = &PublicCard{
		Card: args.DeckCard,
		X:    args.X,
		Y:    args.Y,
		OpID: 0,
		PlID: r.placeCnter,
	}
	r.placeCnter++
	r.broadcast()
}

type CollectArgs struct {
	ID   string `json:"id"`
	OpID uint   `json:"op_id"`
}

func (r *Room) handleCollect(player string, args CollectArgs) {
	card, ok := r.board[args.ID]
	if !ok {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "公共牌不存在"})
		return
	}
	if card.OpID != args.OpID {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "操作超时"})
		return
	}
	if _, ok = r.hole[player][card.Card]; ok {
		r.hole[player][card.Card]++
	} else {
		r.hole[player][card.Card] = 1
	}
	delete(r.board, args.ID)
	r.broadcast()
}

func (r *Room) handleDiscardBoard(player string, args CollectArgs) {
	card, ok := r.board[args.ID]
	if !ok {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "公共牌不存在"})
		return
	}
	if card.OpID != args.OpID {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "操作超时"})
		return
	}
	if card.Card.DeckId >= len(r.decks) {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "牌堆不存在"})
		return
	}
	r.decks[card.Card.DeckId].Return(card.Card.Card)
	delete(r.board, args.ID)
	r.broadcast()
}

func (r *Room) handleDiscardHole(player string, card DeckCard) {
	hole := r.hole[player]
	if i, ok := hole[card]; !ok || i <= 0 {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "手牌余量不足"})
		return
	}
	if card.DeckId >= len(r.decks) {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "牌堆不存在"})
		return
	}
	if hole[card] <= 1 {
		delete(hole, card)
	} else {
		hole[card]--
	}
	r.decks[card.DeckId].Return(card.Card)
	r.broadcast()
}

func (r *Room) handleReset() {
	for id, card := range r.board {
		r.decks[card.Card.DeckId].Return(card.Card.Card)
		delete(r.board, id)
	}
	for _, hole := range r.hole {
		for card, cnt := range hole {
			for i := 0; i < cnt; i++ {
				r.decks[card.DeckId].Return(card.Card)
			}
			delete(hole, card)
		}
	}
	r.broadcast()
}

func (r *Room) handleAllCollect(player string) {
	hole := r.hole[player]
	for id, card := range r.board {
		if _, ok := hole[card.Card]; !ok {
			hole[card.Card] = 0
		}
		hole[card.Card]++
		delete(r.board, id)
	}
	r.broadcast()
}

type AddDeckArgs struct {
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}

func (r *Room) handleAddDeck(player string, args AddDeckArgs) {
	d, err := NewDeck(args.Name, args.Params)
	if err != nil {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: fmt.Sprintf("添加牌堆失败：%v", err.Error())})
		return
	}
	r.decks = append(r.decks, d)
	r.broadcast()
}

type MoveArgs struct {
	ID   string  `json:"id"`
	OpID uint    `json:"op_id"`
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
}

func (r *Room) handleMove(player string, args MoveArgs) {
	card, ok := r.board[args.ID]
	if !ok {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: "公共牌不存在"})
		return
	}
	if card.OpID != args.OpID {
		r.sendMsgTo(player, &ServerMessage{Type: "error", Data: fmt.Sprintf("操作超时, src=%d, dst=%d", args.OpID, card.OpID)})
		return
	}
	card.X = min(max(args.X, .0), 1.0)
	card.Y = min(max(args.Y, .0), 1.0)
	card.OpID++
	card.PlID = r.placeCnter
	r.placeCnter++
	r.broadcast()
}

func (r *Room) sendMsgTo(player string, msg *ServerMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		logrus.Errorf("write message to %s failed due to marshal error, err=%v, msg_type=%s", player, err, msg.Type)
		msg = &ServerMessage{Type: "fatal", Data: err.Error()}
	}
	for _, conn := range r.conn[player] {
		err = conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			logrus.Errorf("write message to %s failed, err=%v, msg=%s", conn.RemoteAddr().String(), err, string(data))
		}
	}
}

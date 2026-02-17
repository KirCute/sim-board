package poker

import (
	"fmt"
	"reflect"

	"github.com/KirCute/sim-board"
)

const Name = "扑克牌"

type Params struct {
	CustomName      string `json:"custom_name" label:"自定义名称" type:"string"`
	Count           int    `json:"count" label:"副数" type:"int" min:"1" default:"1"`
	CountSuit       int    `json:"count_suit" label:"花色种类" type:"int" min:"1" max:"4" default:"4"`
	CountRank       int    `json:"count_rank" label:"数值范围" type:"int" min:"1" max:"13" default:"13"`
	CountRedJoker   int    `json:"count_red_joker" label:"大王总数量" type:"int" default:"1"`
	CountBlackJoker int    `json:"count_black_joker" label:"小王总数量" type:"int" default:"1"`
}

func Create(params *Params) *Poker {
	params.Count = max(params.Count, 1)
	params.CountSuit = min(params.CountSuit, 4)
	params.CountSuit = max(params.CountSuit, 1)
	params.CountRank = min(params.CountRank, 13)
	params.CountRank = max(params.CountRank, 1)
	ret := &Poker{Params: params}
	for k := 0; k < params.Count; k++ {
		for i := 0; i < params.CountSuit; i++ {
			for j := 0; j < params.CountRank; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%d-%d", i, 12-j)))
			}
		}
	}
	for i := 0; i < params.CountRedJoker; i++ {
		ret.rest = append(ret.rest, "rj")
	}
	for i := 0; i < params.CountBlackJoker; i++ {
		ret.rest = append(ret.rest, "bj")
	}
	return ret
}

func init() {
	sim_board.RegisterDeck(Name, reflect.ValueOf(Create), GetHTML)
}

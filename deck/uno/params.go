package uno

import (
	"fmt"
	"reflect"

	"github.com/KirCute/sim-board"
)

const Name = "UNO"

type Params struct {
	CustomName             string `json:"custom_name" label:"自定义名称" type:"string"`
	Count                  int    `json:"count" label:"副数" type:"int" min:"0" default:"1"`
	CountColor             int    `json:"count_color" label:"颜色种类" type:"int" min:"0" max:"4" default:"4"`
	CountRank              int    `json:"count_rank" label:"数值范围" type:"int" min:"0" max:"10" default:"10"`
	CountColoredSkip       int    `json:"count_colored_skip" label:"每种颜色跳过卡数量" type:"int" min:"0" default:"1"`
	CountBlackSkip         int    `json:"count_black_skip" label:"无色跳过总数量" type:"int" min:"0" default:"0"`
	CountColoredReverse    int    `json:"count_colored_reverse" label:"每种颜色反转卡数量" type:"int" min:"0" default:"1"`
	CountBlackReverse      int    `json:"count_black_reverse" label:"无色反转卡总数量" type:"int" min:"0" default:"0"`
	CountTrans             int    `json:"count_trans" label:"万能卡总数量" type:"int" min:"0" default:"4"`
	CountColoredApp2       int    `json:"count_colored_append_2" label:"每种颜色+2卡数量" type:"int" min:"0" default:"1"`
	CountBlackApp2         int    `json:"count_black_append_2" label:"无色+2卡总数量" type:"int" min:"0" default:"0"`
	CountColoredApp4       int    `json:"count_colored_append_4" label:"每种颜色+4卡数量" type:"int" min:"0" default:"0"`
	CountBlackApp4         int    `json:"count_black_append_4" label:"无色+4卡总数量" type:"int" min:"0" default:"4"`
	CountColoredApp6       int    `json:"count_colored_append_6" label:"每种颜色+6卡数量" type:"int" min:"0" default:"0"`
	CountBlackApp6         int    `json:"count_black_append_6" label:"无色+6卡总数量" type:"int" min:"0" default:"0"`
	CountColoredApp8       int    `json:"count_colored_append_8" label:"每种颜色+8卡数量" type:"int" min:"0" default:"0"`
	CountBlackApp8         int    `json:"count_black_append_8" label:"无色+8卡总数量" type:"int" min:"0" default:"0"`
	CountColoredApp10      int    `json:"count_colored_append_10" label:"每种颜色+10卡数量" type:"int" min:"0" default:"0"`
	CountBlackApp10        int    `json:"count_black_append_10" label:"无色+10卡总数量" type:"int" min:"0" default:"0"`
	CountColoredSkipAll    int    `json:"count_colored_skip_all" label:"每种颜色全场跳过卡数量" type:"int" min:"0" default:"0"`
	CountBlackSkipAll      int    `json:"count_black_skip_all" label:"无色全场跳过卡总数量" type:"int" min:"0" default:"0"`
	CountColoredDiscardAll int    `json:"count_colored_discard_all" label:"每种颜色全弃卡数量" type:"int" min:"0" default:"0"`
	CountBlackDiscardAll   int    `json:"count_black_discard_all" label:"无色全弃卡总数量" type:"int" min:"0" default:"0"`
	CountColoredSwap       int    `json:"count_colored_swap" label:"每种颜色交换手牌卡数量" type:"int" min:"0" default:"0"`
	CountBlackSwap         int    `json:"count_black_swap" label:"无色交换手牌卡总数量" type:"int" min:"0" default:"0"`
	CountColoredBlank      int    `json:"count_colored_blank" label:"每种颜色空白卡数量" type:"int" min:"0" default:"0"`
	CountBlackBlank        int    `json:"count_black_blank" label:"无色空白卡总数量" type:"int" min:"0" default:"0"`
}

var COLORS = []string{"#ff5555", "#fcaa04", "#58a858", "#5555fc"}

func Create(params *Params) *Uno {
	params.CountColor = min(params.CountColor, len(COLORS))
	ret := &Uno{Params: params}
	for k := 0; k < params.Count; k++ {
		for i := 0; i < params.CountColor; i++ {
			for j := 0; j < params.CountRank; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-plain-%d", COLORS[i], 10-j)))
			}
			for j := 0; j < params.CountColoredApp2; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-plain-+2", COLORS[i])))
			}
			for j := 0; j < params.CountColoredApp4; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-plain-+4", COLORS[i])))
			}
			for j := 0; j < params.CountColoredApp6; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-plain-+6", COLORS[i])))
			}
			for j := 0; j < params.CountColoredApp8; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-plain-+8", COLORS[i])))
			}
			for j := 0; j < params.CountColoredApp10; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-plain-+10", COLORS[i])))
			}
			for j := 0; j < params.CountColoredSkip; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-skip", COLORS[i])))
			}
			for j := 0; j < params.CountColoredReverse; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-reverse", COLORS[i])))
			}
			for j := 0; j < params.CountColoredSkipAll; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-skipall", COLORS[i])))
			}
			for j := 0; j < params.CountColoredDiscardAll; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-discardall", COLORS[i])))
			}
			for j := 0; j < params.CountColoredSwap; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-swap", COLORS[i])))
			}
			for j := 0; j < params.CountColoredBlank; j++ {
				ret.rest = append(ret.rest, sim_board.Card(fmt.Sprintf("%s-blank", COLORS[i])))
			}
		}
	}
	for j := 0; j < params.CountBlackApp2; j++ {
		ret.rest = append(ret.rest, "black-plain-+2")
	}
	for j := 0; j < params.CountBlackApp4; j++ {
		ret.rest = append(ret.rest, "black-plain-+4")
	}
	for j := 0; j < params.CountBlackApp6; j++ {
		ret.rest = append(ret.rest, "black-plain-+6")
	}
	for j := 0; j < params.CountBlackApp8; j++ {
		ret.rest = append(ret.rest, "black-plain-+8")
	}
	for j := 0; j < params.CountBlackApp10; j++ {
		ret.rest = append(ret.rest, "black-plain-+10")
	}
	for j := 0; j < params.CountBlackSkip; j++ {
		ret.rest = append(ret.rest, "black-skip")
	}
	for j := 0; j < params.CountBlackReverse; j++ {
		ret.rest = append(ret.rest, "black-reverse")
	}
	for j := 0; j < params.CountBlackSkipAll; j++ {
		ret.rest = append(ret.rest, "black-skipall")
	}
	for j := 0; j < params.CountBlackDiscardAll; j++ {
		ret.rest = append(ret.rest, "black-discardall")
	}
	for j := 0; j < params.CountBlackSwap; j++ {
		ret.rest = append(ret.rest, "black-swap")
	}
	for j := 0; j < params.CountBlackBlank; j++ {
		ret.rest = append(ret.rest, "black-blank")
	}
	for j := 0; j < params.CountTrans; j++ {
		ret.rest = append(ret.rest, "black-trans")
	}
	return ret
}

func init() {
	sim_board.RegisterDeck(Name, reflect.ValueOf(Create), GetHTML)
}

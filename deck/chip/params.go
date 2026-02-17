package chip

import (
	"reflect"

	"github.com/KirCute/sim-board"
)

const Name = "筹码"

type Params struct {
	CustomName string `json:"custom_name" label:"自定义名称" type:"string"`
	Count1     int    `json:"count_1" label:"面值1数量" type:"int" min:"0" default:"0"`
	Count5     int    `json:"count_5" label:"面值5数量" type:"int" min:"0" default:"0"`
	Count20    int    `json:"count_20" label:"面值20数量" type:"int" min:"0" default:"0"`
	Count100   int    `json:"count_100" label:"面值100数量" type:"int" min:"0" default:"0"`
	Count500   int    `json:"count_500" label:"面值500数量" type:"int" min:"0" default:"0"`
	Count2k    int    `json:"count_2k" label:"面值2,000数量" type:"int" min:"0" default:"0"`
	Count1w    int    `json:"count_1w" label:"面值10,000数量" type:"int" min:"0" default:"0"`
}

func Create(params *Params) *Chip {
	ret := &Chip{Params: params}
	for i := 0; i < params.Count1; i++ {
		ret.Pool = append(ret.Pool, "1")
	}
	for i := 0; i < params.Count5; i++ {
		ret.Pool = append(ret.Pool, "5")
	}
	for i := 0; i < params.Count20; i++ {
		ret.Pool = append(ret.Pool, "20")
	}
	for i := 0; i < params.Count100; i++ {
		ret.Pool = append(ret.Pool, "100")
	}
	for i := 0; i < params.Count500; i++ {
		ret.Pool = append(ret.Pool, "500")
	}
	for i := 0; i < params.Count2k; i++ {
		ret.Pool = append(ret.Pool, "2000")
	}
	for i := 0; i < params.Count1w; i++ {
		ret.Pool = append(ret.Pool, "10000")
	}
	return ret
}

func init() {
	sim_board.RegisterDeck(Name, reflect.ValueOf(Create), GetHTML)
}

package dice

import (
	"reflect"

	"github.com/KirCute/sim-board"
)

const Name = "骰子"

type Params struct {
	CustomName string `json:"custom_name" label:"自定义名称" type:"string"`
	Face       int    `json:"face" label:"面数" type:"int" min:"1" default:"6"`
}

func Create(params *Params) *Dice {
	return &Dice{Params: params}
}

func init() {
	sim_board.RegisterDeck(Name, reflect.ValueOf(Create), GetHTML)
}

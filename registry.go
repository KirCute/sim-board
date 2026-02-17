package sim_board

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type deckRegistry struct {
	constructor reflect.Value
	paramSchema []map[string]any
	htmlGetter  func(card Card) (string, bool)
}

var decks map[string]*deckRegistry

func init() {
	decks = make(map[string]*deckRegistry)
}

func RegisterDeck(name string, constructor reflect.Value, htmlGetter func(card Card) (string, bool)) {
	ct := constructor.Type()
	if ct.Kind() != reflect.Func {
		panic("deck constructor must be a function")
	}
	if ct.NumIn() != 1 {
		panic("deck constructor must have exactly one input parameter")
	}
	if ct.NumOut() != 1 {
		panic("deck constructor must have exactly one output parameter")
	}
	deckType := ct.Out(0)
	if deckType.Kind() != reflect.Ptr {
		panic("deck constructor must return pointer to deck struct")
	}
	if !deckType.Implements(reflect.TypeOf((*Deck)(nil)).Elem()) {
		panic("deck constructor must return pointer to deck struct")
	}
	paramsType := ct.In(0)
	if paramsType.Kind() == reflect.Ptr {
		paramsType = paramsType.Elem()
	}
	if paramsType.Kind() != reflect.Struct {
		panic("input must be a struct or pointer to struct")
	}
	schema := make([]map[string]any, 0, paramsType.NumField())
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}
		key := strings.Split(jsonTag, ",")[0]
		desc := make(map[string]interface{})
		desc["key"] = key
		if label := field.Tag.Get("label"); label != "" {
			desc["label"] = label
		} else {
			desc["label"] = key
		}
		if typ := field.Tag.Get("type"); typ != "" {
			desc["type"] = typ
		} else {
			desc["type"] = "string"
		}
		if _min := field.Tag.Get("min"); _min != "" {
			if val, err := strconv.Atoi(_min); err == nil {
				desc["min"] = val
			}
		}
		if _max := field.Tag.Get("max"); _max != "" {
			if val, err := strconv.Atoi(_max); err == nil {
				desc["max"] = val
			}
		}
		if _default := field.Tag.Get("default"); _default != "" {
			switch desc["type"] {
			case "bool":
				desc["default"] = _default == "true"
			case "int":
				if val, err := strconv.Atoi(_default); err == nil {
					desc["default"] = val
				}
			default:
				desc["default"] = _default
			}
		}
		schema = append(schema, desc)
	}
	decks[name] = &deckRegistry{
		constructor: constructor,
		paramSchema: schema,
		htmlGetter:  htmlGetter,
	}
}

func NewDeck(name string, param json.RawMessage) (Deck, error) {
	reg, ok := decks[name]
	if !ok {
		return nil, fmt.Errorf("deck '%s' not found", name)
	}
	paramsType := reg.constructor.Type().In(0)
	var arg reflect.Value
	if paramsType.Kind() == reflect.Ptr {
		arg = reflect.New(paramsType.Elem())
	} else {
		arg = reflect.New(paramsType)
	}
	if err := json.Unmarshal(param, arg.Interface()); err != nil {
		return nil, fmt.Errorf("failed to unmarshal deck params: %+v", err)
	}
	var callArg reflect.Value
	if paramsType.Kind() == reflect.Ptr {
		callArg = arg
	} else {
		callArg = arg.Elem()
	}
	ret := reg.constructor.Call([]reflect.Value{callArg})[0].Interface().(Deck)
	return ret, nil
}

func GetCardHTML(deck, card string) (string, bool) {
	d, ok := decks[deck]
	if !ok {
		return "", false
	}
	return d.htmlGetter(Card(card))
}

func GetAllAvailableDecks() map[string][]map[string]any {
	ret := make(map[string][]map[string]any, len(decks))
	for deck, reg := range decks {
		ret[deck] = reg.paramSchema
	}
	return ret
}

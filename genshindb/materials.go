package genshindb

import (
	_ "embed"
	"encoding/json"
	"strings"
)

//go:embed materials.json
var materialsJSON []byte

var Materials MaterialMetaSet

func init() {
	_ = json.Unmarshal(materialsJSON, &Materials)
}

type MaterialMeta struct {
	Icon         string
	DropWeekdays []uint `json:",omitempty"`
}

type MaterialMetaSet map[string]MaterialMeta

func (s MaterialMetaSet) Cost(name string, count uint) (MaterialCostWithMeta, bool) {
	if m, ok := s.Get(name); ok {
		return MaterialCostWithMeta{
			MaterialCost: MaterialCost{
				Name:  name,
				Count: count,
			},
			MaterialMeta: m,
		}, true
	}

	return MaterialCostWithMeta{}, false
}

func (s MaterialMetaSet) Get(name string) (MaterialMeta, bool) {
	if mm, ok := s[name]; ok {
		return mm, ok
	}
	if mm, ok := s[strings.ReplaceAll(name, "之", "的")]; ok {
		return mm, ok
	}
	return MaterialMeta{}, false
}

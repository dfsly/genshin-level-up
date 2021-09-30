package genshindb

import (
	_ "embed"
	"encoding/json"
)

//go:embed reliquaries.json
var reliquariesJSON []byte

var Reliquaries ReliquaryFactory

func init() {
	_ = json.Unmarshal(reliquariesJSON, &Reliquaries)
}

type ReliquaryFactory struct {
	ReliquaryAddProps    [][]FightProps
	ReliquaryAffixDepots map[uint]map[string][]float64
	ReliquarySets        map[uint]ReliquarySet
	Reliquaries          []Reliquary
}

var maxes = map[uint]uint{
	5: 20,
	4: 16,
	3: 12,
	2: 8,
	1: 4,
}

func (f ReliquaryFactory) MainFightProps(rankLevel uint, level uint, fightPropTypes ...string) FightProps {
	if max := maxes[rankLevel]; level > max {
		level = max
	}

	props := f.ReliquaryAddProps[rankLevel-1][level]

	if len(fightPropTypes) > 0 {
		p := FightProps{}

		for _, k := range fightPropTypes {
			p[k] = props[k]
		}

		return p
	}

	return props
}

func (f ReliquaryFactory) Get(name string, rankLevel uint) *Reliquary {
	for i := range f.Reliquaries {
		if f.Reliquaries[i].Name == name && f.Reliquaries[i].RankLevel == rankLevel {
			r := f.Reliquaries[i]

			if set, ok := f.ReliquarySets[r.SetId]; ok {
				r.Set = &set
			}

			if appendPropTypes, ok := f.ReliquaryAffixDepots[r.AppendPropDepotID]; ok {
				r.AppendPropTypes = appendPropTypes
			}

			return &r
		}
	}
	return nil
}

type ReliquarySet struct {
	ActiveNum    uint `json:",omitempty"`
	Name         string
	EquipAffixes []struct {
		NeedNum   uint
		Desc      string
		AddProps  FightProps
		ParamList []float64
	}
}

func (s ReliquarySet) Active() bool {
	for i := range s.EquipAffixes {
		a := s.EquipAffixes[i]

		if s.ActiveNum >= a.NeedNum {
			return true
		}
	}

	return false
}

func (s *ReliquarySet) ReliquaryAffixes() []EquipAffix {
	if s != nil {
		affixes := make([]EquipAffix, 0)

		for i := range s.EquipAffixes {
			a := s.EquipAffixes[i]

			if s.ActiveNum >= a.NeedNum {
				affixes = append(affixes, EquipAffix{
					Name:      s.Name,
					Desc:      a.Desc,
					AddProps:  a.AddProps,
					ParamList: a.ParamList,
				})
			}
		}

		return affixes
	}

	return nil
}

type Reliquary struct {
	Name              string
	Desc              string
	EquipType         string
	RankLevel         uint
	MaxLevel          uint
	MainPropTypes     []string
	AddPropLevels     []uint
	AppendPropNum     uint
	AppendPropTypes   map[string][]float64
	SetId             uint
	AppendPropDepotID uint
	Set               *ReliquarySet `json:",omitempty"`
}

func (r Reliquary) AffixFightProps(level uint, fightPropTypes ...string) FightProps {
	if maxLevel := r.MaxLevel; level > maxLevel {
		level = maxLevel
	}

	c := 1

	for _, l := range r.AddPropLevels {
		if level >= l {
			c++
		}
	}

	p := FightProps{}

	for _, k := range fightPropTypes {
		appendPropTypes := r.AppendPropTypes[k]

		p[k] = appendPropTypes[len(appendPropTypes)-1] * float64(c)
	}

	return p
}

const (
	EQUIP_BRACER   = "EQUIP_BRACER"
	EQUIP_NECKLACE = "EQUIP_NECKLACE"
	EQUIP_SHOES    = "EQUIP_SHOES"
	EQUIP_RING     = "EQUIP_RING"
	EQUIP_DRESS    = "EQUIP_DRESS"
)

func posForEquipType(equipType string) int {
	switch equipType {
	case EQUIP_BRACER:
		return 1
	case EQUIP_NECKLACE:
		return 2
	case EQUIP_SHOES:
		return 3
	case EQUIP_RING:
		return 4
	case EQUIP_DRESS:
		return 5
	}
	return 1
}

func (r *Reliquary) LevelUpCost(current uint, to uint, propTypes ...string) LevelUpCost {
	exp := LevelupExpCosts.LevelUpReliquary(r.RankLevel, current, to)

	levelupCost := LevelUpCost{
		CurrentLevel: current,
		ToLevel:      to,
		Exp:          exp,
	}

	levelupCost.CurrentFightProps = Reliquaries.MainFightProps(r.RankLevel, current, propTypes...)
	levelupCost.ToFightProps = Reliquaries.MainFightProps(r.RankLevel, to, propTypes...)

	return levelupCost
}

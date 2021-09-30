package genshindb

import (
	_ "embed"
	"encoding/json"
)

//go:embed avatars.json
var avatarsJSON []byte

var Avatars AvatarList

func init() {
	_ = json.Unmarshal(avatarsJSON, &Avatars)
}

type AvatarList []Avatar

func (list AvatarList) Get(name string) *Avatar {
	for i := range list {
		if list[i].Name == name {
			return &list[i]
		}
	}
	return nil
}

type Avatar struct {
	Name        string
	RankLevel   uint
	ElementType string

	StaminaRecoverSpeed float64

	Critical         float64
	CriticalHurt     float64
	ChargeEfficiency float64

	Skills         []Skill
	InherentSkills []InherentSkill
	Talents        []Talent

	Promotes       []Promote
	PropGrowCurves map[string]PropGrowCurve
}

type Talent struct {
	Name     string
	Desc     string
	AddProps FightProps
}

type LevelUpCost struct {
	CurrentLevel      uint
	ToLevel           uint
	CurrentFightProps FightProps
	ToFightProps      FightProps
	Exp               uint
	PromoteCosts      []PromoteCost
}

func (a *Avatar) LevelUpCost(current uint, to uint) LevelUpCost {
	exp := LevelupExpCosts.LevelUpAvatar(current, to)

	levelupCost := LevelUpCost{
		CurrentLevel: current,
		ToLevel:      to,
		Exp:          exp,
	}

	levelupCost.CurrentFightProps = a.FightProps(current)
	levelupCost.ToFightProps = a.FightProps(to)

	for i := range a.Promotes {
		p := a.Promotes[i]

		if p.MinLevel >= current && p.MinLevel < to {
			levelupCost.PromoteCosts = append(levelupCost.PromoteCosts, p.PromoteCost())
		}
	}

	return levelupCost
}

func (a *Avatar) FightProps(level uint, addPropsList ...FightProps) FightProps {
	fightProps := FightProps{
		FIGHT_PROP_CRITICAL:          a.Critical,
		FIGHT_PROP_CRITICAL_HURT:     a.CriticalHurt,
		FIGHT_PROP_CHARGE_EFFICIENCY: a.ChargeEfficiency,
	}

	for k := range a.PropGrowCurves {
		fightProps.AddFightProps(FightProps{
			k: a.PropGrowCurves[k].Sum(level),
		})
	}

	fightProps.AddFightProps(a.Promotes[PromoteLevel(level)].AddProps)

	fightProps.AddFightProps(addPropsList...)

	return fightProps
}

type Skill struct {
	Name          string
	Desc          string
	ProudSkills   ProudSkills
	BreakLevels   []uint
	CoinCosts     []float64
	MaterialCosts [][]MaterialCost
}

type InherentSkill struct {
	Name       string
	Desc       string
	BreakLevel uint
	ParamNames []string
	Params     []float64
}

type ProudSkills struct {
	CdTime       uint   `json:",omitempty"`
	CostElemType string `json:",omitempty"`
	CostElemVal  uint   `json:",omitempty"`
	ParamNames   []string
	Params       [][]float64
}

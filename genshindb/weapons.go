package genshindb

import (
	_ "embed"
	"encoding/json"
)

//go:embed weapons.json
var weaponsJSON []byte

var Weapons WeaponList

func init() {
	_ = json.Unmarshal(weaponsJSON, &Weapons)
}

type WeaponList []Weapon

func (list WeaponList) Get(name string) *Weapon {
	for i := range list {
		if list[i].Name == name {
			return &list[i]
		}
	}
	return nil
}

type Weapon struct {
	Name      string
	Desc      string
	RankLevel uint

	Promotes       []Promote
	Affixes        [][]EquipAffix
	PropGrowCurves map[string]PropGrowCurve
}

func (w Weapon) WeaponAffixes(affixLevel uint) []EquipAffix {
	affixes := make([]EquipAffix, len(w.Affixes))

	for i, affixForLevels := range w.Affixes {
		affixes[i] = affixForLevels[affixLevel-1]
	}

	return affixes
}

type EquipAffix struct {
	Name      string
	Desc      string
	Level     uint
	AddProps  FightProps
	ParamList []float64
}

func (w *Weapon) LevelUpCost(current uint, to uint) LevelUpCost {
	exp := LevelupExpCosts.LevelUpWeapon(w.RankLevel, current, to)

	levelupCost := LevelUpCost{
		CurrentLevel: current,
		ToLevel:      to,
		Exp:          exp,
	}

	levelupCost.CurrentFightProps = w.FightProps(current, 0)
	levelupCost.ToFightProps = w.FightProps(to, 0)

	for i := range w.Promotes {
		p := w.Promotes[i]

		if p.MinLevel >= current && p.MinLevel < to {
			levelupCost.PromoteCosts = append(levelupCost.PromoteCosts, p.PromoteCost())
		}
	}

	return levelupCost
}

func (w *Weapon) FightProps(level uint, affixLevel uint, addPropsList ...FightProps) FightProps {
	fightProps := FightProps{}

	for k := range w.PropGrowCurves {
		fightProps.AddFightProps(FightProps{
			k: w.PropGrowCurves[k].Sum(level),
		})
	}

	if affixLevel > 0 {
		for _, affix := range w.WeaponAffixes(affixLevel) {
			fightProps.AddFightProps(affix.AddProps)
		}
	}

	promoteLevel := PromoteLevel(level)

	if maxPromoteLevel := uint(len(w.Promotes)) - 1; promoteLevel > maxPromoteLevel {
		promoteLevel = maxPromoteLevel
	}

	fightProps.AddFightProps(w.Promotes[promoteLevel].AddProps)

	fightProps.AddFightProps(addPropsList...)

	return fightProps
}

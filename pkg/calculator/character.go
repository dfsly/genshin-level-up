package calculator

import (
	"sort"
	"strings"

	"github.com/morlay/genshin-level-up/genshindb"
	"github.com/morlay/genshin-level-up/pkg/gameinfo"
)

func FromCharacters(characters []gameinfo.Character) *Result {
	ret := &Result{}

	for i := range characters {
		c := characters[i]

		ret.Characters = append(ret.Characters, CalcCharacter(c))
	}

	sort.Slice(ret.Characters, func(i, j int) bool {
		return weight(ret.Characters[i]) > weight(ret.Characters[j])
	})

	return ret
}

func weight(c Character) float64 {
	return float64(len(c.Reliquaries))*10.0 + float64(c.Level)/90.0*float64(c.RankLevel) + float64(c.Weapon.Level)/90.0*float64(c.Weapon.RankLevel)
}

func CalcCharacter(c gameinfo.Character) Character {
	cc := Character{
		CharacterIcons:  c.CharacterIcons,
		CharacterLevels: c.CharacterLevels,
		Weapon:          CalcWeapon(c.Weapon),
		Reliquaries:     make([]Reliquary, len(c.Reliquaries)),
	}

	cc.ElementIcon = genshindb.ElementIcons[strings.ToLower(cc.Element)]

	if avatar := genshindb.Avatars.Get(c.Name); avatar != nil {
		cc.Avatar = *avatar
	} else {
		return cc
	}

	reliquarySet := map[string]*genshindb.ReliquarySet{}
	reliquarySetNames := make([]string, 0)

	for i := range c.Reliquaries {
		rr := CalcReliquary(c.Name, c.Reliquaries[i])
		cc.Reliquaries[i] = rr

		if rr.Set != nil {
			if _, ok := reliquarySet[rr.Set.Name]; ok {
				reliquarySet[rr.Set.Name].ActiveNum++
			} else {
				set := *rr.Set
				reliquarySet[rr.Set.Name] = &set
				set.ActiveNum = 1
				reliquarySetNames = append(reliquarySetNames, set.Name)
			}
		}
	}

	for _, name := range reliquarySetNames {
		s := *reliquarySet[name]
		if s.Active() {
			cc.ReliquarySet = append(cc.ReliquarySet, s)
		}
	}

	// 满配圣遗物才可需要升级角色
	if len(cc.Reliquaries) == 5 {
		levels := []uint{61, 81, 86, 89}

		from := cc.Level

		for i := range levels {
			if cc.Level < levels[i] {
				cc.LevelupPlans = append(cc.LevelupPlans, genshindb.Avatars.Get(cc.Name).LevelUpCost(from, levels[i]))
				from = levels[i]
			}
		}
	}

	cc.FightProps = cc.Avatar.FightProps(cc.Level)

	cc.FightProps.AddFightProps(cc.Weapon.FightProps)

	for i := range cc.Reliquaries {
		cc.FightProps.AddFightProps(cc.Reliquaries[i].MainFightProps)
	}

	for i := range cc.ReliquarySet {
		if cc.ReliquarySet[i].Active() {
			for _, affix := range cc.ReliquarySet[i].ReliquaryAffixes() {
				cc.FightProps.AddFightProps(affix.AddProps)
			}
		}
	}

	cc.MaxFightProps = cc.FightProps.Clone()

	for i := range cc.Reliquaries {
		cc.MaxFightProps.AddFightProps(cc.Reliquaries[i].AffixFightProps)
	}

	return cc
}

type Character struct {
	genshindb.Avatar

	gameinfo.CharacterLevels
	gameinfo.CharacterIcons

	Weapon Weapon

	FightProps    genshindb.FightProps
	MaxFightProps genshindb.FightProps

	LevelupPlans []genshindb.LevelUpCost

	Reliquaries  []Reliquary
	ReliquarySet []genshindb.ReliquarySet
}

package calculator

import (
	"github.com/morlay/genshin-level-up/genshindb"
	"github.com/morlay/genshin-level-up/pkg/gameinfo"
)

var ReliquaryLevelupSet = map[uint][]uint{
	4: {8, 12, 16},
	5: {8, 12, 16, 20},
}

func CalcReliquary(name string, r gameinfo.Reliquary) Reliquary {
	rr := Reliquary{
		ReliquaryIcons:  r.ReliquaryIcons,
		ReliquaryLevels: r.ReliquaryLevels,
	}

	if r := genshindb.Reliquaries.Get(r.Name, r.Rarity); r != nil {
		rr.Reliquary = *r
	} else {
		return rr
	}

	mainPropType, affixPropTypes := genshindb.Recommends.For(name, rr.EquipType)

	rr.MainFightProps = genshindb.Reliquaries.MainFightProps(r.Rarity, r.Level, mainPropType)
	rr.AffixFightProps = rr.Reliquary.AffixFightProps(r.Level, affixPropTypes...)

	if levels, ok := ReliquaryLevelupSet[rr.RankLevel]; ok {
		from := rr.Level

		for i := range levels {
			if rr.Level < levels[i] {
				rr.LevelupPlans = append(rr.LevelupPlans, rr.LevelUpCost(from, levels[i], mainPropType))
				from = levels[i]
			}
		}
	}

	return rr
}

type Reliquary struct {
	genshindb.Reliquary

	gameinfo.ReliquaryIcons
	gameinfo.ReliquaryLevels

	MainFightProps  genshindb.FightProps
	AffixFightProps genshindb.FightProps

	LevelupPlans []genshindb.LevelUpCost
}

type Result struct {
	Characters []Character
}

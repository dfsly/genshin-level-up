package calculator

import (
	"github.com/morlay/genshin-level-up/genshindb"
	"github.com/morlay/genshin-level-up/pkg/gameinfo"
)

func normor() {

}

func CalcWeapon(w gameinfo.Weapon) Weapon {
	ww := Weapon{
		WeaponLevels: w.WeaponLevels,
		WeaponIcons:  w.WeaponIcons,
	}

	if found := genshindb.Weapons.Get(w.Name); found != nil {
		ww.Weapon = *found
		ww.FightProps = found.FightProps(w.Level, w.AffixLevel)
	}

	levels := []uint{51, 61, 71, 81, 90}

	from := ww.Level

	for i := range levels {
		if ww.Level < levels[i] {
			ww.LevelupPlans = append(ww.LevelupPlans, genshindb.Weapons.Get(ww.Name).LevelUpCost(from, levels[i]))
			from = levels[i]
		}
	}

	return ww
}

type Weapon struct {
	genshindb.Weapon

	gameinfo.WeaponIcons
	gameinfo.WeaponLevels

	FightProps   genshindb.FightProps
	LevelupPlans []genshindb.LevelUpCost
}

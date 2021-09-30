package genshindb

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	. "github.com/onsi/gomega"
)

func TestWeapons(t *testing.T) {
	a := Weapons.Get("天空之刃")

	t.Run("ToFightProps", func(t *testing.T) {
		hp := map[uint]float64{
			1:  46,
			20: 122,
			40: 235,
			50: 308,
			60: 382,
			70: 457,
			80: 532,
			90: 608,
		}

		for l, v := range hp {
			NewWithT(t).Expect(a.FightProps(l, 5).Round()["FIGHT_PROP_BASE_ATTACK"]).To(Equal(v))
		}
	})

	t.Run("lv.90", func(t *testing.T) {
		spew.Dump(a.FightProps(90, 5).Round())
	})

	t.Run("LevelUpCost", func(t *testing.T) {
		NewWithT(t).Expect(a.LevelUpCost(70, 81).PromoteCosts).To(HaveLen(2))
		NewWithT(t).Expect(a.LevelUpCost(80, 81).PromoteCosts).To(HaveLen(1))
	})
}

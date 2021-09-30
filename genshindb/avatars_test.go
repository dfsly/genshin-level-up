package genshindb

import (
	"math"
	"testing"

	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/gomega"
)

func TestAvatars(t *testing.T) {
	a := Avatars.Get("宵宫")

	t.Run("ToFightProps", func(t *testing.T) {
		hp := map[uint]float64{
			1:  791,
			20: 2053,
			40: 4086,
			50: 5256,
			60: 6593,
			70: 7777,
			80: 8968,
			90: 10164,
		}

		for l, v := range hp {
			NewWithT(t).Expect(math.Round(a.FightProps(l)["FIGHT_PROP_BASE_HP"])).To(Equal(v))
		}
	})

	t.Run("lv.86", func(t *testing.T) {
		spew.Dump(a.FightProps(86).Round())
	})

	t.Run("LevelUpCost", func(t *testing.T) {
		NewWithT(t).Expect(a.LevelUpCost(70, 81).PromoteCosts).To(HaveLen(2))
		NewWithT(t).Expect(a.LevelUpCost(80, 81).PromoteCosts).To(HaveLen(1))
	})
}

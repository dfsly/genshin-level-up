package genshindb

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	. "github.com/onsi/gomega"
)

func TestReliquaries(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		spew.Dump(Reliquaries.Get("祈望之心", 5))
	})

	t.Run("MainFightProps", func(t *testing.T) {
		NewWithT(t).Expect(
			Reliquaries.MainFightProps(5, 20).Round()[FIGHT_PROP_ATTACK],
		).To(Equal(float64(311)))

		NewWithT(t).Expect(
			Reliquaries.MainFightProps(5, 16).Round()[FIGHT_PROP_ATTACK],
		).To(Equal(float64(258)))

		NewWithT(t).Expect(
			Reliquaries.MainFightProps(4, 20).Round()[FIGHT_PROP_ATTACK],
		).To(Equal(float64(232)))
	})

}

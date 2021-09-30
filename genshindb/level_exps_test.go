package genshindb

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestLevelupExpCosts(t *testing.T) {
	t.Run("Meta", func(t *testing.T) {
		NewWithT(t).Expect(LevelupExpCosts.LevelUpAvatar(1, 2)).To(Equal(uint(1000)))
		NewWithT(t).Expect(LevelupExpCosts.LevelUpAvatar(1, 20)).To(Equal(uint(120175)))
		NewWithT(t).Expect(LevelupExpCosts.LevelUpAvatar(20, 40)).To(Equal(uint(578325)))
	})

	t.Run("Weapon", func(t *testing.T) {
		NewWithT(t).Expect(LevelupExpCosts.LevelUpWeapon(5, 1, 20)).To(Equal(uint(121550)))
		NewWithT(t).Expect(LevelupExpCosts.LevelUpWeapon(5, 20, 40)).To(Equal(uint(622800)))
	})

	t.Run("Reliquary", func(t *testing.T) {
		NewWithT(t).Expect(LevelupExpCosts.LevelUpReliquary(5, 0, 4)).To(Equal(uint(16300)))
		NewWithT(t).Expect(LevelupExpCosts.LevelUpReliquary(5, 4, 8)).To(Equal(uint(28425)))
		NewWithT(t).Expect(LevelupExpCosts.LevelUpReliquary(5, 8, 12)).To(Equal(uint(42425)))
		NewWithT(t).Expect(LevelupExpCosts.LevelUpReliquary(5, 12, 16)).To(Equal(uint(66150)))
		NewWithT(t).Expect(LevelupExpCosts.LevelUpReliquary(5, 16, 20)).To(Equal(uint(117175)))
	})

}

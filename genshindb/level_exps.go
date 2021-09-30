package genshindb

import (
	_ "embed"
	"encoding/json"
)

//go:embed level_exps.json
var levelExpsJSON []byte

var LevelupExpCosts LevelupExps

func init() {
	_ = json.Unmarshal(levelExpsJSON, &LevelupExpCosts)
}

type LevelupExps struct {
	Avatar    []uint
	Weapon    [5][]uint
	Reliquary [5][]uint
}

func cost(rng []uint, current uint, to uint) uint {
	if n := len(rng); int(to) > n {
		to = uint(n)
	}

	t := uint(0)

	for toLevel := current + 1; toLevel <= to; toLevel++ {
		t = t + rng[toLevel-1]
	}

	return t
}

func (c LevelupExps) LevelUpAvatar(current uint, to uint) uint {
	return cost(c.Avatar, current-1, to-1)
}

func (c LevelupExps) LevelUpWeapon(rankLevel uint, current uint, to uint) uint {
	return cost(c.Weapon[rankLevel-1], current-1, to-1)
}

func (c LevelupExps) LevelUpReliquary(rankLevel uint, current uint, to uint) uint {
	return cost(c.Reliquary[rankLevel-1], current, to)
}

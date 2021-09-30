package genshindb

func PromoteLevel(level uint) uint {
	l := uint(0)
	for _, maxLevel := range []uint{20, 40, 50, 60, 70, 80} {
		if level > maxLevel {
			l++
		}
	}
	return l
}

type Promote struct {
	MinLevel       uint
	UnlockMaxLevel uint
	CoinCost       uint
	AddProps       map[string]float64
	MaterialCosts  []MaterialCost
}

func (p Promote) PromoteCost() PromoteCost {
	m := PromoteCost{}
	m.CoinCost = p.CoinCost
	m.MaterialCosts = make([]MaterialCostWithMeta, len(p.MaterialCosts))

	for i := range m.MaterialCosts {
		mc := MaterialCostWithMeta{
			MaterialCost: p.MaterialCosts[i],
		}

		if mm, ok := Materials.Get(mc.Name); ok {
			mc.MaterialMeta = mm
		}

		m.MaterialCosts[i] = mc
	}

	return m
}

type PromoteCost struct {
	CoinCost      uint
	MaterialCosts []MaterialCostWithMeta
}

type MaterialCost struct {
	Name      string
	Count     uint
	RankLevel uint `json:",omitempty"`
}

type MaterialCostWithMeta struct {
	MaterialCost
	MaterialMeta
}

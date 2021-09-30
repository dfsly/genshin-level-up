package genshindb

type PropGrowCurve struct {
	Base   float64
	Values []float64
}

func (c PropGrowCurve) Sum(level uint) float64 {
	return c.Base * c.Values[level-1]
}

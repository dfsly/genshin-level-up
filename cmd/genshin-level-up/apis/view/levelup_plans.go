package view

import (
	"context"
	"fmt"

	"github.com/morlay/genshin-level-up/genshindb"
	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func LevelupPlans(levelupPlans []genshindb.LevelUpCost, exchangeExp func(exp uint) []genshindb.MaterialCostWithMeta) Component {
	return func(ctx context.Context) Element {
		if len(levelupPlans) == 0 {
			return nil
		}

		return Div(
			MapSlice(levelupPlans, func(i int) Element {
				if i >= 2 {
					return nil
				}

				lp := levelupPlans[i]

				delta := lp.ToFightProps.Clone()
				delta.Sub(lp.CurrentFightProps)

				return Div(
					CSS{
						"marginTop": "0.5em",
					},
					Panel(
						Div(
							CSS{
								"lineHeight": "1.2",
								"opacity":    fmt.Sprintf("%v", 0.85-0.2*float64(i)),
							},
							Div(
								CSS{
									"padding":      "0.2em 0",
									"marginBottom": "0.3em",
									"textAlign":    "right",
								},
								Div(
									Span(
										CSS{
											"opacity": 0.6,
										},
										Level(lp.CurrentLevel),
									),
									Text(" ~ "),
									Level(lp.ToLevel),
								),
							),
							FightDeltaProps(delta),
							MaterialCosts(func() []genshindb.MaterialCostWithMeta {
								if len(lp.PromoteCosts) > 0 {
									pc := lp.PromoteCosts[0]
									m2, _ := genshindb.Materials.Cost("摩拉", pc.CoinCost)
									return append(append(pc.MaterialCosts, m2), exchangeExp(lp.Exp)...)
								}
								return exchangeExp(lp.Exp)
							}()),
						),
					),
				)
			}),
		)
	}
}

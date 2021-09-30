package view

import (
	"bytes"
	"context"
	"fmt"

	"github.com/morlay/genshin-level-up/genshindb"
	"github.com/morlay/genshin-level-up/pkg/calculator"
	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func Character(l *calculator.Character) Component {
	return func(ctx context.Context) Element {
		return Div(
			Div(
				CSS{
					"display":  "flex",
					"flexWrap": "wrap",
					"margin":   "-0.3em",
					"& > *": CSS{
						"padding": "0.3em",
					},
				},
				Div(
					CSS{
						"width":    "158px",
						"fontSize": "12px",
						"& > * + *": CSS{
							"marginTop": "0.5em",
						},
					},
					Div(
						Avatar(AvatarProps{
							Affix:     uint(l.ActiveConstellationNum),
							RankLevel: l.RankLevel,
							Icon:      l.Icon,
							Name:      l.Name,
							Level:     l.Level,
							Children: Div(
								Div(
									CSS{
										"position": "absolute",
										"top":      "0.5em",
										"left":     "0.5em",
										"width":    "2em",
									},
									Image(l.ElementIcon),
								),
								Div(
									CSS{
										"padding": "0.3em 0.4em",
									},
									ComputedFightProps(l.FightProps, l.MaxFightProps),
								),
							),
						}),
					),
					LevelupPlans(l.LevelupPlans, func(exp uint) []genshindb.MaterialCostWithMeta {
						m, _ := genshindb.Materials.Cost("大英雄的经验", exp/10000)
						m2, _ := genshindb.Materials.Cost("摩拉", exp/5)
						return []genshindb.MaterialCostWithMeta{m, m2}
					}),
					Panel(
						CSS{
							"fontSize": "8px",
						},
						MapSlice(l.Talents, func(i int) Element {
							t := l.Talents[i]

							return Div(
								CSS{
									"opacity": func() float64 {
										if l.ActiveConstellationNum > i {
											return 1
										}
										return 0.6
									}(),
								},
								NameWithDesc(t.Name, t.Desc),
							)
						}),
					),
				),
				CharacterWeapon(l.Weapon),
				Div(
					CSS{
						"minWidth": fmt.Sprintf("%dpx", (78+2+2)*3+(4+2)*2),
						"maxWidth": fmt.Sprintf("%dpx", (78+2+2)*5+(4+2)*2),
						"& > * + *": CSS{
							"marginTop": "0.5em",
						},
					},
					CharacterReliquaries(l),
					Skills(l),
				),
			),
		)
	}
}

func CharacterWeapon(w calculator.Weapon) Component {
	return func(ctx context.Context) Element {
		affixes := w.WeaponAffixes(w.AffixLevel)

		return Div(
			CSS{
				"width":    "122px",
				"fontSize": "11px",
			},
			Avatar(AvatarProps{
				Affix:     w.AffixLevel,
				RankLevel: w.RankLevel,
				Level:     w.Level,
				Icon:      w.Icon,
				Name:      w.Name,
				Children: Div(
					Div(
						CSS{
							"padding": "0.3em 0.4em",
						},
						MapSlice(affixes, func(i int) Element {
							return Div(
								CSS{
									"fontSize":     "8px",
									"opacity":      0.8,
									"marginBottom": "0.5em",
								},
								Text(affixes[i].Desc),
							)
						}),
						FightDeltaProps(w.FightProps),
					),
				),
			}),
			LevelupPlans(w.LevelupPlans, func(exp uint) []genshindb.MaterialCostWithMeta {
				m, _ := genshindb.Materials.Cost("精锻用魔矿", exp/10000)
				m2, _ := genshindb.Materials.Cost("摩拉", exp)
				return []genshindb.MaterialCostWithMeta{m, m2}
			}),
		)
	}
}

func CharacterReliquaries(c *calculator.Character) Component {
	return func(ctx context.Context) Element {
		if len(c.Reliquaries) == 0 {
			return nil
		}

		return Panel(
			Div(
				CSS{
					"backgroundColor": "rgba(255,255,255,0.2)",
					"display":         "flex",
					"margin":          "0 -2px",
					"& > *": CSS{
						"flex":    "1",
						"padding": "2px",
					},
				},
				MapSlice(c.ReliquarySet, func(i int) Element {
					rs := c.ReliquarySet[i]
					affixes := rs.ReliquaryAffixes()

					return Div(
						NameWithDesc(
							fmt.Sprintf("%s %d 件套", rs.Name, rs.ActiveNum),
							func() string {
								b := bytes.NewBuffer(nil)
								for i := range affixes {
									if i > 0 {
										b.WriteString("\\n")
									}
									b.WriteString(affixes[i].Desc)
								}
								return b.String()
							}(),
						),
					)
				}),
			),
			Div(
				CSS{
					"display":  "flex",
					"flexWrap": "wrap",
					"margin":   "0 -2px",
					"& > *": CSS{
						"padding": "2px",
					},
				},
				MapSlice(c.Reliquaries, func(i int) Element {
					r := c.Reliquaries[i]

					return Div(
						CSS{
							"width": "82px",
						},
						Avatar(AvatarProps{
							RankLevel: r.RankLevel,
							Icon:      r.Icon,
							Name:      r.Name,
							Level:     r.Level,
							Children: Div(
								Div(
									CSS{
										"padding": "0.3em 0.4em",
									},
									FightDeltaProps(r.MainFightProps),
									Div(
										CSS{
											"opacity": "0.6",
										},
										FightDeltaProps(r.AffixFightProps),
									),
								),
							),
						}),
						LevelupPlans(r.LevelupPlans, func(exp uint) []genshindb.MaterialCostWithMeta {
							m, _ := genshindb.Materials.Cost("角色经验", exp)
							m2, _ := genshindb.Materials.Cost("摩拉", exp)
							m.Name = "经验"
							return []genshindb.MaterialCostWithMeta{m, m2}
						}),
					)
				}),
			),
		)
	}
}

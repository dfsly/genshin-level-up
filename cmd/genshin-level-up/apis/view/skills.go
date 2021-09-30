package view

import (
	"context"
	"fmt"
	"strings"

	"github.com/morlay/genshin-level-up/pkg/calculator"

	"github.com/morlay/genshin-level-up/genshindb"
	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func Skills(c *calculator.Character) Component {
	pl := genshindb.PromoteLevel(c.Level)
	skillLevels := make([]uint, 0, len(c.Skills[0].BreakLevels))

	for i, l := range c.Skills[0].BreakLevels {
		skillLevel := i + 1

		if pl == l && skillLevel <= 10 {
			skillLevels = append(skillLevels, uint(skillLevel))
		}
	}

	if len(skillLevels) == 0 {
		skillLevels = []uint{1}
	}

	return func(ctx context.Context) Element {
		return Panel(
			SkillCosts(c.Skills[0], skillLevels),
			SkillTables(c, pl, skillLevels),
		)
	}
}

func SkillTables(c *calculator.Character, promoteLevel uint, skillLevels []uint) Component {
	return func(ctx context.Context) Element {
		trs := make([]Element, 0)

		if len(skillLevels) > 0 && skillLevels[0] != 1 {
			skillLevels = append([]uint{skillLevels[0] - 1}, skillLevels...)
		}

		for i := range c.Skills {
			skill := c.Skills[i]

			trs = append(trs,
				Tr(Td(
					Attr("colspan", len(skillLevels)+1),
					NameWithDesc(skill.Name, skill.Desc),
				)),
			)

			for i := range skill.ProudSkills.ParamNames {
				paramName := skill.ProudSkills.ParamNames[i]
				parts := strings.Split(paramName, "|")

				tds := []Element{
					Td(
						CSS{
							"maxWidth": "5em",
						},
						Text(parts[0]),
					),
				}

				for _, skillLevel := range skillLevels {
					skillLevelIdx := skillLevel - 1
					if n := uint(len(skill.ProudSkills.Params)) - 1; skillLevelIdx > n {
						skillLevelIdx = n
					}

					params := genshindb.SliceFloatToSliceInterface(skill.ProudSkills.Params[skillLevelIdx])

					text := ""

					ret := genshindb.ParseParamName(paramName, params...)
					if len(ret) > 1 {
						text = ret[1]
					}

					tds = append(tds, Td(
						CSS{
							"textAlign":  "right",
							"fontWeight": "bold",
							"wordBreak":  "break-word",
						},
						Span(
							Text(text),
						),
						Div(
							CSS{
								"opacity": 0.6,
							},
							SkillNum(c, parts[0], parts[1], params...),
						),
					))
				}

				trs = append(trs, Tr(tds...))
			}

		}

		for i := range c.InherentSkills {
			inherentSkills := c.InherentSkills[i]

			activeCSS := CSS{
				"opacity": func() float64 {
					if promoteLevel >= inherentSkills.BreakLevel {
						return 1
					}
					return 0.6
				}(),
			}

			trs = append(trs,
				Tr(activeCSS, Td(
					Attr("colspan", len(skillLevels)+1),
					NameWithDesc(inherentSkills.Name, inherentSkills.Desc),
				)),
			)
		}

		return Table(
			CSS{
				"color":         "inherit",
				"width":         "100%",
				"fontSize":      "8px",
				"borderSpacing": 0,
				"& td": CSS{
					"padding":       "2px 0",
					"verticalAlign": "top",
				},
			},
			Thead(
				Tr(
					Td(),
					MapSlice(skillLevels, func(i int) Element {
						return Th(
							CSS{
								"textAlign": "right",
							},
							Level(skillLevels[i]),
						)
					}),
				),
			),
			Tbody(trs...),
		)
	}
}

func HurtValue(computed genshindb.FightProps, hurtAddType string, paramTemplate string, params ...interface{}) Component {
	base := computed.Get(genshindb.FIGHT_PROP_ATTACK)

	return Div(
		Span(Text(fmt.Sprintf(genshindb.Formatter(func(v interface{}, fnName string) string {
			hurt := v.(float64) * base * (1 + computed.Get(hurtAddType))
			return fmt.Sprintf("%.0f(%.0f)", hurt, hurt*(1+computed.Get(genshindb.FIGHT_PROP_CRITICAL_HURT)))
		}).Exec(paramTemplate, params...)))),
	)
}

func HealValue(computed genshindb.FightProps, baseOnPropType string, paramTemplate string, params ...interface{}) Component {
	return Div(
		CSS{
			"color": "green",
		},
		Span(Text(fmt.Sprintf(genshindb.Formatter(func(v interface{}, fnName string) string {
			if fnName == "I" {
				return fmt.Sprintf("%.0f", v.(float64)*(1+computed.Get(genshindb.FIGHT_PROP_HEAL_ADD)))
			}

			return fmt.Sprintf(
				"%.0f",
				v.(float64)*(computed.Get(baseOnPropType))*(1+computed.Get(genshindb.FIGHT_PROP_HEAL_ADD)),
			)
		}).Exec(paramTemplate, params...)))),
	)
}

func ShieldValue(computed genshindb.FightProps, baseOnPropType string, paramTemplate string, params ...interface{}) Component {
	return Div(
		CSS{
			"color": "blue",
		},
		Span(Text(fmt.Sprintf(genshindb.Formatter(func(v interface{}, fnName string) string {
			if fnName == "I" {
				return fmt.Sprintf("%.0f", v.(float64))
			}

			return fmt.Sprintf(
				"%.0f",
				v.(float64)*(computed.Get(baseOnPropType)),
			)
		}).Exec(paramTemplate, params...)))),
	)
}

func AttackAddValue(computed genshindb.FightProps, baseOnPropType string, paramTemplate string, params ...interface{}) Component {
	return Div(
		Span(Text(fmt.Sprintf(genshindb.Formatter(func(v interface{}, fnName string) string {
			return fmt.Sprintf(
				"+%.0f",
				v.(float64)*(computed.Get(baseOnPropType)),
			)
		}).Exec(paramTemplate, params...)))),
	)
}

func SkillNum(c *calculator.Character, label string, paramTemplate string, params ...interface{}) Component {
	return func(ctx context.Context) Element {
		if strings.Contains(paramTemplate, "普通攻击伤害") {
			return nil
		}

		computed := c.FightProps.Compute()

		if strings.HasSuffix(label, "攻击力加成比例") {
			return AttackAddValue(computed, genshindb.FIGHT_PROP_BASE_ATTACK, paramTemplate, params...)
		}

		if strings.HasSuffix(label, "吸收量") {
			if strings.Contains(paramTemplate, "生命值上限") {
				return ShieldValue(computed, genshindb.FIGHT_PROP_HP, strings.ReplaceAll(paramTemplate, "生命值上限", ""), params...)
			}
			if strings.Contains(paramTemplate, "最大生命值") {
				return ShieldValue(computed, genshindb.FIGHT_PROP_HP, strings.ReplaceAll(paramTemplate, "最大生命值", ""), params...)
			}
		}

		if strings.Contains(label, "治疗") {
			if strings.Contains(paramTemplate, "生命值上限") {
				return HealValue(computed, genshindb.FIGHT_PROP_HP, strings.ReplaceAll(paramTemplate, "生命值上限", ""), params...)
			}

			if strings.Contains(paramTemplate, "攻击力") {
				return HealValue(computed, genshindb.FIGHT_PROP_ATTACK, strings.ReplaceAll(paramTemplate, "攻击力", ""), params...)
			}
		}

		if strings.Contains(label, "伤害") || strings.Contains(label, "射击") {
			fightPropTypes := []string{
				genshindb.FIGHT_PROP_FIRE_ADD_HURT,
				genshindb.FIGHT_PROP_ELEC_ADD_HURT,
				genshindb.FIGHT_PROP_WATER_ADD_HURT,
				genshindb.FIGHT_PROP_ICE_ADD_HURT,
				genshindb.FIGHT_PROP_WIND_ADD_HURT,
				genshindb.FIGHT_PROP_ROCK_ADD_HURT,
				genshindb.FIGHT_PROP_GRASS_ADD_HURT,
			}

			return Div(
				func() Element {
					if strings.HasSuffix(label, "段伤害") || strings.HasSuffix(label, "重击伤害") || strings.HasSuffix(label, "下坠期间伤害") || strings.HasSuffix(label, "坠地冲击伤害") {
						return HurtValue(computed, genshindb.FIGHT_PROP_PHYSICAL_ADD_HURT, paramTemplate, params...)
					}
					return nil
				}(),
				MapSlice(fightPropTypes, func(i int) Element {
					elementHurtAdd := computed.Get(fightPropTypes[i])

					if elementHurtAdd == 0 && genshindb.FightPropElements[fightPropTypes[i]] != c.ElementType {
						return nil
					}

					return Div(
						CSS{
							"color": genshindb.FightPropColors[fightPropTypes[i]],
						},
						HurtValue(computed, fightPropTypes[i], paramTemplate, params...),
					)
				}),
			)
		}
		return nil
	}
}

func SkillCosts(skill genshindb.Skill, skillLevels []uint) Component {
	return func(ctx context.Context) Element {
		return Div(
			CSS{
				"padding": "0.5em 0",
			},
			MapSlice(skillLevels, func(j int) Element {
				skillLevel := skillLevels[j]
				costs := skill.MaterialCosts[skillLevel-1]

				if len(costs) == 0 {
					return nil
				}

				costWithMetas := make([]genshindb.MaterialCostWithMeta, len(costs))

				for i := range costWithMetas {
					m, _ := genshindb.Materials.Get(costs[i].Name)

					costWithMetas[i] = genshindb.MaterialCostWithMeta{
						MaterialCost: costs[i],
						MaterialMeta: m,
					}
				}

				mc, _ := genshindb.Materials.Cost("摩拉", uint(skill.CoinCosts[skillLevel-1]))

				costWithMetas = append(costWithMetas, mc)

				return Div(
					Div(
						CSS{
							"display": "flex",
						},
						Div(
							Text("技能突破花费"),
						),
						Div(
							CSS{
								"flex":      1,
								"textAlign": "right",
							},
							Span(
								CSS{
									"opacity": 0.6,
								},
								Level(uint(skillLevel-1)),
							),
							Text(" ~ "),
							Level(uint(skillLevel)),
						),
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
						MaterialCosts(costWithMetas),
					),
				)
			}),
		)
	}
}

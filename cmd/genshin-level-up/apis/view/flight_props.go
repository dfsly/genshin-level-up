package view

import (
	"context"
	"fmt"
	"strings"

	"github.com/morlay/genshin-level-up/genshindb"
	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func ComputedFightProps(fightProps genshindb.FightProps, maxFightProps genshindb.FightProps) Component {
	return func(ctx context.Context) Element {
		tds := Fragment{}

		computedFightProps := fightProps.Compute().Format()
		maxComputedFightProps := maxFightProps.Compute().Format()

		bases := map[string]string{}

		for _, key := range genshindb.FightPropKeys {
			if ret, ok := computedFightProps[key]; ok {
				if len(ret) >= 2 {
					if key == genshindb.FIGHT_PROP_BASE_HP || key == genshindb.FIGHT_PROP_BASE_ATTACK || key == genshindb.FIGHT_PROP_BASE_DEFENSE {
						bases[strings.ReplaceAll(key, "_BASE_", "_")] = ret[1]
						continue
					}

					tds = append(tds, Tr(
						Td(
							Text(func(v string) string {
								if _, ok := bases[key]; ok {
									return "(基础)" + v
								}
								return v
							}(ret[0])),
						),
						Td(
							CSS{
								"textAlign":  "right",
								"fontWeight": "400",
								"fontFamily": "monospace",
							},
							Span(
								CSS{
									"display": "block",
								},
								Text(func(v string) string {
									if base, ok := bases[key]; ok {
										return fmt.Sprintf("(%s)%s", base, v)
									}
									return v
								}(ret[1])),
							),
							func() Element {
								maxV := maxComputedFightProps[key][1]
								if maxV != ret[1] {
									return Span(
										CSS{
											"display": "block",
											"opacity": 0.6,
										},
										Text(fmt.Sprintf("≤ %s", maxV)),
									)
								}
								return nil
							}(),
						),
					))
				}

			}
		}

		return Table(
			CSS{
				"color":         "inherit",
				"width":         "100%",
				"fontSize":      "8px",
				"borderSpacing": 0,
				"& tr": CSS{
					"padding": "2px 0",
				},
				"& td": CSS{
					"verticalAlign": "top",
				},
			},
			Tbody(tds...),
		)
	}
}

func FightDeltaProps(fightProps genshindb.FightProps) Component {
	return func(ctx context.Context) Element {
		tds := Fragment{}

		for _, key := range genshindb.FightPropKeys {
			if v, ok := fightProps[key]; ok {
				if v == 0 {
					continue
				}

				ret := genshindb.ParseParamName(genshindb.FightPropParamNames[key], v)

				if len(ret) >= 2 {
					tds = append(tds, Tr(
						Td(
							Text(ret[0]),
						),
						Td(
							CSS{
								"textAlign":  "right",
								"fontFamily": "monospace",
								"fontWeight": "bold",
							},
							Text(func() string {
								if v > 0 {
									return "+" + ret[1]
								}
								return "-" + ret[1]
							}()),
						),
					))
				}

			}
		}

		return Table(
			CSS{
				"color":         "inherit",
				"width":         "100%",
				"fontSize":      "8px",
				"borderSpacing": 0,
			},
			Tbody(tds...),
		)
	}
}

package view

import (
	"context"
	"fmt"

	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

type AvatarProps struct {
	Name      string
	Icon      string
	RankLevel uint
	Level     uint
	Affix     uint
	Desc      Element
	Children  Element
}

func Avatar(l AvatarProps) Component {
	return func(ctx context.Context) Element {
		fromColor, toColor := RankColorsForRankLevel(l.RankLevel)

		return Div(
			CSS{
				"borderRadius":    "0.4em",
				"backgroundColor": "#ECE3D5",
				"border":          "1px solid rgba(0,0,0,0.05)",
				"position":        "relative",
			},
			Div(
				CSS{
					"position":                "relative",
					"margin":                  "-1px",
					"borderTopLeftRadius":     "0.4em",
					"borderTopRightRadius":    "0.4em",
					"borderBottomRightRadius": "30%",
					"overflow":                "hidden",
					"backgroundImage":         fmt.Sprintf("linear-gradient(0, %s, %s)", fromColor, toColor),
				},
				Image(l.Icon),
				Div(
					CSS{
						"position":     "absolute",
						"left":         0,
						"right":        0,
						"padding":      "0px 0.2em",
						"bottom":       "0.3em",
						"wordBreak":    "keep-all",
						"textOverflow": "ellipsis",
						"color":        "white",
						"textShadow":   fmt.Sprintf("0 0 2px %s,1px 1px 10px %s", toColor, toColor),
						"lineHeight":   "1.1",
						"textAlign":    "left",
					},
					l.Desc,
				),
				Component(func(ctx context.Context) Element {
					if l.Affix > 0 {
						return Div(
							CSS{
								"borderRadius":    "0.2em",
								"lineHeight":      "1.2",
								"padding":         "0 0.2em",
								"textAlign":       "center",
								"position":        "absolute",
								"fontWeight":      "bold",
								"fontFamily":      "monospace",
								"top":             "0.4em",
								"right":           "0.4em",
								"color":           "white",
								"backgroundColor": "rgba(0,0,0,0.2)",
							},
							Text(func() string {
								if l.Affix > 10000 {
									return fmt.Sprintf("%.0fw", float64(l.Affix)/10000)
								}
								return fmt.Sprintf("%d", l.Affix)
							}()),
						)
					}
					return nil
				}),
			),
			Div(
				CSS{
					"display":    "flex",
					"padding":    "0.3em 0.3em 0.2em",
					"alignItems": "flex-end",
					"lineHeight": "1.25",
				},
				Div(
					CSS{
						"flex":         1,
						"fontWeight":   "500",
						"wordBreak":    "keep-all",
						"textOverflow": "ellipsis",
						"overflow":     "hidden",
						"textAlign":    "left",
					},
					Attr("title", l.Name),
					Text(l.Name),
				),
				func() Element {
					if l.Level > 0 {
						return Div(
							CSS{
								"marginLeft": "0.5em",
								"fontSize":   "1.2em",
							},
							Level(l.Level),
						)
					}
					return nil
				}(),
			),
			Div(
				l.Children,
			),
		)
	}
}

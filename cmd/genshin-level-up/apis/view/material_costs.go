package view

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/morlay/genshin-level-up/genshindb"
	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func Level(l uint) Component {
	return func(ctx context.Context) Element {
		return Span(
			CSS{
				"fontWeight": "bold",
				"fontFamily": "monospace",
			},
			Span(
				Attr("style", CSS{
					"fontSize": "0.7em",
				}),
				Text("Lv."),
			),
			Text(fmt.Sprintf("%d", l)),
		)
	}
}

var (
	bgs = [][]string{
		{"transparent", "transparent"},
		{"transparent", "transparent"},
		{"#4E87A6", "#5A7E99"},
		{"#8164AA", "#6A5489"},
		{"rgb(198, 133, 77)", "rgb(167, 113, 67)"},
	}
)

func MaterialCosts(materials []genshindb.MaterialCostWithMeta) Component {
	return func(ctx context.Context) Element {
		if len(materials) != 0 {
			return Div(
				CSS{
					"fontSize": "8px",
				},
				Div(
					CSS{
						"display":  "flex",
						"flexWrap": "wrap",
						"margin":   "0 -2px",
						"& > *": CSS{
							"padding": "2px",
						},
					},
					MapSlice(materials, func(i int) Element {
						m := materials[i]

						w, active := toWeekdays(m.DropWeekdays, time.Now())

						return Div(
							Div(
								CSS{
									"width":    "32px",
									"fontSize": "8px",
									"position": "relative",
									"opacity": func() string {
										if active || len(m.DropWeekdays) == 0 {
											return "1"
										}
										return "0.6"
									}(),
								},
								Avatar(AvatarProps{
									Name:      m.Name,
									Icon:      m.Icon,
									RankLevel: m.RankLevel,
									Affix:     m.Count,
									Desc:      Text(w),
								}),
							),
						)
					}),
				),
			)
		}
		return nil
	}
}

var TZShanghai, _ = time.LoadLocation("Asia/Shanghai")

func toWeekdays(list []uint, t time.Time) (string, bool) {
	b := bytes.NewBuffer(nil)

	active := false

	for i := range list {
		w := list[i]
		if uint(t.In(TZShanghai).Weekday()) == w {
			active = true
		}
		if i == 0 {
			b.WriteString("周")
			continue
		}
		if i > 1 {
			b.WriteString("/")
		}
		b.WriteString(weekdays[w])
	}

	return b.String(), active
}

var weekdays = strings.Split("日一二三四五六", "")

func RankColorsForRankLevel(r uint) (string, string) {
	if r == 0 {
		r = 1
	}
	// 联动角色
	if r > 100 {
		r = r - 100
	}
	return bgs[r-1][0], bgs[r-1][1]
}

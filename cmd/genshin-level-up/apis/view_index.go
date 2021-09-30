package apis

import (
	"context"
	_ "time/tzdata"

	"github.com/morlay/genshin-level-up/cmd/genshin-level-up/apis/view"
	"github.com/morlay/genshin-level-up/pkg/calculator"
	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func ViewIndex(l *calculator.Result) Component {
	return func(ctx context.Context) Element {
		return view.Container(
			Div(
				CSS{
					"lineHeight": "1.4",
					"fontSize":   "10px",
					"& ul": CSS{
						"paddingLeft": "2em",
					},
				},
				Ul(
					Li(
						Span(Text("由于账号「角色天赋等级」和「圣遗物词条」无法获取")),
						Ul(
							Li(Text("「角色天赋等级」跟随「角色突破等级」")),
							Li(
								Span(Text("角色面板计算不包括「圣遗物副词条」，并按「圣遗物有效副词条」计算各自相关属性的理论最大值（不可能同时满足）")),
								Ul(
									Li(Text("「圣遗物主词条」和「圣遗物有效副词条」来自社区通用推荐 https://bbs.nga.cn/read.php?tid=27859119")),
									Li(Text("技能的伤害，盾量，生命，攻击加成等，也以该面板为基础计算")),
								),
							),
						),
					),
					Li(
						A(
							CSS{
								"color": "inherit",
							},
							Attr("href", "https://github.com/morlay/genshin-level-up"),
							Span(Text("源码 & 反馈 & 讨论")),
						),
					),
				),
			),
			Div(
				CSS{
					"display":  "flex",
					"flexWrap": "wrap",
					"& > *": CSS{
						"padding": "1em",
					},
				},
				MapSlice(l.Characters, func(i int) Element {
					return view.Character(&l.Characters[i])
				}),
			),
		)
	}
}

package view

import (
	"context"
	_ "embed"
	"sort"

	. "github.com/morlay/genshin-level-up/pkg/jsx"
	"github.com/morlay/genshin-level-up/pkg/jsx/css"
)

//go:embed normalize.css
var normalize string

func Container(children ...Element) Component {
	return func(ctx context.Context) Element {
		cache := css.NewCache()

		return Provide(func(ctx context.Context) context.Context {
			return css.ContextWithCache(ctx, cache)
		})(
			Head(
				Meta(
					Attributes{
						"name":    "viewport",
						"content": "width=device-width",
					},
				),
				Style(
					RawText(normalize),
				),
				Style(
					RawText("* { box-sizing: border-box }"),
				),
				Defer(func() Element {
					keys := make([]string, 0, len(cache.Registered))

					for k := range cache.Registered {
						keys = append(keys, k)
					}

					sort.Strings(keys)

					list := make(Fragment, len(keys))

					for i := range list {
						k := keys[i]
						s := cache.Registered[k]

						list[i] = Style(
							Attr("data-css", s.Name),
							RawText(s.ToCSS(cache.Key)),
						)
					}

					return list
				}),
			),
			Body(
				CSS{
					"fontFamily": "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif, 'Apple Color Emoji', 'Segoe UI Emoji', 'Segoe UI Symbol', 'Noto Color Emoji'",
				},
				Div(
					CSS{
						"margin":          "0 auto",
						"padding":         "0.5em",
						"overflowX":       "auto",
						"overflowY":       "auto",
						"fontSize":        "8px",
						"color":           "#986f51",
						"backgroundColor": "#FAF3EB",
					},
					Fragment(children),
				),
			),
		)
	}
}

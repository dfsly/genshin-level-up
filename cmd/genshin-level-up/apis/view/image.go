package view

import (
	"context"

	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func Image(src string) Component {
	return func(ctx context.Context) Element {
		return Img(
			CSS{
				"display":   "block",
				"margin":    "0",
				"width":     "100%",
				"minHeight": "50%",
			},
			Attr("src", src),
		)
	}
}

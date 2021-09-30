package view

import (
	"context"

	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func Panel(children ...Element) Component {
	return Div(
		append([]Element{
			CSS{
				"position":     "relative",
				"borderRadius": "4px",
				"padding":      "4px",
				"border":       "1px solid rgba(0,0,0,0.05)",
			},
		}, children...)...,
	)
}

func NameWithDesc(name string, desc string) Component {
	return func(ctx context.Context) Element {
		return Span(
			CSS{
				"padding": "0.5em 0",
				"display": "block",
			},
			Span(
				CSS{
					"fontSize":     "1.2em",
					"lineHeight":   1,
					"fontWeight":   "bold",
					"marginBottom": "0.5em",
					"display":      "block",
				},
				Text(name),
			),
			Span(
				CSS{
					"display":    "block",
					"lineHeight": 1.4,
					"padding":    "0 1em",
				},
				Text(desc),
			),
		)
	}
}

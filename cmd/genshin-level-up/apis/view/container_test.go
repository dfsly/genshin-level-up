package view

import (
	"context"
	"os"
	"testing"

	. "github.com/morlay/genshin-level-up/pkg/jsx"
)

func TestContainer(t *testing.T) {
	_ = Render(
		context.Background(),
		Container(
			Span(
				Text("test"),
			),
		),
		os.Stdout,
	)
}

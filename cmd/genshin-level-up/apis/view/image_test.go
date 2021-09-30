package view

import (
	"context"
	"os"
	"testing"

	"github.com/morlay/genshin-level-up/pkg/jsx"
)

func TestImage(i *testing.T) {
	_ = jsx.Render(context.Background(), Image("//xxx.com/x.png"), os.Stdout)
}

package jsx

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type RawText string

func (t RawText) RenderTo(ctx context.Context, w io.Writer) (err error) {
	_, err = fmt.Fprint(w, string(t))
	return
}

type Text string

func (t Text) RenderTo(ctx context.Context, w io.Writer) (err error) {
	s := string(t)
	s = strings.ReplaceAll(s, "\\n", "<br/>")
	s = strings.ReplaceAll(s, "\n", "<br/>")
	_, err = fmt.Fprint(w, s)
	return
}

type Fragment []Element

func (f Fragment) RenderTo(ctx context.Context, w io.Writer) error {
	for i := range f {
		if err := render(ctx, f[i], w); err != nil {
			return err
		}
	}
	return nil
}

func MapSlice(v interface{}, m func(i int) Element) Fragment {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Slice {
		panic("not a slice Value")
	}

	rv := reflect.ValueOf(v)

	children := make(Fragment, rv.Len())

	r := func(i int) Element {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("MapSlice, render %d failed: %v\n", i, r)
			}
		}()
		return m(i)
	}

	for i := range children {
		children[i] = r(i)
	}

	return children
}

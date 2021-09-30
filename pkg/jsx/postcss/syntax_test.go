package postcss

import (
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/matchers"
)

func TestNode(t *testing.T) {
	node := &Root{
		Nodes: []Node{
			&AtRule{
				Name:   "media",
				Params: "screen and (min-width: 480px)",
				Nodes: []Node{
					&Rule{
						Selector: "body",
						Nodes: []Node{
							&Declaration{Prop: "color", Value: "blue"},
						},
					},
				},
			},
			&Rule{
				Selector: "body",
				Nodes: []Node{
					&Declaration{Prop: "color", Value: "red"},
					&Rule{
						Selector: "__HOLDER__:hover",
						Nodes: []Node{
							&Declaration{Prop: "color", Value: "yellow"},
						},
					},
					&AtRule{
						Name:   "media",
						Params: "screen and (min-width: 480px)",
						Nodes: []Node{
							&Rule{
								Selector: "__HOLDER__",
								Nodes: []Node{
									&Declaration{Prop: "color", Value: "white"},
								},
							},
						},
					},
				},
			},
		},
	}

	NewWithT(t).Expect(node).To(BeCSS(`
@media screen and (min-width: 480px){body{color:blue;}}
body{color:red;}
body:hover{color:yellow;}
@media screen and (min-width: 480px){body{color:white;}}
`))
}

func BeCSS(s string) OmegaMatcher {
	return &BeCSSMatcher{
		EqualMatcher: matchers.EqualMatcher{
			Expected: strings.TrimSpace(s),
		},
	}
}

type BeCSSMatcher struct {
	EqualMatcher matchers.EqualMatcher
}

func (m *BeCSSMatcher) process(actual interface{}) string {
	b := &strings.Builder{}

	if node, ok := actual.(Node); ok {
		node.FormatTo(b, &FormatOpt{OneLine: true, Indent: "  "})
	}

	return strings.TrimSpace(b.String())
}

func (m *BeCSSMatcher) Match(actual interface{}) (success bool, err error) {
	return m.EqualMatcher.Match(m.process(actual))
}

func (m *BeCSSMatcher) FailureMessage(actual interface{}) (message string) {
	spew.Dump(actual)
	fmt.Println(m.process(actual))
	return m.EqualMatcher.FailureMessage(m.process(actual))
}

func (m *BeCSSMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return m.EqualMatcher.NegatedFailureMessage(m.process(actual))
}

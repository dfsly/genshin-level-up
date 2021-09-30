package genshindb

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRecommendSet_For(t *testing.T) {
	spew.Dump(Recommends.For("宵宫", EQUIP_BRACER))
	spew.Dump(Recommends.For("托马", EQUIP_DRESS))
}

func TestGenRecommends(t *testing.T) {
	group := strings.Split("火雷水冰风岩草", "")

	grouped := map[string][]string{}

	for _, a := range Avatars {
		if a.Name == "旅行者" {
			continue
		}
		grouped[a.ElementType] = append(grouped[a.ElementType], a.Name)
	}

	buf := bytes.NewBuffer(nil)

	buf.WriteString(`package genshindb

// https://bbs.nga.cn/read.php?tid=27859119
var Recommends = RecommendSet{
`)

	for _, e := range group {
		names := grouped[e]

		if len(names) == 0 {
			continue
		}

		sort.Strings(names)

		_, _ = fmt.Fprintf(buf, `// %s
`, e)

		for _, name := range names {
			_, _ = fmt.Fprintf(buf, `"%s": {
		%v
		%v
	},
`, name, func() string {
				if v, ok := Recommends[name]; ok {
					if len(v) > 2 {
						return strings.Join(v[0:3], ", ") + ","
					}
				}
				return ""
			}(),
				func() string {
					if v, ok := Recommends[name]; ok {
						if len(v) > 3 {
							return strings.Join(v[3:], ", ") + ","
						}
					}
					return ""
				}(),
			)
		}
	}

	buf.WriteString("}")

	_ = os.WriteFile("recommends_config.go", buf.Bytes(), os.ModePerm)
}

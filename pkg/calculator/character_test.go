package calculator

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/morlay/genshin-level-up/pkg/httputil"

	"github.com/morlay/genshin-level-up/pkg/gameinfo"
)

func TestGetUserInfo(t *testing.T) {
	data, _ := os.ReadFile("../../cookie")
	cookie, _ := base64.StdEncoding.DecodeString(string(data))

	c := &gameinfo.Client{
		Transports: []httputil.Transport{
			gameinfo.NewCommonTransport(string(cookie)),
		},
	}

	t.Run("Calc", func(t *testing.T) {
		characters, err := c.GetAllCharacters(context.Background(), 194435467)
		if err != nil {
			t.Fatal(err)
		}

		ret := FromCharacters(characters)

		for i := range ret.Characters {
			c := ret.Characters[i]

			fmt.Println(c.Name, c.Element, c.Weapon.Level, c.Weapon.PromoteLevel)
			//for i := range c.Reliquaries {
			//	r := c.Reliquaries[i]
			//	fmt.Println("\t", r.Name, r.ID)
			//}
		}
	})
}

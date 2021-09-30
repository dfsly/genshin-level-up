package gameinfo

import (
	"context"
	"encoding/base64"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGetUserInfo(t *testing.T) {
	data, _ := os.ReadFile("../../cookie")
	cookie, _ := base64.StdEncoding.DecodeString(string(data))

	c := NewClient(string(cookie))

	t.Run("GetAllCharacters", func(t *testing.T) {
		characters, err := c.GetAllCharacters(context.Background(), 194435467)
		if err != nil {
			t.Fatal(err)
		}
		spew.Dump(characters)
	})
}

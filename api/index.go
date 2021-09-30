package api

import (
	"net/http"

	"github.com/morlay/genshin-level-up/cmd/genshin-level-up/apis"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	apis.RootRouter.ServeHTTP(w, r)
}

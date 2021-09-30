package apis

import (
	"github.com/gorilla/mux"
)

var RootRouter = mux.NewRouter()

func init() {
	RootRouter.HandleFunc("/{uid}", Index)
	RootRouter.HandleFunc("/", Index)
}

package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/morlay/genshin-level-up/cmd/genshin-level-up/apis"
)

type Options struct {
	Cookie string
	Uid    int
}

func main() {
	opts := Options{}

	flag.StringVar(&opts.Cookie, "cookie", "", "encoded cookie string, could run `btoa(document.cookie)` on https://bbs.mihoyo.com")
	flag.IntVar(&opts.Uid, "uid", 0, "uid of game account")

	flag.Parse()

	srv := &http.Server{
		Addr:    ":8888",
		Handler: apis.RootRouter,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

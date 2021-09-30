package apis

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/morlay/genshin-level-up/pkg/httputil"

	"github.com/gorilla/mux"

	"github.com/morlay/genshin-level-up/pkg/calculator"

	"github.com/morlay/genshin-level-up/pkg/gameinfo"

	"github.com/morlay/genshin-level-up/pkg/jsx"
)

var SharedCookie = os.Getenv("COOKIE")

func Index(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p := mux.Vars(r)

	qCookie := q.Get("cookie")
	if qCookie == "" {
		qCookie = SharedCookie
	}

	qUid := q.Get("uid")
	if qUid == "" {
		qUid = p["uid"]
	}

	uid, _ := strconv.ParseInt(qUid, 10, 64)
	cookie, _ := base64.StdEncoding.DecodeString(qCookie)

	if uid == 0 {
		writeErr(rw, http.StatusInternalServerError, fmt.Errorf("missing `uid`"))
		return
	}

	if len(cookie) == 0 {
		writeErr(rw, http.StatusInternalServerError, fmt.Errorf("invalid `cookie`"))
		return
	}

	characters, err := gameinfo.NewClient(string(cookie)).GetAllCharacters(httputil.ContextWithRequest(r.Context(), r), int(uid))
	if err != nil {
		writeErr(rw, http.StatusInternalServerError, err)
		return
	}

	ret := calculator.FromCharacters(characters)

	jsx.RenderToResponse(r.Context(), rw, ViewIndex(ret))
}

func writeErr(rw http.ResponseWriter, code int, err error) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	_ = json.NewEncoder(rw).Encode(map[string]interface{}{
		"code": code,
		"msg":  err.Error(),
	})
}

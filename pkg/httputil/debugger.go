package httputil

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func NewDebuggerTransport() Transport {
	return func(roundTripper http.RoundTripper) http.RoundTripper {
		return &debuggerRoundTripper{
			next: roundTripper,
		}
	}
}

type debuggerRoundTripper struct {
	next http.RoundTripper
}

func (rt *debuggerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	d, _ := httputil.DumpRequestOut(req, true)
	fmt.Println(string(d))
	return rt.next.RoundTrip(req.WithContext(req.Context()))
}

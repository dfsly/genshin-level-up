package httputil

import (
	"context"
	"net"
	"net/http"
	"time"
)

type contextKeyRequest struct{}

func ContextWithRequest(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, contextKeyRequest{}, req)
}

func RequestFromContext(ctx context.Context) *http.Request {
	return ctx.Value(contextKeyRequest{}).(*http.Request)
}

func GetShortConnClientContext(ctx context.Context, timeout time.Duration, transports ...Transport) *http.Client {
	t := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 0,
		}).DialContext,
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: t,
	}

	for i := range transports {
		client.Transport = transports[i](client.Transport)
	}

	return client
}

type Transport = func(next http.RoundTripper) http.RoundTripper

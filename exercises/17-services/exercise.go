//go:build !solution

package service

import (
	"context"
	"net"
	"net/http"
)

// NewHandler returns the API's routes.
func NewHandler(quotes QuoteStore) http.Handler {
	mux := http.NewServeMux()
	// TODO: register the routes
	return mux
}

// Serve serves h on ln until ctx is cancelled, then shuts down gracefully.
func Serve(ctx context.Context, ln net.Listener, h http.Handler) error {
	// TODO: graceful shutdown
	return http.Serve(ln, h)
}

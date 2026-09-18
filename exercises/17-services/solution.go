//go:build solution

package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type api struct {
	quotes QuoteStore
	nextID atomic.Int64
}

// NewHandler returns the API's routes.
func NewHandler(quotes QuoteStore) http.Handler {
	a := &api{quotes: quotes}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /quotes/{symbol}", a.getQuote)
	mux.HandleFunc("POST /orders", a.postOrder)
	return mux
}

func (a *api) getQuote(w http.ResponseWriter, r *http.Request) {
	sym := r.PathValue("symbol")
	price, ok := a.quotes.Price(sym)
	if !ok {
		writeError(w, http.StatusNotFound, fmt.Sprintf("unknown symbol %q", sym))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"symbol": sym, "price": price})
}

func (a *api) postOrder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req OrderRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if req.Symbol == "" || req.Qty <= 0 || (req.Side != "buy" && req.Side != "sell") {
		writeError(w, http.StatusBadRequest, "symbol, side (buy|sell) and positive qty are required")
		return
	}
	id := a.nextID.Add(1)
	writeJSON(w, http.StatusCreated, map[string]string{"id": fmt.Sprintf("ord-%d", id)})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// Serve serves h on ln until ctx is cancelled, then shuts down gracefully.
func Serve(ctx context.Context, ln net.Listener, h http.Handler) error {
	srv := &http.Server{
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second, // guards against slowloris
	}
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()

	select {
	case err := <-errc:
		return err // failed to serve at all
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

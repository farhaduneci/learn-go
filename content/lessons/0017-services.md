---
title: Building services: net/http, JSON, graceful shutdown, gRPC
section: Part 3 · Senior Go
summary: ServeMux routing, defensive JSON, server/client timeouts, graceful shutdown, gRPC basics.
description: Go's standard library contains a production-grade HTTP server. This lesson builds a JSON API the way reviewers expect it (routing, validation, timeouts, graceful shutdown) and introduces gRPC, the usual protocol between internal services.
source: <a href="https://pkg.go.dev/net/http">net/http docs</a>, <a href="https://go.dev/blog/routing-enhancements">Routing Enhancements for Go 1.22</a>, <a href="https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/">How I write HTTP services in Go after 13 years</a> (Mat Ryer), and <a href="https://grpc.io/docs/languages/go/quickstart/">gRPC Go Quick start</a>.
quiz: Which request method returns a path wildcard like {symbol}? || PathValue
quiz: Which server method stops accepting connections and waits for in-flight requests? || Shutdown
quiz: Which json.Decoder method rejects unexpected JSON keys? || DisallowUnknownFields
quiz: Which http wrapper caps how many body bytes a handler reads? || MaxBytesReader
quiz: gRPC messages are defined in which IDL? || protobuf
---

## Routing with the standard library

Since Go 1.22, `http.ServeMux` supports methods and wildcards, so a third-party router is optional:

```go
mux := http.NewServeMux()
mux.HandleFunc("GET /healthz", health)
mux.HandleFunc("GET /quotes/{symbol}", a.getQuote)   // r.PathValue("symbol")
mux.HandleFunc("POST /orders", a.postOrder)
// wrong method → 405 automatically; unknown path → 404
```

A handler is anything with `ServeHTTP(w http.ResponseWriter, r *http.Request)`. `net/http` runs **each request in its own goroutine**, so handlers must be safe for concurrent use. Shared state needs a mutex or atomics (the order-ID counter in exercise 17 is an `atomic.Int64`).

## JSON in and out, defensively

```go
func (a *api) postOrder(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // never trust body size
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()                     // typos in field names become 400s
	var req OrderRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if err := req.Validate(); err != nil { ... }
	writeJSON(w, http.StatusCreated, resp)
}
```

Struct tags control encoding (`json:"symbol"`, `json:"price,omitempty"`). Set `Content-Type`, then `WriteHeader`, then write the body. Headers set after the first write are silently ignored. Map domain errors to status codes **in one place**, not in every handler.

## Timeouts are mandatory

The zero-value `http.Server` has **no timeouts**, so slow clients can hold connections forever. Also, `http.Get` uses `http.DefaultClient`, which has **no timeout**. Production settings:

```go
srv := &http.Server{
	Addr:              ":8080",
	Handler:           mux,
	ReadHeaderTimeout: 5 * time.Second,
	ReadTimeout:       10 * time.Second,
	WriteTimeout:      10 * time.Second,
	IdleTimeout:       60 * time.Second,
}
client := &http.Client{Timeout: 3 * time.Second} // reuse one client: it pools connections
```

Always `defer resp.Body.Close()` **and** read the body to completion, or connections can't be reused.

## Graceful shutdown

On Kubernetes, a pod gets SIGTERM and then a grace period (30s by default). The service should stop taking new requests, finish in-flight ones, flush buffers and close resources:

```go
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func Serve(ctx context.Context, ln net.Listener, h http.Handler) error {
	srv := &http.Server{Handler: h, ReadHeaderTimeout: 5 * time.Second}
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	sctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // fresh ctx: the old one is done
	defer cancel()
	if err := srv.Shutdown(sctx); err != nil {
		return err
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
```

The `main` → `run(ctx) error` split keeps `main` to three lines and makes the whole service testable.

## gRPC

Internal service-to-service traffic in trading systems usually runs over **gRPC**: HTTP/2, binary **protobuf** messages, generated clients, deadlines propagated across services, and **streaming**.

```protobuf
service OrderService {
  rpc PlaceOrder(PlaceOrderRequest) returns (PlaceOrderResponse);
  rpc StreamTrades(StreamTradesRequest) returns (stream Trade); // server streaming
}
```

```sh
protoc --go_out=. --go-grpc_out=. order.proto   # or `buf generate`
```

The generated server interface is implemented like any Go interface. Every RPC receives a `ctx` carrying the caller's deadline. Middleware becomes **interceptors**. Know the four RPC kinds (unary, server-streaming, client-streaming, bidirectional) and why protobuf beats JSON on the hot path: smaller, faster to parse, and schema-evolved through field numbers. Capstone stage 3 has a stretch goal to expose the engine over gRPC.

## Exercise

[`exercises/17-services`](../exercises/17-services/README.md): a standard-library JSON API with method routing, body limits, strict decoding, concurrency-safe IDs, and a graceful shutdown that the test proves lets an in-flight request finish.

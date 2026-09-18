# Exercise 17 — An HTTP service with graceful shutdown

Lesson: [17 · Building services: net/http, JSON, graceful shutdown, gRPC](../../lessons/0017-services.html)

```sh
go test -race ./exercises/17-services/
```

Standard library only. Go 1.22+ `ServeMux` supports methods and path wildcards (`"GET /quotes/{symbol}"`, `r.PathValue("symbol")`).

1. `NewHandler(quotes)` routes:
   - `GET /healthz` returns `200` with body `ok`.
   - `GET /quotes/{symbol}` returns `200` and `{"symbol":"AAPL","price":18950}`, or `404` and `{"error":"..."}`.
   - `POST /orders` decodes `{"symbol":"AAPL","side":"buy","qty":10}`. Reject unknown fields (`DisallowUnknownFields`), bodies over 1 MiB (`http.MaxBytesReader`), and invalid values with `400`. On success return `201` and `{"id":"ord-1"}` (IDs increment and must be safe under concurrent requests).
   - Any other method on a known path gets `405`, which `ServeMux` does for free.
2. `Serve(ctx, ln, h)` serves until `ctx` is cancelled, then calls `srv.Shutdown` with a 5s timeout, so **in-flight requests finish** before it returns. It returns `nil` after a clean shutdown.

Reference solution: `solution.go`.

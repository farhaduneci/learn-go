# Exercise 19 — Structured logs & latency histograms

Lesson: [19 · Observability: slog, metrics, tracing, profiling](../../lessons/0019-observability.html)

```sh
go test -race ./exercises/19-observability/
```

1. `NewLogger(w, level)` returns a JSON `*slog.Logger` that always carries `service=order-gateway`. Records below `level` are dropped.
2. `LogFill(ctx, logger, f)` logs `"order filled"` at Info, with the fill as a **group** (`"fill":{"symbol":..,"qty":..,"price":..}`) plus `request_id` read from the context when one is present. Use `logger.LogAttrs` with typed attrs (`slog.String`, `slog.Int64`), which avoids the `any` boxing that the `logger.Info("msg", k, v)` form pays for.
3. `Histogram` is a fixed-bucket latency histogram that is **safe for concurrent use without a mutex** (one `atomic.Int64` per bucket). `Quantile(0.99)` returns the upper bound of the bucket holding the 99th-percentile observation. Prometheus histograms work the same way, and latency SLOs are stated as p99, never as averages.
4. `Timed(h, next)` is middleware that records every call's duration into the histogram.

Reference solution: `solution.go`.

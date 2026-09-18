---
title: Observability: slog, metrics, tracing, profiling
section: Part 3 · Senior Go
summary: slog structured logs, Prometheus histograms, OpenTelemetry traces, continuous profiling.
description: A service nobody can see into can't be operated. This lesson covers the Go observability stack (structured logs with slog, Prometheus metrics, OpenTelemetry traces, continuous profiling) and what to measure in a trading system.
source: <a href="https://go.dev/blog/slog">Structured Logging with slog</a> (go.dev blog), <a href="https://opentelemetry.io/docs/languages/go/getting-started/">OpenTelemetry Go — Getting started</a>, <a href="https://prometheus.io/docs/guides/go-application/">Instrumenting a Go application for Prometheus</a>, and Google's <a href="https://sre.google/sre-book/monitoring-distributed-systems/">SRE book — Monitoring Distributed Systems</a> (the four golden signals).
quiz: Which slog method takes typed Attr values and avoids any-boxing? || LogAttrs
quiz: Which Prometheus metric type is right for request latency? || histogram
quiz: The four golden signals are latency, traffic, errors and what? || saturation
quiz: Logs, metrics, traces, and which fourth signal does continuous profiling add? || profiles
---

## Structured logging with log/slog

Since Go 1.21 the standard library has structured, levelled logging:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).
	With(slog.String("service", "order-gateway"))

logger.LogAttrs(ctx, slog.LevelInfo, "order filled",
	slog.Group("fill", slog.String("symbol", f.Symbol), slog.Int64("qty", f.Qty)),
	slog.String("request_id", id),
)
// {"time":"…","level":"INFO","msg":"order filled","service":"order-gateway","fill":{"symbol":"AAPL","qty":10},"request_id":"req-7"}
```

- Log **events with fields**, not formatted sentences. Fields are queryable, and `fmt.Sprintf` messages aren't.
- Use `LogAttrs` with typed attrs on hot paths. The `logger.Info("msg", "k", v)` form boxes every value into `any`.
- Pass request IDs and trace IDs via context, and use a custom `slog.Handler` that reads them from `ctx`, so no call site forgets them.
- Don't log secrets or PII, and don't log every order at Info on a high-volume path. Sample, or rely on metrics.
- Pass the logger into constructors. `slog.SetDefault` is fine in `main` but is global state.

## Metrics

Metrics are cheap, aggregated numbers, and alerting is built on them. With `prometheus/client_golang`:

```go
var orderLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "engine_submit_seconds",
	Buckets: prometheus.ExponentialBuckets(0.00001, 2, 16), // 10µs … ~0.3s
}, []string{"symbol", "result"})

orderLatency.WithLabelValues("AAPL", "filled").Observe(time.Since(start).Seconds())
```

- **Counter**: only goes up (orders total). **Gauge**: goes up and down (book depth, queue length). **Histogram**: distributions such as latency and sizes, which are what p99 is computed from.
- The **four golden signals**: latency, traffic, errors, saturation. For a matching engine, also track queue depth into the engine goroutine (saturation), event-publish lag, and `go_gc_duration_seconds` / `go_goroutines`.
- **Label cardinality**: never label with order IDs or user IDs. Each unique label set becomes a separate time series and can take down the metrics backend.

A fixed-bucket histogram with one atomic counter per bucket (exercise 19) is how these libraries work internally, and it's why recording a metric can be lock-free and nanoseconds-cheap.

## Distributed tracing

A trace follows one request across services as a tree of **spans**. OpenTelemetry is the standard:

```go
ctx, span := tracer.Start(ctx, "engine.Submit")
defer span.End()
span.SetAttributes(attribute.String("symbol", sym))
```

Instrumented HTTP and gRPC clients and servers (`otelhttp`, `otelgrpc`) **propagate the trace context** in headers, so spans across services join up. This is why `ctx` has to be passed everywhere: it carries the span.

## Profiling as the fourth signal

Continuous profiling (lesson 18) answers "why is this slower than last week?" down to a line of code. Together, logs, metrics, traces and profiles give coverage that no single one provides.

## Exercise

[`exercises/19-observability`](../exercises/19-observability/README.md): a JSON slog logger with a service tag, a grouped log record carrying a request ID from context, a lock-free latency histogram with quantiles, and timing middleware.

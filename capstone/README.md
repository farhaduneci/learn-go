# Capstone — a matching engine in Go

One project that grows across the track. By the end it covers most of what a core-trading Go role touches: in-memory data structures on the hot path, a single-writer concurrency model, a network API, an event stream, and measured performance.

Stages 1–2 come with tests (`capstone/orderbook/`). Stages 3–6 are specs with acceptance criteria. Build them in new packages under `capstone/`, write the tests first, and follow the same stub/solution convention if the work is meant to be shared.

```sh
go test -race ./capstone/...                           # stubs: failing until stages 1–2 are done
go test -race -tags solution ./capstone/...            # reference implementation
go test -run '^$' -bench . -benchmem ./capstone/...    # stage 5 baseline
```

## Stage 1 — The order book (after lessons 1–7)

`orderbook.Book`: a single-symbol limit order book.

- **Price-time priority**: best price first, then first-come-first-served within a price level.
- Trades execute at the **maker's** (resting order's) price.
- `Limit` orders rest their remainder; `Market` orders drop their remainder.
- `Cancel(id)` must be fast. Keep an index from order ID to its resting position.
- Money and prices are `int64` cents, never floats.

Design notes: a price level is a FIFO queue (`container/list`, or a slice-backed ring). Each side is a sorted slice of levels with the **best price at the end**, so taking liquidity pops the end in O(1). A red-black tree or skip list is the textbook alternative. The sorted slice wins on cache locality when the number of levels is small, and that's worth measuring in stage 5.

Tests: `book_test.go`.

## Stage 2 — The engine: one writer, many clients (after lessons 8–13)

`orderbook.Engine`: many goroutines submit orders concurrently, and **one goroutine owns all the books**. Books then need no locks, and every order gets a single, total sequence. This is the LMAX/"single-writer principle" design that real matching engines use.

- Commands travel over a channel. Each carries a reply channel, or, as in the reference solution, a closure run on the loop.
- Every public method honours `context.Context`.
- `Close()` is idempotent, and later calls return `ErrClosed` without blocking.

Tests: `engine_test.go`. They check conservation of quantity under 8 concurrent writers with `-race`, symbol isolation, and shutdown.

## Stage 3 — HTTP/JSON API (after lessons 14–17)

Package `capstone/api` plus `cmd/engine/main.go`.

- `POST /v1/orders/{symbol}` → `201 {"id":…,"trades":[…]}`; `DELETE /v1/orders/{symbol}/{id}`; `GET /v1/book/{symbol}?depth=10`.
- Validation errors map to `400`, unknown order to `404`, engine closed to `503`. The domain package defines errors, and the HTTP layer maps them in one place.
- Graceful shutdown on `SIGINT`/`SIGTERM` (`signal.NotifyContext`): stop accepting requests, drain in-flight ones, then `Engine.Close()`.
- Accepted when: `httptest`-based table tests cover every route and status; a test proves an in-flight request completes during shutdown.
- Stretch: the same API over **gRPC** (`protoc-gen-go-grpc`), including a server-streaming `StreamTrades` RPC.

## Stage 4 — Event stream (after lesson 20)

Publish every trade and book change to a broker. **NATS JetStream** (`nats-io/nats.go`) is the simplest to run locally with `docker run -p 4222:4222 nats -js`. Kafka (`twmb/franz-go`) is the alternative.

- Subject per symbol: `trades.AAPL`. Message ID = `<orderID>-<seq>`, so consumers can deduplicate.
- **Outbox rule**: publishing must never block matching. The engine writes events to a bounded channel, and a separate goroutine publishes. On overflow, pick a policy (block, drop, or halt) and document it. Trading systems usually prefer halting to silently losing fills.
- A consumer service builds a positions ledger from the trade stream using the idempotent consumer from exercise 20.
- Accepted when: an integration test (`//go:build integration`) starts NATS with `testcontainers-go` and asserts the consumer's positions match the engine's fills after a restart mid-stream.

## Stage 5 — Performance (after lesson 18)

- Baseline with `BenchmarkBookSubmit` and `BenchmarkEngineSubmitParallel`, and save it with `benchstat`.
- Profile (`-cpuprofile`, `-memprofile`) and cut allocations on the hot path. Targets to aim for on a laptop: `Book.Submit` **0–1 allocs/op**, under 150 ns/op. Ideas: pool `resting` structs, replace `container/list` with an intrusive or slice-backed queue, preallocate the trades slice or let callers pass one in (`SubmitInto(o, dst []Trade)`).
- Measure the **latency distribution** (p50/p99/p99.9) of `Engine.Submit` with the histogram from exercise 19, not just throughput. Tail latency is what trading clients see.
- Try `GOGC` and `GOMEMLIMIT` settings and record their effect on p99.

## Stage 6 — Observability (after lesson 19)

- `slog` JSON logs with a request ID carried from HTTP into the engine call.
- Prometheus metrics (`prometheus/client_golang`): orders by result, a match-latency histogram, and book depth as a gauge.
- OpenTelemetry traces across HTTP → engine → publish.
- `net/http/pprof` on a separate admin port, never the public one.

## Interview use

Be ready to explain in 5 minutes: why a single writer, what happens when the event publisher is slow, how cancel is O(1), and what the p99 is and why.

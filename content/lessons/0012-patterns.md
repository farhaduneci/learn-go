---
title: Concurrency patterns
section: Part 2 · Concurrency
summary: Worker pools, errgroup with SetLimit, semaphores, pipelines, rate limiting, choosing a pattern.
description: Most concurrent Go is built from a handful of patterns (worker pools, bounded parallelism, errgroup, pipelines, rate limiting). Knowing them lets each problem be recognised as one of them instead of designed from scratch.
source: <a href="https://pkg.go.dev/golang.org/x/sync/errgroup">errgroup docs</a>, <a href="https://go.dev/talks/2013/advconc.slide">Advanced Go Concurrency Patterns</a> (Sameer Ajmani), and the book <em>Concurrency in Go</em> (Katherine Cox-Buday, O'Reilly), chapters 3–5.
quiz: Which errgroup method caps the number of goroutines in flight? || SetLimit
quiz: What does errgroup.WithContext do to the context when one goroutine returns an error? || cancels it
quiz: A buffered channel of size N used to cap concurrency is acting as what? || semaphore
quiz: Which golang.org/x/time package implements a token-bucket limiter? || rate
---

## Worker pool

A fixed number of goroutines pull jobs from one channel. It bounds concurrency and memory, and keeps downstream systems from being overloaded.

```go
func WorkerPool[J, R any](ctx context.Context, jobs <-chan J, n int, fn func(J) R) <-chan R {
	out := make(chan R)
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}
					select {                       // the send must also be cancellable,
					case out <- fn(j):             // or a worker blocks forever when
					case <-ctx.Done():             // the consumer stops reading
						return
					}
				}
			}
		})
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
```

How many workers? For CPU-bound work, `runtime.GOMAXPROCS(0)`. For I/O-bound work, whatever the downstream (a DB pool, an API rate limit) can take.

## errgroup: bounded parallelism with first-error cancellation

`golang.org/x/sync/errgroup` is the standard answer to "run these N things concurrently, stop on the first failure, and wait for all of them":

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(8)                       // at most 8 in flight; g.Go blocks when full
results := make([]Quote, len(syms))
for i, s := range syms {
	g.Go(func() error {
		q, err := fetch(ctx, s)     // use the group's ctx: it is cancelled on first error
		if err != nil {
			return fmt.Errorf("fetch %s: %w", s, err)
		}
		results[i] = q
		return nil
	})
}
if err := g.Wait(); err != nil {    // first non-nil error
	return nil, err
}
```

In an interview, this is the code to write for "fetch prices for 1,000 symbols without overwhelming the API".

## Semaphores

A buffered channel **is** a counting semaphore:

```go
sem := make(chan struct{}, 10)
sem <- struct{}{}        // acquire (blocks when 10 are held)
defer func() { <-sem }() // release
```

`golang.org/x/sync/semaphore` adds weighted acquisition and context support.

## Pipelines

Stages connected by channels (lesson 9), where **each stage owns and closes its output**. Every stage should take a `ctx` or `done` channel so a downstream failure can stop upstream work. Otherwise upstream goroutines block forever on sends nobody reads.

For fan-out within a stage, run N copies reading the same input channel. Order is lost. Restore it by carrying a sequence number, or process per-key in order (lesson 20's partitioned dispatcher).

## Rate limiting

```go
lim := rate.NewLimiter(rate.Limit(50), 10) // golang.org/x/time/rate: 50/s, burst 10
if err := lim.Wait(ctx); err != nil { return err }
```

A `time.Ticker` gives simple smoothing, and a token bucket allows bursts (implemented by hand in exercise 21). Broker and exchange APIs enforce limits, and exceeding them gets a client throttled or banned.

## Picking a pattern

| Problem | Pattern |
| --- | --- |
| N independent calls, fail fast, bounded | `errgroup` + `SetLimit` |
| Continuous stream of jobs | worker pool over a channel |
| Transform stream in stages | pipeline with context |
| Protect a shared resource's concurrency | semaphore |
| Stay under an external rate | `x/time/rate` |
| Same key must be processed in order | partition by key hash to one goroutine each |
| First answer wins | buffered results + cancel (lesson 11) |

## Exercise

[`exercises/12-patterns`](../exercises/12-patterns/README.md): a cancellable worker pool (verified leak-free with `goleak`), `errgroup` with a limit and first-error cancellation, and a ticker-based throttle.

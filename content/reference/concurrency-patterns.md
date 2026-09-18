---
title: Concurrency patterns & rules
description: ownership rules, channel semantics table, and copy-ready snippets for every pattern in Part 2.
source: Compressed from lessons 8–13. Sources of truth: <a href="https://go.dev/ref/mem">The Go Memory Model</a> and <a href="https://go.dev/blog/pipelines">Pipelines and cancellation</a>.
---

## The rules

1. Every goroutine has a known **exit condition** and a known **waiter**.
2. The **sender closes**. With many senders, a coordinator closes after `wg.Wait()`.
3. Every blocking send or receive in a long-lived goroutine sits in a `select` with `<-ctx.Done()`.
4. A result channel that may be abandoned is **buffered to the number of senders**.
5. Shared mutable state: a mutex, **or** one owner goroutine, **or** per-goroutine slots. Never "probably fine".
6. `ctx` first, `defer cancel()` always, and never store a context in a struct.
7. `go test -race` in CI, plus `goleak` in concurrent packages.

## Channel semantics

| Operation | nil | open | closed |
| --- | --- | --- | --- |
| `ch <- v` | blocks forever | blocks until received or buffered | **panic** |
| `<-ch` | blocks forever | blocks until value | zero value, `ok == false` |
| `close(ch)` | **panic** | ok | **panic** |

## Wait for N goroutines

```go
var wg sync.WaitGroup
for i, x := range xs {
	wg.Go(func() { out[i] = f(x) }) // Go 1.25+
}
wg.Wait()
```

## Bounded, fail-fast parallelism

```go
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(n)
for i, x := range xs {
	g.Go(func() error { r, err := f(ctx, x); out[i] = r; return err })
}
err := g.Wait()
```

## Generator / pipeline stage

```go
func stage(ctx context.Context, in <-chan T) <-chan U {
	out := make(chan U)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case out <- transform(v):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}
```

## Fan-in

```go
func merge[T any](cs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Go(func() { for v := range c { out <- v } })
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
```

## Timeout without a leak

```go
done := make(chan T, 1) // buffered!
go func() { done <- slow() }()
select {
case v := <-done:
case <-time.After(d):
}
```

## First response wins

```go
ctx, cancel := context.WithCancel(ctx)
defer cancel()
res := make(chan result, len(backends))
for _, b := range backends {
	go func() { v, err := b(ctx); res <- result{v, err} }()
}
```

## Semaphore

```go
sem := make(chan struct{}, n)
sem <- struct{}{}
defer func() { <-sem }()
```

## Single-writer owner goroutine

```go
cmds := make(chan func(*State), 1024)
go func() {
	st := newState() // only this goroutine touches st: no locks
	for cmd := range cmds { cmd(st) }
}()
```

## Per-key ordering, cross-key parallelism

```go
idx := fnv32(key) % uint32(len(workers))
workers[idx] <- msg // each worker processes its channel sequentially
```

## Idempotent close / stop

```go
var once sync.Once
stop := func() { once.Do(func() { close(done) }) }
```

## Glossary anchors

| Term | Meaning |
| --- | --- |
| GMP | Goroutine, Machine (OS thread), Processor (run queue + right to run Go code) |
| happens-before | ordering guarantee from sync ops; without it, no visibility guarantee |
| data race | concurrent access, ≥1 write, no sync: undefined behaviour |
| goroutine leak | goroutine blocked forever; its memory is never freed |
| backpressure | slow consumer slows producer (bounded buffer) instead of unbounded growth |

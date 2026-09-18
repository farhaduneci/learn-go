# Exercise 12 — Concurrency patterns

Lesson: [12 · Concurrency patterns](../../lessons/0012-patterns.html)

```sh
go test -race ./exercises/12-patterns/
```

1. **Worker pool.** `WorkerPool(ctx, jobs, n, fn)` starts `n` workers that read from `jobs` and send results to the returned channel. The output closes once all workers finish. Workers stop early when `ctx` is cancelled, and no goroutine may be left blocked on a send.
2. **Bounded parallel with first-error cancel.** `ProcessAll(ctx, items, limit, fn)` runs `fn` on every item with **at most `limit` in flight**. Results come back in input order. The first error cancels the rest. Use `golang.org/x/sync/errgroup` with `SetLimit`.
3. **Rate limiter.** `Ticker`-based limiting: `Throttle(ctx, in, perSecond)` forwards values no faster than `perSecond`. This is how client code stays under a broker or exchange API rate limit.

Reference solution: `solution.go`.

# Exercise 10 — sync, atomics & the race detector

Lesson: [10 · Shared memory: sync, atomics & the memory model](../../lessons/0010-sync.html)

```sh
go test -race ./exercises/10-sync/
```

The starter code **works in a single goroutine and is broken under concurrency**. Run the tests without `-race` first, then with it, and read the race report: it prints both conflicting accesses with their stack traces.

1. `Positions` tracks net quantity per symbol, updated by many goroutines. Protect it with a `sync.Mutex`.
2. `Counter` counts events. Use `sync/atomic` (`atomic.Int64`) with no lock.
3. `QuoteCache` is read far more than written. Use `sync.RWMutex` and let readers proceed in parallel.
4. `Config` is loaded once, lazily, by whichever goroutine asks first. Use `sync.Once` (or `sync.OnceValue`) so the loader runs exactly once.

Reference solution: `solution.go`.

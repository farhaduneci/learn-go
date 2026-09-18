# Exercise 08 — Goroutines & WaitGroup

Lesson: [08 · Goroutines & the scheduler](../../lessons/0008-goroutines.html)

```sh
go test -race ./exercises/08-goroutines/
```

Always run concurrency exercises with `-race`.

1. `ParallelSum(nums, workers)` splits `nums` into `workers` parts, sums each part in its own goroutine, and combines the results. Give each goroutine its **own slot** in a results slice, so no locking is needed.
2. `FetchAll(symbols, fetch)` calls `fetch` for every symbol **concurrently** and returns the results **in input order**. The test fails if the calls run one after another.
3. `FetchAll` must not leak goroutines, and it must not return before every goroutine has finished (`sync.WaitGroup`, or `wg.Go` on Go ≥ 1.25).

Reference solution: `solution.go`.

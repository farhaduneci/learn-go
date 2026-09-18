# Exercise 11 — context: cancellation, deadlines, values

Lesson: [11 · context.Context](../../lessons/0011-context.html)

```sh
go test -race ./exercises/11-context/
```

1. `Sleep(ctx, d)` waits for `d` or returns `ctx.Err()` early if the context is cancelled. This is the basic building block for any context-aware wait.
2. `FirstSuccess(ctx, venues...)` queries every venue concurrently and returns the first successful price. It then **cancels the others** through a derived context. If all venues fail, it returns every error joined with `errors.Join`. Once the parent context is done, it returns `ctx.Err()`.
3. `WithRequestID` / `RequestID` store and read a request ID in the context. Use an **unexported key type** so no other package can collide with it.
4. `Retry(ctx, attempts, backoff, op)` retries `op` with a doubling backoff. It stops early when the context is cancelled.

Reference solution: `solution.go`.

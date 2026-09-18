# Exercise 09 — Channels & select

Lesson: [09 · Channels & select](../../lessons/0009-channels.html)

```sh
go test -race ./exercises/09-channels/
```

Rule for every function here: **the goroutine that sends on a channel is the one that closes it.**

1. `Gen(nums...)` returns a receive-only channel that emits `nums` and then closes.
2. `Square(in)` reads from `in` and emits squares. It closes its output when `in` closes.
3. `Merge(chans...)` fans many channels into one. It closes the output only after **all** inputs are drained (a WaitGroup plus a closer goroutine).
4. `WithTimeout(f, d)` runs `f` in a goroutine and returns `(result, true)`, or `(0, false)` if `d` passes first. Use a **buffered** channel of size 1, so the goroutine can still finish its send (and exit) after a timeout.
5. `Batch(in, size, wait)` groups items into slices of up to `size`. It emits a partial batch after `wait` if no more items arrive. Market-data and order-flush pipelines use this pattern: flush on size **or** time.

Tests fail with a timeout message instead of hanging when a channel is never closed.

Reference solution: `solution.go`.

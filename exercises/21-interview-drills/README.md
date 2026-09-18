# Exercise 21 — Interview drills

Lesson: [21 · Go interview drills](../../lessons/0021-interview-drills.html)

```sh
go test -race ./exercises/21-interview-drills/
```

Timebox each drill to **25 minutes**, the length of a live-coding round. Talk through the design out loud before typing, and state the complexity when finished.

1. **LRU cache** with O(1) `Get`/`Put` (`container/list` + map). Evict the least recently used item when full. Make it safe for concurrent use.
2. **Token-bucket rate limiter.** `NewLimiter(rate, burst, now)`. `Allow()` refills based on elapsed time, then spends one token. The clock is injected so tests are deterministic.
3. **Ping-pong.** Two goroutines take turns appending `"ping"` and `"pong"`, `n` rounds, using only unbuffered channels. The result is `ping pong ping pong ...`.
4. **Top-K** most frequent symbols with a size-k min-heap (`container/heap`), O(n log k). Break ties by symbol ascending.

Reference solution: `solution.go`.

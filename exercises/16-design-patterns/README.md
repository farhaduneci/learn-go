# Exercise 16 — Design patterns, the Go way

Lesson: [16 · Design patterns in Go](../../lessons/0016-design-patterns.html)

```sh
go test -race ./exercises/16-design-patterns/
```

1. **Functional options.** `NewClient(baseURL, opts...)` applies defaults (timeout 5s, 3 retries), then each `Option`. `WithTimeout` and `WithRetries` return an error for invalid values, and `NewClient` returns the first one.
2. **Middleware / decorator.** `Handler` is `func(context.Context, Order) error`. Write `Chain(h, mws...)` so `Chain(h, A, B)` runs A, then B, then h (the first middleware is the outermost). Then write two middlewares: `Validate` rejects `Qty <= 0` without calling the next handler, and `Recover` turns a panic in the next handler into an error.
3. **Pub/sub with backpressure.** `Bus.Subscribe(topic, buf)` returns a channel and an unsubscribe func. `Publish` must **never block** the publisher: if a subscriber's buffer is full, drop the message for that subscriber and count it in `Dropped()`. Market-data fan-out has to work this way; one slow consumer can't stall the feed.

Reference solution: `solution.go`.

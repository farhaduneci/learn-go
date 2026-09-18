---
title: Channels & select
section: Part 2 · Concurrency
summary: Unbuffered vs buffered, channel ownership, pipelines, fan-in, select, the timeout-leak trap.
description: Channels let goroutines communicate by passing values, and moving data through them transfers ownership. Getting the ownership rules right (who sends, who closes, who waits) is what separates working concurrent Go from code that deadlocks.
source: <a href="https://go.dev/blog/pipelines">Go Concurrency Patterns: Pipelines and cancellation</a> (go.dev blog, the essential read) and <a href="https://go.dev/doc/effective_go#channels">Effective Go — Channels</a>. Talk: <a href="https://www.youtube.com/watch?v=f6kdp27TYZs">Go Concurrency Patterns</a> (Rob Pike, 2012).
quiz: Who should close a channel: the sender or the receiver? || sender
quiz: What does a receive from a closed, drained channel return? || zero value
quiz: What happens on a send to a closed channel? || panic
quiz: What does a send or receive on a nil channel do? || blocks forever
quiz: Which select clause makes a channel operation non-blocking? || default
---

## Don't communicate by sharing memory; share memory by communicating

```go
ch := make(chan int)      // unbuffered
ch <- 42                  // send: blocks until a receiver takes it
v := <-ch                 // receive: blocks until a sender provides one
v, ok := <-ch             // ok == false once the channel is closed and drained
close(ch)                 // no more sends; receivers drain what is left, then get zero values
for v := range ch { }     // receive until closed
```

An **unbuffered** channel is a synchronisation point: the send completes only when a receiver takes the value, so both goroutines meet. A **buffered** channel `make(chan int, 64)` lets the sender run ahead until the buffer is full. The buffer is a **queue with backpressure**: when consumers fall behind, producers block. That's often exactly what's wanted.

## Direction types document ownership

```go
func Gen(nums ...int) <-chan int          // caller can only receive
func Square(in <-chan int) <-chan int      // reads in, owns and closes its output
func sink(out chan<- Trade)                // can only send
```

The rules that prevent most channel bugs:

1. **The sending side closes the channel**, and only when it's certain nothing else will send. Receivers never close.
2. With **several senders**, none of them can close safely. A coordinator closes after a `WaitGroup` says all senders are done (see `Merge` below).
3. Closing is a **broadcast**: every receiver sees it. That's why `close(done)` is the idiomatic stop signal.

| Operation | nil channel | open channel | closed channel |
| --- | --- | --- | --- |
| send | blocks forever | blocks until received or buffered | **panic** |
| receive | blocks forever | blocks until a value arrives | zero value, `ok=false` |
| close | **panic** | succeeds | **panic** |

## Pipelines and fan-in

```go
func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)       // this goroutine is the only sender, so it closes
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

func Merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range chans {
		wg.Go(func() {
			for v := range c { out <- v }
		})
	}
	go func() { wg.Wait(); close(out) }() // close once ALL senders are finished
	return out
}

for v := range Merge(Square(Gen(1, 2)), Square(Gen(3))) { ... }
```

## select

`select` waits on several channel operations and runs whichever is ready first (at random if several are):

```go
select {
case q := <-quotes:
	handle(q)
case <-time.After(50 * time.Millisecond):
	return ErrTimeout
case <-ctx.Done():
	return ctx.Err()
}
```

- A `default` case makes it **non-blocking**: `select { case ch <- v: default: dropped++ }` is how a publisher avoids stalling on a slow consumer (exercise 16).
- Setting a channel variable to `nil` **disables** its case, which is handy when one of two inputs has closed.
- In loops, use one `time.NewTimer` and `Reset` it instead of calling `time.After` on every iteration.

## The timeout-leak trap

```go
func WithTimeout(f func() int, d time.Duration) (int, bool) {
	done := make(chan int, 1)       // BUFFERED: if the timeout wins,
	go func() { done <- f() }()     // this send still succeeds and the goroutine exits
	select {
	case v := <-done: return v, true
	case <-time.After(d): return 0, false
	}
}
```

With an unbuffered channel, a timed-out call leaves its goroutine blocked on `done <- f()` forever. This exact leak comes up in code reviews and interviews.

## Channels or mutex?

Use **channels** to transfer ownership of data, distribute work, or signal events. Use a **mutex** to protect state that several goroutines read and update in place, like a cache or a counter. *"Channels orchestrate; mutexes serialise."* A channel used as a lock is slower and harder to read than `sync.Mutex`.

## Exercise

[`exercises/09-channels`](../exercises/09-channels/README.md): generator, pipeline stage, fan-in, a leak-free timeout, and a size-or-time batcher (the flush pattern used for market data and order writes).

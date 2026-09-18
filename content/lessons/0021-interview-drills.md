---
title: Go interview drills
section: Part 4 · Interview
summary: Rapid-fire senior Go questions with short answers, plus four timed live-coding drills.
description: The questions and live-coding tasks that come up in senior Go interviews, grouped by topic, each with the short answer a strong candidate gives. Use this lesson for self-testing and come back to it before every interview.
source: <a href="https://go.dev/doc/faq">The Go FAQ</a> (answers many "why is Go like this" questions directly), <a href="https://go-proverbs.github.io/">Go Proverbs</a>, and <a href="https://go.dev/wiki/CodeReviewComments">Code Review Comments</a>. Practice live coding at <a href="https://exercism.org/tracks/go">Exercism Go track</a>.
quiz: Which standard package provides a doubly linked list for an O(1) LRU? || container/list
quiz: A token bucket has two parameters: rate and what? || burst
quiz: Which standard package provides heap operations for a priority queue? || container/heap
quiz: An interface value is nil only when both its value and its what are nil? || type
---

## How to use this lesson

For each question, answer **out loud in under a minute** before reading the answer. If the answer drifts, redo the linked lesson. Then do the four timed drills in the exercise.

## Language

**Value vs pointer receivers?** Pointer receivers when the method mutates, the struct is large, or it holds a mutex. Stay consistent across a type's methods. Value receivers suit small immutable types. (Lesson 2)

**What's in a slice, and what happens on append?** A pointer, a length and a capacity. Append writes into spare capacity or reallocates about 2× and copies, so two slices can alias until a reallocation. (Lesson 3)

**Why is `err != nil` true when the function returned a nil `*MyError`?** An interface is nil only if both its type and value are nil. Returning a typed nil pointer gives a non-nil interface. Return a literal `nil`. (Lesson 5)

**Map iteration order? Concurrent map writes?** Randomised on purpose. Concurrent writes crash with a fatal error. Use a mutex (or `sync.Map` for its specific cases). (Lesson 3)

**Embedding vs inheritance?** Embedding promotes methods, but there's no virtual dispatch. The embedded type never calls the outer type's methods. (Lesson 4)

**When to use generics vs interfaces?** Generics for the same logic over different types (containers, algorithms). Interfaces for different behaviour behind one API. (Lesson 6)

**What does `defer` evaluate immediately?** Its arguments. The call runs at function return, in LIFO order, and can change named results. (Lesson 5)

## Concurrency

**Goroutine vs thread?** ~2 KB growable stack, scheduled by the runtime onto GOMAXPROCS threads (GMP) with work stealing, cheap user-space switches, and the netpoller for I/O. (Lesson 8)

**Buffered vs unbuffered channel?** Unbuffered synchronises sender and receiver. Buffered decouples them up to the capacity, which gives backpressure once full. (Lesson 9)

**Who closes a channel? What if there are many senders?** The sender does. With many senders, a coordinator closes after a WaitGroup. Receivers never close. (Lesson 9)

**How do goroutines leak, and how are leaks found?** Blocked forever on a send or receive (abandoned receivers, unclosed channels, missing ctx cases). Find them with goleak in tests and goroutine profiles and metrics in production. (Lesson 13)

**Mutex vs channel?** Mutex to protect shared state. Channels to transfer ownership or signal. (Lessons 9–10)

**What is the race detector, and what can't it do?** Runtime instrumentation that reports conflicting unsynchronised accesses, but only on code paths that actually ran. It costs ~5–10× overhead, so run it in CI. (Lesson 10)

**Implement "fetch 1,000 URLs, at most 10 at a time, stop on the first error".** `errgroup.WithContext` + `SetLimit(10)` + per-index result slots. (Lesson 12)

**How does context cancellation propagate?** Down a tree of derived contexts. Code has to check `ctx.Done()` or `ctx.Err()`, and `cancel` must always be called. (Lesson 11)

## Production and design

**What does graceful shutdown involve?** `signal.NotifyContext`, then `srv.Shutdown` with a timeout, then drain workers and consumers, flush, and close. (Lesson 17)

**A Go service's p99 latency doubled. What next?** Check metrics first (GC CPU, goroutines, lock wait, queue depth), then the CPU and alloc profiles and the execution trace. Look for new allocations, lock contention or GC assists. Tune GOGC/GOMEMLIMIT only after cutting allocations. (Lesson 18)

**What makes a Kafka/NATS consumer correct?** At-least-once delivery plus idempotent processing (dedupe in the same transaction), a DLQ for poison messages, per-key ordering through partitioning, and acking after processing. (Lesson 20)

**Design a matching engine's concurrency.** A single writer per book (or per symbol shard) behind a channel, so the book itself needs no locks and gets a total order. Events leave through a bounded outbox. Explain the backpressure policy and the p99 measurements. (Capstone)

**Where do interfaces belong?** In the consuming package, kept small. Accept interfaces, return structs. (Lesson 15)

## Live-coding etiquette

1. Restate the problem and ask about edge cases (empty input, concurrency, sizes).
2. Say the approach and its complexity **before** typing.
3. Write the simplest correct version first, then improve it.
4. Test with a few cases out loud, including an edge case.
5. For concurrent code, say who closes each channel and how every goroutine ends.

## Exercise

[`exercises/21-interview-drills`](../exercises/21-interview-drills/README.md): four 25-minute drills, done in order: an O(1) thread-safe LRU, a deterministic token-bucket limiter, a two-goroutine ping-pong, and heap-based top-K.

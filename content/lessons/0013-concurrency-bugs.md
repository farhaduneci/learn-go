---
title: Concurrency bugs: leaks, deadlocks, races
section: Part 2 · Concurrency
summary: Goroutine leaks, deadlocks, data races, misuse panics, and the tools that expose each.
description: Senior engineers are hired partly to find bugs like these before production does. This lesson catalogues the failure modes behind nearly every Go concurrency incident and the tools that expose each one.
source: <a href="https://songlh.github.io/paper/go-study.pdf">Understanding Real-World Concurrency Bugs in Go</a> (ASPLOS 2019; read sections 1, 5 and 6) and <a href="https://github.com/uber-go/goleak">uber-go/goleak</a>. Also <a href="https://www.uber.com/blog/data-race-patterns-in-go/">Data Race Patterns in Go</a> (Uber engineering).
quiz: When every goroutine is blocked, the runtime's fatal error reads: all goroutines are asleep - what? || deadlock!
quiz: Which signal makes a Go program dump all goroutine stacks? || SIGQUIT
quiz: Which Uber library fails a test if goroutines are left running? || goleak
quiz: The classic fix for two locks acquired in opposite orders is a global lock what? || ordering
---

## 1. Goroutine leaks

A goroutine blocked forever on a channel operation is never collected, and neither is everything it references. Leaks grow slowly until the pod gets OOM-killed. The research study above found that **blocking bugs are more common than races** in real Go code, and most of them involve channels.

Typical causes:

- A sender whose receiver gave up (the unbuffered-timeout bug in lesson 9, or first-result-wins with an unbuffered channel).
- A receiver waiting on a channel nobody will ever close.
- A worker without a `ctx.Done()` case, after the consumer stopped reading.

Detection:

```go
func TestMain(m *testing.M) { goleak.VerifyTestMain(m) } // whole package
defer goleak.VerifyNone(t)                              // one test
```

In production, watch the `go_goroutines` metric, which should be flat under steady load, and pull `/debug/pprof/goroutine?debug=1` to see where goroutines are parked.

## 2. Deadlocks

- **Lock ordering:** goroutine A locks X then Y, while B locks Y then X. Fix it with one global order, for example by account ID (exercise 13).
- **Re-locking:** Go mutexes aren't reentrant. A method holding `mu` that calls another method which also locks `mu` hangs.
- **Lock held across a channel send or a callback:** the receiver needs the same lock, so both wait on each other.
- **Unbuffered send with no receiver** in the same goroutine: `ch <- v; <-ch` blocks forever.

If **every** goroutine is blocked, the runtime detects it: `fatal error: all goroutines are asleep - deadlock!`. Partial deadlocks, where some goroutines still run, aren't detected, so the process just stops making progress. Send `SIGQUIT` (Ctrl-\\) or hit `/debug/pprof/goroutine?debug=2` and read the stacks: two goroutines each in `sync.(*Mutex).Lock`, waiting on each other's lock, is the signature.

## 3. Data races

- Shared variables written from goroutines (`err = ...` inside `wg.Go`). Use a mutex, per-goroutine slots, or `errgroup`.
- Maps written concurrently. These crash outright.
- Slices appended from several goroutines. `append` isn't atomic.
- Reading a struct while another goroutine updates it field by field.
- Capturing a loop variable before Go 1.22, still common in older codebases.

`go test -race` in CI is non-negotiable. Stress tests that run the concurrent path many times (`-count=100`) make rare interleavings show up.

## 4. Misuse panics

- **Double close** of a channel: use `sync.Once`, or restructure so exactly one owner closes.
- **Send on a closed channel:** a sender outlived its owner's shutdown. Fix the ownership so the channel's owner decides when sends stop.
- **`sync.WaitGroup` negative counter:** `Done` called more times than `Add`.
- **Copying a mutex or WaitGroup** by value: `go vet` catches it.

## A debugging checklist

1. Can it be reproduced under `-race`, with `-count=N` and `GOMAXPROCS=1` or `=many`?
2. Take a goroutine dump. Which goroutines exist, and what is each one blocked on?
3. For every channel involved: who sends, who closes, and who stops reading first?
4. For every lock: which locks does this goroutine already hold here?
5. Write the failing test **first**, then fix, then keep the test.

## Exercise

[`exercises/13-concurrency-bugs`](../exercises/13-concurrency-bugs/README.md): four functions, four classic bugs (leak, lock-order deadlock, data race, double close). Diagnose each from its tool output, fix it minimally, and name it in a comment.

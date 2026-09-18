---
title: Goroutines & the scheduler
section: Part 2 · Concurrency
summary: Goroutines vs threads and asyncio, the GMP scheduler, WaitGroup, rules for every goroutine.
description: Goroutines are the reason teams choose Go for high-throughput backends. This lesson gives the mental model (cheap tasks, M:N scheduling, real parallelism, no GIL) and the discipline around starting and stopping them.
source: <a href="https://go.dev/tour/concurrency/1">A Tour of Go — Concurrency</a>, then Rob Pike's talk <a href="https://go.dev/blog/waza-talk">Concurrency is not Parallelism</a> (video, 30 min). For the scheduler: <a href="https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html">Scheduling In Go, Part II</a> (Ardan Labs).
quiz: In the GMP model, what does P stand for? || processor
quiz: Which environment variable sets how many Ps (parallel OS threads running Go code) exist? || GOMAXPROCS
quiz: Roughly how big is a new goroutine's initial stack? || 2 KB
quiz: Which sync type waits for a collection of goroutines to finish? || WaitGroup
quiz: What happens to the program if a goroutine panics without recovering? || crashes
---

## A goroutine is a function call that runs concurrently

```go
go fetch("AAPL") // returns immediately; fetch runs concurrently
```

That's the whole syntax. A goroutine starts with a ~2 KB stack that grows and shrinks as needed, so **hundreds of thousands** of them are normal. OS threads cost roughly 1–8 MB and a kernel context switch each. Goroutine switches happen in user space and cost ~100 ns.

## Compared with Python

| | Python threads | Python asyncio | Go goroutines |
| --- | --- | --- | --- |
| Parallel CPU work | No (GIL; free-threaded 3.13+ is experimental) | No | **Yes, on all cores** |
| Blocking I/O | OK | Must be `await`ed, or it blocks the loop | OK, the runtime parks the goroutine |
| Function colouring | No | Yes, `async` spreads through the codebase | **No**, any function can be run with `go` |
| Cost per task | ~MBs | small | ~2 KB |

Go gets asyncio's cheap concurrency with plain blocking code. Write `resp, err := http.Get(url)` in a goroutine; while it waits on the network, the runtime runs other goroutines on that thread. No `async`/`await` and no event loop to avoid blocking.

## The scheduler: G, M, P

- **G**: a goroutine (stack, instruction pointer, state).
- **M**: an OS thread (machine).
- **P**: a logical processor, a resource needed to run Go code, with a **local run queue** of Gs. There are `GOMAXPROCS` of them (default: CPU cores, or the container's CPU limit since Go 1.25).

Each M must hold a P to run Gs. When a G blocks in a syscall, its M detaches from the P and another M picks the P up, so other goroutines keep running. Idle Ps **steal work** from busy Ps' queues. Since Go 1.14 the scheduler can **preempt** long-running goroutines asynchronously, so a tight loop can't starve the others.

Network I/O goes through the **netpoller** (epoll/kqueue). A goroutine waiting on a socket takes no thread at all. That's how one Go process handles 100k connections.

For latency-sensitive code, the same facts become questions to ask: how many runnable goroutines compete for `GOMAXPROCS` Ps? Is a hot goroutine being preempted or migrated between cores? `GODEBUG=schedtrace=1000` and the execution tracer (`go tool trace`) answer them (lesson 18).

## Waiting for goroutines

`main` returning kills every goroutine. Waiting has to be explicit:

```go
var wg sync.WaitGroup
for i, s := range symbols {
	wg.Go(func() {            // Go 1.25+: Add(1), go, defer Done() in one call
		out[i] = fetch(s)     // each goroutine writes its OWN index: no data race
	})
}
wg.Wait()
```

Before Go 1.25 the same code reads `wg.Add(1); go func() { defer wg.Done(); ... }()`. Call `Add` **before** `go`, never inside the goroutine, or `Wait` can return before the goroutine has registered.

**Loop variables:** since Go 1.22 every iteration gets a fresh `i` and `s`, so capturing them in a closure is safe. Older code (and older interview questions) show `go func(i int){...}(i)` to work around the pre-1.22 shared variable.

## Rules for every goroutine

1. **Know how it ends.** Every `go` statement needs an answer to "what makes this goroutine return?", whether that's the input channel closing, context cancellation, or the work finishing. A goroutine that never returns is a **leak**: memory that can't be collected (lesson 13).
2. **Know who waits for it.** Fire-and-forget goroutines outlive their request, swallow errors, and race with shutdown.
3. **Recover at the top if it runs untrusted work.** A panic in any goroutine crashes the whole process, and the parent can't catch it.
4. **Don't share memory without synchronisation.** Either give each goroutine its own data (as above) or use channels and locks (lessons 9–10).

## Concurrency is not parallelism

Concurrency is **structuring** a program as independently executing pieces. Parallelism is **executing** them at the same time on several cores. Go programs are written concurrently, and the runtime makes them parallel when cores are available. A concurrent design stays correct with `GOMAXPROCS=1`, which is a useful test.

## Exercise

[`exercises/08-goroutines`](../exercises/08-goroutines/README.md): a parallel sum with per-goroutine result slots, and a concurrent, order-preserving `FetchAll` whose test measures real parallelism and checks for leaked goroutines.

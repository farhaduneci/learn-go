---
title: Performance: benchmarks, pprof, allocations & the GC
section: Part 3 · Senior Go
summary: benchstat, pprof, escape analysis, allocation budgets, cache-friendly layout, GOGC/GOMEMLIMIT, PGO.
description: Low latency and high throughput come up in every trading-platform job posting. In Go the work is mostly measurement, allocation control and working with the garbage collector. This lesson is the toolkit.
source: <a href="https://go.dev/doc/diagnostics">Diagnostics</a> (go.dev), <a href="https://go.dev/doc/gc-guide">A Guide to the Go Garbage Collector</a> (essential), <a href="https://go.dev/blog/pprof">Profiling Go Programs</a>, and <a href="https://go.dev/doc/pgo">Profile-guided optimization</a>. Book: <em>Efficient Go</em> (Bartłomiej Płotka, O'Reilly).
quiz: Which tool compares benchmark results statistically? || benchstat
quiz: Which environment variable sets a soft memory limit for the GC? || GOMEMLIMIT
quiz: Which build flag prints escape analysis decisions? || -gcflags=-m
quiz: Which sync type reuses temporary objects to reduce allocations? || sync.Pool
quiz: Using a CPU profile to guide compiler optimizations is abbreviated as what? || PGO
---

## Rule zero: measure

Intuition about Go performance is often wrong. Changes that look faster frequently aren't, and the real cost often turns out to be an allocation nobody noticed. The loop is: **benchmark → profile → change one thing → benchmark again → benchstat.**

```sh
go test -run '^$' -bench BenchmarkSubmit -benchmem -count 10 ./orderbook > old.txt
# ...change...
go test -run '^$' -bench BenchmarkSubmit -benchmem -count 10 ./orderbook > new.txt
benchstat old.txt new.txt
```

Three numbers matter: **ns/op**, **B/op** and **allocs/op**. On a hot path, allocs/op is often the one to drive down first.

## pprof

```sh
go test -bench . -cpuprofile cpu.out -memprofile mem.out ./orderbook
go tool pprof -http=:8080 cpu.out          # flame graph, top, source view
go tool pprof -sample_index=alloc_space mem.out
```

In a running service, `import _ "net/http/pprof"` on an **internal admin port** exposes `/debug/pprof/`. Profiles to know:

| Profile | Answers |
| --- | --- |
| CPU | where time is spent |
| heap / allocs | what allocates, and what stays live |
| goroutine | what every goroutine is doing (leaks, deadlocks) |
| mutex / block | where goroutines wait on locks or channels |
| execution trace (`go tool trace`) | scheduler latency, GC pauses, per-goroutine timelines |

**Continuous profiling** (Pyroscope, Grafana, Google Cloud Profiler, Datadog) samples production all the time, so a regression can be compared to last week's profile. It's listed as a plus in trading-platform job postings.

## Allocations and escape analysis

Stack allocation is essentially free and disappears on return. Heap allocation costs at allocation time **and** later, in GC work. Values escape to the heap when:

- a pointer to them outlives the function (returned, or stored in a heap object);
- they're stored in an interface (`any`, `fmt.Println(x)`, `slog.Info(..., "k", v)`);
- they're captured by a closure that escapes, such as a goroutine;
- their size is unknown at compile time or too large.

```sh
go build -gcflags=-m ./orderbook 2>&1 | grep -E 'escapes|moved to heap'
```

Common fixes on hot paths:

- **Preallocate**: `make([]T, 0, n)`, `make(map[K]V, n)`, `strings.Builder.Grow(n)`.
- **Reuse buffers**: take a `dst []T` parameter and `append` to it (`strconv.AppendInt`, `SubmitInto(o, dst)`).
- **Avoid interface boxing** in tight loops. Use concrete types, typed `slog.Attr`s, and no `fmt` on hot paths.
- **`sync.Pool`** for short-lived objects allocated per request (buffers, encoders).
- Pass small structs **by value**. A pointer can force the struct to escape.

## Memory layout and the CPU cache

At low latency, **cache misses** often cost more than instructions. A main-memory load is ~100 ns; an L1 hit is ~1 ns.

- Order struct fields **largest to smallest** to cut padding (exercise 18: 32 → 24 bytes). The `fieldalignment` analyzer flags this.
- Prefer **slices of values** over slices of pointers for data iterated in bulk. Contiguous memory prefetches well, and pointer chasing doesn't. The capstone's sorted-slice price levels are chosen for this.
- Keep hot and cold fields apart, and pad contended atomics to avoid **false sharing** (lesson 10).

## The garbage collector

Go's GC is **concurrent, tri-colour mark-and-sweep, non-generational and non-moving**. Stop-the-world pauses are typically well under a millisecond. The cost shows up instead as **CPU** spent marking (roughly 25% of GOMAXPROCS while a cycle runs) and as **mark assists**, which make allocating goroutines do GC work. Those assists hurt tail latency. Go 1.26 turns on the **Green Tea** collector by default, which scans memory in a more cache-friendly way and cuts GC CPU cost noticeably on many workloads.

The two settings:

- `GOGC=100` (default): the next GC starts when the heap grows 100% over the live heap after the last one. Higher means fewer GCs and more memory.
- `GOMEMLIMIT=2GiB` (Go 1.19+): a soft limit. The GC runs more often as the limit gets close. In containers, set it to ~90% of the memory limit. A common low-latency setup is `GOGC=off` or a high GOGC, **plus** GOMEMLIMIT, so the GC only runs when memory really needs reclaiming.

The biggest GC improvement is still **allocating less**: less garbage means fewer cycles.

## Other levers

- **PGO**: put a production CPU profile at `default.pgo` in the main package, and `go build` uses it for inlining and devirtualisation decisions. That's typically 2–14% for free.
- **`GOMAXPROCS`** is container-aware since Go 1.25. `runtime.LockOSThread` plus CPU pinning are extreme measures for dedicated hot loops.
- **cgo** (calling C/C++) costs tens of nanoseconds per call and complicates builds and profiling. Use it only when a C/C++ library is genuinely required, and batch the calls. Trading shops that mix Go and C++ often keep the ultra-low-latency core in C++ and use Go around it for gateways, risk, the OMS and APIs.

## Latency is a distribution

Report **p50, p99 and p99.9**, never the average. One GC assist or lock convoy that hits 1% of orders is invisible in a mean and very visible to a client. Record latencies in a histogram (exercise 19) and load-test at realistic concurrency.

## Exercise

[`exercises/18-performance`](../exercises/18-performance/README.md): fix struct padding, get a VWAP hot path to zero allocations, and rewrite string building from 4,000 allocations to one, with budgets enforced by `testing.AllocsPerRun`. Then profile it with pprof and compare with benchstat.

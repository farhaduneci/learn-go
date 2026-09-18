# Exercise 18 — Measure, then make it fast

Lesson: [18 · Performance: benchmarks, pprof, allocations & the GC](../../lessons/0018-performance.html)

The code in `exercise.go` is correct and slow. The tests check correctness **and** a budget: a struct size and allocations per call. Measure first:

```sh
go test -run '^$' -bench . -benchmem ./exercises/18-performance/ | tee old.txt
go test -run '^$' -bench BenchmarkFormat -cpuprofile cpu.out -memprofile mem.out ./exercises/18-performance/
go tool pprof -top cpu.out             # or: go tool pprof -http=:8080 cpu.out
go build -gcflags=-m ./exercises/18-performance/ 2>&1 | grep escapes   # escape analysis
```

Then fix each item and compare with `benchstat old.txt new.txt` (`go install golang.org/x/perf/cmd/benchstat@latest`).

1. **Struct layout.** `Tick` is 32 bytes on 64-bit platforms because of alignment padding. Reorder its fields (without changing any types) to get it to 24 bytes, so a slice of ticks packs 33% more into each cache line. Check with `unsafe.Sizeof`, or run the `fieldalignment` analyzer.
2. **Zero-allocation hot path.** `VWAP` allocates an intermediate slice it doesn't need. Get it to **0 allocs/op**.
3. **Build strings without garbage.** `FormatTrades` uses `+=` and `fmt.Sprintf` in a loop. Use `strings.Builder` with `Grow` and `strconv.AppendInt` to get it to 2 or fewer allocs/op.

Budgets are enforced with `testing.AllocsPerRun`, so a regression fails CI, not just a benchmark.

Reference solution: `solution.go`.

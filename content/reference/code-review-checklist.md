---
title: Go code-review checklist
description: what a senior reviewer checks in a Go PR, grouped by category.
source: Compressed from lessons 5, 7, 10–18. Sources of truth: <a href="https://go.dev/wiki/CodeReviewComments">Go Code Review Comments</a>, <a href="https://google.github.io/styleguide/go/">Google Go Style Guide</a>, <a href="https://github.com/uber-go/guide/blob/master/style.md">Uber Go Style Guide</a>.
---

## Errors

- [ ] Every error is checked or explicitly discarded (`_ =`) with a reason.
- [ ] Errors are wrapped with context that says **what was being done**: `fmt.Errorf("reserve cash for %s: %w", id, err)`.
- [ ] `errors.Is`/`As` are used, never `==` on wrapped errors or string matching.
- [ ] Success returns a literal `nil`, never a typed nil pointer.
- [ ] No `panic` for expected failures. `recover` only at boundaries.
- [ ] Errors are either logged or returned, not both (double logging).

## API and naming

- [ ] Package names are short and specific; no `util`, `common` or `models`.
- [ ] No stutter (`orderbook.Book`, not `orderbook.OrderBook`).
- [ ] Interfaces are small and declared by the consumer. Constructors return concrete types.
- [ ] Exported identifiers have doc comments starting with the name.
- [ ] No exported mutable package-level variables. No I/O in `init()`.
- [ ] Internal slices and maps aren't returned without copying.
- [ ] Zero values are useful, or a constructor is required and documented.

## Concurrency

- [ ] Every `go` statement: how does it end, and who waits for it?
- [ ] Blocking channel operations in long-lived goroutines sit in a `select` with `ctx.Done()`.
- [ ] Channels are closed only by their sender, and only once.
- [ ] Shared state has a documented guard (`mu sync.Mutex // guards x`).
- [ ] No lock is held across I/O, a channel send, or a callback.
- [ ] Mutex-holding structs aren't copied (pointer receivers).
- [ ] `ctx` is the first parameter, is passed along rather than stored, and `cancel` is deferred.
- [ ] Tests run with `-race`, and concurrent packages use `goleak`.

## Performance (hot paths only)

- [ ] Slices and maps are preallocated when the size is known.
- [ ] No `fmt.Sprintf`, `+=` string building, or `any` boxing in tight loops.
- [ ] Allocation budgets are enforced by a test or benchmark (`AllocsPerRun`, `-benchmem`).
- [ ] Claims backed by `benchstat`, not a single run.
- [ ] `http.Client` and `http.Server` have timeouts, and response bodies are closed and drained.

## Testing

- [ ] Table-driven, with named subtests and `got X, want Y` messages.
- [ ] Fakes over mocks. Clock and randomness are injected.
- [ ] Parsers of external input have fuzz tests.
- [ ] No sleeps for synchronisation in tests. Use channels, or poll with a deadline.

## Readability

- [ ] Happy path left-aligned, with early returns for errors.
- [ ] Functions fit on a screen; no boolean flag parameters.
- [ ] Comments explain **why**, not what.
- [ ] `gofmt`, `go vet`, `staticcheck` / `golangci-lint` are clean.

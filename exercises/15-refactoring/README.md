# Exercise 15 — Refactor legacy Go

Lesson: [15 · Refactoring Go code](../../lessons/0015-refactoring.html)

```sh
go test -race ./exercises/15-refactoring/
```

`legacy.go` holds `ProcessOrders`, a working function that can't be tested. It reads the global clock, mutates a package-level map, prints to stdout, and reports errors as strings. The tests describe the **target design**. Move the logic into it step by step, keeping behaviour the same.

Target design (the stubs in `exercise.go`):

1. `ParseOrder(line)` is a **pure function** returning `Order` or an error wrapping `ErrBadFormat` / `ErrBadQty` / `ErrBadSide`.
2. Dependencies are **small interfaces declared by the consumer** (`PositionStore`, `Notifier`), plus the clock as a `func() time.Time`. Nothing global.
3. `NewService(store, now, notifier)` wires them together. `Service.Process(input)` returns `[]Result`, one per non-blank line, with **typed errors** (`ErrMarketClosed`, `ErrInsufficientPosition`) instead of strings.
4. `MemStore` is a concrete, exported in-memory `PositionStore`. Return the struct, accept the interface.
5. Delete `legacy.go` at the end. Nothing should depend on it.

Suggested order (each step keeps the build green): extract `ParseOrder`, then introduce the `Order` type, then inject the clock, then replace the global map with `PositionStore`, then replace `fmt.Println` with `Notifier`, then replace strings with sentinel errors.

Reference solution: `solution.go`.

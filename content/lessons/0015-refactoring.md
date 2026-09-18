---
title: Refactoring Go code
section: Part 3 · Senior Go
summary: Characterization tests, pure functions, consumer-side interfaces, removing globals, go fix.
description: Senior Go work is mostly changing existing code safely. This lesson covers the recurring refactors that turn untestable Go into small, well-bounded packages, in steps that keep the build green.
source: <a href="https://go.dev/wiki/CodeReviewComments">Go Code Review Comments</a> (the checklist reviewers use), <a href="https://dave.cheney.net/practical-go/presentations/qcon-china.html">Practical Go</a> (Dave Cheney), and <em>Refactoring</em> (Martin Fowler, 2nd ed.) for the step-by-step mechanics.
quiz: Where should an interface usually be defined: by the producer or the consumer package? || consumer
quiz: Complete the proverb: accept interfaces, return what? || structs
quiz: Tests written to pin down existing behaviour before a refactor are called what? || characterization tests
quiz: Which gofmt flag applies a pattern rewrite rule? || -r
---

## Start with a safety net

Before changing untested code, write **characterization tests** that pin down what it does today, bugs included. Then refactor in **small steps**, running `go test ./...` after each one, with each step committed separately. If a step breaks a test, revert it rather than debugging a large diff.

In Go the compiler is part of the safety net: rename a method, and every call site that breaks is listed. gopls can rename, extract functions or variables, inline calls and change signatures safely.

## Refactor 1: extract pure functions

Separate **decisions** from **effects**. Parsing, validation and pricing rules shouldn't touch clocks, databases or stdout:

```go
// before: parsing mixed into a 60-line loop that also writes to a global map
// after:
func ParseOrder(line string) (Order, error) // pure: trivially table-testable
```

## Refactor 2: inject dependencies, and declare interfaces at the consumer

The legacy pattern: `time.Now()`, a package-level `db`, `fmt.Println` for notifications. Replace each **hidden input** with an explicit dependency:

```go
// Declared in the package that USES them: small, only the methods needed.
type PositionStore interface {
	Position(symbol string) int64
	Add(symbol string, delta int64)
}
type Notifier interface{ Filled(o Order) }

type Service struct {
	store  PositionStore
	now    func() time.Time    // a func is the smallest possible "clock interface"
	notify Notifier
}

func NewService(store PositionStore, now func() time.Time, n Notifier) *Service
```

In Java or C# the interface usually lives next to the implementation. In Go it lives **where it's consumed**, because implicit satisfaction means the producer doesn't need to know about it. The concrete `*PostgresStore` satisfies `PositionStore` without importing the package that declares it, so packages stay decoupled and fakes need two methods instead of forty.

**Accept interfaces, return structs:** constructors return concrete types (`*MemStore`), and functions accept the narrowest interface they use. Returning an interface hides useful methods and forces type assertions.

## Refactor 3: replace string errors with typed errors

`out += "line 3: market closed"` becomes `ErrMarketClosed`, returned per line in a `Result`. Callers can now branch with `errors.Is` instead of parsing strings (lesson 5).

## Refactor 4: remove global state

Package-level mutable variables mean tests interfere with each other and parallel tests are impossible. Move the state into a struct created by a constructor, and wire it up in `main`. The same goes for `init()` that connects to things: move that into an explicit `Open`/`New`.

## Refactor 5: shrink and split

Signals that a function or package has grown too much:

- Functions longer than a screen, or nested more than 3 levels deep. Use early returns and guard clauses: *handle the error, return, keep the happy path left-aligned*.
- Boolean parameters (`Process(o, true, false)`). Split into two functions, or use an options struct.
- More than ~4 parameters. Use a config struct or functional options (lesson 16).
- A package whose name doesn't describe everything in it, or that imports half the codebase.
- Interfaces with many methods and one implementation. Delete the interface, or split it at the consumers.

## Tools that do the refactoring

```sh
gofmt -r 'a[b:len(a)] -> a[b:]' -w .   # pattern rewrite
go fix ./...                            # apply modernisations (Go 1.26 "modernizers")
gopls rename / extract / inline         # via the editor
staticcheck ./...                       # dead code, deprecated APIs, simplifications
golangci-lint run                       # runs dozens of linters in one pass
```

`go fix` in recent Go versions rewrites old idioms (`interface{}` → `any`, manual loops → `slices.Contains`, `for i := 0; i < n; i++` → `for i := range n`) automatically. Running it on an old codebase is a cheap, reviewable win.

## Exercise

[`exercises/15-refactoring`](../exercises/15-refactoring/README.md): take a working-but-untestable `ProcessOrders` (global map, `time.Now`, `fmt.Println`, string errors) and move it into a pure parser plus a service with injected, consumer-defined interfaces, step by step.

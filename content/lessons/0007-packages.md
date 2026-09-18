---
title: Packages, visibility & project layout
section: Part 1 · Go core
summary: Exported vs unexported, naming, internal/, cmd/, package-by-domain, external test packages.
description: In Go, package boundaries are the main design tool. There are no classes to hide behind, so what a package exports is its whole contract.
source: <a href="https://go.dev/doc/modules/layout">Organizing a Go module</a> (official) and <a href="https://go.dev/blog/package-names">Package names</a> (go.dev blog). The canonical style guide is <a href="https://google.github.io/styleguide/go/">Google's Go Style Guide</a>.
quiz: What makes an identifier exported? || capital letter
quiz: Packages under which directory name can only be imported by code rooted at its parent? || internal
quiz: A test file with package foo_test is called what kind of test package? || external
quiz: Which special function runs automatically when a package is initialised? || init
---

## Visibility is one rule

**Capitalised = exported**, visible to other packages. Lower-case = package-private. That covers types, functions, methods, struct fields and constants. There's no `protected` and no `_underscore` convention, and the compiler enforces it.

Because the **package** is the encapsulation unit, keep internal state unexported and expose behaviour:

```go
type Ledger struct {
	balances map[string]int64 // callers can't reach in and break invariants
}

func New() *Ledger { return &Ledger{balances: make(map[string]int64)} }
func (l *Ledger) Accounts() []string { return slices.Sorted(maps.Keys(l.balances)) } // a copy
```

Returning an internal slice or map directly hands callers a way to mutate state. Return a copy.

## Naming

- Package names are **short, lower-case, one word**: `ledger`, `orderbook`, `http`. Never `utils`, `common`, `helpers`, `models`.
- Names are read **with the package**, so avoid stutter: `ledger.New`, not `ledger.NewLedger`; `http.Server`, not `http.HTTPServer`.
- Initialisms keep one case: `ID`, `URL`, `HTTPClient`, `userID`.
- Getters drop `Get`: `acct.Balance()`, not `acct.GetBalance()`.
- Receivers are one or two letters (`l *Ledger`), never `self` or `this`.

## A real service layout

```
learn-go/
├── go.mod
├── cmd/
│   └── engine/main.go        # tiny: parse config, wire dependencies, run
├── internal/                 # importable only inside this module
│   ├── orderbook/            # domain logic, no I/O
│   ├── api/                  # HTTP handlers
│   └── store/                # database access
└── pkg/                      # optional: code meant for OTHER modules (often skipped)
```

- **Start flat.** A small service can be a single `main` package. Split when a boundary appears, not before.
- **`internal/`** is enforced by the compiler: `a/b/internal/c` can only be imported from `a/b/...`. It's the right default for application code.
- **`cmd/<name>/main.go`** holds one directory per binary. Keep `main` thin so everything else is testable.
- Package by **domain** (`orderbook`, `ledger`) rather than by layer (`models`, `controllers`).
- **No import cycles.** The compiler forbids them. A cycle means the boundary is wrong. Extract the shared piece, or define an interface on the consumer side.

## Test packages

A `_test.go` file can use **the same package** (`package ledger`: white-box, sees unexported names) or an **external test package** (`package ledger_test`: black-box, only the exported API). External tests prove the API is usable as designed and are good at catching awkward APIs early.

## init() and package-level state

`func init()` runs automatically when a package loads. Use it sparingly: it makes behaviour depend on import order and makes tests harder to write. Package-level mutable variables (like the global map in exercise 15's legacy code) are the main source of untestable Go. Pass dependencies in explicitly.

## Exercise

[`exercises/07-packages`](../exercises/07-packages/README.md): a double-entry ledger package with an unexported representation, exported sentinel errors, a copy-returning accessor, and black-box tests in `ledger_test`.

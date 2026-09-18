---
title: Errors as values
section: Part 1 · Go core
summary: error values, sentinels, typed errors, %w wrapping, errors.Is/As, the nil-interface trap, panic/recover.
description: Go has no exceptions for expected failures. Errors are ordinary return values, and that turns error handling into visible, reviewable control flow. Doing it well (wrapping, sentinels, typed errors) marks senior Go code.
source: <a href="https://go.dev/blog/go1.13-errors">Working with Errors in Go 1.13</a> (go.dev blog) and <a href="https://go.dev/blog/errors-are-values">Errors are values</a> (Rob Pike). The <a href="https://pkg.go.dev/errors">errors package docs</a> are short.
quiz: Which fmt verb wraps an error so errors.Is can see through it? || %w
quiz: Which function checks whether any error in the chain equals a target value? || errors.Is
quiz: Which function extracts a specific error type from the chain? || errors.As
quiz: Which built-in stops a panic when called inside a deferred function? || recover
quiz: Which function combines several errors into one (Go 1.20)? || errors.Join
---

## The pattern

```go
f, err := os.Open(path)
if err != nil {
	return fmt.Errorf("load config %s: %w", path, err)
}
defer f.Close()
```

`error` is a one-method interface: `Error() string`. Functions that can fail return it **last**. Callers check it **immediately**, then handle it or return it with added context. The happy path stays left-aligned and every failure point is visible, so no code path can throw from somewhere unexpected.

Python translation: `raise` becomes `return ..., err`, and `try/except` becomes `if err != nil`. `finally` becomes `defer`.

## Three kinds of errors

**1. Opaque**: `errors.New("...")` or `fmt.Errorf("...")`. Callers only log or return it. This is the default.

**2. Sentinel**: a package-level value callers compare against:

```go
var ErrInsufficientFunds = errors.New("insufficient funds")
// io.EOF, sql.ErrNoRows, context.Canceled are all sentinels
```

**3. Typed**: a struct carrying data:

```go
type ValidationError struct{ Field, Reason string }
func (e *ValidationError) Error() string { return "invalid " + e.Field + ": " + e.Reason }
```

Each exported sentinel or type becomes part of the package's API, so export them only when callers genuinely need to branch on them.

## Wrapping, Is and As

`%w` wraps an error while adding context. The chain can still be inspected:

```go
err := fmt.Errorf("place order %s: %w", id, ErrInsufficientFunds)

errors.Is(err, ErrInsufficientFunds) // true: walks the chain comparing values

var ve *ValidationError
if errors.As(err, &ve) {             // finds the first *ValidationError in the chain
	log.Println("bad field:", ve.Field)
}
```

Never compare with `err == ErrX` once wrapping is involved, and never match on `err.Error()` strings. Use `%v` instead of `%w` when the underlying error should **not** become part of the API. Wrapping exposes it.

Good context reads like a path: `place order o-17: reserve cash: db: connection refused`. Each layer adds **what it was doing**, not "error occurred".

## The nil interface trap

An interface is nil only when **both** its type and value are nil:

```go
func validate() error {
	var ve *ValidationError // nil pointer
	if bad { ve = &ValidationError{...} }
	return ve               // BUG: returns (type=*ValidationError, value=nil), which != nil
}
```

Return a literal `nil` on success. This is one of the most common Go interview questions.

## panic and recover

`panic` is for **programmer errors and impossible states**: index out of range, nil map write, broken invariants. Never use it for expected failures like bad input or network errors. `recover` inside a `defer` stops a panic, but only in the same goroutine:

```go
func SafeDiv(a, b int) (q int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("safe div: %v", r) // named result lets defer set it
		}
	}()
	return a / b, nil
}
```

Legitimate uses are narrow: HTTP servers recover per request (`net/http` does this), and worker pools recover per job so one bad message doesn't kill the process. **A panic in any goroutine without a recover crashes the entire program**, which matters for lesson 8.

## defer

`defer` runs a call when the surrounding function returns, in LIFO order. The arguments are evaluated immediately. It's the tool for cleanup: `defer f.Close()`, `defer mu.Unlock()`, `defer cancel()`. A deferred closure can read and modify **named** results, as `SafeDiv` does.

## Exercise

[`exercises/05-errors`](../exercises/05-errors/README.md): a typed validation error, a sentinel, wrapping with context, `errors.Is`/`As` through the chain, and turning a panic into an error.

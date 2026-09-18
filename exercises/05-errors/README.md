# Exercise 05 — Errors

Lesson: [05 · Errors as values](../../lessons/0005-errors.html)

```sh
go test -race ./exercises/05-errors/
```

1. `Validate(o)` returns a `*ValidationError` naming the first bad field: `Symbol` empty, `Qty <= 0`, or `Side` not `"buy"`/`"sell"`. Return `nil` for a valid order. Make sure a valid order gives a **nil `error` interface**, not a typed nil pointer (see the lesson's nil-interface trap).
2. `Account.Place(o)` validates, then returns `ErrInsufficientFunds` (a sentinel) when `Qty*Price > Cash`. Wrap every error with context: `fmt.Errorf("place order %s: %w", o.ID, err)`.
3. Callers must be able to use `errors.Is(err, ErrInsufficientFunds)` and `errors.As(err, &ve)` through the wrapping.
4. `SafeDiv(a, b)` turns a runtime panic (integer divide by zero) into an error with `defer` + `recover`. This is a boundary pattern; it isn't normal control flow.

Reference solution: `solution.go`.

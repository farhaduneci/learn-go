---
title: context.Context
section: Part 2 · Concurrency
summary: Cancellation trees, deadlines, request-scoped values, first-result-wins, retry with backoff.
description: Context carries cancellation, deadlines and request-scoped values across API boundaries and goroutines. Every production Go function that does I/O takes one, and interviewers check that it's used correctly.
source: <a href="https://go.dev/blog/context">Go Concurrency Patterns: Context</a> (go.dev blog) and the <a href="https://pkg.go.dev/context">context package docs</a>. Also <a href="https://go.dev/blog/context-and-structs">Contexts and structs</a>.
quiz: What must always be called after context.WithTimeout returns, typically via defer? || cancel
quiz: Which method returns a channel closed when the context is cancelled? || Done
quiz: Which error does ctx.Err() return when a deadline passes? || context.DeadlineExceeded
quiz: By convention, context is which parameter of a function? || first
quiz: Should a context be stored in a struct field? || no
---

## The interface

```go
type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}  // closed on cancel or deadline
	Err() error             // nil, context.Canceled, or context.DeadlineExceeded
	Value(key any) any
}
```

Contexts form a **tree**. Derive children from a parent, and cancelling a parent cancels every descendant:

```go
ctx := context.Background()                              // root, in main or tests
ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
defer cancel()                                           // ALWAYS: releases the timer and children
ctx, cancel2 := context.WithCancel(ctx)                  // cancel manually, e.g. first result wins
ctx = context.WithValue(ctx, requestIDKey{}, "req-42")
```

`net/http` gives each request `r.Context()`, which is cancelled when the client disconnects. Database drivers, gRPC and HTTP clients all take a context and stop work when it's done.

## The conventions

1. `ctx context.Context` is the **first parameter**, named `ctx`: `func (s *Svc) Place(ctx context.Context, o Order) error`.
2. **Don't store contexts in structs.** Pass them through calls. A context belongs to one operation, not to an object's lifetime.
3. **Never pass `nil`.** Use `context.TODO()` while plumbing is unfinished.
4. **Always call `cancel`**, even on success. `go vet` (lostcancel) catches a missing call.
5. **Values are for request-scoped metadata** such as request IDs, trace spans and auth identity. Don't use them for optional parameters or dependencies like DB handles or loggers configured per service. Key them with an **unexported type** so packages can't collide:

```go
type requestIDKey struct{}
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}
```

## Honouring cancellation

A context only works if long operations **check** it:

```go
func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

for _, item := range batch {        // CPU-bound loops: check between units of work
	if err := ctx.Err(); err != nil {
		return err
	}
	process(item)
}
```

## First result wins, cancel the rest

A classic pattern: query several venues or replicas and take the fastest answer ("hedged requests"):

```go
func FirstSuccess(ctx context.Context, venues ...Venue) (int64, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()                               // losers are cancelled as soon as this returns
	results := make(chan result, len(venues))    // buffered: losers can always finish their send
	for _, v := range venues {
		go func() { p, err := v(ctx); results <- result{p, err} }()
	}
	...
}
```

## Newer helpers

- `context.WithoutCancel(ctx)` (1.21): keeps the values but detaches cancellation. Use it for work that must finish after the request ends, such as writing an audit log.
- `context.AfterFunc(ctx, f)` (1.21): runs `f` when the context is done.
- `context.WithCancelCause` / `context.Cause` (1.20): records **why** something was cancelled.
- `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)`: a context cancelled by SIGINT/SIGTERM, which is the start of graceful shutdown (lesson 17).

## Exercise

[`exercises/11-context`](../exercises/11-context/README.md): a cancellable sleep, first-success with loser cancellation, request-ID values with a private key type, and retry with exponential backoff that respects deadlines.

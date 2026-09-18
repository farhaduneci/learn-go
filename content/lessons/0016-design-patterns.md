---
title: Design patterns in Go
section: Part 3 · Senior Go
summary: Functional options, middleware chains, strategy as func, pub/sub backpressure, iterators.
description: Classic patterns look different in Go. With first-class functions, implicit interfaces and no inheritance, most Gang-of-Four patterns shrink to a few lines, and a few Go-specific ones come up in every codebase.
source: <a href="https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis">Functional options for friendly APIs</a> (Dave Cheney), <a href="https://go.dev/blog/errors-are-values">Errors are values</a>, and <a href="https://github.com/uber-go/guide/blob/master/style.md">Uber Go Style Guide</a> (sections on functional options, embedding, and interfaces).
quiz: A variadic list of func(*T) arguments used to configure a constructor is which pattern? || functional options
quiz: In Chain(h, A, B), which middleware runs first? || A
quiz: When a pub/sub subscriber is slow, a non-blocking publisher uses select with which clause? || default
quiz: Which GoF pattern is a type implementing an interface by wrapping another implementation of it? || decorator
---

## Functional options

The answer to "constructors with many optional settings" when there are no default arguments or overloads:

```go
type Option func(*Client) error

func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d <= 0 {
			return errors.New("timeout must be positive")
		}
		c.Timeout = d
		return nil
	}
}

func NewClient(baseURL string, opts ...Option) (*Client, error) {
	c := &Client{BaseURL: baseURL, Timeout: 5 * time.Second, Retries: 3} // defaults
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

c, err := NewClient(url, WithTimeout(time.Second), WithRetries(0))
```

It's backwards compatible (adding an option breaks no caller), self-documenting, and has sensible defaults. The simpler alternative is a **config struct** whose zero values mean "default". Both are idiomatic: use options for libraries, and a config struct when settings come from a file.

## Middleware (decorator)

Wrap a function type in another of the same type:

```go
type Handler func(ctx context.Context, o Order) error
type Middleware func(next Handler) Handler

func Chain(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- { // wrap in reverse so mws[0] is outermost
		h = mws[i](h)
	}
	return h
}

h := Chain(placeOrder, Recover, Logging, Validate, RateLimit)
```

`net/http` middleware is exactly this with `http.Handler`. So are gRPC interceptors. Cross-cutting concerns such as logging, metrics, auth, validation and panic recovery live here instead of in every handler.

## Strategy = interface or func

A strategy is an interface (`FeeModel`, lesson 4) or just a function value (`func(int64) int64`). When the "interface" would have one method and no state, **a func type is simpler**: `sort.Slice` takes `less func(i, j int) bool`, and `http.HandlerFunc` adapts a function to an interface.

## Observer / pub-sub with backpressure

Go's version uses channels, and the important design question is **what happens when a subscriber is slow**:

```go
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs[e.Topic] {
		select {
		case ch <- e:        // delivered into the subscriber's buffer
		default:             // buffer full: drop for THIS subscriber, count it
			b.dropped.Add(1)
		}
	}
}
```

The choices are block (one slow consumer stalls everyone), drop (with metrics), coalesce (keep only the latest price per symbol), or disconnect the slow client. Market-data fan-out usually drops or coalesces. Order and fill streams must never drop, so they block or halt. Saying this out loud in a system-design interview scores points.

## Other patterns that show up

- **Constructor + unexported struct fields** for invariants (lesson 7).
- **Adapter:** `http.HandlerFunc(f)` turns a func into an interface.
- **Singleton:** a package-level `sync.OnceValue`, used sparingly because it's global state.
- **Builder:** rare, since functional options or struct literals cover it.
- **Generator/iterator:** `iter.Seq[T]` (Go 1.23) lets custom collections work with `for range`:

```go
func (b *Book) Levels() iter.Seq[Level] {
	return func(yield func(Level) bool) {
		for _, l := range b.levels {
			if !yield(l) { return }
		}
	}
}
for l := range book.Levels() { ... }
```

## What not to port from Python or Java

- Abstract base classes and deep hierarchies. Use small interfaces plus composition instead.
- Getters and setters for every field. Export the field, or give the type real behaviour.
- `Manager`/`Helper`/`Util` types. They're named after nothing in the domain.
- Interfaces declared before a second implementation (or a test fake) exists.

## Exercise

[`exercises/16-design-patterns`](../exercises/16-design-patterns/README.md): functional options with validation, an ordered middleware chain with validation and panic recovery, and a non-blocking pub/sub bus with drop accounting.

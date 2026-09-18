---
title: Methods, interfaces & composition
section: Part 1 · Go core
summary: Implicit interfaces, small interfaces, type switches, decorators, embedding without inheritance.
description: Go has no classes and no inheritance. Small, implicitly satisfied interfaces plus struct embedding take their place, and the design style that follows is the core of idiomatic Go.
source: <a href="https://go.dev/doc/effective_go#interfaces_and_types">Effective Go — Interfaces and other types</a> and <a href="https://go.dev/doc/effective_go#embedding">Embedding</a>. Also <a href="https://go.dev/tour/methods/1">A Tour of Go — Methods and interfaces</a>.
quiz: Does a Go type need to declare which interfaces it implements? || no
quiz: What is the name of the empty interface alias added in Go 1.18? || any
quiz: Which construct branches on an interface value's dynamic type? || type switch
quiz: Promoting the methods of a field into the outer struct is called what? || embedding
quiz: Go proverb: the bigger the interface, the weaker the what? || abstraction
---

## Interfaces are satisfied implicitly

An interface is a set of method signatures. **Any type with those methods satisfies it.** There's no `implements` keyword.

```go
type FeeModel interface {
	Fee(notional int64) int64
}

type FlatFee int64
func (f FlatFee) Fee(int64) int64 { return int64(f) }

type BpsFee int64
func (b BpsFee) Fee(n int64) int64 { return n * int64(b) / 10_000 }

var m FeeModel = BpsFee(5) // compiles because BpsFee has Fee(int64) int64
```

This is Python's duck typing checked at compile time. A package can define an interface that a type from another package already satisfies, without touching that type. That changes where interfaces live (lesson 15).

A compile-time assertion documents intent and catches drift:

```go
var _ FeeModel = FlatFee(0)
```

## Small interfaces

The standard library's most powerful interfaces have one method:

```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
type Stringer interface { String() string }
type error interface { Error() string }
```

Files, sockets, buffers, gzip streams, HTTP bodies and hashes all compose through `io.Reader` and `io.Writer`. *"The bigger the interface, the weaker the abstraction"* is one of the Go proverbs. Build bigger interfaces by embedding small ones: `type ReadWriter interface { Reader; Writer }`.

## Interface values, type assertions, type switches

An interface value holds a (**dynamic type**, **value**) pair. To get the concrete type back:

```go
bf, ok := m.(BpsFee) // type assertion, with "comma ok" to avoid a panic

switch v := x.(type) { // type switch
case FlatFee:
	fmt.Println("flat", int64(v))
case BpsFee, Capped:
	fmt.Println("bps or capped")
default:
	fmt.Println("unknown")
}
```

`any` (= `interface{}`) accepts every type. It's fine at true boundaries (JSON decoding, logging) and a smell in domain code, where it throws away the type checker.

## Composition over inheritance

**Wrap** instead of subclassing. A decorator holds another implementation of the same interface:

```go
type Capped struct {
	Inner FeeModel
	Max   int64
}

func (c Capped) Fee(n int64) int64 { return min(c.Inner.Fee(n), c.Max) }
```

`Capped{Inner: BpsFee(5), Max: 1000}` is itself a `FeeModel`, so decorators stack. HTTP middleware (lesson 16) is the same idea.

## Embedding

Putting a type in a struct **without a field name** promotes its fields and methods:

```go
type AuditedAccount struct {
	*Ledger       // embedded: a.Post(...) and a.Total work directly
	Log []int64
}

func (a *AuditedAccount) Post(amount int64) { // "overrides" by shadowing
	a.Log = append(a.Log, amount)
	a.Ledger.Post(amount) // explicit call to the embedded one: there is no super()
}
```

This is **not inheritance**. The embedded `Ledger` doesn't know about the outer type, so a `Ledger` method that calls `Post` internally still calls `Ledger.Post`. There's no virtual dispatch through embedding. Expect that point in interviews.

## Exercise

[`exercises/04-interfaces`](../exercises/04-interfaces/README.md): three fee models behind one interface, a decorator, a type switch, and an embedded ledger with an overriding method.

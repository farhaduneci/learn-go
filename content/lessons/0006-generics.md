---
title: Generics
section: Part 1 · Go core
summary: Type parameters, constraints, ~underlying types, generic containers, when not to use generics.
description: Type parameters (Go 1.18+) remove the copy-paste and interface{} casts behind containers and algorithms. Knowing when not to use them matters just as much.
source: <a href="https://go.dev/doc/tutorial/generics">Tutorial: Getting started with generics</a> and <a href="https://go.dev/blog/when-generics">When To Use Generics</a> (Ian Lance Taylor, go.dev blog).
quiz: Which constraint allows == and != and map keys? || comparable
quiz: Which standard package defines the Ordered constraint? || cmp
quiz: In a constraint, what does the ~ in ~int mean? || underlying type
quiz: Methods cannot declare their own type parameters: true or false? || true
---

## Syntax

```go
func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		out = append(out, f(v))
	}
	return out
}

strs := Map([]int{1, 2}, strconv.Itoa) // T=int, U=string inferred
```

Type parameters go in square brackets. Each has a **constraint**, which is an interface.

## Constraints

```go
any                       // no requirements
comparable                // supports == (usable as a map key)
cmp.Ordered               // supports < <= > >= (ints, floats, strings)

type Number interface {   // a type set
	~int | ~int64 | ~float64
}
```

`~int64` means "any type whose **underlying type** is int64", so a custom `type Cents int64` satisfies it. Without the tilde, only `int64` itself does.

## Generic types

```go
type Set[T comparable] struct{ m map[T]struct{} }

func (s *Set[T]) Add(v T) {
	if s.m == nil { s.m = make(map[T]struct{}) } // useful zero value
	s.m[v] = struct{}{}
}

type Stack[T any] struct{ items []T }
func (s *Stack[T]) Pop() (v T, ok bool) { ... } // "var zero T" gives the zero value
```

Methods can use the type's parameters but **can't introduce new ones**. A method `func (s *Set[T]) Map[U any]()` is illegal, so write a top-level function instead. The same goes for stricter constraints: `SortedItems[T cmp.Ordered](s *Set[T])` has to be a function.

## When to use them

Use generics for:

- **Containers**: sets, stacks, queues, caches, trees, ring buffers.
- **Algorithms over any element type**: map/filter/reduce, sorting helpers, `slices.*` and `maps.*`.
- Removing duplicated code that differs **only in type**.

Prefer an **interface** when behaviour differs per type. `FeeModel` (lesson 4) has different logic for each implementation, so a type parameter would add nothing. The rule from the Go team: *write code, don't design types*. Reach for generics when the same code is being written a second time for a different type, not up front.

Performance note: Go compiles generics with **GC shape stenciling**, which shares code between types of the same shape and passes a dictionary. Generic code is usually as fast as hand-written code, but calling methods through a type-parameter constraint can be slower than on a concrete type in hot paths. Benchmark before assuming.

## Exercise

[`exercises/06-generics`](../exercises/06-generics/README.md): Map/Filter/Reduce, `MaxOf` with `cmp.Ordered`, a zero-value-ready `Set`, a `Stack` that clears popped references, and `GroupBy`.

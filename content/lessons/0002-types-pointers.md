---
title: Types, zero values, structs & pointers
section: Part 1 · Go core
summary: Static types, zero values, structs, methods, pointer vs value receivers, escape analysis, money as int64.
description: Static types, zero values and explicit pointers are the biggest mental shift from Python. They're also where Go's performance and predictability come from.
source: <a href="https://go.dev/tour/basics/1">A Tour of Go — Basics</a> and <a href="https://go.dev/tour/moretypes/1">More types: pointers and structs</a>. For depth, <a href="https://go.dev/ref/spec#Types">the spec's Types section</a> is short and readable.
quiz: What is the zero value of a string? || ""
quiz: What is the zero value of a pointer, slice, map or channel? || nil
quiz: Which operator takes the address of a variable? || &
quiz: To let a method modify its receiver, the receiver must be a what? || pointer
quiz: What compiler analysis decides whether a variable lives on the stack or the heap? || escape analysis
---

## Static types, inferred where obvious

```go
var count int          // declared, zero value 0
price := 189.5         // inferred float64
var qty int64 = 10
name, ok := "AAPL", true
```

Basic types: `bool`, `string`, `int`/`int8`…`int64`, `uint`…`uint64`, `float32`/`float64`, `byte` (= `uint8`), `rune` (= `int32`, a Unicode code point). `int` is 64-bit on 64-bit platforms. **No implicit conversions**: `int64(qty) * price` won't compile until both sides have the same type. That feels strict for a day and prevents whole classes of bugs.

Constants are untyped until used, and `iota` builds enums:

```go
type Side int

const (
	Buy Side = iota + 1 // 1
	Sell                // 2
)
```

## Zero values instead of None

Every type has a **zero value**, and every variable starts at it: `0`, `""`, `false`, `nil` for pointers/slices/maps/channels/interfaces/functions, and a struct whose fields are all zero. There is no "uninitialized" state and no `None` sprinkled through the code.

Idiomatic Go designs types so the **zero value is useful**: `var mu sync.Mutex` is an unlocked mutex, `var b strings.Builder` is an empty builder, `var s []int` can be appended to. Aim for the same in custom types.

## Structs, not classes

```go
type Account struct {
	ID      string
	Balance Cents // named type over int64
}

a := Account{ID: "acc-1"}   // Balance is 0
p := &Account{ID: "acc-2"}  // pointer to a new Account
p.Balance = 100             // automatic dereference: no p->Balance
```

There are no constructors. By convention a function named `NewAccount(...) *Account` builds one when setup is needed. No inheritance either; lesson 4 shows composition instead.

## Methods and receivers

```go
type Cents int64

func (c Cents) String() string { ... }         // value receiver: gets a copy
func (a *Account) Deposit(amount Cents) {     // pointer receiver: can modify
	a.Balance += amount
}
```

Rules of thumb:

- Use a **pointer receiver** if the method mutates, if the struct is large, or if it contains a `sync.Mutex` (copying a mutex is a bug that `go vet` catches).
- Keep receivers consistent. If one method needs a pointer, make them all pointers.
- Value receivers suit small immutable types like `Cents`, `time.Time` and `Side`.

## Pointers: explicit references

Python passes references to objects implicitly. Go passes **everything by value** (a copy), and uses pointers when sharing is wanted:

```go
func Swap(a, b *int) { *a, *b = *b, *a }

x, y := 1, 2
Swap(&x, &y) // & = address-of, * = dereference
```

There's no pointer arithmetic (outside `unsafe`), and the garbage collector manages lifetimes. Returning a pointer to a local variable is fine: the compiler notices it **escapes** and allocates it on the heap.

```sh
go build -gcflags=-m ./...   # prints escape-analysis decisions: "moved to heap: x"
```

Stack allocation is nearly free, and heap allocation costs GC work later. This is the first lever for low-latency Go, covered in lesson 18.

## Money is never a float

`0.1 + 0.2 != 0.3` in every IEEE-754 language. Trading and ledger systems represent money as **integers in the smallest unit** (cents, or micro-units for crypto and FX) or as a decimal type (`shopspring/decimal`). This track uses `int64` cents throughout, and interviewers at trading firms notice when someone reaches for `float64` for prices.

## Closures

Functions are values, and closures capture variables by reference, just like Python:

```go
func NewCounter() func() int {
	n := 0
	return func() int { n++; return n }
}
```

## Exercise

[`exercises/02-types-pointers`](../exercises/02-types-pointers/README.md): a `Cents` type with `String()`, an `Account` with a deliberately wrong value receiver to fix, pointer swapping, and a closure counter.

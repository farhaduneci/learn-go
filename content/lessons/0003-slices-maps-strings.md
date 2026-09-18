---
title: Slices, maps, strings & runes
section: Part 1 · Go core
summary: Slice headers, append aliasing, full slice expressions, maps as sets, bytes vs runes.
description: Slices look like Python lists and behave differently in ways that cause real bugs. The slice header, append and aliasing are standard interview material.
source: <a href="https://go.dev/blog/slices-intro">Go Slices: usage and internals</a> (go.dev blog, ~10 min). Then <a href="https://go.dev/blog/strings">Strings, bytes, runes and characters in Go</a>.
quiz: A slice header holds a pointer, a length and a what? || capacity
quiz: Iterating a string with range yields bytes or runes? || runes
quiz: Is map iteration order in Go deterministic? || no
quiz: Which full slice expression caps capacity so append must copy? || s[lo:hi:max]
quiz: What does len() of a string count? || bytes
---

## Arrays are values, slices are views

An **array** has its length in its type (`[4]int`) and is copied on assignment. It's rarely used directly. A **slice** is a small header of three words that points into an array:

```
slice header:  ptr ──► [ 1 | 2 | 3 | 4 | _ | _ ]   backing array
               len = 4
               cap = 6
```

```go
s := []int{1, 2, 3}             // literal
s = make([]int, 0, 100)         // len 0, cap 100: preallocate when the size is known
t := s[1:3]                     // a new header onto the SAME array
```

Because `t` shares the backing array, **writing `t[0]` changes `s[1]`**. Python's `lst[1:3]` makes a copy; Go's doesn't.

## append and the aliasing trap

`append` writes into spare capacity if there is any, and otherwise allocates a bigger array (roughly doubling) and copies. Always reassign: `s = append(s, v)`.

```go
base := make([]int, 3, 10)
a := append(base, 1) // fits in cap: writes base's array at index 3
b := append(base, 2) // ALSO writes index 3, so a[3] is now 2
```

Two slices from the same parent can silently clobber each other. The fixes:

- The **full slice expression** `s[lo:hi:hi]` sets `cap == len`, so the next append must copy.
- An explicit copy: `c := slices.Clone(s)`, or `make` + `copy(dst, src)`.

This bug is common in chunking and batching code, and it's what exercise 03's `Chunk` test catches.

## nil vs empty

`var s []int` is **nil** (len 0, no array). `s := []int{}` is empty but non-nil. Both work with `len`, `range` and `append`. The difference shows in JSON (`null` vs `[]`) and in `== nil` checks. Return `nil` for "nothing" unless an API contract says otherwise.

## The slices and maps packages

Since Go 1.21 the standard library covers most of what Python's builtins do:

```go
slices.Sort(s); slices.Contains(s, v); slices.Index(s, v)
slices.Reverse(s); slices.Max(s); slices.Equal(a, b)
slices.SortFunc(people, func(a, b P) int { return cmp.Compare(a.Age, b.Age) })
keys := slices.Sorted(maps.Keys(m)) // iterators (Go 1.23)
```

## Maps

```go
vol := map[string]int{"AAPL": 50}
vol["TSLA"] += 10          // missing key reads as the zero value, so this works
v, ok := vol["MSFT"]       // "comma ok" distinguishes missing from zero
delete(vol, "AAPL")
set := map[string]struct{}{} // a set: struct{} takes zero bytes
```

- **Iteration order is randomized on purpose.** Sort keys when output must be deterministic.
- A nil map can be read but **panics on write**. Use `make` or a literal first.
- Maps are **not safe for concurrent writes**. The runtime crashes with `fatal error: concurrent map writes`. Lesson 10 fixes that.
- Keys must be comparable: no slices, maps or funcs.

## Strings, bytes and runes

A Go string is an **immutable sequence of bytes**, conventionally UTF-8. `len("héllo")` is 6 (bytes), not 5.

```go
s := "héllo"
s[1]                // a byte (0xc3), not 'é'
for i, r := range s // r is a rune (code point); i is its byte offset
[]rune(s)           // decode to code points: reverse/index by character
utf8.RuneCountInString(s) // 5
```

Building strings in a loop with `+=` allocates a new string each time, O(n²). Use `strings.Builder` (lesson 18). `strings` has the usual helpers: `Fields`, `Split`, `Cut`, `TrimSpace`, `ToLower`, `HasPrefix`, `Join`.

## Exercise

[`exercises/03-slices-maps-strings`](../exercises/03-slices-maps-strings/README.md): reversing without mutating, order-preserving dedupe with a set, word frequency, rune-aware reversal, aliasing-safe chunking, and top-k with a deterministic sort.

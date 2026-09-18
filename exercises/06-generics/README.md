# Exercise 06 — Generics

Lesson: [06 · Generics](../../lessons/0006-generics.html)

```sh
go test -race ./exercises/06-generics/
```

1. `Map`, `Filter`, `Reduce` over slices. Go has no built-in versions; these are the classic first generic functions.
2. `MaxOf[T cmp.Ordered]` returns the largest element and `false` for an empty slice.
3. `Set[T comparable]` with `Add`, `Has`, `Len`. Sorting needs the stricter `cmp.Ordered` constraint, so the sorted listing is a separate function: `SortedItems[T cmp.Ordered](*Set[T])`.
4. `Stack[T]` with `Push`, `Pop` (returns `T, bool`) and `Peek`. The zero value of `Stack[T]` must be ready to use.
5. `GroupBy` groups items into a `map[K][]T` by a key function.

The lesson covers when **not** to use generics. The fee models in exercise 04 are the right place for an interface, not a type parameter.

Reference solution: `solution.go`.

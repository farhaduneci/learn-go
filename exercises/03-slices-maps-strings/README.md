# Exercise 03 — Slices, maps & strings

Lesson: [03 · Slices, maps, strings & runes](../../lessons/0003-slices-maps-strings.html)

```sh
go test -race ./exercises/03-slices-maps-strings/
```

1. `Reversed(s)` returns a reversed **copy**. The input must not change (slices share their backing array).
2. `Dedupe(s)` removes duplicates and keeps first-seen order. Use a `map[string]struct{}` as a set.
3. `WordFreq(text)` counts lower-cased words split on whitespace (`strings.Fields`).
4. `ReverseString(s)` reverses by **rune**, not byte: `"héllo"` → `"olléh"`, `"日本"` → `"本日"`.
5. `Chunk(s, n)` splits into chunks of at most `n`. Appending to one chunk must not overwrite the next one. Hint: the full slice expression `s[lo:hi:hi]`.
6. `TopSymbols(volumes, k)` returns the `k` symbols with the highest volume, ties broken alphabetically. Map iteration order is random, so sort the result (`slices.SortFunc`).

Reference solution: `solution.go`.

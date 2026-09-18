# Exercise 01 — Toolchain

Lesson: [01 · Toolchain & modules](../../lessons/0001-toolchain.html)

Implement the functions in `exercise.go` until the tests pass.

```sh
go test ./exercises/01-toolchain/          # run the tests
go test -run TestHello -v ./exercises/01-toolchain/
go vet ./exercises/01-toolchain/           # static checks
gofmt -l ./exercises/01-toolchain/         # lists unformatted files (expect no output)
```

Tasks:

1. `Hello(name)` returns `"Hello, <name>!"`, or `"Hello, world!"` when `name` is empty.
2. `Sum(nums ...int)` returns the sum of any number of ints (a variadic function — like `*args` in Python).
3. `Stats(nums []int)` returns the min, max, and whether the slice was non-empty (multiple return values instead of a tuple).

Reference solution: `solution.go` (run it with `go test -tags solution ./exercises/01-toolchain/`). Try the exercise before opening it.

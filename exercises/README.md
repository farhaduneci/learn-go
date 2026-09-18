# Exercises

One directory per lesson. Each has:

- `README.md`: the tasks.
- `exercise.go`: stubs to implement (build tag `!solution`).
- `solution.go`: a reference implementation (build tag `solution`). Attempt the exercise first.
- `*_test.go`: tests shared by both.

```sh
go test -race ./exercises/03-slices-maps-strings/          # run one exercise against the stubs
go test -race -tags solution ./exercises/...               # run every reference solution
```

Some stubs are deliberately incomplete types, so the package **doesn't compile** until the right methods exist. That compile error is the first failing test.

Always use `-race` for exercises 08 onward. Exercise 14 is different: the code is given and the tests are the task.

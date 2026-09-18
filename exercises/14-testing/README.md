# Exercise 14 — Write the tests (and find the bug)

Lesson: [14 · Testing like a Go team](../../lessons/0014-testing.html)

This time the code is written and the **tests are the task**. `ParsePrice` and `FormatPrice` convert between strings like `"12.34"` and integer cents. There is **one bug** in `exercise.go`. Don't hunt for it by reading; let the tests find it.

Fill in `exercise_test.go`:

1. **Table-driven test** for `ParsePrice`: valid inputs (`"0"`, `"12"`, `"12.3"`, `"12.34"`, `"-5.00"`) and invalid ones (`""`, `"abc"`, `"1.234"`, `"1.2.3"`). Use `t.Run` subtests with descriptive names.
2. **Fuzz test** `FuzzRoundTrip`: for any `int64` cents `c`, `ParsePrice(FormatPrice(c))` must return `c`. Add a few seeds with `f.Add`, then run:
   ```sh
   go test -run '^$' -fuzz FuzzRoundTrip -fuzztime 20s ./exercises/14-testing/
   ```
   A failure saves the input under `testdata/fuzz/`, where it becomes a permanent regression test.
3. **Benchmark** `BenchmarkParsePrice` using `b.Loop()` (Go ≥ 1.24), and run it with `-benchmem`:
   ```sh
   go test -run '^$' -bench . -benchmem ./exercises/14-testing/
   ```
4. Fix the bug in `exercise.go`. Keep the failing input as a table case.

Reference: `solution_test.go` and `solution.go` (`go test -tags solution ./exercises/14-testing/`).

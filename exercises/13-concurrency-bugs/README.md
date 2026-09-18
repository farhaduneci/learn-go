# Exercise 13 — Find and fix concurrency bugs

Lesson: [13 · Concurrency bugs: leaks, deadlocks, races](../../lessons/0013-concurrency-bugs.html)

```sh
go test -race ./exercises/13-concurrency-bugs/
```

`exercise.go` compiles and looks fine, but each function has a classic concurrency bug. For each one, read the failing test output, **name the bug**, then fix it with the smallest change that works.

| Function | Symptom | Tool that catches it |
| --- | --- | --- |
| `FirstQuote` | goroutines left behind after it returns | `goleak` |
| `Transfer` | program hangs under concurrent opposite transfers | test timeout / `SIGQUIT` goroutine dump |
| `ValidateAll` | wrong or missing error, sometimes | `-race` |
| `Feed.Stop` | `panic: close of closed channel` | the panic itself |

Write down the bug name next to each fix, e.g. `// fix: lock ordering`. Interviewers ask about these four more than any others.

Reference solution: `solution.go`.

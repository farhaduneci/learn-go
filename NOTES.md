# Notes

- Toolchain: Go 1.26 (uses `wg.Go` from 1.25, `b.Loop` from 1.24, range-over-int and per-iteration loop variables from 1.22). Older Go versions won't compile the solutions.
- Build: lesson HTML is generated. Edit `content/**/*.md`, then run `python3 scripts/build.py` (needs `pandoc`). The index's learning path comes from each lesson's front matter (`status`, `summary`), so update `status` there as lessons are completed.
- Exercises: stubs use `//go:build !solution`, reference solutions `//go:build solution`, and both share one test file. `scripts/check.sh` verifies every solution under `-race`. Stubs are expected to fail (or, in 04 and 15, not compile) until implemented.
- Dependencies: `golang.org/x/sync` (errgroup), `go.uber.org/goleak`. The capstone's later stages add NATS, Prometheus and OpenTelemetry clients when reached.
- Teaching style: Python comparisons first, then Go idiom, then the senior/interview angle. One primary source per lesson.

# Go Resources

## Knowledge

- [A Tour of Go](https://go.dev/tour/): the official interactive introduction. Use it for syntax in the first sessions.
- [Effective Go](https://go.dev/doc/effective_go): the canonical idioms document. Partly dated (pre-generics, pre-modules) but still the base vocabulary.
- [The Go Programming Language Specification](https://go.dev/ref/spec): short and readable. The final word on semantics.
- [The Go Memory Model](https://go.dev/ref/mem): required for concurrency correctness (lesson 10).
- [Go blog](https://go.dev/blog/): primary-source articles on slices, errors, pipelines, context, generics, slog, PGO and more. Each lesson cites the relevant post.
- [A Guide to the Go Garbage Collector](https://go.dev/doc/gc-guide): how GOGC and GOMEMLIMIT actually work (lesson 18).
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments): what Go reviewers flag. The base of the checklist reference.
- [Google Go Style Guide](https://google.github.io/styleguide/go/) and [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md): the two most-cited production style guides.
- Book: _Learning Go_, 2nd ed. (Jon Bodner, O'Reilly): the best single book for an experienced developer coming from another language.
- Book: _Concurrency in Go_ (Katherine Cox-Buday, O'Reilly): patterns and pitfalls in depth (lessons 8–13).
- Book: _100 Go Mistakes and How to Avoid Them_ (Teiva Harsanyi, Manning): the practical catalogue of bugs senior reviewers catch. Good interview prep.
- Book: _Efficient Go_ (Bartłomiej Płotka, O'Reilly): performance methodology, benchmarks and profiling (lesson 18).
- [Understanding Real-World Concurrency Bugs in Go](https://songlh.github.io/paper/go-study.pdf) (ASPLOS 2019): empirical study behind lesson 13.
- Talks: [Concurrency is not Parallelism](https://go.dev/blog/waza-talk) and [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs) (Rob Pike). [Go Proverbs](https://go-proverbs.github.io/) (Rob Pike).

## Practice

- [Exercism Go track](https://exercism.org/tracks/go): small problems with optional mentor review. Good for live-coding fluency.
- [Go by Example](https://gobyexample.com/): quick annotated snippets for lookup.

## Wisdom (Communities)

- [Gophers Slack](https://invite.slack.golangbridge.org/): `#newbies`, `#performance`, `#concurrency` channels.
- [r/golang](https://reddit.com/r/golang): active. Useful for "how do teams actually do X" questions.
- [Go Forum](https://forum.golangbridge.org/): slower, longer-form Q&A.

## Gaps

- No curated video course. The official docs plus the books above cover the material better; revisit if a visual walkthrough of the scheduler or GC is needed.
- Broker internals, system design, trading-domain knowledge and C/C++ are deliberately left to sibling tracks (see `MISSION.md`).

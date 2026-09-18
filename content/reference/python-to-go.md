---
title: Python → Go cheat sheet
description: side-by-side translations of everyday Python idioms into Go.
source: Compressed from lessons 1–7. Source of truth: <a href="https://go.dev/ref/spec">The Go Programming Language Specification</a>.
---

## Basics

| Python | Go |
| --- | --- |
| `x = 5` | `x := 5` (declare) / `x = 5` (assign) |
| `x: int = 5` | `var x int = 5` |
| `None` | zero value; `nil` for pointers, slices, maps, chans, funcs, interfaces |
| `def f(a, b=2):` | no defaults: options struct or functional options |
| `def f(*args):` | `func f(args ...int)`; call `f(s...)` |
| `return a, b` (tuple) | `return a, b` (multiple results): `func f() (int, error)` |
| `lambda x: x*2` | `func(x int) int { return x * 2 }` |
| `f"{name} {n:.2f}"` | `fmt.Sprintf("%s %.2f", name, n)`; `%v` any value, `%+v` with field names, `%q` quoted, `%T` type |
| `if x in d:` | `if _, ok := d[x]; ok {` |
| `for i, v in enumerate(s):` | `for i, v := range s {` |
| `for k, v in d.items():` | `for k, v := range m {` (random order) |
| `for i in range(10):` | `for i := range 10 {` (Go 1.22+) |
| `while cond:` | `for cond {` |
| `x if c else y` | no ternary: `if`/`else` |
| `with open(p) as f:` | `f, err := os.Open(p); if err != nil {...}; defer f.Close()` |
| `try/except/finally` | `if err != nil { return ..., fmt.Errorf("ctx: %w", err) }`, `defer` |
| `raise ValueError` | `return errors.New(...)` / a sentinel / a typed error |
| `isinstance(x, T)` | `v, ok := x.(T)` or a type switch |

## Collections

| Python | Go |
| --- | --- |
| `[]`, `list()` | `var s []int` (nil, usable) / `make([]int, 0, n)` |
| `s.append(v)` | `s = append(s, v)`, always reassign |
| `s[1:3]` (copy) | `s[1:3]` (**shared view**); copy with `slices.Clone` |
| `len(s)` | `len(s)`, `cap(s)` |
| `sorted(s)` | `slices.Sorted(...)` / `slices.Sort(s)` (in place) |
| `sorted(s, key=f)` | `slices.SortFunc(s, func(a, b T) int { return cmp.Compare(...) })` |
| `v in s` | `slices.Contains(s, v)` |
| `{}`, `dict()` | `map[string]int{}` / `make(map[string]int, n)` |
| `d.get(k, 0)` | `d[k]` (missing → zero value) |
| `set()` | `map[T]struct{}` |
| `list(d.keys())` | `slices.Collect(maps.Keys(m))` |
| `collections.deque` | slice, or `container/list` |
| `heapq` | `container/heap` |
| `"".join(parts)` | `strings.Join(parts, "")`; build with `strings.Builder` |
| `s.split()` | `strings.Fields(s)` |
| `len("héllo") == 5` | `len("héllo") == 6` (bytes); `utf8.RuneCountInString` |

## Types and OOP

| Python | Go |
| --- | --- |
| `class A:` | `type A struct { ... }` |
| `__init__` | `func NewA(...) *A` |
| `self.x` | receiver: `func (a *A) M() { a.x }` |
| `_private` | lower-case identifier (package-private) |
| `__str__` | `func (a A) String() string` (fmt.Stringer) |
| inheritance | embedding + interfaces (composition) |
| ABC / Protocol | interface, satisfied implicitly |
| `@dataclass` | a struct literal: `A{X: 1}` |
| `Enum` | `type Side int; const ( Buy Side = iota + 1; Sell )` |
| `TypeVar` generics | `func F[T any](x T) T`, `cmp.Ordered`, `comparable` |

## Concurrency

| Python | Go |
| --- | --- |
| `threading.Thread(target=f).start()` | `go f()` |
| `asyncio.gather(*tasks)` | `sync.WaitGroup` / `errgroup.Group` |
| `asyncio.Queue` / `queue.Queue` | `make(chan T, n)` |
| `asyncio.wait_for(x, timeout)` | `context.WithTimeout` + `select` |
| `threading.Lock()` | `sync.Mutex` |
| `concurrent.futures.ThreadPoolExecutor(10)` | worker pool / `errgroup.SetLimit(10)` |
| `functools.lru_cache` | hand-rolled map + mutex, or `sync.OnceValue` |
| the GIL | none: goroutines run in parallel, so **use `-race`** |

## Tooling

| Python | Go |
| --- | --- |
| `pip install x` / `uv add x` | `go get example.com/x@latest` |
| `requirements.txt` / lock | `go.mod` + `go.sum` |
| `black` / `ruff` | `gofmt` / `go vet`, `staticcheck`, `golangci-lint` |
| `pytest` | `go test ./...` |
| `pytest -k name` | `go test -run 'TestName/sub' ./...` |
| `hypothesis` | native fuzzing: `go test -fuzz` |
| `cProfile` / py-spy | `pprof` |
| `python -m venv` | nothing: modules are per-project |
| PyInstaller | `go build` (static binary; `GOOS/GOARCH` to cross-compile) |

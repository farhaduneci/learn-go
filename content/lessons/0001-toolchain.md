---
title: Toolchain & modules
section: Part 1 · Go core
status: next
summary: go run/build/test/vet/fmt, modules, go.mod/go.sum, the shape of a Go file.
description: Go ships one tool that builds, tests, formats, vets and manages dependencies. Learning it first removes most of the friction that Python's tooling ecosystem trains people to expect.
source: <a href="https://go.dev/doc/tutorial/getting-started">Tutorial: Get started with Go</a> and <a href="https://go.dev/doc/modules/managing-dependencies">Managing dependencies</a> (go.dev, ~20 min). Then the first ten slides of <a href="https://go.dev/tour/welcome/1">A Tour of Go</a> for syntax.
quiz: Which command adds missing and removes unused module requirements? || go mod tidy
quiz: Which file records the checksums of every dependency? || go.sum
quiz: Go's equivalent of Python's black is which command? || gofmt
quiz: Which flag enables the data race detector? || -race
quiz: A test function's signature takes which parameter type? || *testing.T
---

## One binary, every job

In Python, a working setup often means pyenv + venv + pip/poetry/uv + black + ruff + mypy + pytest. Go folds nearly all of that into the `go` command:

| Job | Python | Go |
| --- | --- | --- |
| Run a script | `python main.py` | `go run .` |
| Build an artifact | wheel / PyInstaller | `go build` → one static binary |
| Dependencies | `pyproject.toml` + lock file | `go.mod` + `go.sum` |
| Format | black | `gofmt` / `go fmt ./...` (no options, no debates) |
| Lint | ruff / pylint | `go vet ./...` (+ `staticcheck`, `golangci-lint`) |
| Type check | mypy | the compiler |
| Test | pytest | `go test ./...` |
| Docs | Sphinx | `go doc`, pkg.go.dev from comments |

`./...` means "this directory and every package below it". It shows up everywhere.

## Modules and packages

A **module** is a versioned unit of code with a `go.mod` at its root. A **package** is a directory of `.go` files that share one `package` clause. A module contains packages.

```sh
mkdir hello && cd hello
go mod init github.com/example/hello   # module path = where it would be fetched from
```

```go
// main.go
package main // package main + func main = an executable

import "fmt"

func main() {
	fmt.Println("hello")
}
```

```sh
go run .                     # compile to a temp dir and run
go build -o bin/hello .      # static binary, no interpreter needed at deploy time
GOOS=linux GOARCH=amd64 go build -o bin/hello-linux .   # cross-compile, no toolchain setup
```

Adding a dependency is just importing it and running `go mod tidy`, or explicitly `go get example.com/pkg@v1.2.3`. `go.sum` pins cryptographic hashes; commit both files. There is no global site-packages and no virtualenv. Each module's dependency graph is resolved with **minimal version selection**: the oldest version that satisfies every requirement, which makes builds reproducible without a lock-file solver.

## The shape of a Go file

```go
package pricing

import (
	"errors"
	"fmt"
)

// Exported names start with an upper-case letter. This comment is the doc.
func Hello(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

func helper() {} // lower-case: private to the package
```

Differences a Python developer hits in the first hour:

- **Braces, not indentation**, and the opening brace must be on the same line.
- **No unused imports or variables.** They're compile errors, not warnings.
- **`:=` declares and assigns** (`x := 5`). `=` only assigns to an existing variable.
- **One loop keyword:** `for`. `for i := 0; i < n; i++`, `for cond {}`, `for {}`, `for i, v := range s`, and since Go 1.22, `for i := range 10`.
- **Multiple return values** instead of tuples: `func Stats(s []int) (lo, hi int, ok bool)`.
- **Variadic** `func Sum(nums ...int)` is `*args`; call with a slice as `Sum(s...)`.

## Tests live next to the code

`foo_test.go` files in the same directory are compiled only by `go test`. No framework is needed:

```go
func TestHello(t *testing.T) {
	if got := Hello(""); got != "Hello, world!" {
		t.Errorf("Hello(\"\") = %q, want %q", got, "Hello, world!")
	}
}
```

```sh
go test ./...                   # all packages
go test -run TestHello -v ./... # one test, verbose
go test -race ./...             # with the race detector (use it always from lesson 8 on)
go test -cover ./...            # coverage
```

The "got X, want Y" message format is a Go convention. Lesson 14 covers testing in depth.

## Editor setup

Install **gopls**, the official language server (VS Code's Go extension does it automatically). It gives completion, go-to-definition, rename, and format + organize imports on save. With format-on-save on, style never needs to be thought about again, and every Go codebase reads the same. That's a large part of why Go teams move fast.

## Exercise

[`exercises/01-toolchain`](../exercises/01-toolchain/README.md): three warm-up functions, plus running `go test`, `go vet` and `gofmt -l` on them. Also create a throwaway module outside this repo with `go mod init`, add a dependency with `go get`, and read the resulting `go.mod`.

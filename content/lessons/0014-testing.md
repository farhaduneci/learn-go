---
title: Testing like a Go team
section: Part 3 · Senior Go
summary: Table-driven tests, subtests, fakes over mocks, fuzzing, benchmarks with benchstat.
description: Go's testing culture is standard-library first. Table-driven tests, subtests, fakes built on small interfaces, fuzzing and benchmarks are what reviewers expect in every PR.
source: <a href="https://go.dev/wiki/TableDrivenTests">TableDrivenTests</a> (Go wiki), <a href="https://go.dev/doc/tutorial/fuzz">Tutorial: Getting started with fuzzing</a>, and the <a href="https://pkg.go.dev/testing">testing package docs</a>. Style: <a href="https://google.github.io/styleguide/go/decisions#useful-test-failures">Google style guide — Useful test failures</a>.
quiz: Which method runs a named subtest? || t.Run
quiz: Which method marks a function as a test helper so failures report the caller's line? || t.Helper
quiz: Which benchmark loop form was added in Go 1.24? || b.Loop
quiz: Which flag prints allocations per operation in benchmarks? || -benchmem
quiz: Failing fuzz inputs are saved under which directory? || testdata/fuzz
---

## Table-driven tests

The default shape of a Go test:

```go
func TestParsePrice(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		{name: "whole", in: "12", want: 1200},
		{name: "two decimals", in: "12.34", want: 1234},
		{name: "letters", in: "abc", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePrice(tt.in)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParsePrice(%q) err = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("ParsePrice(%q) = %d, want %d", tt.in, got, tt.want)
			}
		})
	}
}
```

- Adding a case is one line, and each subtest runs by name: `go test -run 'TestParsePrice/letters'`.
- Failure messages say **function(input) = got, want X**, so a failure can be understood without opening the test.
- `t.Errorf` records a failure and continues. `t.Fatalf` stops the current (sub)test. Use Fatal when later checks would be meaningless.
- `t.Helper()` in assertion helpers makes failures point at the caller.
- `t.Parallel()` runs subtests concurrently. `t.Cleanup(f)`, `t.TempDir()`, `t.Setenv()` and `t.Context()` (1.24) handle setup and teardown.

The standard library is enough. Many teams add `github.com/google/go-cmp/cmp` for readable struct diffs (`cmp.Diff(want, got)`). `testify` is popular, and some teams avoid it. Match the codebase.

## Fakes over mocks

Go's implicit interfaces make hand-written fakes trivial. Define the **small interface the code needs** (lesson 15), then:

```go
type fakeNotifier struct{ fills []Order }
func (f *fakeNotifier) Filled(o Order) { f.fills = append(f.fills, o) }

type fakeClock struct{ t time.Time }
func (c *fakeClock) Now() time.Time { return c.t }
```

Fakes are plain code: readable, debuggable, and checked for state rather than call sequences. Generated mocks (`gomock`, `mockery`) help for large interfaces, but a large interface is usually the actual problem.

## Fuzzing

Fuzzing (Go 1.18+) generates inputs to break **properties**:

```go
func FuzzRoundTrip(f *testing.F) {
	f.Add(int64(1234))                 // seed corpus
	f.Fuzz(func(t *testing.T, c int64) {
		got, err := ParsePrice(FormatPrice(c))
		if err != nil || got != c {
			t.Errorf("round trip %d -> %q -> %d, %v", c, FormatPrice(c), got, err)
		}
	})
}
```

```sh
go test -run '^$' -fuzz FuzzRoundTrip -fuzztime 30s ./pricing/
```

Good properties to fuzz: round-trips (parse ↔ format, encode ↔ decode), "never panics on any input", and agreement between a fast implementation and a simple reference one. Parsers of external input, like FIX messages or JSON orders, should always have a fuzz test.

## Benchmarks

```go
func BenchmarkParsePrice(b *testing.B) {
	for b.Loop() {                   // Go 1.24+: handles timer reset and keeps results alive
		ParsePrice("12345.67")
	}
}
```

```sh
go test -run '^$' -bench . -benchmem -count 10 ./... > new.txt
benchstat old.txt new.txt             # statistically compares runs
```

A single benchmark run is noise. Use `-count` of 6–10 and `benchstat`. Lesson 18 goes further.

## Integration tests and golden files

- Guard slow tests with a build tag (`//go:build integration`) or `testing.Short()`.
- Use `testcontainers-go` to run real Postgres, Kafka or NATS in Docker from a test.
- For large outputs, compare against `testdata/*.golden` files with an `-update` flag to regenerate them.
- `httptest.NewRecorder` and `httptest.NewServer` test HTTP code without a network (lesson 17).

## Exercise

[`exercises/14-testing`](../exercises/14-testing/README.md): the code is given, the tests are the task. Write table tests, a fuzz round-trip and a benchmark, then let the fuzzer find the planted bug.

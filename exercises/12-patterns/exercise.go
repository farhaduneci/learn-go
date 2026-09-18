//go:build !solution

// Package patterns practises worker pools, bounded parallelism and throttling.
package patterns

import "context"

// WorkerPool processes jobs with n workers and emits results in any order.
func WorkerPool[J, R any](ctx context.Context, jobs <-chan J, n int, fn func(J) R) <-chan R {
	// TODO: implement
	return make(chan R)
}

// ProcessAll applies fn to every item with at most limit calls in flight,
// returning results in input order or the first error.
func ProcessAll[T, R any](ctx context.Context, items []T, limit int, fn func(context.Context, T) (R, error)) ([]R, error) {
	// TODO: implement with errgroup
	out := make([]R, len(items))
	for i, it := range items {
		r, err := fn(ctx, it)
		if err != nil {
			return nil, err
		}
		out[i] = r
	}
	return out, nil
}

// Throttle forwards values from in at most perSecond per second.
// The output closes when in closes or ctx is done.
func Throttle[T any](ctx context.Context, in <-chan T, perSecond int) <-chan T {
	// TODO: implement with time.Ticker
	return in
}

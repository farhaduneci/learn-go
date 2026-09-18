//go:build solution

// Package patterns practises worker pools, bounded parallelism and throttling.
package patterns

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// WorkerPool processes jobs with n workers and emits results in any order.
func WorkerPool[J, R any](ctx context.Context, jobs <-chan J, n int, fn func(J) R) <-chan R {
	out := make(chan R)
	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case j, ok := <-jobs:
					if !ok {
						return
					}
					select {
					case out <- fn(j):
					case <-ctx.Done(): // never block forever on a send nobody reads
						return
					}
				}
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// ProcessAll applies fn to every item with at most limit calls in flight,
// returning results in input order or the first error.
func ProcessAll[T, R any](ctx context.Context, items []T, limit int, fn func(context.Context, T) (R, error)) ([]R, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(limit) // g.Go blocks while limit goroutines are running
	out := make([]R, len(items))
	for i, it := range items {
		g.Go(func() error {
			r, err := fn(ctx, it)
			if err != nil {
				return err
			}
			out[i] = r
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil, err
	}
	return out, nil
}

// Throttle forwards values from in at most perSecond per second.
// The output closes when in closes or ctx is done.
func Throttle[T any](ctx context.Context, in <-chan T, perSecond int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		tick := time.NewTicker(time.Second / time.Duration(perSecond))
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-tick.C:
				case <-ctx.Done():
					return
				}
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

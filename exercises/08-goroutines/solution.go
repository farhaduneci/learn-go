//go:build solution

package goroutines

import "sync"

// ParallelSum sums nums using the given number of goroutines.
func ParallelSum(nums []int, workers int) int {
	if workers < 1 {
		workers = 1
	}
	partial := make([]int, workers) // one slot per goroutine: no shared writes
	size := (len(nums) + workers - 1) / workers

	var wg sync.WaitGroup
	for w := range workers {
		lo := min(w*size, len(nums))
		hi := min(lo+size, len(nums))
		wg.Go(func() {
			for _, n := range nums[lo:hi] {
				partial[w] += n
			}
		})
	}
	wg.Wait()

	total := 0
	for _, p := range partial {
		total += p
	}
	return total
}

// FetchAll fetches every symbol concurrently and returns quotes in input order.
func FetchAll(symbols []string, fetch Fetcher) []Quote {
	out := make([]Quote, len(symbols))
	var wg sync.WaitGroup
	for i, s := range symbols {
		wg.Go(func() { // Go 1.22+: i and s are per-iteration, safe to capture
			p, err := fetch(s)
			out[i] = Quote{Symbol: s, Price: p, Err: err}
		})
	}
	wg.Wait()
	return out
}

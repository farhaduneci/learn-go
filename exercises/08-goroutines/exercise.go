//go:build !solution

package goroutines

// ParallelSum sums nums using the given number of goroutines.
func ParallelSum(nums []int, workers int) int {
	// TODO: implement with goroutines
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// FetchAll fetches every symbol concurrently and returns quotes in input order.
func FetchAll(symbols []string, fetch Fetcher) []Quote {
	// TODO: this runs sequentially — make it concurrent.
	out := make([]Quote, 0, len(symbols))
	for _, s := range symbols {
		p, err := fetch(s)
		out = append(out, Quote{Symbol: s, Price: p, Err: err})
	}
	return out
}

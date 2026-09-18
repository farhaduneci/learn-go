// Package perf practises measuring and removing allocations.
package perf

// Trade is one execution.
type Trade struct {
	Symbol string
	Price  int64 // cents
	Qty    int64
}

// Package goroutines practises starting goroutines and waiting for them.
package goroutines

// Quote is the result of fetching one symbol.
type Quote struct {
	Symbol string
	Price  int64
	Err    error
}

// Fetcher looks up the price of a symbol. It may be slow.
type Fetcher func(symbol string) (int64, error)

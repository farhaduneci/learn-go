//go:build !solution

package orderbook

import "context"

// Book is a single-symbol limit order book. It is NOT safe for concurrent use;
// the Engine is the only writer.
type Book struct {
	// TODO (stage 1): bids and asks as price levels, each a FIFO queue of resting
	// orders, plus an index from order ID to its resting position for O(1) cancel.
}

// NewBook returns an empty book.
func NewBook() *Book {
	// TODO (stage 1)
	return &Book{}
}

// Submit matches o against the opposite side with price-time priority and
// returns the resulting trades. A Limit remainder rests; a Market remainder is dropped.
func (b *Book) Submit(o Order) ([]Trade, error) {
	// TODO (stage 1)
	return nil, nil
}

// Cancel removes a resting order and reports whether it was found.
func (b *Book) Cancel(id uint64) bool {
	// TODO (stage 1)
	return false
}

// Depth returns up to n aggregated levels on side s, best price first.
func (b *Book) Depth(s Side, n int) []Level {
	// TODO (stage 1)
	return nil
}

// Engine owns one Book per symbol and serialises all access through one goroutine.
type Engine struct {
	// TODO (stage 2)
}

// NewEngine starts the engine's event loop.
func NewEngine() *Engine {
	// TODO (stage 2)
	return &Engine{}
}

// Submit sends o to the book for symbol and waits for the trades.
func (e *Engine) Submit(ctx context.Context, symbol string, o Order) ([]Trade, error) {
	// TODO (stage 2)
	return nil, nil
}

// Cancel cancels a resting order on symbol's book.
func (e *Engine) Cancel(ctx context.Context, symbol string, id uint64) (bool, error) {
	// TODO (stage 2)
	return false, nil
}

// Depth returns up to n levels of symbol's book on side s.
func (e *Engine) Depth(ctx context.Context, symbol string, s Side, n int) ([]Level, error) {
	// TODO (stage 2)
	return nil, nil
}

// Close stops the event loop. Later calls return ErrClosed.
func (e *Engine) Close() {
	// TODO (stage 2)
}

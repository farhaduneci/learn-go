// Package refactor turns an untestable legacy function into a testable service.
package refactor

import "errors"

// Errors reported per line.
var (
	ErrBadFormat            = errors.New("bad format")
	ErrBadQty               = errors.New("bad qty")
	ErrBadSide              = errors.New("bad side")
	ErrMarketClosed         = errors.New("market closed")
	ErrInsufficientPosition = errors.New("insufficient position")
)

// Side is buy or sell.
type Side int

// Sides.
const (
	Buy Side = iota + 1
	Sell
)

// Order is one parsed instruction.
type Order struct {
	Side   Side
	Symbol string
	Qty    int64
}

// Result is the outcome for one input line (1-based). Err is nil on a fill.
type Result struct {
	Line int
	Err  error
}

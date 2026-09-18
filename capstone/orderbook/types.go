// Package orderbook is a limit order book and a single-writer matching engine.
package orderbook

import "errors"

// Side is the direction of an order.
type Side uint8

// Sides.
const (
	Buy Side = iota + 1
	Sell
)

// Kind is the order type.
type Kind uint8

// Kinds.
const (
	Limit  Kind = iota + 1 // rests on the book at Price if not fully filled
	Market                 // takes whatever liquidity exists; any remainder is cancelled
)

// Order is an instruction to trade. Price is in cents and ignored for Market orders.
type Order struct {
	ID    uint64
	Side  Side
	Kind  Kind
	Price int64
	Qty   int64
}

// Trade is one fill between an incoming (taker) and a resting (maker) order.
// It always executes at the maker's price.
type Trade struct {
	TakerID uint64
	MakerID uint64
	Price   int64
	Qty     int64
}

// Level is aggregated resting quantity at one price.
type Level struct {
	Price int64
	Qty   int64
}

// Errors returned by Submit.
var (
	ErrInvalidOrder = errors.New("orderbook: invalid order")
	ErrDuplicateID  = errors.New("orderbook: duplicate order id")
	ErrClosed       = errors.New("orderbook: engine closed")
)

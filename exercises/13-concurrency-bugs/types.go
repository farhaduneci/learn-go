// Package bugs contains classic concurrency bugs to diagnose and fix.
package bugs

import (
	"errors"
	"sync"
)

// Quote is a price from one venue.
type Quote struct {
	Venue string
	Price int64
}

// Account is a lockable balance.
type Account struct {
	ID      int
	mu      sync.Mutex
	Balance int64
}

// ErrInvalid is returned for an invalid order quantity.
var ErrInvalid = errors.New("invalid quantity")

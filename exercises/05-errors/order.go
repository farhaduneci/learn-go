package errs

import (
	"errors"
	"fmt"
)

// ErrInsufficientFunds is returned when an order costs more than the available cash.
var ErrInsufficientFunds = errors.New("insufficient funds")

// Order is a request to trade.
type Order struct {
	ID     string
	Symbol string
	Side   string // "buy" or "sell"
	Qty    int64
	Price  int64 // cents
}

// ValidationError reports an invalid order field.
type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

// Account holds cash available for trading.
type Account struct {
	Cash int64
}

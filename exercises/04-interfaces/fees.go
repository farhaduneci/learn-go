package interfaces

// FeeModel computes the fee for a trade of the given notional value (in cents).
type FeeModel interface {
	Fee(notional int64) int64
}

// Ledger records a running total.
type Ledger struct {
	Total int64
}

// Post adds amount to the total.
func (l *Ledger) Post(amount int64) { l.Total += amount }

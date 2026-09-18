//go:build solution

package ledger

import (
	"maps"
	"slices"
)

// Ledger holds account balances in cents.
type Ledger struct {
	balances map[string]int64
}

// New returns an empty ledger.
func New() *Ledger {
	return &Ledger{balances: make(map[string]int64)}
}

// Post moves amount from one account to another.
func (l *Ledger) Post(from, to string, amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if from == to {
		return ErrSameAccount
	}
	l.balances[from] -= amount
	l.balances[to] += amount
	return nil
}

// Balance returns the balance of account.
func (l *Ledger) Balance(account string) int64 {
	return l.balances[account]
}

// Accounts returns every account name, sorted.
func (l *Ledger) Accounts() []string {
	return slices.Sorted(maps.Keys(l.balances))
}

// Sum returns the sum of all balances (always zero for a consistent ledger).
func (l *Ledger) Sum() int64 {
	var sum int64
	for _, b := range l.balances {
		sum += b
	}
	return sum
}

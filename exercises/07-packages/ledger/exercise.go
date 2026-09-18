//go:build !solution

package ledger

// Ledger holds account balances in cents.
type Ledger struct {
	// TODO: unexported fields
}

// New returns an empty ledger.
func New() *Ledger {
	// TODO: implement
	return &Ledger{}
}

// Post moves amount from one account to another.
func (l *Ledger) Post(from, to string, amount int64) error {
	// TODO: implement
	return nil
}

// Balance returns the balance of account.
func (l *Ledger) Balance(account string) int64 {
	// TODO: implement
	return 0
}

// Accounts returns every account name, sorted.
func (l *Ledger) Accounts() []string {
	// TODO: implement
	return nil
}

// Sum returns the sum of all balances (always zero for a consistent ledger).
func (l *Ledger) Sum() int64 {
	// TODO: implement
	return -1
}

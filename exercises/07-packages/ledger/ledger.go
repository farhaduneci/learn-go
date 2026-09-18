// Package ledger is a minimal double-entry ledger.
//
// Every posting debits one account and credits another by the same amount,
// so the sum of all balances is always zero.
package ledger

import "errors"

// Sentinel errors returned by Post.
var (
	ErrInvalidAmount = errors.New("ledger: amount must be positive")
	ErrSameAccount   = errors.New("ledger: from and to must differ")
)

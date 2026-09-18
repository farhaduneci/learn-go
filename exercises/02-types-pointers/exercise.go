//go:build !solution

// Package types practises value types, pointer receivers and closures.
package types

// Cents is an amount of money in hundredths of a dollar.
type Cents int64

// String formats c as "$12.34" or "-$0.05".
func (c Cents) String() string {
	// TODO: implement
	return ""
}

// Account is a cash account.
type Account struct {
	ID      string
	Balance Cents
}

// NewAccount returns a new account with a zero balance.
func NewAccount(id string) *Account {
	// TODO: implement
	return nil
}

// Deposit adds amount to the balance.
// TODO: this receiver is a value — fix it so the change sticks.
func (a Account) Deposit(amount Cents) {
	a.Balance += amount
}

// Withdraw removes amount if funds allow and reports whether it did.
func (a *Account) Withdraw(amount Cents) bool {
	// TODO: implement
	return false
}

// Swap swaps the values a and b point to.
func Swap(a, b *int) {
	// TODO: implement
}

// NewCounter returns a function that returns 1, 2, 3, ... on successive calls.
func NewCounter() func() int {
	// TODO: implement
	return func() int { return 0 }
}

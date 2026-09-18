//go:build solution

// Package types practises value types, pointer receivers and closures.
package types

import "fmt"

// Cents is an amount of money in hundredths of a dollar.
type Cents int64

// String formats c as "$12.34" or "-$0.05".
func (c Cents) String() string {
	sign := ""
	if c < 0 {
		sign = "-"
		c = -c
	}
	return fmt.Sprintf("%s$%d.%02d", sign, c/100, c%100)
}

// Account is a cash account.
type Account struct {
	ID      string
	Balance Cents
}

// NewAccount returns a new account with a zero balance.
func NewAccount(id string) *Account {
	return &Account{ID: id}
}

// Deposit adds amount to the balance.
func (a *Account) Deposit(amount Cents) {
	a.Balance += amount
}

// Withdraw removes amount if funds allow and reports whether it did.
func (a *Account) Withdraw(amount Cents) bool {
	if amount > a.Balance {
		return false
	}
	a.Balance -= amount
	return true
}

// Swap swaps the values a and b point to.
func Swap(a, b *int) {
	*a, *b = *b, *a
}

// NewCounter returns a function that returns 1, 2, 3, ... on successive calls.
func NewCounter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

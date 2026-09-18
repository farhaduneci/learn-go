//go:build solution

package errs

import "fmt"

// Validate reports the first invalid field of o, or nil.
func Validate(o Order) error {
	switch {
	case o.Symbol == "":
		return &ValidationError{Field: "Symbol", Reason: "empty"}
	case o.Qty <= 0:
		return &ValidationError{Field: "Qty", Reason: "must be positive"}
	case o.Side != "buy" && o.Side != "sell":
		return &ValidationError{Field: "Side", Reason: fmt.Sprintf("unknown side %q", o.Side)}
	}
	return nil // untyped nil: the error interface is truly nil
}

// Place validates o and checks funds; errors are wrapped with the order ID.
func (a *Account) Place(o Order) error {
	if err := Validate(o); err != nil {
		return fmt.Errorf("place order %s: %w", o.ID, err)
	}
	if o.Qty*o.Price > a.Cash {
		return fmt.Errorf("place order %s: %w", o.ID, ErrInsufficientFunds)
	}
	a.Cash -= o.Qty * o.Price
	return nil
}

// SafeDiv divides a by b, converting a panic into an error.
func SafeDiv(a, b int) (q int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("safe div: %v", r)
		}
	}()
	return a / b, nil
}

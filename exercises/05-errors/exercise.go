//go:build !solution

package errs

// Validate reports the first invalid field of o, or nil.
func Validate(o Order) error {
	// TODO: implement. Return a *ValidationError for a bad field.
	return nil
}

// Place validates o and checks funds; errors are wrapped with the order ID.
func (a *Account) Place(o Order) error {
	// TODO: implement
	return nil
}

// SafeDiv divides a by b, converting a panic into an error.
func SafeDiv(a, b int) (q int, err error) {
	// TODO: add a deferred recover
	return a / b, nil
}

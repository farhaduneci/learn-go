//go:build !solution

package messaging

// Dedupe runs next at most once per successful message ID.
func Dedupe(next Handler) Handler {
	// TODO
	return next
}

// WithDeadLetter nacks failures until maxAttempts, then sends to dlq and acks.
func WithDeadLetter(next Handler, maxAttempts int, dlq func(Message, error)) Handler {
	// TODO
	return next
}

// Dispatcher routes messages to per-key-ordered partitions.
type Dispatcher struct {
	// TODO
}

// NewDispatcher starts partitions goroutines, each running h sequentially.
func NewDispatcher(partitions int, h Handler) *Dispatcher {
	// TODO
	return &Dispatcher{}
}

// Dispatch enqueues m on the partition for m.Key.
func (d *Dispatcher) Dispatch(m Message) {
	// TODO
}

// Close stops accepting messages, drains every partition and waits.
func (d *Dispatcher) Close() {
	// TODO
}

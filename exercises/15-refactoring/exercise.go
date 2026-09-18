//go:build !solution

package refactor

import "time"

// ParseOrder parses "buy aapl 10" into an Order.
func ParseOrder(line string) (Order, error) {
	// TODO: extract from ProcessOrders
	return Order{}, nil
}

// PositionStore holds net positions per symbol.
type PositionStore interface {
	// TODO: the two methods Service needs, and no more
}

// Notifier is told about every fill.
type Notifier interface {
	// TODO
}

// MemStore is an in-memory PositionStore.
type MemStore struct{}

// NewMemStore returns an empty MemStore.
func NewMemStore() *MemStore { return &MemStore{} }

// Service processes order batches.
type Service struct{}

// NewService wires a Service from its dependencies.
func NewService(store PositionStore, now func() time.Time, n Notifier) *Service {
	// TODO
	return &Service{}
}

// Process handles every non-blank line of input.
func (s *Service) Process(input string) []Result {
	// TODO: move the ProcessOrders logic here
	return nil
}

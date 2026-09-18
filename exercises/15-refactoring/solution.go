//go:build solution

package refactor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseOrder parses "buy aapl 10" into an Order.
func ParseOrder(line string) (Order, error) {
	parts := strings.Fields(line)
	if len(parts) != 3 {
		return Order{}, fmt.Errorf("%w: %q", ErrBadFormat, line)
	}
	var side Side
	switch strings.ToUpper(parts[0]) {
	case "BUY":
		side = Buy
	case "SELL":
		side = Sell
	default:
		return Order{}, fmt.Errorf("%w: %q", ErrBadSide, parts[0])
	}
	qty, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil || qty <= 0 {
		return Order{}, fmt.Errorf("%w: %q", ErrBadQty, parts[2])
	}
	return Order{Side: side, Symbol: strings.ToUpper(parts[1]), Qty: qty}, nil
}

// PositionStore holds net positions per symbol.
type PositionStore interface {
	Position(symbol string) int64
	Add(symbol string, delta int64)
}

// Notifier is told about every fill.
type Notifier interface {
	Filled(o Order)
}

// MemStore is an in-memory PositionStore.
type MemStore struct {
	pos map[string]int64
}

// NewMemStore returns an empty MemStore.
func NewMemStore() *MemStore { return &MemStore{pos: make(map[string]int64)} }

// Position implements PositionStore.
func (m *MemStore) Position(symbol string) int64 { return m.pos[symbol] }

// Add implements PositionStore.
func (m *MemStore) Add(symbol string, delta int64) { m.pos[symbol] += delta }

// Service processes order batches.
type Service struct {
	store  PositionStore
	now    func() time.Time
	notify Notifier
}

// NewService wires a Service from its dependencies.
func NewService(store PositionStore, now func() time.Time, n Notifier) *Service {
	return &Service{store: store, now: now, notify: n}
}

// Process handles every non-blank line of input.
func (s *Service) Process(input string) []Result {
	var results []Result
	for i, line := range strings.Split(input, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		results = append(results, Result{Line: i + 1, Err: s.processLine(line)})
	}
	return results
}

func (s *Service) processLine(line string) error {
	o, err := ParseOrder(line)
	if err != nil {
		return err
	}
	if h := s.now().Hour(); h < 9 || h >= 16 {
		return ErrMarketClosed
	}
	delta := o.Qty
	if o.Side == Sell {
		if s.store.Position(o.Symbol) < o.Qty {
			return ErrInsufficientPosition
		}
		delta = -o.Qty
	}
	s.store.Add(o.Symbol, delta)
	s.notify.Filled(o)
	return nil
}

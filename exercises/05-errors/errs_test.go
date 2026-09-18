package errs

import (
	"errors"
	"strings"
	"testing"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name  string
		order Order
		field string // "" means valid
	}{
		{"valid", Order{Symbol: "AAPL", Side: "buy", Qty: 1}, ""},
		{"no symbol", Order{Side: "buy", Qty: 1}, "Symbol"},
		{"zero qty", Order{Symbol: "AAPL", Side: "buy"}, "Qty"},
		{"bad side", Order{Symbol: "AAPL", Side: "hold", Qty: 1}, "Side"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.order)
			if tt.field == "" {
				if err != nil {
					t.Fatalf("Validate = %v (%T), want nil", err, err)
				}
				return
			}
			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("Validate = %v, want *ValidationError", err)
			}
			if ve.Field != tt.field {
				t.Errorf("Field = %q, want %q", ve.Field, tt.field)
			}
		})
	}
}

func TestPlace(t *testing.T) {
	a := &Account{Cash: 10_000}

	err := a.Place(Order{ID: "o1", Symbol: "AAPL", Side: "buy", Qty: 10, Price: 5_000})
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf("Place = %v, want wrapped ErrInsufficientFunds", err)
	}
	if !strings.Contains(err.Error(), "o1") {
		t.Errorf("error %q lacks order ID context", err)
	}

	err = a.Place(Order{ID: "o2", Symbol: "", Side: "buy", Qty: 1, Price: 1})
	var ve *ValidationError
	if !errors.As(err, &ve) || !strings.Contains(err.Error(), "o2") {
		t.Fatalf("Place = %v, want wrapped *ValidationError mentioning o2", err)
	}

	if err := a.Place(Order{ID: "o3", Symbol: "AAPL", Side: "buy", Qty: 2, Price: 1_000}); err != nil {
		t.Fatalf("Place = %v, want nil", err)
	}
	if a.Cash != 8_000 {
		t.Errorf("Cash = %d, want 8000", a.Cash)
	}
}

func TestSafeDiv(t *testing.T) {
	if q, err := SafeDiv(7, 2); err != nil || q != 3 {
		t.Errorf("SafeDiv(7, 2) = %d, %v", q, err)
	}
	if _, err := SafeDiv(1, 0); err == nil {
		t.Error("SafeDiv(1, 0) err = nil, want error")
	}
}

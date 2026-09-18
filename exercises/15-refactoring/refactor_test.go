package refactor

import (
	"errors"
	"testing"
	"time"
)

type fakeNotifier struct{ fills []Order }

func (f *fakeNotifier) Filled(o Order) { f.fills = append(f.fills, o) }

func at(hour int) func() time.Time {
	return func() time.Time { return time.Date(2026, 1, 5, hour, 30, 0, 0, time.UTC) }
}

func TestParseOrder(t *testing.T) {
	tests := []struct {
		in      string
		want    Order
		wantErr error
	}{
		{"buy aapl 10", Order{Buy, "AAPL", 10}, nil},
		{"SELL tsla 3", Order{Sell, "TSLA", 3}, nil},
		{"buy aapl", Order{}, ErrBadFormat},
		{"buy aapl -1", Order{}, ErrBadQty},
		{"buy aapl x", Order{}, ErrBadQty},
		{"hold aapl 1", Order{}, ErrBadSide},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := ParseOrder(tt.in)
			if !errors.Is(err, tt.wantErr) || got != tt.want {
				t.Errorf("ParseOrder(%q) = %+v, %v; want %+v, %v", tt.in, got, err, tt.want, tt.wantErr)
			}
		})
	}
}

func TestProcess(t *testing.T) {
	store := NewMemStore()
	n := &fakeNotifier{}
	svc := NewService(store, at(10), n)

	got := svc.Process("buy aapl 10\n\nsell aapl 4\nsell aapl 100\nnonsense\n")
	want := []Result{{1, nil}, {3, nil}, {4, ErrInsufficientPosition}, {5, ErrBadFormat}}
	if len(got) != len(want) {
		t.Fatalf("got %d results, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i].Line != want[i].Line || !errors.Is(got[i].Err, want[i].Err) {
			t.Errorf("result %d = %+v, want %+v", i, got[i], want[i])
		}
	}
	if p := store.Position("AAPL"); p != 6 {
		t.Errorf("AAPL position = %d, want 6", p)
	}
	if len(n.fills) != 2 {
		t.Errorf("notified %d fills, want 2", len(n.fills))
	}
}

func TestProcessMarketClosed(t *testing.T) {
	svc := NewService(NewMemStore(), at(20), &fakeNotifier{})
	got := svc.Process("buy aapl 1")
	if len(got) != 1 || !errors.Is(got[0].Err, ErrMarketClosed) {
		t.Errorf("Process at 20:30 = %+v, want ErrMarketClosed", got)
	}
}

func TestServicesAreIndependent(t *testing.T) {
	a, b := NewMemStore(), NewMemStore()
	NewService(a, at(10), &fakeNotifier{}).Process("buy aapl 5")
	if b.Position("AAPL") != 0 {
		t.Error("state leaked between stores (global map still in use?)")
	}
}

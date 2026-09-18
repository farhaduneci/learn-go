package orderbook

import (
	"errors"
	"slices"
	"testing"
)

func limit(id uint64, s Side, price, qty int64) Order {
	return Order{ID: id, Side: s, Kind: Limit, Price: price, Qty: qty}
}

func market(id uint64, s Side, qty int64) Order {
	return Order{ID: id, Side: s, Kind: Market, Qty: qty}
}

func mustSubmit(t *testing.T, b *Book, o Order) []Trade {
	t.Helper()
	trades, err := b.Submit(o)
	if err != nil {
		t.Fatalf("Submit(%+v) = %v", o, err)
	}
	return trades
}

func TestRestingOrdersAndDepth(t *testing.T) {
	b := NewBook()
	mustSubmit(t, b, limit(1, Buy, 100, 5))
	mustSubmit(t, b, limit(2, Buy, 101, 3))
	mustSubmit(t, b, limit(3, Buy, 100, 2))
	mustSubmit(t, b, limit(4, Sell, 105, 4))
	mustSubmit(t, b, limit(5, Sell, 103, 1))

	if got := b.Depth(Buy, 10); !slices.Equal(got, []Level{{101, 3}, {100, 7}}) {
		t.Errorf("bids = %v, want [{101 3} {100 7}] (best first)", got)
	}
	if got := b.Depth(Sell, 1); !slices.Equal(got, []Level{{103, 1}}) {
		t.Errorf("asks top = %v, want [{103 1}]", got)
	}
}

func TestPriceTimePriority(t *testing.T) {
	b := NewBook()
	mustSubmit(t, b, limit(1, Sell, 101, 5)) // worse price
	mustSubmit(t, b, limit(2, Sell, 100, 3)) // best price, earlier
	mustSubmit(t, b, limit(3, Sell, 100, 3)) // best price, later

	got := mustSubmit(t, b, limit(10, Buy, 101, 8))
	want := []Trade{
		{TakerID: 10, MakerID: 2, Price: 100, Qty: 3}, // best price first
		{TakerID: 10, MakerID: 3, Price: 100, Qty: 3}, // then time priority
		{TakerID: 10, MakerID: 1, Price: 101, Qty: 2}, // then next level, partial
	}
	if !slices.Equal(got, want) {
		t.Fatalf("trades =\n%v\nwant\n%v", got, want)
	}
	if got := b.Depth(Sell, 5); !slices.Equal(got, []Level{{101, 3}}) {
		t.Errorf("asks after = %v, want [{101 3}]", got)
	}
	if got := b.Depth(Buy, 5); len(got) != 0 {
		t.Errorf("fully filled taker should not rest: bids = %v", got)
	}
}

func TestLimitDoesNotCrossBeyondPrice(t *testing.T) {
	b := NewBook()
	mustSubmit(t, b, limit(1, Sell, 102, 5))
	if got := mustSubmit(t, b, limit(2, Buy, 101, 5)); len(got) != 0 {
		t.Fatalf("buy@101 traded against ask@102: %v", got)
	}
	if got := b.Depth(Buy, 1); !slices.Equal(got, []Level{{101, 5}}) {
		t.Errorf("bid should rest: %v", got)
	}
	got := mustSubmit(t, b, limit(3, Sell, 100, 2)) // aggressive sell trades at maker's 101
	if !slices.Equal(got, []Trade{{TakerID: 3, MakerID: 2, Price: 101, Qty: 2}}) {
		t.Errorf("trades = %v, want fill at maker price 101", got)
	}
}

func TestMarketOrder(t *testing.T) {
	b := NewBook()
	mustSubmit(t, b, limit(1, Buy, 99, 2))
	mustSubmit(t, b, limit(2, Buy, 98, 2))
	got := mustSubmit(t, b, market(3, Sell, 10))
	if len(got) != 2 || got[0].Price != 99 || got[1].Price != 98 {
		t.Fatalf("market sweep = %v", got)
	}
	if d := b.Depth(Sell, 5); len(d) != 0 {
		t.Errorf("market remainder must not rest: asks = %v", d)
	}
	if got := mustSubmit(t, b, market(4, Buy, 1)); len(got) != 0 {
		t.Errorf("market into empty book traded: %v", got)
	}
}

func TestCancel(t *testing.T) {
	b := NewBook()
	mustSubmit(t, b, limit(1, Buy, 100, 5))
	mustSubmit(t, b, limit(2, Buy, 100, 5))
	if !b.Cancel(1) {
		t.Fatal("Cancel(1) = false")
	}
	if b.Cancel(1) {
		t.Error("second Cancel(1) = true")
	}
	if got := b.Depth(Buy, 1); !slices.Equal(got, []Level{{100, 5}}) {
		t.Errorf("bids = %v, want [{100 5}]", got)
	}
	b.Cancel(2)
	if got := b.Depth(Buy, 1); len(got) != 0 {
		t.Errorf("empty level not removed: %v", got)
	}
	got := mustSubmit(t, b, market(3, Sell, 5))
	if len(got) != 0 {
		t.Errorf("traded against cancelled orders: %v", got)
	}
}

func TestPartiallyFilledMakerCanCancel(t *testing.T) {
	b := NewBook()
	mustSubmit(t, b, limit(1, Sell, 100, 10))
	mustSubmit(t, b, limit(2, Buy, 100, 4))
	if !b.Cancel(1) {
		t.Fatal("partially filled maker should still be cancellable")
	}
	if d := b.Depth(Sell, 1); len(d) != 0 {
		t.Errorf("asks = %v after cancel", d)
	}
}

func TestSubmitErrors(t *testing.T) {
	b := NewBook()
	bad := []Order{
		limit(1, Buy, 100, 0),
		limit(2, Buy, 0, 1),
		{ID: 3, Side: 9, Kind: Limit, Price: 1, Qty: 1},
		{ID: 4, Side: Buy, Kind: 9, Price: 1, Qty: 1},
	}
	for _, o := range bad {
		if _, err := b.Submit(o); !errors.Is(err, ErrInvalidOrder) {
			t.Errorf("Submit(%+v) = %v, want ErrInvalidOrder", o, err)
		}
	}
	mustSubmit(t, b, limit(5, Buy, 100, 1))
	if _, err := b.Submit(limit(5, Buy, 100, 1)); !errors.Is(err, ErrDuplicateID) {
		t.Errorf("duplicate id err = %v, want ErrDuplicateID", err)
	}
}

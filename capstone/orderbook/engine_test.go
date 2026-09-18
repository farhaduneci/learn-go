package orderbook

import (
	"context"
	"errors"
	"math/rand/v2"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestEngineConcurrentConservation(t *testing.T) {
	e := NewEngine()
	defer e.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var nextID atomic.Uint64
	var buyQty, sellQty, traded atomic.Int64
	var wg sync.WaitGroup
	for w := range 8 {
		wg.Go(func() {
			r := rand.New(rand.NewPCG(uint64(w), 42))
			for range 500 {
				o := limit(nextID.Add(1), Side(1+r.IntN(2)), int64(95+r.IntN(11)), int64(1+r.IntN(10)))
				trades, err := e.Submit(ctx, "AAPL", o)
				if err != nil {
					t.Error(err)
					return
				}
				if o.Side == Buy {
					buyQty.Add(o.Qty)
				} else {
					sellQty.Add(o.Qty)
				}
				for _, tr := range trades {
					traded.Add(tr.Qty)
				}
			}
		})
	}
	wg.Wait()

	sum := func(lv []Level) (n int64) {
		for _, l := range lv {
			n += l.Qty
		}
		return n
	}
	bids, err := e.Depth(ctx, "AAPL", Buy, 1_000)
	if err != nil {
		t.Fatal(err)
	}
	asks, _ := e.Depth(ctx, "AAPL", Sell, 1_000)
	if traded.Load() == 0 {
		t.Fatal("no trades happened")
	}
	if got := traded.Load() + sum(bids); got != buyQty.Load() {
		t.Errorf("buy side: traded %d + resting %d != submitted %d", traded.Load(), sum(bids), buyQty.Load())
	}
	if got := traded.Load() + sum(asks); got != sellQty.Load() {
		t.Errorf("sell side: traded %d + resting %d != submitted %d", traded.Load(), sum(asks), sellQty.Load())
	}
	if len(bids) > 0 && len(asks) > 0 && bids[0].Price >= asks[0].Price {
		t.Errorf("crossed book: best bid %d >= best ask %d", bids[0].Price, asks[0].Price)
	}
}

func TestEngineSymbolsAreIsolated(t *testing.T) {
	e := NewEngine()
	defer e.Close()
	ctx := context.Background()
	if _, err := e.Submit(ctx, "AAPL", limit(1, Sell, 100, 5)); err != nil {
		t.Fatal(err)
	}
	trades, err := e.Submit(ctx, "TSLA", limit(2, Buy, 100, 5))
	if err != nil || len(trades) != 0 {
		t.Errorf("TSLA order matched AAPL liquidity: %v, %v", trades, err)
	}
	if ok, err := e.Cancel(ctx, "AAPL", 1); !ok || err != nil {
		t.Errorf("Cancel = %v, %v", ok, err)
	}
}

func TestEngineClose(t *testing.T) {
	e := NewEngine()
	e.Close()
	e.Close() // idempotent
	done := make(chan error, 1)
	go func() {
		_, err := e.Submit(context.Background(), "AAPL", limit(1, Buy, 1, 1))
		done <- err
	}()
	select {
	case err := <-done:
		if !errors.Is(err, ErrClosed) {
			t.Errorf("Submit after Close = %v, want ErrClosed", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Submit after Close blocked")
	}
}

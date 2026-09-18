package patterns

import (
	"context"
	"errors"
	"slices"
	"sync"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("https://api.example.com")
	if err != nil || c.Timeout != 5*time.Second || c.Retries != 3 {
		t.Fatalf("defaults = %+v, %v", c, err)
	}
	c, err = NewClient("u", WithTimeout(time.Second), WithRetries(0))
	if err != nil || c.Timeout != time.Second || c.Retries != 0 {
		t.Fatalf("options = %+v, %v", c, err)
	}
	if _, err := NewClient("u", WithTimeout(-1)); err == nil {
		t.Error("negative timeout accepted")
	}
	if _, err := NewClient("u", WithRetries(-1)); err == nil {
		t.Error("negative retries accepted")
	}
}

func TestChainOrder(t *testing.T) {
	var trace []string
	tag := func(name string) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, o Order) error {
				trace = append(trace, name)
				return next(ctx, o)
			}
		}
	}
	h := Chain(func(context.Context, Order) error { trace = append(trace, "h"); return nil }, tag("A"), tag("B"))
	_ = h(context.Background(), Order{Qty: 1})
	if !slices.Equal(trace, []string{"A", "B", "h"}) {
		t.Errorf("trace = %v, want [A B h]", trace)
	}
}

func TestValidateAndRecover(t *testing.T) {
	called := false
	h := Chain(func(context.Context, Order) error { called = true; return nil }, Validate)
	if err := h(context.Background(), Order{Qty: 0}); !errors.Is(err, ErrInvalidOrder) || called {
		t.Errorf("Validate: err=%v called=%v", err, called)
	}
	boom := Chain(func(context.Context, Order) error { panic("boom") }, Recover)
	if err := boom(context.Background(), Order{Qty: 1}); err == nil {
		t.Error("Recover returned nil for a panic")
	}
}

func TestBusDelivers(t *testing.T) {
	b := NewBus()
	ch, unsub := b.Subscribe("AAPL", 4)
	other, unsubOther := b.Subscribe("TSLA", 4)
	defer unsubOther()
	b.Publish(Event{"AAPL", "100.5"})
	select {
	case e := <-ch:
		if e.Payload != "100.5" {
			t.Errorf("payload = %q", e.Payload)
		}
	case <-time.After(time.Second):
		t.Fatal("no event delivered")
	}
	select {
	case e := <-other:
		t.Errorf("TSLA subscriber got %v", e)
	default:
	}
	unsub()
	unsub() // idempotent
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("unexpected event after unsubscribe")
		}
	case <-time.After(time.Second):
		t.Fatal("channel not closed after unsubscribe")
	}
	b.Publish(Event{"AAPL", "x"}) // must not panic on a closed subscriber
}

func TestBusNeverBlocks(t *testing.T) {
	b := NewBus()
	_, unsub := b.Subscribe("AAPL", 1) // slow consumer: never reads
	defer unsub()
	done := make(chan struct{})
	go func() {
		for range 100 {
			b.Publish(Event{"AAPL", "tick"})
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Publish blocked on a slow subscriber")
	}
	if b.Dropped() != 99 {
		t.Errorf("Dropped = %d, want 99", b.Dropped())
	}
}

func TestBusConcurrent(t *testing.T) {
	b := NewBus()
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			ch, unsub := b.Subscribe("X", 8)
			for range 10 {
				b.Publish(Event{"X", "p"})
			}
			unsub()
			for range ch {
			}
		})
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("subscribers never finished (channels not closed on unsubscribe?)")
	}
}

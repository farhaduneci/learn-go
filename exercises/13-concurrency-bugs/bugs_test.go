package bugs

import (
	"errors"
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestFirstQuoteNoLeak(t *testing.T) {
	defer goleak.VerifyNone(t)
	venues := []func() Quote{
		func() Quote { return Quote{"A", 100} },
		func() Quote { time.Sleep(5 * time.Millisecond); return Quote{"B", 101} },
		func() Quote { time.Sleep(10 * time.Millisecond); return Quote{"C", 102} },
	}
	if q := FirstQuote(venues); q.Venue != "A" {
		t.Errorf("FirstQuote = %+v, want venue A", q)
	}
	time.Sleep(30 * time.Millisecond) // let the losers finish, if they can
}

func TestTransferNoDeadlock(t *testing.T) {
	a := &Account{ID: 1, Balance: 1_000}
	b := &Account{ID: 2, Balance: 1_000}
	done := make(chan struct{})
	go func() {
		var wg sync.WaitGroup
		for range 1_000 {
			wg.Go(func() { Transfer(a, b, 1) })
			wg.Go(func() { Transfer(b, a, 1) })
		}
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("deadlock: transfers did not finish")
	}
	if a.Balance+b.Balance != 2_000 {
		t.Errorf("money created or destroyed: %d + %d", a.Balance, b.Balance)
	}
}

func TestValidateAll(t *testing.T) {
	for range 50 {
		if err := ValidateAll([]int64{1, 2, -1, 4, 0, 6}); !errors.Is(err, ErrInvalid) {
			t.Fatalf("ValidateAll = %v, want ErrInvalid", err)
		}
	}
	if err := ValidateAll([]int64{1, 2, 3}); err != nil {
		t.Errorf("ValidateAll(valid) = %v", err)
	}
}

func TestFeedStopTwice(t *testing.T) {
	f := NewFeed()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("second Stop panicked: %v", r)
		}
	}()
	f.Stop()
	f.Stop()
	select {
	case <-f.Done():
	default:
		t.Error("Done not closed after Stop")
	}
}

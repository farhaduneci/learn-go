package patterns

import (
	"context"
	"errors"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) { goleak.VerifyTestMain(m) }

func feed(n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := range n {
			ch <- i
		}
	}()
	return ch
}

func TestWorkerPool(t *testing.T) {
	out := WorkerPool(context.Background(), feed(100), 4, func(n int) int { return n * 2 })
	var got []int
	timeout := time.After(2 * time.Second)
	for done := false; !done; {
		select {
		case v, ok := <-out:
			if !ok {
				done = true
				break
			}
			got = append(got, v)
		case <-timeout:
			t.Fatal("output never closed")
		}
	}
	slices.Sort(got)
	if len(got) != 100 || got[0] != 0 || got[99] != 198 {
		t.Errorf("got %d results, first %v", len(got), got[:min(3, len(got))])
	}
}

func TestWorkerPoolCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	jobs := make(chan int) // never closed: only cancellation can stop the workers
	out := WorkerPool(ctx, jobs, 3, func(n int) int { return n })
	select {
	case jobs <- 1:
	case <-time.After(time.Second):
		t.Fatal("no worker is reading jobs")
	}
	cancel() // a worker may be blocked sending 1 to out; it must still exit
	select {
	case <-out:
	case <-time.After(time.Second):
		t.Fatal("pool did not shut down after cancel")
	}
	for range out {
	}
}

func TestProcessAllLimitAndOrder(t *testing.T) {
	var inFlight, peak atomic.Int32
	items := []int{1, 2, 3, 4, 5, 6, 7, 8}
	got, err := ProcessAll(context.Background(), items, 3, func(_ context.Context, n int) (int, error) {
		cur := inFlight.Add(1)
		for p := peak.Load(); cur > p && !peak.CompareAndSwap(p, cur); p = peak.Load() {
		}
		time.Sleep(20 * time.Millisecond)
		inFlight.Add(-1)
		return n * n, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []int{1, 4, 9, 16, 25, 36, 49, 64}) {
		t.Errorf("results = %v (order must match input)", got)
	}
	if p := peak.Load(); p > 3 || p < 2 {
		t.Errorf("peak in-flight = %d, want 2..3 (limit 3, and must be parallel)", p)
	}
}

func TestProcessAllFirstErrorCancels(t *testing.T) {
	boom := errors.New("boom")
	var cancelled atomic.Int32
	_, err := ProcessAll(context.Background(), []int{0, 1, 2, 3}, 4, func(ctx context.Context, n int) (int, error) {
		if n == 0 {
			return 0, boom
		}
		select {
		case <-ctx.Done():
			cancelled.Add(1)
			return 0, ctx.Err()
		case <-time.After(time.Second):
			return n, nil
		}
	})
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
	if cancelled.Load() != 3 {
		t.Errorf("%d siblings saw cancellation, want 3", cancelled.Load())
	}
}

func TestThrottle(t *testing.T) {
	start := time.Now()
	n := 0
	for range Throttle(context.Background(), feed(5), 50) { // 20ms apart
		n++
	}
	if n != 5 {
		t.Fatalf("forwarded %d values, want 5", n)
	}
	if el := time.Since(start); el < 80*time.Millisecond {
		t.Errorf("5 values at 50/s took %v, want >= ~100ms", el)
	}
}

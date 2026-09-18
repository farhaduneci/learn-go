package goroutines

import (
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelSum(t *testing.T) {
	nums := make([]int, 10_001)
	want := 0
	for i := range nums {
		nums[i] = i
		want += i
	}
	for _, w := range []int{1, 3, 8, 64} {
		if got := ParallelSum(nums, w); got != want {
			t.Errorf("ParallelSum(workers=%d) = %d, want %d", w, got, want)
		}
	}
}

func TestFetchAllConcurrentAndOrdered(t *testing.T) {
	var inFlight, peak atomic.Int32
	fetch := func(s string) (int64, error) {
		n := inFlight.Add(1)
		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
		inFlight.Add(-1)
		if s == "BAD" {
			return 0, errors.New("unknown symbol")
		}
		return int64(len(s)), nil
	}
	syms := []string{"A", "BB", "BAD", "DDDD", "EEEEE"}

	before := runtime.NumGoroutine()
	start := time.Now()
	got := FetchAll(syms, fetch)
	elapsed := time.Since(start)

	if peak.Load() < 2 {
		t.Fatalf("peak concurrent fetches = %d; calls ran sequentially (took %v)", peak.Load(), elapsed)
	}
	for i, q := range got {
		if q.Symbol != syms[i] {
			t.Fatalf("result %d is %q, want %q (order lost)", i, q.Symbol, syms[i])
		}
	}
	if got[2].Err == nil || got[1].Price != 2 {
		t.Errorf("unexpected results: %+v", got)
	}
	time.Sleep(10 * time.Millisecond)
	if after := runtime.NumGoroutine(); after > before {
		t.Errorf("goroutines before=%d after=%d: leaked", before, after)
	}
}

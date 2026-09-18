package drills

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestLRU(t *testing.T) {
	c := NewLRU(2)
	c.Put("a", 1)
	c.Put("b", 2)
	if v, ok := c.Get("a"); !ok || v != 1 { // a is now most recent
		t.Fatalf("Get(a) = %d, %v", v, ok)
	}
	c.Put("c", 3) // evicts b
	if _, ok := c.Get("b"); ok {
		t.Error("b should have been evicted")
	}
	c.Put("a", 10) // update, no eviction
	if v, _ := c.Get("a"); v != 10 {
		t.Errorf("Get(a) = %d, want 10", v)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Errorf("Get(c) = %d, %v", v, ok)
	}
}

func TestLRUConcurrent(t *testing.T) {
	c := NewLRU(16)
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			for j := range 200 {
				k := fmt.Sprint((i + j) % 32)
				c.Put(k, j)
				c.Get(k)
			}
		})
	}
	wg.Wait()
}

func TestLimiter(t *testing.T) {
	now := time.Unix(0, 0)
	clock := func() time.Time { return now }
	l := NewLimiter(2, 3, clock) // 2/s, burst 3
	for i := range 3 {
		if !l.Allow() {
			t.Fatalf("burst call %d denied", i)
		}
	}
	if l.Allow() {
		t.Fatal("4th immediate call allowed; bucket should be empty")
	}
	now = now.Add(500 * time.Millisecond) // +1 token
	if !l.Allow() || l.Allow() {
		t.Error("after 500ms exactly one call should be allowed")
	}
	now = now.Add(10 * time.Second) // refill caps at burst
	allowed := 0
	for range 10 {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 3 {
		t.Errorf("after long idle allowed %d, want burst 3", allowed)
	}
}

func TestPingPong(t *testing.T) {
	got := PingPong(3)
	if strings.Join(got, " ") != "ping pong ping pong ping pong" {
		t.Errorf("PingPong(3) = %v", got)
	}
	if len(PingPong(0)) != 0 {
		t.Error("PingPong(0) should be empty")
	}
}

func TestTopK(t *testing.T) {
	syms := strings.Fields("AAPL TSLA AAPL MSFT TSLA AAPL NVDA MSFT TSLA")
	if got := TopK(syms, 2); !slices.Equal(got, []string{"AAPL", "TSLA"}) {
		t.Errorf("TopK(2) = %v, want [AAPL TSLA]", got)
	}
	if got := TopK(syms, 4); !slices.Equal(got, []string{"AAPL", "TSLA", "MSFT", "NVDA"}) {
		t.Errorf("TopK(4) = %v", got)
	}
}

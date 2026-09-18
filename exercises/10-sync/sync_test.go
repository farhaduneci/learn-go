package syncx

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

const goroutines = 50

func TestPositions(t *testing.T) {
	p := NewPositions()
	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Go(func() {
			for range 100 {
				p.Apply("AAPL", 1)
				p.Apply(fmt.Sprint("S", i%5), -1)
				_ = p.Get("AAPL")
			}
		})
	}
	wg.Wait()
	if got := p.Get("AAPL"); got != goroutines*100 {
		t.Errorf("AAPL = %d, want %d", got, goroutines*100)
	}
}

func TestCounter(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	for range goroutines {
		wg.Go(func() {
			for range 1000 {
				c.Inc()
			}
		})
	}
	wg.Wait()
	if got := c.Value(); got != goroutines*1000 {
		t.Errorf("Counter = %d, want %d (lost updates)", got, goroutines*1000)
	}
}

func TestQuoteCache(t *testing.T) {
	c := NewQuoteCache()
	var wg sync.WaitGroup
	for i := range goroutines {
		wg.Go(func() {
			if i%10 == 0 {
				c.Set("TSLA", int64(i))
			}
			c.Get("TSLA")
		})
	}
	wg.Wait()
	if _, ok := c.Get("TSLA"); !ok {
		t.Error("TSLA missing")
	}
}

func TestConfigLoadsOnce(t *testing.T) {
	var loads atomic.Int32
	cfg := NewConfig(func() map[string]string {
		loads.Add(1)
		return map[string]string{"venue": "XNAS"}
	})
	var wg sync.WaitGroup
	for range goroutines {
		wg.Go(func() {
			if got := cfg.Get("venue"); got != "XNAS" {
				t.Errorf("Get = %q", got)
			}
		})
	}
	wg.Wait()
	if n := loads.Load(); n != 1 {
		t.Errorf("loader ran %d times, want 1", n)
	}
}

package obs

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger(&buf, slog.LevelInfo)
	logger.Debug("dropped")
	ctx := WithRequestID(context.Background(), "req-7")
	LogFill(ctx, logger, Fill{Symbol: "AAPL", Qty: 10, Price: 18950})

	var rec map[string]any
	if err := json.Unmarshal(buf.Bytes(), &rec); err != nil {
		t.Fatalf("want exactly one JSON record, got %q: %v", buf.String(), err)
	}
	if rec["msg"] != "order filled" || rec["service"] != "order-gateway" || rec["request_id"] != "req-7" {
		t.Errorf("record = %v", rec)
	}
	fill, ok := rec["fill"].(map[string]any)
	if !ok || fill["symbol"] != "AAPL" || fill["qty"] != float64(10) || fill["price"] != float64(18950) {
		t.Errorf("fill group = %v", rec["fill"])
	}
}

func TestLogFillWithoutRequestID(t *testing.T) {
	var buf bytes.Buffer
	LogFill(context.Background(), NewLogger(&buf, slog.LevelInfo), Fill{Symbol: "A"})
	var rec map[string]any
	_ = json.Unmarshal(buf.Bytes(), &rec)
	if _, ok := rec["request_id"]; ok {
		t.Error("request_id present without one in context")
	}
}

func TestHistogram(t *testing.T) {
	h := NewHistogram(time.Millisecond, 5*time.Millisecond, 10*time.Millisecond)
	for range 90 {
		h.Observe(500 * time.Microsecond)
	}
	for range 9 {
		h.Observe(4 * time.Millisecond)
	}
	h.Observe(time.Second)
	if h.Count() != 100 {
		t.Fatalf("Count = %d", h.Count())
	}
	for q, want := range map[float64]time.Duration{0.5: time.Millisecond, 0.9: time.Millisecond, 0.99: 5 * time.Millisecond, 1: -1} {
		if got := h.Quantile(q); got != want {
			t.Errorf("Quantile(%v) = %v, want %v", q, got, want)
		}
	}
}

func TestHistogramConcurrent(t *testing.T) {
	h := NewHistogram(time.Millisecond)
	var wg sync.WaitGroup
	for range 20 {
		wg.Go(func() {
			for range 500 {
				h.Observe(time.Microsecond)
			}
		})
	}
	wg.Wait()
	if h.Count() != 10_000 {
		t.Errorf("Count = %d, want 10000", h.Count())
	}
}

func TestTimed(t *testing.T) {
	h := NewHistogram(time.Millisecond, 100*time.Millisecond)
	f := Timed(h, func(context.Context) error { time.Sleep(5 * time.Millisecond); return nil })
	_ = f(context.Background())
	if h.Count() != 1 || h.Quantile(1) != 100*time.Millisecond {
		t.Errorf("Count=%d p100=%v", h.Count(), h.Quantile(1))
	}
}

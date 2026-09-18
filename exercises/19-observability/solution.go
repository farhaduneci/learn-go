//go:build solution

package obs

import (
	"context"
	"io"
	"log/slog"
	"math"
	"sync/atomic"
	"time"
)

// NewLogger returns a JSON logger at level, tagged with service=order-gateway.
func NewLogger(w io.Writer, level slog.Level) *slog.Logger {
	h := slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level})
	return slog.New(h).With(slog.String("service", "order-gateway"))
}

// LogFill logs f as a group, with the request ID from ctx if present.
func LogFill(ctx context.Context, logger *slog.Logger, f Fill) {
	attrs := []slog.Attr{
		slog.Group("fill",
			slog.String("symbol", f.Symbol),
			slog.Int64("qty", f.Qty),
			slog.Int64("price", f.Price),
		),
	}
	if id, ok := RequestID(ctx); ok {
		attrs = append(attrs, slog.String("request_id", id))
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "order filled", attrs...)
}

// Histogram counts durations into fixed buckets. Safe for concurrent use.
type Histogram struct {
	bounds []time.Duration
	counts []atomic.Int64 // len(bounds)+1; the last is the overflow (+Inf) bucket
}

// NewHistogram returns a histogram with ascending upper bounds.
func NewHistogram(bounds ...time.Duration) *Histogram {
	return &Histogram{bounds: bounds, counts: make([]atomic.Int64, len(bounds)+1)}
}

// Observe records d in the first bucket whose bound is >= d.
func (h *Histogram) Observe(d time.Duration) {
	i := 0
	for i < len(h.bounds) && d > h.bounds[i] {
		i++
	}
	h.counts[i].Add(1)
}

// Count returns the total number of observations.
func (h *Histogram) Count() int64 {
	var n int64
	for i := range h.counts {
		n += h.counts[i].Load()
	}
	return n
}

// Quantile returns the upper bound of the bucket containing quantile q (0 < q <= 1).
// Observations beyond the last bound report -1 (unknown, "+Inf").
func (h *Histogram) Quantile(q float64) time.Duration {
	total := h.Count()
	if total == 0 {
		return 0
	}
	rank := int64(math.Ceil(q * float64(total)))
	var seen int64
	for i, b := range h.bounds {
		seen += h.counts[i].Load()
		if seen >= rank {
			return b
		}
	}
	return -1
}

// Timed wraps next, recording its duration in h.
func Timed(h *Histogram, next func(context.Context) error) func(context.Context) error {
	return func(ctx context.Context) error {
		start := time.Now()
		defer func() { h.Observe(time.Since(start)) }()
		return next(ctx)
	}
}

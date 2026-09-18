//go:build !solution

package obs

import (
	"context"
	"io"
	"log/slog"
	"time"
)

// NewLogger returns a JSON logger at level, tagged with service=order-gateway.
func NewLogger(w io.Writer, level slog.Level) *slog.Logger {
	// TODO
	return slog.New(slog.NewTextHandler(w, nil))
}

// LogFill logs f as a group, with the request ID from ctx if present.
func LogFill(ctx context.Context, logger *slog.Logger, f Fill) {
	// TODO
}

// Histogram counts durations into fixed buckets. Safe for concurrent use.
type Histogram struct {
	// TODO: bounds and one atomic counter per bucket (+1 overflow bucket)
}

// NewHistogram returns a histogram with ascending upper bounds.
func NewHistogram(bounds ...time.Duration) *Histogram {
	// TODO
	return &Histogram{}
}

// Observe records d in the first bucket whose bound is >= d.
func (h *Histogram) Observe(d time.Duration) {
	// TODO
}

// Count returns the total number of observations.
func (h *Histogram) Count() int64 {
	// TODO
	return 0
}

// Quantile returns the upper bound of the bucket containing quantile q (0 < q <= 1).
// Observations beyond the last bound report -1 (unknown, "+Inf").
func (h *Histogram) Quantile(q float64) time.Duration {
	// TODO
	return 0
}

// Timed wraps next, recording its duration in h.
func Timed(h *Histogram, next func(context.Context) error) func(context.Context) error {
	// TODO
	return next
}

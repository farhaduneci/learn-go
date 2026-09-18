//go:build !solution

package ctxx

import (
	"context"
	"time"
)

// Sleep waits for d, returning ctx.Err() early if ctx is done.
func Sleep(ctx context.Context, d time.Duration) error {
	// TODO: implement with select
	time.Sleep(d)
	return nil
}

// FirstSuccess returns the first successful quote and cancels the other venues.
func FirstSuccess(ctx context.Context, venues ...Venue) (int64, error) {
	// TODO: implement
	return 0, nil
}

// WithRequestID returns a copy of ctx carrying id.
func WithRequestID(ctx context.Context, id string) context.Context {
	// TODO: implement with an unexported key type
	return ctx
}

// RequestID returns the request ID stored in ctx, if any.
func RequestID(ctx context.Context) (string, bool) {
	// TODO: implement
	return "", false
}

// Retry calls op up to attempts times, sleeping backoff, 2*backoff, ... between
// failures. It returns nil on the first success, ctx.Err() if cancelled while
// waiting, or the last error.
func Retry(ctx context.Context, attempts int, backoff time.Duration, op func(context.Context) error) error {
	// TODO: implement
	return op(ctx)
}

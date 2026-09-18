//go:build solution

package ctxx

import (
	"context"
	"errors"
	"time"
)

// Sleep waits for d, returning ctx.Err() early if ctx is done.
func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// FirstSuccess returns the first successful quote and cancels the other venues.
func FirstSuccess(ctx context.Context, venues ...Venue) (int64, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // stops the losers once a winner returns

	type result struct {
		price int64
		err   error
	}
	results := make(chan result, len(venues)) // buffered: losers never block
	for _, v := range venues {
		go func() {
			p, err := v(ctx)
			results <- result{p, err}
		}()
	}

	var errs []error
	for range venues {
		select {
		case r := <-results:
			if r.err == nil {
				return r.price, nil
			}
			errs = append(errs, r.err)
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
	return 0, errors.Join(errs...)
}

type requestIDKey struct{}

// WithRequestID returns a copy of ctx carrying id.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID returns the request ID stored in ctx, if any.
func RequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	return id, ok
}

// Retry calls op up to attempts times, sleeping backoff, 2*backoff, ... between
// failures. It returns nil on the first success, ctx.Err() if cancelled while
// waiting, or the last error.
func Retry(ctx context.Context, attempts int, backoff time.Duration, op func(context.Context) error) error {
	var err error
	for i := range attempts {
		if err = op(ctx); err == nil {
			return nil
		}
		if i == attempts-1 {
			break
		}
		if serr := Sleep(ctx, backoff); serr != nil {
			return serr
		}
		backoff *= 2
	}
	return err
}

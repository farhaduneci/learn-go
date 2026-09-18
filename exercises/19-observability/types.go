// Package obs practises structured logging and latency metrics.
package obs

import "context"

// Fill is an executed order.
type Fill struct {
	Symbol string
	Qty    int64
	Price  int64
}

type requestIDKey struct{}

// WithRequestID stores id in ctx.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, id)
}

// RequestID reads the id stored by WithRequestID.
func RequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	return id, ok
}

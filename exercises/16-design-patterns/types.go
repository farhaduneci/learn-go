// Package patterns implements idiomatic Go design patterns.
package patterns

import (
	"context"
	"errors"
	"time"
)

// Order is a minimal order for the middleware chain.
type Order struct {
	Symbol string
	Qty    int64
}

// Handler handles an order.
type Handler func(ctx context.Context, o Order) error

// Middleware wraps a Handler with extra behaviour.
type Middleware func(next Handler) Handler

// ErrInvalidOrder is returned by Validate.
var ErrInvalidOrder = errors.New("invalid order")

// Client is an API client configured with functional options.
type Client struct {
	BaseURL string
	Timeout time.Duration
	Retries int
}

// Event is a message on the bus.
type Event struct {
	Topic   string
	Payload string
}

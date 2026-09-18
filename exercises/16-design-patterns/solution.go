//go:build solution

package patterns

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Option configures a Client.
type Option func(*Client) error

// WithTimeout sets the request timeout; it must be positive.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) error {
		if d <= 0 {
			return errors.New("timeout must be positive")
		}
		c.Timeout = d
		return nil
	}
}

// WithRetries sets the retry count; it must be >= 0.
func WithRetries(n int) Option {
	return func(c *Client) error {
		if n < 0 {
			return errors.New("retries must be >= 0")
		}
		c.Retries = n
		return nil
	}
}

// NewClient builds a Client with defaults overridden by opts.
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	c := &Client{BaseURL: baseURL, Timeout: 5 * time.Second, Retries: 3}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("new client: %w", err)
		}
	}
	return c, nil
}

// Chain wraps h so that mws[0] runs first.
func Chain(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// Validate rejects orders with Qty <= 0.
func Validate(next Handler) Handler {
	return func(ctx context.Context, o Order) error {
		if o.Qty <= 0 {
			return fmt.Errorf("%w: qty %d", ErrInvalidOrder, o.Qty)
		}
		return next(ctx, o)
	}
}

// Recover converts a panic in next into an error.
func Recover(next Handler) Handler {
	return func(ctx context.Context, o Order) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("handler panic: %v", r)
			}
		}()
		return next(ctx, o)
	}
}

// Bus is a non-blocking topic-based pub/sub.
type Bus struct {
	mu      sync.RWMutex
	subs    map[string]map[chan Event]struct{}
	dropped atomic.Int64
}

// NewBus returns an empty Bus.
func NewBus() *Bus { return &Bus{subs: make(map[string]map[chan Event]struct{})} }

// Subscribe returns a channel of events for topic with buffer size buf,
// and a function that unsubscribes and closes the channel.
func (b *Bus) Subscribe(topic string, buf int) (<-chan Event, func()) {
	ch := make(chan Event, buf)
	b.mu.Lock()
	if b.subs[topic] == nil {
		b.subs[topic] = make(map[chan Event]struct{})
	}
	b.subs[topic][ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	return ch, func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subs[topic], ch)
			b.mu.Unlock()
			close(ch) // safe: Publish sends only while holding the read lock
		})
	}
}

// Publish delivers e to every subscriber of e.Topic without blocking.
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs[e.Topic] {
		select {
		case ch <- e:
		default:
			b.dropped.Add(1)
		}
	}
}

// Dropped returns how many deliveries were dropped because a buffer was full.
func (b *Bus) Dropped() int64 { return b.dropped.Load() }

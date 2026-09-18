//go:build solution

package messaging

import (
	"context"
	"hash/fnv"
	"sync"
)

// Dedupe runs next at most once per successful message ID.
func Dedupe(next Handler) Handler {
	var (
		mu   sync.Mutex
		seen = make(map[string]struct{}) // production: a DB table or TTL cache, not unbounded memory
	)
	return func(ctx context.Context, m Message) error {
		mu.Lock()
		_, dup := seen[m.ID]
		mu.Unlock()
		if dup {
			return nil
		}
		if err := next(ctx, m); err != nil {
			return err
		}
		mu.Lock()
		seen[m.ID] = struct{}{}
		mu.Unlock()
		return nil
	}
}

// WithDeadLetter nacks failures until maxAttempts, then sends to dlq and acks.
func WithDeadLetter(next Handler, maxAttempts int, dlq func(Message, error)) Handler {
	return func(ctx context.Context, m Message) error {
		err := next(ctx, m)
		if err == nil || m.Attempt < maxAttempts {
			return err
		}
		dlq(m, err)
		return nil
	}
}

// Dispatcher routes messages to per-key-ordered partitions.
type Dispatcher struct {
	parts []chan Message
	wg    sync.WaitGroup
}

// NewDispatcher starts partitions goroutines, each running h sequentially.
func NewDispatcher(partitions int, h Handler) *Dispatcher {
	d := &Dispatcher{parts: make([]chan Message, partitions)}
	for i := range d.parts {
		ch := make(chan Message, 64)
		d.parts[i] = ch
		d.wg.Go(func() {
			for m := range ch {
				_ = h(context.Background(), m)
			}
		})
	}
	return d
}

// Dispatch enqueues m on the partition for m.Key.
func (d *Dispatcher) Dispatch(m Message) {
	h := fnv.New32a()
	_, _ = h.Write([]byte(m.Key))
	d.parts[h.Sum32()%uint32(len(d.parts))] <- m
}

// Close stops accepting messages, drains every partition and waits.
func (d *Dispatcher) Close() {
	for _, ch := range d.parts {
		close(ch)
	}
	d.wg.Wait()
}

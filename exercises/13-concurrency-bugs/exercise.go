//go:build !solution

package bugs

import "sync"

// FirstQuote returns the first quote produced by any venue.
func FirstQuote(venues []func() Quote) Quote {
	ch := make(chan Quote)
	for _, v := range venues {
		go func() { ch <- v() }()
	}
	return <-ch
}

// Transfer moves amount from one account to another atomically.
func Transfer(from, to *Account, amount int64) {
	from.mu.Lock()
	defer from.mu.Unlock()
	to.mu.Lock()
	defer to.mu.Unlock()
	from.Balance -= amount
	to.Balance += amount
}

// ValidateAll checks every quantity concurrently and returns an error if any is invalid.
func ValidateAll(qtys []int64) error {
	var firstErr error
	var wg sync.WaitGroup
	for _, q := range qtys {
		wg.Go(func() {
			if q <= 0 {
				firstErr = ErrInvalid
			}
		})
	}
	wg.Wait()
	return firstErr
}

// Feed is a stream that can be stopped from several places.
type Feed struct {
	done chan struct{}
}

// NewFeed returns a running feed.
func NewFeed() *Feed { return &Feed{done: make(chan struct{})} }

// Stop stops the feed. It is safe to call more than once.
func (f *Feed) Stop() { close(f.done) }

// Done is closed once the feed stops.
func (f *Feed) Done() <-chan struct{} { return f.done }

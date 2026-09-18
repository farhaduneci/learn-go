//go:build solution

package bugs

import "sync"

// FirstQuote returns the first quote produced by any venue.
func FirstQuote(venues []func() Quote) Quote {
	// fix: goroutine leak. With an unbuffered channel the losers block forever
	// on send. A buffer sized to the number of senders lets every one finish.
	ch := make(chan Quote, len(venues))
	for _, v := range venues {
		go func() { ch <- v() }()
	}
	return <-ch
}

// Transfer moves amount from one account to another atomically.
func Transfer(from, to *Account, amount int64) {
	// fix: lock-ordering deadlock. A->B and B->A each held one lock and waited
	// for the other. Always acquire locks in one global order (by ID).
	first, second := from, to
	if second.ID < first.ID {
		first, second = second, first
	}
	first.mu.Lock()
	defer first.mu.Unlock()
	second.mu.Lock()
	defer second.mu.Unlock()
	from.Balance -= amount
	to.Balance += amount
}

// ValidateAll checks every quantity concurrently and returns an error if any is invalid.
func ValidateAll(qtys []int64) error {
	// fix: data race on firstErr. Guard the shared write with a mutex
	// (errgroup would also handle this).
	var (
		mu       sync.Mutex
		firstErr error
		wg       sync.WaitGroup
	)
	for _, q := range qtys {
		wg.Go(func() {
			if q <= 0 {
				mu.Lock()
				if firstErr == nil {
					firstErr = ErrInvalid
				}
				mu.Unlock()
			}
		})
	}
	wg.Wait()
	return firstErr
}

// Feed is a stream that can be stopped from several places.
type Feed struct {
	done chan struct{}
	once sync.Once
}

// NewFeed returns a running feed.
func NewFeed() *Feed { return &Feed{done: make(chan struct{})} }

// Stop stops the feed. It is safe to call more than once.
func (f *Feed) Stop() {
	// fix: double close panics. sync.Once makes the close idempotent.
	f.once.Do(func() { close(f.done) })
}

// Done is closed once the feed stops.
func (f *Feed) Done() <-chan struct{} { return f.done }

//go:build solution

// Package channels practises channel ownership, fan-in and select.
package channels

import (
	"sync"
	"time"
)

// Gen emits nums on the returned channel, then closes it.
func Gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Square emits the square of every value received from in.
func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// Merge fans in all chans into one channel, closed when every input is drained.
func Merge(chans ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range chans {
		wg.Go(func() {
			for v := range c {
				out <- v
			}
		})
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// WithTimeout returns f's result, or ok=false if d elapses first.
func WithTimeout(f func() int, d time.Duration) (result int, ok bool) {
	done := make(chan int, 1) // buffered: the sender never blocks, so it never leaks
	go func() { done <- f() }()
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case v := <-done:
		return v, true
	case <-timer.C:
		return 0, false
	}
}

// Batch groups values from in into slices of at most size, flushing a partial
// batch when wait passes without a new value. The output closes after in closes
// and any remaining partial batch is sent.
func Batch(in <-chan int, size int, wait time.Duration) <-chan []int {
	out := make(chan []int)
	go func() {
		defer close(out)
		var buf []int
		timer := time.NewTimer(wait)
		defer timer.Stop()
		flush := func() {
			if len(buf) > 0 {
				out <- buf
				buf = nil
			}
		}
		for {
			select {
			case v, ok := <-in:
				if !ok {
					flush()
					return
				}
				buf = append(buf, v)
				if len(buf) == size {
					flush()
				}
				timer.Reset(wait)
			case <-timer.C:
				flush()
				timer.Reset(wait)
			}
		}
	}()
	return out
}

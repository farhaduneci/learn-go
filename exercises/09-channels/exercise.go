//go:build !solution

// Package channels practises channel ownership, fan-in and select.
package channels

import "time"

// Gen emits nums on the returned channel, then closes it.
func Gen(nums ...int) <-chan int {
	// TODO: implement
	return make(chan int)
}

// Square emits the square of every value received from in.
func Square(in <-chan int) <-chan int {
	// TODO: implement
	return make(chan int)
}

// Merge fans in all chans into one channel, closed when every input is drained.
func Merge(chans ...<-chan int) <-chan int {
	// TODO: implement
	return make(chan int)
}

// WithTimeout returns f's result, or ok=false if d elapses first.
func WithTimeout(f func() int, d time.Duration) (result int, ok bool) {
	// TODO: implement
	return f(), true
}

// Batch groups values from in into slices of at most size, flushing a partial
// batch when wait passes without a new value. The output closes after in closes
// and any remaining partial batch is sent.
func Batch(in <-chan int, size int, wait time.Duration) <-chan []int {
	// TODO: implement
	return make(chan []int)
}

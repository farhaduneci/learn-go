//go:build !solution

// Package drills holds classic Go live-coding interview problems.
package drills

import "time"

// LRU is a fixed-capacity least-recently-used cache. Safe for concurrent use.
type LRU struct{}

// NewLRU returns an empty cache holding at most capacity items.
func NewLRU(capacity int) *LRU {
	// TODO
	return &LRU{}
}

// Get returns the value for key and marks it most recently used.
func (c *LRU) Get(key string) (int, bool) {
	// TODO
	return 0, false
}

// Put inserts or updates key, evicting the least recently used item if full.
func (c *LRU) Put(key string, value int) {
	// TODO
}

// Limiter is a token-bucket rate limiter.
type Limiter struct{}

// NewLimiter allows rate tokens per second with bursts of up to burst.
func NewLimiter(rate float64, burst int, now func() time.Time) *Limiter {
	// TODO
	return &Limiter{}
}

// Allow reports whether one event may happen now, spending a token if so.
func (l *Limiter) Allow() bool {
	// TODO
	return true
}

// PingPong returns 2*n alternating "ping"/"pong" strings produced by two goroutines.
func PingPong(n int) []string {
	// TODO
	return nil
}

// TopK returns the k most frequent symbols, most frequent first, ties by symbol.
func TopK(symbols []string, k int) []string {
	// TODO
	return nil
}

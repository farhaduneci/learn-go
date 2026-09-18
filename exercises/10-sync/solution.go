//go:build solution

// Package syncx practises mutexes, atomics and once-initialisation.
package syncx

import (
	"sync"
	"sync/atomic"
)

// Positions tracks net quantity per symbol.
type Positions struct {
	mu  sync.Mutex // guards qty
	qty map[string]int64
}

// NewPositions returns empty positions.
func NewPositions() *Positions { return &Positions{qty: make(map[string]int64)} }

// Apply adds delta to symbol's position.
func (p *Positions) Apply(symbol string, delta int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.qty[symbol] += delta
}

// Get returns symbol's position.
func (p *Positions) Get(symbol string) int64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.qty[symbol]
}

// Counter counts events.
type Counter struct {
	n atomic.Int64
}

// Inc adds one.
func (c *Counter) Inc() { c.n.Add(1) }

// Value returns the count.
func (c *Counter) Value() int64 { return c.n.Load() }

// QuoteCache maps symbols to last prices; read-heavy.
type QuoteCache struct {
	mu     sync.RWMutex // guards prices
	prices map[string]int64
}

// NewQuoteCache returns an empty cache.
func NewQuoteCache() *QuoteCache { return &QuoteCache{prices: make(map[string]int64)} }

// Set stores price for symbol.
func (c *QuoteCache) Set(symbol string, price int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prices[symbol] = price
}

// Get returns the price for symbol.
func (c *QuoteCache) Get(symbol string) (int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.prices[symbol]
	return p, ok
}

// Config lazily loads settings exactly once.
type Config struct {
	settings func() map[string]string
}

// NewConfig returns a Config that calls load on first use.
func NewConfig(load func() map[string]string) *Config {
	return &Config{settings: sync.OnceValue(load)}
}

// Get returns the setting for key.
func (c *Config) Get(key string) string { return c.settings()[key] }

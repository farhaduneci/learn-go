//go:build !solution

// Package syncx practises mutexes, atomics and once-initialisation.
package syncx

// Positions tracks net quantity per symbol.
type Positions struct {
	qty map[string]int64 // TODO: guard with a mutex
}

// NewPositions returns empty positions.
func NewPositions() *Positions { return &Positions{qty: make(map[string]int64)} }

// Apply adds delta to symbol's position.
func (p *Positions) Apply(symbol string, delta int64) { p.qty[symbol] += delta }

// Get returns symbol's position.
func (p *Positions) Get(symbol string) int64 { return p.qty[symbol] }

// Counter counts events.
type Counter struct {
	n int64 // TODO: make this atomic
}

// Inc adds one.
func (c *Counter) Inc() { c.n++ }

// Value returns the count.
func (c *Counter) Value() int64 { return c.n }

// QuoteCache maps symbols to last prices; read-heavy.
type QuoteCache struct {
	prices map[string]int64 // TODO: guard with an RWMutex
}

// NewQuoteCache returns an empty cache.
func NewQuoteCache() *QuoteCache { return &QuoteCache{prices: make(map[string]int64)} }

// Set stores price for symbol.
func (c *QuoteCache) Set(symbol string, price int64) { c.prices[symbol] = price }

// Get returns the price for symbol.
func (c *QuoteCache) Get(symbol string) (int64, bool) {
	p, ok := c.prices[symbol]
	return p, ok
}

// Config lazily loads settings exactly once.
type Config struct {
	load     func() map[string]string
	settings map[string]string // TODO: initialise exactly once, safely
}

// NewConfig returns a Config that calls load on first use.
func NewConfig(load func() map[string]string) *Config { return &Config{load: load} }

// Get returns the setting for key.
func (c *Config) Get(key string) string {
	if c.settings == nil {
		c.settings = c.load()
	}
	return c.settings[key]
}

// Package service is a small JSON HTTP API with graceful shutdown.
package service

// QuoteStore looks up the last price for a symbol.
type QuoteStore interface {
	Price(symbol string) (int64, bool)
}

// MapQuotes is a QuoteStore backed by a map.
type MapQuotes map[string]int64

// Price implements QuoteStore.
func (m MapQuotes) Price(symbol string) (int64, bool) {
	p, ok := m[symbol]
	return p, ok
}

// OrderRequest is the POST /orders body.
type OrderRequest struct {
	Symbol string `json:"symbol"`
	Side   string `json:"side"`
	Qty    int64  `json:"qty"`
}

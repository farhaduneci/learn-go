//go:build !solution

package perf

import "fmt"

// Tick is a compact market-data update.
type Tick struct {
	Halted  bool
	Price   int64
	IsBid   bool
	Size    int32
	IsFinal bool
	Seq     int32
}

// VWAP returns the volume-weighted average price, or 0 for no volume.
func VWAP(trades []Trade) int64 {
	notionals := make([]int64, 0)
	var volume int64
	for _, t := range trades {
		notionals = append(notionals, t.Price*t.Qty)
		volume += t.Qty
	}
	if volume == 0 {
		return 0
	}
	var total int64
	for _, n := range notionals {
		total += n
	}
	return total / volume
}

// FormatTrades renders one "SYMBOL qty@price" line per trade.
func FormatTrades(trades []Trade) string {
	out := ""
	for _, t := range trades {
		out += fmt.Sprintf("%s %d@%d\n", t.Symbol, t.Qty, t.Price)
	}
	return out
}

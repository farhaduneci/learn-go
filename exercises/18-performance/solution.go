//go:build solution

package perf

import (
	"strconv"
	"strings"
)

// Tick is a compact market-data update.
// Largest fields first: 8 + 4 + 4 + 1 + 1 + 1 = 19 bytes, padded to 24
// (the struct's alignment is 8). The original order padded after every bool: 32.
type Tick struct {
	Price   int64
	Size    int32
	Seq     int32
	Halted  bool
	IsBid   bool
	IsFinal bool
}

// VWAP returns the volume-weighted average price, or 0 for no volume.
func VWAP(trades []Trade) int64 {
	var notional, volume int64
	for _, t := range trades {
		notional += t.Price * t.Qty
		volume += t.Qty
	}
	if volume == 0 {
		return 0
	}
	return notional / volume
}

// FormatTrades renders one "SYMBOL qty@price" line per trade.
func FormatTrades(trades []Trade) string {
	var b strings.Builder
	b.Grow(len(trades) * 32) // one up-front allocation for typical lines
	var num [20]byte         // scratch space on the stack for number formatting
	for _, t := range trades {
		b.WriteString(t.Symbol)
		b.WriteByte(' ')
		b.Write(strconv.AppendInt(num[:0], t.Qty, 10))
		b.WriteByte('@')
		b.Write(strconv.AppendInt(num[:0], t.Price, 10))
		b.WriteByte('\n')
	}
	return b.String()
}

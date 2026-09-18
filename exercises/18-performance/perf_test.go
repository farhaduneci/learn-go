package perf

import (
	"strings"
	"testing"
	"unsafe"
)

func sampleTrades(n int) []Trade {
	syms := []string{"AAPL", "TSLA", "MSFT", "NVDA"}
	out := make([]Trade, n)
	for i := range out {
		out[i] = Trade{Symbol: syms[i%len(syms)], Price: int64(10_000 + i), Qty: int64(1 + i%100)}
	}
	return out
}

func TestTickSize(t *testing.T) {
	if unsafe.Sizeof(uintptr(0)) != 8 {
		t.Skip("layout budget assumes a 64-bit platform")
	}
	if got := unsafe.Sizeof(Tick{}); got > 24 {
		t.Errorf("sizeof(Tick) = %d, want <= 24", got)
	}
}

func TestVWAP(t *testing.T) {
	trades := []Trade{{"A", 100, 1}, {"A", 200, 3}}
	if got := VWAP(trades); got != 175 {
		t.Errorf("VWAP = %d, want 175", got)
	}
	if got := VWAP(nil); got != 0 {
		t.Errorf("VWAP(nil) = %d", got)
	}
	big := sampleTrades(1_000)
	if allocs := testing.AllocsPerRun(100, func() { VWAP(big) }); allocs != 0 {
		t.Errorf("VWAP allocates %.0f times per call, want 0", allocs)
	}
}

func TestFormatTrades(t *testing.T) {
	got := FormatTrades([]Trade{{"AAPL", 18950, 10}, {"TSLA", 25000, 1}})
	if want := "AAPL 10@18950\nTSLA 1@25000\n"; got != want {
		t.Errorf("FormatTrades = %q, want %q", got, want)
	}
	big := sampleTrades(1_000)
	if n := strings.Count(FormatTrades(big), "\n"); n != 1_000 {
		t.Fatalf("got %d lines", n)
	}
	if allocs := testing.AllocsPerRun(20, func() { FormatTrades(big) }); allocs > 2 {
		t.Errorf("FormatTrades allocates %.0f times per call, want <= 2", allocs)
	}
}

func BenchmarkVWAP(b *testing.B) {
	trades := sampleTrades(1_000)
	for b.Loop() {
		VWAP(trades)
	}
}

func BenchmarkFormatTrades(b *testing.B) {
	trades := sampleTrades(1_000)
	for b.Loop() {
		FormatTrades(trades)
	}
}

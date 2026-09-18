package orderbook

import (
	"context"
	"math/rand/v2"
	"sync/atomic"
	"testing"
)

// BenchmarkBookSubmit measures the single-threaded matching hot path.
func BenchmarkBookSubmit(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	book := NewBook()
	var id uint64
	for b.Loop() {
		id++
		o := Order{ID: id, Side: Side(1 + r.IntN(2)), Kind: Limit, Price: int64(9_990 + r.IntN(21)), Qty: int64(1 + r.IntN(100))}
		if _, err := book.Submit(o); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkEngineSubmitParallel measures throughput through the channel-based engine.
func BenchmarkEngineSubmitParallel(b *testing.B) {
	e := NewEngine()
	defer e.Close()
	ctx := context.Background()
	var nextID atomic.Uint64
	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewPCG(rand.Uint64(), 2))
		for pb.Next() {
			o := Order{ID: nextID.Add(1), Side: Side(1 + r.IntN(2)), Kind: Limit, Price: int64(9_990 + r.IntN(21)), Qty: int64(1 + r.IntN(100))}
			if _, err := e.Submit(ctx, "AAPL", o); err != nil {
				b.Error(err)
				return
			}
		}
	})
}

//go:build solution

// Package drills holds classic Go live-coding interview problems.
package drills

import (
	"container/heap"
	"container/list"
	"slices"
	"sync"
	"time"
)

type entry struct {
	key   string
	value int
}

// LRU is a fixed-capacity least-recently-used cache. Safe for concurrent use.
type LRU struct {
	mu    sync.Mutex // Get mutates recency, so even reads need the write lock
	cap   int
	order *list.List // front = most recent
	items map[string]*list.Element
}

// NewLRU returns an empty cache holding at most capacity items.
func NewLRU(capacity int) *LRU {
	return &LRU{cap: capacity, order: list.New(), items: make(map[string]*list.Element)}
}

// Get returns the value for key and marks it most recently used.
func (c *LRU) Get(key string) (int, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.items[key]
	if !ok {
		return 0, false
	}
	c.order.MoveToFront(el)
	return el.Value.(*entry).value, true
}

// Put inserts or updates key, evicting the least recently used item if full.
func (c *LRU) Put(key string, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[key]; ok {
		el.Value.(*entry).value = value
		c.order.MoveToFront(el)
		return
	}
	if c.order.Len() == c.cap {
		oldest := c.order.Back()
		c.order.Remove(oldest)
		delete(c.items, oldest.Value.(*entry).key)
	}
	c.items[key] = c.order.PushFront(&entry{key, value})
}

// Limiter is a token-bucket rate limiter.
type Limiter struct {
	mu     sync.Mutex
	rate   float64 // tokens per second
	burst  float64
	tokens float64
	last   time.Time
	now    func() time.Time
}

// NewLimiter allows rate tokens per second with bursts of up to burst.
func NewLimiter(rate float64, burst int, now func() time.Time) *Limiter {
	return &Limiter{rate: rate, burst: float64(burst), tokens: float64(burst), last: now(), now: now}
}

// Allow reports whether one event may happen now, spending a token if so.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	t := l.now()
	l.tokens = min(l.burst, l.tokens+t.Sub(l.last).Seconds()*l.rate)
	l.last = t
	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

// PingPong returns 2*n alternating "ping"/"pong" strings produced by two goroutines.
func PingPong(n int) []string {
	var out []string
	ping, pong := make(chan struct{}), make(chan struct{})
	done := make(chan struct{})
	go func() {
		for range n {
			<-ping // wait for this goroutine's turn: the channel handoff is the synchronisation
			out = append(out, "ping")
			pong <- struct{}{}
		}
	}()
	go func() {
		for i := range n {
			<-pong
			out = append(out, "pong")
			if i == n-1 {
				close(done)
				return
			}
			ping <- struct{}{}
		}
	}()
	if n == 0 {
		return nil
	}
	ping <- struct{}{}
	<-done
	return out
}

type symCount struct {
	sym string
	n   int
}

// minHeap keeps the "worst" candidate at the root so it is cheap to evict.
type minHeap []symCount

func (h minHeap) Len() int { return len(h) }
func (h minHeap) Less(i, j int) bool {
	if h[i].n != h[j].n {
		return h[i].n < h[j].n
	}
	return h[i].sym > h[j].sym // on ties, the later symbol is "worse"
}
func (h minHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)   { *h = append(*h, x.(symCount)) }
func (h *minHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

// TopK returns the k most frequent symbols, most frequent first, ties by symbol.
func TopK(symbols []string, k int) []string {
	freq := make(map[string]int)
	for _, s := range symbols {
		freq[s]++
	}
	h := &minHeap{}
	for s, n := range freq {
		heap.Push(h, symCount{s, n})
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	out := make([]string, h.Len())
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = heap.Pop(h).(symCount).sym
	}
	return slices.Clip(out)
}

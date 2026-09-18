//go:build solution

package orderbook

import (
	"container/list"
	"context"
	"fmt"
	"slices"
	"sync"
)

type resting struct {
	id  uint64
	qty int64
}

type level struct {
	price  int64
	total  int64
	orders list.List // of *resting, oldest at the front (time priority)
}

// bookSide keeps levels sorted so the BEST price is the LAST element:
// taking liquidity pops from the end, which is O(1).
type bookSide struct {
	levels []*level
	// worse reports whether price a is worse than b for this side.
	worse func(a, b int64) bool
}

func (s *bookSide) best() *level {
	if len(s.levels) == 0 {
		return nil
	}
	return s.levels[len(s.levels)-1]
}

func (s *bookSide) search(price int64) (int, bool) {
	return slices.BinarySearchFunc(s.levels, price, func(l *level, p int64) int {
		switch {
		case l.price == p:
			return 0
		case s.worse(l.price, p):
			return -1
		default:
			return 1
		}
	})
}

func (s *bookSide) levelAt(price int64) *level {
	i, found := s.search(price)
	if found {
		return s.levels[i]
	}
	l := &level{price: price}
	s.levels = slices.Insert(s.levels, i, l)
	return l
}

func (s *bookSide) remove(l *level) {
	if i, found := s.search(l.price); found {
		s.levels = slices.Delete(s.levels, i, i+1)
	}
}

type ref struct {
	el   *list.Element
	lvl  *level
	side *bookSide
}

// Book is a single-symbol limit order book. It is NOT safe for concurrent use;
// the Engine is the only writer.
type Book struct {
	bids, asks bookSide
	index      map[uint64]ref
}

// NewBook returns an empty book.
func NewBook() *Book {
	return &Book{
		bids:  bookSide{worse: func(a, b int64) bool { return a < b }}, // lower bid is worse
		asks:  bookSide{worse: func(a, b int64) bool { return a > b }}, // higher ask is worse
		index: make(map[uint64]ref),
	}
}

func validate(o Order) error {
	switch {
	case o.Qty <= 0:
		return fmt.Errorf("%w: qty %d", ErrInvalidOrder, o.Qty)
	case o.Side != Buy && o.Side != Sell:
		return fmt.Errorf("%w: side %d", ErrInvalidOrder, o.Side)
	case o.Kind != Limit && o.Kind != Market:
		return fmt.Errorf("%w: kind %d", ErrInvalidOrder, o.Kind)
	case o.Kind == Limit && o.Price <= 0:
		return fmt.Errorf("%w: limit price %d", ErrInvalidOrder, o.Price)
	}
	return nil
}

// Submit matches o against the opposite side with price-time priority and
// returns the resulting trades. A Limit remainder rests; a Market remainder is dropped.
func (b *Book) Submit(o Order) ([]Trade, error) {
	if err := validate(o); err != nil {
		return nil, err
	}
	if _, dup := b.index[o.ID]; dup {
		return nil, fmt.Errorf("%w: %d", ErrDuplicateID, o.ID)
	}

	own, opp := &b.bids, &b.asks
	if o.Side == Sell {
		own, opp = opp, own
	}

	var trades []Trade
	for o.Qty > 0 {
		lvl := opp.best()
		if lvl == nil || (o.Kind == Limit && opp.worse(lvl.price, o.Price)) {
			break // no liquidity, or the best price is beyond the limit
		}
		for o.Qty > 0 && lvl.orders.Len() > 0 {
			front := lvl.orders.Front()
			r := front.Value.(*resting)
			q := min(o.Qty, r.qty)
			trades = append(trades, Trade{TakerID: o.ID, MakerID: r.id, Price: lvl.price, Qty: q})
			o.Qty -= q
			r.qty -= q
			lvl.total -= q
			if r.qty == 0 {
				lvl.orders.Remove(front)
				delete(b.index, r.id)
			}
		}
		if lvl.orders.Len() == 0 {
			opp.levels = opp.levels[:len(opp.levels)-1]
		}
	}

	if o.Qty > 0 && o.Kind == Limit {
		lvl := own.levelAt(o.Price)
		el := lvl.orders.PushBack(&resting{id: o.ID, qty: o.Qty})
		lvl.total += o.Qty
		b.index[o.ID] = ref{el: el, lvl: lvl, side: own}
	}
	return trades, nil
}

// Cancel removes a resting order and reports whether it was found.
func (b *Book) Cancel(id uint64) bool {
	r, ok := b.index[id]
	if !ok {
		return false
	}
	delete(b.index, id)
	r.lvl.total -= r.el.Value.(*resting).qty
	r.lvl.orders.Remove(r.el)
	if r.lvl.orders.Len() == 0 {
		r.side.remove(r.lvl)
	}
	return true
}

// Depth returns up to n aggregated levels on side s, best price first.
func (b *Book) Depth(s Side, n int) []Level {
	side := &b.bids
	if s == Sell {
		side = &b.asks
	}
	out := make([]Level, 0, min(n, len(side.levels)))
	for i := len(side.levels) - 1; i >= 0 && len(out) < n; i-- {
		out = append(out, Level{Price: side.levels[i].price, Qty: side.levels[i].total})
	}
	return out
}

// Engine owns one Book per symbol and serialises all access through one goroutine.
type Engine struct {
	cmds      chan func(map[string]*Book)
	quit      chan struct{}
	stopped   chan struct{}
	closeOnce sync.Once
}

// NewEngine starts the engine's event loop.
func NewEngine() *Engine {
	e := &Engine{
		cmds:    make(chan func(map[string]*Book), 1024),
		quit:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
	go e.loop()
	return e
}

// loop is the single writer: books are only ever touched here, so they need no locks.
func (e *Engine) loop() {
	defer close(e.stopped)
	books := make(map[string]*Book)
	for {
		select {
		case cmd := <-e.cmds:
			cmd(books)
		case <-e.quit:
			return
		}
	}
}

// do runs fn on the event loop and waits for it to finish.
func (e *Engine) do(ctx context.Context, fn func(map[string]*Book)) error {
	done := make(chan struct{})
	cmd := func(books map[string]*Book) {
		fn(books)
		close(done)
	}
	select {
	case e.cmds <- cmd:
	case <-e.quit:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case <-done:
		return nil
	case <-e.stopped:
		return ErrClosed
	case <-ctx.Done():
		return ctx.Err() // the command may still run; callers must treat this as "unknown"
	}
}

func bookFor(books map[string]*Book, symbol string) *Book {
	b, ok := books[symbol]
	if !ok {
		b = NewBook()
		books[symbol] = b
	}
	return b
}

// Submit sends o to the book for symbol and waits for the trades.
func (e *Engine) Submit(ctx context.Context, symbol string, o Order) ([]Trade, error) {
	var (
		trades []Trade
		err    error
	)
	if derr := e.do(ctx, func(books map[string]*Book) {
		trades, err = bookFor(books, symbol).Submit(o)
	}); derr != nil {
		return nil, derr
	}
	return trades, err
}

// Cancel cancels a resting order on symbol's book.
func (e *Engine) Cancel(ctx context.Context, symbol string, id uint64) (bool, error) {
	var ok bool
	err := e.do(ctx, func(books map[string]*Book) { ok = bookFor(books, symbol).Cancel(id) })
	return ok, err
}

// Depth returns up to n levels of symbol's book on side s.
func (e *Engine) Depth(ctx context.Context, symbol string, s Side, n int) ([]Level, error) {
	var lv []Level
	err := e.do(ctx, func(books map[string]*Book) { lv = bookFor(books, symbol).Depth(s, n) })
	return lv, err
}

// Close stops the event loop. Later calls return ErrClosed.
func (e *Engine) Close() {
	e.closeOnce.Do(func() { close(e.quit) })
	<-e.stopped
}

package messaging

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

// fakeBroker delivers every message at least once: it duplicates some, and
// redelivers (Attempt+1) whenever the handler returns an error.
func fakeBroker(msgs []Message, h Handler, maxRedeliveries int) {
	for _, m := range msgs {
		for attempt := 1; attempt <= maxRedeliveries; attempt++ {
			m.Attempt = attempt
			err := h(context.Background(), m)
			if m.Seq%3 == 0 { // duplicate delivery after an ack was lost
				_ = h(context.Background(), m)
			}
			if err == nil {
				break
			}
		}
	}
}

func TestDedupeAppliesOnce(t *testing.T) {
	var balance int64
	apply := func(_ context.Context, m Message) error { balance += m.Amount; return nil }
	var msgs []Message
	for i := 1; i <= 9; i++ {
		msgs = append(msgs, Message{ID: fmt.Sprint("m", i), Seq: i, Amount: 10})
	}
	fakeBroker(msgs, Dedupe(apply), 5)
	if balance != 90 {
		t.Errorf("balance = %d, want 90 (duplicates applied?)", balance)
	}
}

func TestDedupeRetriesFailures(t *testing.T) {
	calls := 0
	flaky := func(context.Context, Message) error {
		calls++
		if calls == 1 {
			return errors.New("db timeout")
		}
		return nil
	}
	h := Dedupe(flaky)
	m := Message{ID: "x", Attempt: 1}
	if err := h(context.Background(), m); err == nil {
		t.Fatal("first attempt should fail")
	}
	if err := h(context.Background(), m); err != nil || calls != 2 {
		t.Fatalf("retry: err=%v calls=%d, want nil, 2", err, calls)
	}
	_ = h(context.Background(), m)
	if calls != 2 {
		t.Errorf("duplicate after success re-ran handler (calls=%d)", calls)
	}
}

func TestDedupeConcurrent(t *testing.T) {
	var mu sync.Mutex
	count := map[string]int{}
	h := Dedupe(func(_ context.Context, m Message) error {
		mu.Lock()
		count[m.ID]++
		mu.Unlock()
		return nil
	})
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() { _ = h(context.Background(), Message{ID: fmt.Sprint(i)}) })
	}
	wg.Wait()
	if len(count) != 100 {
		t.Errorf("processed %d IDs, want 100", len(count))
	}
}

func TestDeadLetter(t *testing.T) {
	poison := errors.New("cannot parse")
	var dead []Message
	h := WithDeadLetter(func(context.Context, Message) error { return poison }, 3, func(m Message, err error) {
		if !errors.Is(err, poison) {
			t.Errorf("dlq got err %v", err)
		}
		dead = append(dead, m)
	})
	for attempt := 1; attempt <= 2; attempt++ {
		if err := h(context.Background(), Message{ID: "p", Attempt: attempt}); err == nil {
			t.Fatalf("attempt %d acked; want nack before maxAttempts", attempt)
		}
	}
	if err := h(context.Background(), Message{ID: "p", Attempt: 3}); err != nil {
		t.Fatalf("final attempt err = %v, want ack after dead-lettering", err)
	}
	if len(dead) != 1 {
		t.Errorf("dlq received %d messages, want 1", len(dead))
	}
	ok := WithDeadLetter(func(context.Context, Message) error { return nil }, 1, func(Message, error) { t.Error("dlq called on success") })
	_ = ok(context.Background(), Message{Attempt: 1})
}

func TestDispatcherOrderingAndParallelism(t *testing.T) {
	var mu sync.Mutex
	last := map[string]int{}
	var outOfOrder []string
	processed := 0
	d := NewDispatcher(4, func(_ context.Context, m Message) error {
		time.Sleep(100 * time.Microsecond)
		mu.Lock()
		defer mu.Unlock()
		if m.Seq <= last[m.Key] {
			outOfOrder = append(outOfOrder, fmt.Sprintf("%s#%d after #%d", m.Key, m.Seq, last[m.Key]))
		}
		last[m.Key] = m.Seq
		processed++
		return nil
	})
	keys := []string{"acct-1", "acct-2", "acct-3", "acct-4", "acct-5", "acct-6"}
	for seq := 1; seq <= 50; seq++ {
		for _, k := range keys {
			d.Dispatch(Message{Key: k, Seq: seq})
		}
	}
	done := make(chan struct{})
	go func() { d.Close(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Close did not return")
	}
	if len(outOfOrder) > 0 {
		t.Errorf("per-key order violated: %v", outOfOrder[:min(3, len(outOfOrder))])
	}
	if processed != 300 {
		t.Errorf("processed %d, want 300 (Close must drain)", processed)
	}
}

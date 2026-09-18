package ctxx

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestSleepCancelled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := Sleep(ctx, time.Second)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Sleep err = %v, want DeadlineExceeded", err)
	}
	if time.Since(start) > 300*time.Millisecond {
		t.Error("Sleep ignored the deadline")
	}
	if err := Sleep(context.Background(), time.Millisecond); err != nil {
		t.Errorf("Sleep = %v, want nil", err)
	}
}

func venue(price int64, delay time.Duration, err error, cancelled *atomic.Int32) Venue {
	return func(ctx context.Context) (int64, error) {
		select {
		case <-time.After(delay):
			return price, err
		case <-ctx.Done():
			if cancelled != nil {
				cancelled.Add(1)
			}
			return 0, ctx.Err()
		}
	}
}

func TestFirstSuccess(t *testing.T) {
	var cancelled atomic.Int32
	price, err := FirstSuccess(context.Background(),
		venue(101, 300*time.Millisecond, nil, &cancelled),
		venue(100, 10*time.Millisecond, nil, &cancelled),
		venue(0, 1*time.Millisecond, errors.New("venue down"), &cancelled),
	)
	if err != nil || price != 100 {
		t.Fatalf("FirstSuccess = %d, %v; want 100, nil", price, err)
	}
	time.Sleep(50 * time.Millisecond)
	if cancelled.Load() != 1 {
		t.Errorf("slow venue cancelled %d times, want 1 (losers must be cancelled)", cancelled.Load())
	}
}

func TestFirstSuccessAllFail(t *testing.T) {
	errA, errB := errors.New("a down"), errors.New("b down")
	_, err := FirstSuccess(context.Background(),
		venue(0, time.Millisecond, errA, nil),
		venue(0, 2*time.Millisecond, errB, nil),
	)
	if !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Errorf("err = %v, want both venue errors joined", err)
	}
}

func TestFirstSuccessParentCancelled(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := FirstSuccess(ctx, venue(1, time.Second, nil, nil))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("err = %v, want DeadlineExceeded", err)
	}
}

func TestRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-42")
	if id, ok := RequestID(ctx); !ok || id != "req-42" {
		t.Errorf("RequestID = %q, %v", id, ok)
	}
	if _, ok := RequestID(context.Background()); ok {
		t.Error("RequestID on empty ctx ok = true")
	}
}

func TestRetry(t *testing.T) {
	var calls int
	err := Retry(context.Background(), 5, time.Millisecond, func(context.Context) error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Errorf("Retry = %v after %d calls, want nil after 3", err, calls)
	}

	calls = 0
	permanent := errors.New("permanent")
	err = Retry(context.Background(), 3, time.Millisecond, func(context.Context) error { calls++; return permanent })
	if !errors.Is(err, permanent) || calls != 3 {
		t.Errorf("Retry = %v after %d calls, want permanent after 3", err, calls)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	err = Retry(ctx, 10, 20*time.Millisecond, func(context.Context) error { return permanent })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Retry with deadline = %v, want DeadlineExceeded", err)
	}
}

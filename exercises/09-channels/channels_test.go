package channels

import (
	"slices"
	"testing"
	"time"
)

// collect drains ch, failing the test instead of hanging if ch never closes.
func collect[T any](t *testing.T, ch <-chan T) []T {
	t.Helper()
	var out []T
	deadline := time.After(2 * time.Second)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return out
			}
			out = append(out, v)
		case <-deadline:
			t.Fatalf("channel not closed within 2s (got %v so far)", out)
		}
	}
}

func TestPipeline(t *testing.T) {
	got := collect(t, Square(Gen(1, 2, 3)))
	if !slices.Equal(got, []int{1, 4, 9}) {
		t.Errorf("Square(Gen(1,2,3)) = %v, want [1 4 9]", got)
	}
}

func TestMerge(t *testing.T) {
	got := collect(t, Merge(Gen(1, 2), Gen(3), Gen(), Gen(4, 5, 6)))
	slices.Sort(got)
	if !slices.Equal(got, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Merge = %v", got)
	}
}

func TestWithTimeout(t *testing.T) {
	if v, ok := WithTimeout(func() int { return 42 }, time.Second); !ok || v != 42 {
		t.Errorf("fast f: got %d, %v", v, ok)
	}
	start := time.Now()
	_, ok := WithTimeout(func() int { time.Sleep(500 * time.Millisecond); return 1 }, 20*time.Millisecond)
	if ok {
		t.Error("slow f: ok = true, want timeout")
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Errorf("WithTimeout waited %v, should return at the deadline", time.Since(start))
	}
}

func TestBatchBySize(t *testing.T) {
	got := collect(t, Batch(Gen(1, 2, 3, 4, 5), 2, time.Second))
	want := [][]int{{1, 2}, {3, 4}, {5}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("Batch = %v, want %v", got, want)
	}
}

func TestBatchByTime(t *testing.T) {
	in := make(chan int)
	out := Batch(in, 100, 30*time.Millisecond)
	for _, v := range []int{1, 2} {
		select {
		case in <- v:
		case <-time.After(time.Second):
			t.Fatal("Batch is not reading its input")
		}
	}
	select {
	case b := <-out:
		if !slices.Equal(b, []int{1, 2}) {
			t.Errorf("timed flush = %v, want [1 2]", b)
		}
	case <-time.After(time.Second):
		t.Fatal("partial batch not flushed after wait")
	}
	close(in)
	if rest := collect(t, out); len(rest) != 0 {
		t.Errorf("unexpected trailing batches %v", rest)
	}
}

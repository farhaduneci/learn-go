package collections

import (
	"maps"
	"slices"
	"testing"
)

func TestReversed(t *testing.T) {
	in := []int{1, 2, 3}
	got := Reversed(in)
	if !slices.Equal(got, []int{3, 2, 1}) {
		t.Errorf("Reversed = %v, want [3 2 1]", got)
	}
	if !slices.Equal(in, []int{1, 2, 3}) {
		t.Errorf("input modified to %v", in)
	}
	if got := Reversed([]int{}); len(got) != 0 {
		t.Errorf("Reversed(empty) = %v", got)
	}
}

func TestDedupe(t *testing.T) {
	got := Dedupe([]string{"AAPL", "TSLA", "AAPL", "MSFT", "TSLA"})
	if want := []string{"AAPL", "TSLA", "MSFT"}; !slices.Equal(got, want) {
		t.Errorf("Dedupe = %v, want %v", got, want)
	}
}

func TestWordFreq(t *testing.T) {
	got := WordFreq("Buy buy SELL\n hold  buy")
	want := map[string]int{"buy": 3, "sell": 1, "hold": 1}
	if !maps.Equal(got, want) {
		t.Errorf("WordFreq = %v, want %v", got, want)
	}
}

func TestReverseString(t *testing.T) {
	for in, want := range map[string]string{"héllo": "olléh", "日本": "本日", "": "", "go": "og"} {
		if got := ReverseString(in); got != want {
			t.Errorf("ReverseString(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestChunk(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	got := Chunk(s, 2)
	if len(got) != 3 || !slices.Equal(got[0], []int{1, 2}) || !slices.Equal(got[2], []int{5}) {
		t.Fatalf("Chunk = %v, want [[1 2] [3 4] [5]]", got)
	}
	_ = append(got[0], 99) // must not clobber got[1][0]
	if got[1][0] != 3 {
		t.Errorf("append to chunk 0 overwrote chunk 1: %v", got[1])
	}
}

func TestTopSymbols(t *testing.T) {
	vol := map[string]int{"AAPL": 50, "TSLA": 90, "MSFT": 50, "NVDA": 10}
	got := TopSymbols(vol, 3)
	if want := []string{"TSLA", "AAPL", "MSFT"}; !slices.Equal(got, want) {
		t.Errorf("TopSymbols = %v, want %v", got, want)
	}
	if got := TopSymbols(vol, 10); len(got) != 4 {
		t.Errorf("TopSymbols(k>len) returned %d items, want 4", len(got))
	}
}

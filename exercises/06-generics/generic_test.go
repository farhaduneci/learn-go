package generic

import (
	"slices"
	"strconv"
	"testing"
)

func TestMapFilterReduce(t *testing.T) {
	strs := Map([]int{1, 2, 3}, strconv.Itoa)
	if !slices.Equal(strs, []string{"1", "2", "3"}) {
		t.Errorf("Map = %v", strs)
	}
	even := Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
	if !slices.Equal(even, []int{2, 4}) {
		t.Errorf("Filter = %v", even)
	}
	sum := Reduce([]int{1, 2, 3}, 0.5, func(acc float64, n int) float64 { return acc + float64(n) })
	if sum != 6.5 {
		t.Errorf("Reduce = %v, want 6.5", sum)
	}
}

func TestMaxOf(t *testing.T) {
	if m, ok := MaxOf([]string{"b", "c", "a"}); !ok || m != "c" {
		t.Errorf("MaxOf = %q, %v", m, ok)
	}
	if _, ok := MaxOf([]float64{}); ok {
		t.Error("MaxOf(empty) ok = true")
	}
}

func TestSet(t *testing.T) {
	var s Set[string] // zero value must work
	s.Add("TSLA")
	s.Add("AAPL")
	s.Add("TSLA")
	if s.Len() != 2 || !s.Has("AAPL") || s.Has("MSFT") {
		t.Errorf("Set state wrong: len=%d", s.Len())
	}
	if got := SortedItems(&s); !slices.Equal(got, []string{"AAPL", "TSLA"}) {
		t.Errorf("SortedItems = %v", got)
	}
	var empty Set[int]
	if empty.Has(1) || empty.Len() != 0 {
		t.Error("zero-value Set should be empty")
	}
}

func TestStack(t *testing.T) {
	var s Stack[int]
	if _, ok := s.Pop(); ok {
		t.Fatal("Pop on empty stack ok = true")
	}
	s.Push(1)
	s.Push(2)
	if v, _ := s.Peek(); v != 2 {
		t.Errorf("Peek = %d, want 2", v)
	}
	if v, _ := s.Pop(); v != 2 {
		t.Errorf("Pop = %d, want 2", v)
	}
	if v, ok := s.Pop(); !ok || v != 1 {
		t.Errorf("Pop = %d, %v, want 1, true", v, ok)
	}
}

func TestGroupBy(t *testing.T) {
	type trade struct {
		sym string
		qty int
	}
	trades := []trade{{"A", 1}, {"B", 2}, {"A", 3}}
	g := GroupBy(trades, func(tr trade) string { return tr.sym })
	if len(g) != 2 || len(g["A"]) != 2 || g["A"][1].qty != 3 {
		t.Errorf("GroupBy = %v", g)
	}
}

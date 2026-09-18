//go:build solution

// Package generic practises type parameters and constraints.
package generic

import (
	"cmp"
	"slices"
)

// Map returns f applied to every element of s.
func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		out = append(out, f(v))
	}
	return out
}

// Filter returns the elements of s for which keep returns true.
func Filter[T any](s []T, keep func(T) bool) []T {
	var out []T
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// Reduce folds s into a single value starting from init.
func Reduce[T, A any](s []T, init A, f func(A, T) A) A {
	acc := init
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

// MaxOf returns the largest element of s; ok is false when s is empty.
func MaxOf[T cmp.Ordered](s []T) (m T, ok bool) {
	if len(s) == 0 {
		return m, false
	}
	return slices.Max(s), true
}

// Set is an unordered collection of unique values.
type Set[T comparable] struct {
	m map[T]struct{}
}

// Add inserts v. It must work on a zero-value Set.
func (s *Set[T]) Add(v T) {
	if s.m == nil {
		s.m = make(map[T]struct{})
	}
	s.m[v] = struct{}{}
}

// Has reports whether v is in the set.
func (s *Set[T]) Has(v T) bool {
	_, ok := s.m[v] // reading a nil map is fine
	return ok
}

// Len returns the number of elements.
func (s *Set[T]) Len() int { return len(s.m) }

// SortedItems returns the elements of s in ascending order.
func SortedItems[T cmp.Ordered](s *Set[T]) []T {
	out := make([]T, 0, len(s.m))
	for v := range s.m {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

// Stack is a LIFO stack. Its zero value is an empty stack.
type Stack[T any] struct {
	items []T
}

// Push adds v to the top.
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop removes and returns the top element; ok is false when empty.
func (s *Stack[T]) Pop() (v T, ok bool) {
	v, ok = s.Peek()
	if ok {
		var zero T
		s.items[len(s.items)-1] = zero // drop the reference so the GC can collect it
		s.items = s.items[:len(s.items)-1]
	}
	return v, ok
}

// Peek returns the top element without removing it.
func (s *Stack[T]) Peek() (v T, ok bool) {
	if len(s.items) == 0 {
		return v, false
	}
	return s.items[len(s.items)-1], true
}

// GroupBy groups the elements of s by key(element), keeping order within each group.
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	out := make(map[K][]T)
	for _, v := range s {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

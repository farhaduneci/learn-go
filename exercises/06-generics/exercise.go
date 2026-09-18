//go:build !solution

// Package generic practises type parameters and constraints.
package generic

import "cmp"

// Map returns f applied to every element of s.
func Map[T, U any](s []T, f func(T) U) []U {
	// TODO: implement
	return nil
}

// Filter returns the elements of s for which keep returns true.
func Filter[T any](s []T, keep func(T) bool) []T {
	// TODO: implement
	return nil
}

// Reduce folds s into a single value starting from init.
func Reduce[T, A any](s []T, init A, f func(A, T) A) A {
	// TODO: implement
	return init
}

// MaxOf returns the largest element of s; ok is false when s is empty.
func MaxOf[T cmp.Ordered](s []T) (m T, ok bool) {
	// TODO: implement
	return m, false
}

// Set is an unordered collection of unique values.
type Set[T comparable] struct {
	m map[T]struct{}
}

// Add inserts v. It must work on a zero-value Set.
func (s *Set[T]) Add(v T) {
	// TODO: implement
}

// Has reports whether v is in the set.
func (s *Set[T]) Has(v T) bool {
	// TODO: implement
	return false
}

// Len returns the number of elements.
func (s *Set[T]) Len() int {
	// TODO: implement
	return 0
}

// SortedItems returns the elements of s in ascending order.
func SortedItems[T cmp.Ordered](s *Set[T]) []T {
	// TODO: implement
	return nil
}

// Stack is a LIFO stack. Its zero value is an empty stack.
type Stack[T any] struct {
	items []T
}

// Push adds v to the top.
func (s *Stack[T]) Push(v T) {
	// TODO: implement
}

// Pop removes and returns the top element; ok is false when empty.
func (s *Stack[T]) Pop() (v T, ok bool) {
	// TODO: implement
	return v, false
}

// Peek returns the top element without removing it.
func (s *Stack[T]) Peek() (v T, ok bool) {
	// TODO: implement
	return v, false
}

// GroupBy groups the elements of s by key(element), keeping order within each group.
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	// TODO: implement
	return nil
}

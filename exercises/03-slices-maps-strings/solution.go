//go:build solution

// Package collections practises slices, maps and strings.
package collections

import (
	"cmp"
	"slices"
	"strings"
)

// Reversed returns a reversed copy of s without modifying s.
func Reversed(s []int) []int {
	out := make([]int, len(s))
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

// Dedupe returns s without duplicates, keeping first-seen order.
func Dedupe(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	out := make([]string, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

// WordFreq counts whitespace-separated, lower-cased words in text.
func WordFreq(text string) map[string]int {
	freq := make(map[string]int)
	for _, w := range strings.Fields(text) {
		freq[strings.ToLower(w)]++
	}
	return freq
}

// ReverseString reverses s by rune.
func ReverseString(s string) string {
	r := []rune(s)
	slices.Reverse(r)
	return string(r)
}

// Chunk splits s into consecutive chunks of at most n elements.
// Appending to a returned chunk must never modify another chunk.
func Chunk(s []int, n int) [][]int {
	if n <= 0 {
		return nil
	}
	var out [][]int
	for lo := 0; lo < len(s); lo += n {
		hi := min(lo+n, len(s))
		out = append(out, s[lo:hi:hi]) // cap == len, so append must reallocate
	}
	return out
}

// TopSymbols returns the k symbols with the highest volume,
// highest first, ties broken by symbol ascending.
func TopSymbols(volumes map[string]int, k int) []string {
	syms := make([]string, 0, len(volumes))
	for s := range volumes {
		syms = append(syms, s)
	}
	slices.SortFunc(syms, func(a, b string) int {
		if c := cmp.Compare(volumes[b], volumes[a]); c != 0 {
			return c
		}
		return cmp.Compare(a, b)
	})
	return syms[:min(k, len(syms))]
}

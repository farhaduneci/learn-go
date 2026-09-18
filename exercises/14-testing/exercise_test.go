//go:build !solution

package pricing

import "testing"

func TestParsePrice(t *testing.T) {
	t.Skip("TODO: write a table-driven test with t.Run subtests")
}

func FuzzRoundTrip(f *testing.F) {
	// TODO: f.Add a few seed values, then assert
	// ParsePrice(FormatPrice(c)) == c inside f.Fuzz.
	f.Fuzz(func(t *testing.T, c int64) {
		t.Skip("TODO")
	})
}

// TODO: BenchmarkParsePrice using b.Loop().

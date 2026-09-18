//go:build solution

package pricing

import (
	"errors"
	"math"
	"testing"
)

func TestParsePrice(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    int64
		wantErr bool
	}{
		{name: "zero", in: "0", want: 0},
		{name: "whole", in: "12", want: 1200},
		{name: "one decimal", in: "12.3", want: 1230},
		{name: "two decimals", in: "12.34", want: 1234},
		{name: "negative", in: "-5.00", want: -500},
		{name: "negative under one", in: "-0.50", want: -50}, // found by FuzzRoundTrip
		{name: "empty", in: "", wantErr: true},
		{name: "letters", in: "abc", wantErr: true},
		{name: "three decimals", in: "1.234", wantErr: true},
		{name: "two dots", in: "1.2.3", wantErr: true},
		{name: "trailing dot", in: "1.", wantErr: true},
		{name: "negative fraction", in: "1.-5", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePrice(tt.in)
			if tt.wantErr {
				if !errors.Is(err, ErrBadPrice) {
					t.Fatalf("ParsePrice(%q) err = %v, want ErrBadPrice", tt.in, err)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Errorf("ParsePrice(%q) = %d, %v; want %d", tt.in, got, err, tt.want)
			}
		})
	}
}

func FuzzRoundTrip(f *testing.F) {
	for _, seed := range []int64{0, 1, -1, 99, 1234, math.MaxInt64 / 100} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, c int64) {
		if c == math.MinInt64 {
			t.Skip("no positive counterpart")
		}
		s := FormatPrice(c)
		got, err := ParsePrice(s)
		if err != nil || got != c {
			t.Errorf("ParsePrice(FormatPrice(%d) = %q) = %d, %v", c, s, got, err)
		}
	})
}

func BenchmarkParsePrice(b *testing.B) {
	for b.Loop() {
		_, _ = ParsePrice("12345.67")
	}
}

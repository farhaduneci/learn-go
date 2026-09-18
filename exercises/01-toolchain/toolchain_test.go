package toolchain

import "testing"

func TestHello(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"named", "Gopher", "Hello, Gopher!"},
		{"empty", "", "Hello, world!"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hello(tt.in); got != tt.want {
				t.Errorf("Hello(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSum(t *testing.T) {
	if got := Sum(); got != 0 {
		t.Errorf("Sum() = %d, want 0", got)
	}
	if got := Sum(1, 2, 3); got != 6 {
		t.Errorf("Sum(1, 2, 3) = %d, want 6", got)
	}
	nums := []int{10, -4}
	if got := Sum(nums...); got != 6 {
		t.Errorf("Sum(nums...) = %d, want 6", got)
	}
}

func TestStats(t *testing.T) {
	lo, hi, ok := Stats([]int{3, -1, 7, 2})
	if !ok || lo != -1 || hi != 7 {
		t.Errorf("Stats = (%d, %d, %v), want (-1, 7, true)", lo, hi, ok)
	}
	if _, _, ok := Stats(nil); ok {
		t.Error("Stats(nil) ok = true, want false")
	}
}

//go:build solution

// Package toolchain is the first exercise: a warm-up for go test, go vet and gofmt.
package toolchain

import "fmt"

// Hello returns a greeting for name, or "Hello, world!" when name is empty.
func Hello(name string) string {
	if name == "" {
		name = "world"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Sum returns the sum of nums.
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Stats returns the smallest and largest value in nums.
// ok is false when nums is empty.
func Stats(nums []int) (lo, hi int, ok bool) {
	if len(nums) == 0 {
		return 0, 0, false
	}
	lo, hi = nums[0], nums[0]
	for _, n := range nums[1:] {
		lo = min(lo, n)
		hi = max(hi, n)
	}
	return lo, hi, true
}

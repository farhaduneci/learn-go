//go:build solution

// Package pricing converts between decimal price strings and integer cents.
package pricing

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrBadPrice is returned for malformed input.
var ErrBadPrice = errors.New("bad price")

// ParsePrice parses "12", "12.3" or "12.34" (optionally negative) into cents.
func ParsePrice(s string) (int64, error) {
	// fix: "-0.50" parsed as +50 because the sign lived on w, and w == 0.
	// Take the sign off the string first, then parse unsigned parts.
	neg := strings.HasPrefix(s, "-")
	s = strings.TrimPrefix(s, "-")

	whole, frac, hasFrac := strings.Cut(s, ".")
	if whole == "" || (hasFrac && (frac == "" || len(frac) > 2)) {
		return 0, fmt.Errorf("%w: %q", ErrBadPrice, s)
	}
	w, err := strconv.ParseUint(whole, 10, 63)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrBadPrice, s)
	}
	var f uint64
	if hasFrac {
		if len(frac) == 1 {
			frac += "0"
		}
		f, err = strconv.ParseUint(frac, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("%w: %q", ErrBadPrice, s)
		}
	}
	cents := w*100 + f
	if neg {
		return -int64(cents), nil
	}
	return int64(cents), nil
}

// FormatPrice formats cents as "12.34" or "-0.05".
func FormatPrice(c int64) string {
	sign := ""
	u := uint64(c)
	if c < 0 {
		sign = "-"
		u = -u
	}
	return fmt.Sprintf("%s%d.%02d", sign, u/100, u%100)
}

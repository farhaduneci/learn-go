//go:build !solution

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
	whole, frac, hasFrac := strings.Cut(s, ".")
	if whole == "" || (hasFrac && (frac == "" || len(frac) > 2)) {
		return 0, fmt.Errorf("%w: %q", ErrBadPrice, s)
	}
	w, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrBadPrice, s)
	}
	var f int64
	if hasFrac {
		if len(frac) == 1 {
			frac += "0"
		}
		f, err = strconv.ParseInt(frac, 10, 64)
		if err != nil || f < 0 {
			return 0, fmt.Errorf("%w: %q", ErrBadPrice, s)
		}
	}
	if w < 0 {
		return w*100 - f, nil
	}
	return w*100 + f, nil
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

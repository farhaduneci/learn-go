//go:build !solution

package refactor

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var positions = map[string]int64{}

// ProcessOrders is the legacy entry point. It works, and it can't be tested.
func ProcessOrders(input string) string {
	var out string
	now := time.Now()
	lines := strings.Split(input, "\n")
	for i := 0; i < len(lines); i++ {
		l := strings.TrimSpace(lines[i])
		if l == "" {
			continue
		}
		parts := strings.Fields(l)
		if len(parts) != 3 {
			out += "line " + strconv.Itoa(i+1) + ": bad format\n"
			continue
		}
		side := strings.ToUpper(parts[0])
		sym := strings.ToUpper(parts[1])
		qty, err := strconv.Atoi(parts[2])
		if err != nil || qty <= 0 {
			out += "line " + strconv.Itoa(i+1) + ": bad qty\n"
			continue
		}
		if now.Hour() < 9 || now.Hour() >= 16 {
			out += "line " + strconv.Itoa(i+1) + ": market closed\n"
			continue
		}
		if side == "BUY" {
			positions[sym] += int64(qty)
		} else if side == "SELL" {
			if positions[sym] < int64(qty) {
				out += "line " + strconv.Itoa(i+1) + ": insufficient position\n"
				continue
			}
			positions[sym] -= int64(qty)
		} else {
			out += "line " + strconv.Itoa(i+1) + ": bad side\n"
			continue
		}
		fmt.Println("NOTIFY: filled", side, sym, qty)
		out += "line " + strconv.Itoa(i+1) + ": ok\n"
	}
	return out
}

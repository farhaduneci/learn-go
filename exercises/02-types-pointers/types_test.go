package types

import (
	"fmt"
	"testing"
)

func TestCentsString(t *testing.T) {
	tests := []struct {
		in   Cents
		want string
	}{
		{1234, "$12.34"},
		{5, "$0.05"},
		{0, "$0.00"},
		{-5, "-$0.05"},
		{100000, "$1000.00"},
	}
	for _, tt := range tests {
		if got := fmt.Sprint(tt.in); got != tt.want {
			t.Errorf("fmt.Sprint(Cents(%d)) = %q, want %q", int64(tt.in), got, tt.want)
		}
	}
}

func TestAccount(t *testing.T) {
	a := NewAccount("acc-1")
	if a == nil {
		t.Fatal("NewAccount returned nil")
	}
	if a.ID != "acc-1" || a.Balance != 0 {
		t.Fatalf("NewAccount = %+v, want ID acc-1 and zero balance", *a)
	}
	a.Deposit(500)
	if a.Balance != 500 {
		t.Fatalf("after Deposit(500) balance = %d, want 500 (value receiver bug?)", a.Balance)
	}
	if a.Withdraw(600) {
		t.Error("Withdraw(600) on 500 succeeded, want refusal")
	}
	if !a.Withdraw(200) || a.Balance != 300 {
		t.Errorf("Withdraw(200) -> balance %d, want 300", a.Balance)
	}
}

func TestSwap(t *testing.T) {
	x, y := 1, 2
	Swap(&x, &y)
	if x != 2 || y != 1 {
		t.Errorf("after Swap x=%d y=%d, want 2 1", x, y)
	}
}

func TestCounter(t *testing.T) {
	c1, c2 := NewCounter(), NewCounter()
	for want := 1; want <= 3; want++ {
		if got := c1(); got != want {
			t.Fatalf("c1() = %d, want %d", got, want)
		}
	}
	if got := c2(); got != 1 {
		t.Errorf("c2() = %d, want 1 (counters must be independent)", got)
	}
}

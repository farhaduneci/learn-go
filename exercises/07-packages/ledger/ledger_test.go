package ledger_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/farhaduneci/learn-go/exercises/07-packages/ledger"
)

func TestPostAndBalance(t *testing.T) {
	l := ledger.New()
	if err := l.Post("cash", "alice", 1_000); err != nil {
		t.Fatal(err)
	}
	if err := l.Post("alice", "bob", 250); err != nil {
		t.Fatal(err)
	}
	for acct, want := range map[string]int64{"cash": -1_000, "alice": 750, "bob": 250, "nobody": 0} {
		if got := l.Balance(acct); got != want {
			t.Errorf("Balance(%q) = %d, want %d", acct, got, want)
		}
	}
	if got := l.Sum(); got != 0 {
		t.Errorf("Sum = %d, want 0", got)
	}
}

func TestPostErrors(t *testing.T) {
	l := ledger.New()
	if err := l.Post("a", "b", 0); !errors.Is(err, ledger.ErrInvalidAmount) {
		t.Errorf("Post amount 0 = %v, want ErrInvalidAmount", err)
	}
	if err := l.Post("a", "a", 5); !errors.Is(err, ledger.ErrSameAccount) {
		t.Errorf("Post same account = %v, want ErrSameAccount", err)
	}
	if got := l.Sum(); got != 0 {
		t.Errorf("rejected postings changed balances: Sum = %d", got)
	}
}

func TestAccountsIsSortedCopy(t *testing.T) {
	l := ledger.New()
	_ = l.Post("zeta", "alpha", 1)
	got := l.Accounts()
	if !slices.Equal(got, []string{"alpha", "zeta"}) {
		t.Fatalf("Accounts = %v", got)
	}
	got[0] = "mutated"
	if l.Accounts()[0] != "alpha" {
		t.Error("Accounts must return a copy")
	}
}

package interfaces

import (
	"slices"
	"testing"
)

// Compile-time checks: these lines fail to build until the types satisfy FeeModel.
var (
	_ FeeModel = FlatFee(0)
	_ FeeModel = BpsFee(0)
	_ FeeModel = Capped{}
)

func TestFees(t *testing.T) {
	tests := []struct {
		name     string
		model    FeeModel
		notional int64
		want     int64
	}{
		{"flat", FlatFee(100), 1_000_000, 100},
		{"bps", BpsFee(5), 1_000_000, 500},
		{"bps rounds down", BpsFee(1), 9_999, 0},
		{"capped under", Capped{Inner: BpsFee(5), Max: 1_000}, 1_000_000, 500},
		{"capped over", Capped{Inner: BpsFee(5), Max: 300}, 1_000_000, 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.model.Fee(tt.notional); got != tt.want {
				t.Errorf("Fee(%d) = %d, want %d", tt.notional, got, tt.want)
			}
		})
	}
}

func TestTotalFee(t *testing.T) {
	models := []FeeModel{FlatFee(100), BpsFee(10), Capped{Inner: FlatFee(50), Max: 20}}
	if got := TotalFee(models, 100_000); got != 100+100+20 {
		t.Errorf("TotalFee = %d, want 220", got)
	}
}

func TestDescribe(t *testing.T) {
	cases := map[string]any{"flat": FlatFee(1), "bps": BpsFee(1), "capped": Capped{}, "unknown": 42}
	for want, v := range cases {
		if got := Describe(v); got != want {
			t.Errorf("Describe(%T) = %q, want %q", v, got, want)
		}
	}
}

func TestAuditedAccount(t *testing.T) {
	a := &AuditedAccount{Ledger: &Ledger{}}
	a.Post(10)
	a.Post(-3)
	if a.Total != 7 {
		t.Errorf("Total = %d, want 7 (promoted field from embedded Ledger)", a.Total)
	}
	if !slices.Equal(a.Log, []int64{10, -3}) {
		t.Errorf("Log = %v, want [10 -3]", a.Log)
	}
}

//go:build !solution

package interfaces

// FlatFee charges a fixed fee per trade.
type FlatFee int64

// TODO: implement Fee for FlatFee.

// BpsFee charges basis points of notional (1 bps = 1/10000), rounded down.
type BpsFee int64

// TODO: implement Fee for BpsFee.

// Capped caps the fee of an inner model at Max.
type Capped struct {
	Inner FeeModel
	Max   int64
}

// TODO: implement Fee for Capped.

// TotalFee sums the fee of every model.
func TotalFee(models []FeeModel, notional int64) int64 {
	// TODO: implement
	return 0
}

// Describe names the concrete fee model type.
func Describe(v any) string {
	// TODO: implement with a type switch
	return ""
}

// AuditedAccount is a Ledger that also keeps a log of every posting.
type AuditedAccount struct {
	*Ledger
	Log []int64
}

// TODO: implement Post on *AuditedAccount so it logs and then delegates to the embedded Ledger.

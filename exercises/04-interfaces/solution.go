//go:build solution

package interfaces

// FlatFee charges a fixed fee per trade.
type FlatFee int64

// Fee implements FeeModel.
func (f FlatFee) Fee(int64) int64 { return int64(f) }

// BpsFee charges basis points of notional (1 bps = 1/10000), rounded down.
type BpsFee int64

// Fee implements FeeModel.
func (b BpsFee) Fee(notional int64) int64 { return notional * int64(b) / 10_000 }

// Capped caps the fee of an inner model at Max.
type Capped struct {
	Inner FeeModel
	Max   int64
}

// Fee implements FeeModel.
func (c Capped) Fee(notional int64) int64 { return min(c.Inner.Fee(notional), c.Max) }

// TotalFee sums the fee of every model.
func TotalFee(models []FeeModel, notional int64) int64 {
	var total int64
	for _, m := range models {
		total += m.Fee(notional)
	}
	return total
}

// Describe names the concrete fee model type.
func Describe(v any) string {
	switch v.(type) {
	case FlatFee:
		return "flat"
	case BpsFee:
		return "bps"
	case Capped:
		return "capped"
	default:
		return "unknown"
	}
}

// AuditedAccount is a Ledger that also keeps a log of every posting.
type AuditedAccount struct {
	*Ledger
	Log []int64
}

// Post logs amount and then delegates to the embedded Ledger.
func (a *AuditedAccount) Post(amount int64) {
	a.Log = append(a.Log, amount)
	a.Ledger.Post(amount)
}

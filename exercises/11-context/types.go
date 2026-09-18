// Package ctxx practises context cancellation, deadlines and values.
package ctxx

import "context"

// Venue returns a price quote, honouring ctx cancellation.
type Venue func(ctx context.Context) (int64, error)

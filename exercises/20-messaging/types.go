// Package messaging builds at-least-once consumer primitives.
package messaging

import "context"

// Message is one delivery from a broker. Attempt starts at 1.
type Message struct {
	ID      string
	Key     string
	Seq     int
	Amount  int64
	Attempt int
}

// Handler processes a message. A nil error acks it; an error nacks it.
type Handler func(ctx context.Context, m Message) error

---
title: Messaging from Go: Kafka, NATS & delivery guarantees
section: Part 3 · Senior Go
summary: Kafka/NATS from Go, at-least-once, idempotent consumers, DLQs, per-key ordering, outbox.
description: Job postings for trading platforms ask for experience with message queues and brokers. This lesson covers the Go-side mechanics (clients, consumer loops, acks, idempotency, ordering by key) that make a consumer correct under at-least-once delivery. Broker internals and system design belong to a separate system-design track.
source: <a href="https://docs.nats.io/nats-concepts/jetstream">NATS JetStream concepts</a>, <a href="https://kafka.apache.org/documentation/#semantics">Kafka — Message Delivery Semantics</a>, and the Go clients <a href="https://github.com/nats-io/nats.go">nats.go</a> and <a href="https://github.com/twmb/franz-go">franz-go</a>. Background: <em>Designing Data-Intensive Applications</em>, chapter 11.
quiz: Most brokers guarantee at-most-once, at-least-once, or exactly-once by default? || at-least-once
quiz: Making a handler safe to run twice on the same message is called what? || idempotency
quiz: In Kafka, ordering is guaranteed only within a what? || partition
quiz: A queue where repeatedly failing messages are parked is abbreviated what? || DLQ
quiz: Writing events to a DB table in the same transaction as the state change is which pattern? || outbox
---

## Why brokers show up everywhere in trading

Order entry, the matching engine, risk checks, market data, the ledger and client notifications are separate services. Brokers decouple them in time (a consumer can be down briefly), in rate (bursts are buffered), and in fan-out (one fill feeds positions, P&L, notifications and audit). The common ones:

| | Kafka | NATS / JetStream | Google Pub/Sub | RabbitMQ |
| --- | --- | --- | --- | --- |
| Model | partitioned, replicated log | subjects; JetStream adds persisted streams | managed topics/subscriptions | queues + exchanges |
| Ordering | per partition | per stream/consumer | optional per ordering key | per queue |
| Replay | yes (offsets) | yes (JetStream) | limited (seek) | no |
| Go client | franz-go, confluent-kafka-go (cgo), segmentio/kafka-go | nats.go (official) | cloud.google.com/go/pubsub | amqp091-go |

NATS is itself written in Go and is common in low-latency Go shops. Kafka is the most common event log overall.

## Delivery guarantees

- **At-most-once:** ack before processing. A crash loses the message.
- **At-least-once:** ack **after** processing succeeds. A crash between processing and ack means **redelivery**. This is the usual default.
- **Exactly-once:** not something the network provides end to end. It's built from at-least-once delivery plus **idempotent processing**, and Kafka transactions only cover Kafka-to-Kafka flows.

So the practical rule is: **every consumer must be idempotent.**

## Idempotent consumers

```go
func Dedupe(next Handler) Handler {
	return func(ctx context.Context, m Message) error {
		if seen(m.ID) { return nil }   // duplicate: ack without re-applying
		if err := next(ctx, m); err != nil { return err }
		markSeen(m.ID)                 // only after success
		return nil
	}
}
```

In production, "seen" is a unique constraint in the same **database transaction** as the state change (`INSERT INTO processed(msg_id) ...` plus `UPDATE positions ...`), so the check and the effect can't diverge. An in-memory map is fine for learning and wrong across restarts.

## Retries and dead letters

A message that always fails (a "poison" message) must not block its partition forever. Retry with backoff for **transient** errors, and after N attempts, park it in a **dead-letter queue** with the error attached, alert on it, and move on. Distinguish error kinds: a validation error will never succeed on retry, so dead-letter it immediately.

## Ordering: per key, not global

Global ordering doesn't scale. Order is only needed **per entity**: all events for account 42 in sequence, while accounts can proceed in parallel. Kafka does this by hashing the message key to a partition. Inside a Go consumer the same idea keeps parallelism without breaking order:

```go
func (d *Dispatcher) Dispatch(m Message) {
	h := fnv.New32a()
	h.Write([]byte(m.Key))
	d.parts[h.Sum32()%uint32(len(d.parts))] <- m // one goroutine per partition = in-order per key
}
```

## The dual-write problem and the outbox

"Update the DB, then publish the event" fails when the process dies between the two steps. **Transactional outbox**: write the event to an `outbox` table in the **same transaction** as the state change, and have a relay publish from that table (polling or CDC). Publishing becomes at-least-once, and consumers are idempotent anyway. This is standard system-design material for any ledger or order system.

## Consumer loop hygiene in Go

- One goroutine reads from the client, and processing fans out by key (as above) with bounded buffers for backpressure.
- Commit offsets or ack only **after** processing. With batching, commit the highest contiguous processed offset.
- Handle rebalances: stop processing revoked partitions and flush before giving them up.
- Shutdown on context cancel: stop fetching, drain in-flight work, commit, close. The same shape as HTTP graceful shutdown.
- Export consumer lag as a metric. It's the first thing to check in an incident.

## Exercise

[`exercises/20-messaging`](../exercises/20-messaging/README.md): an idempotent-consumer middleware, retry with a dead-letter queue, and a key-partitioned dispatcher that preserves per-key order, tested against a fake broker that duplicates and redelivers. Capstone stage 4 then wires this to real NATS.

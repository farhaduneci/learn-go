# Exercise 20 — At-least-once consumers

Lesson: [20 · Messaging from Go: Kafka, NATS & delivery guarantees](../../lessons/0020-messaging.html)

```sh
go test -race ./exercises/20-messaging/
```

Brokers such as Kafka, NATS JetStream and Pub/Sub deliver **at least once**, so duplicates and retries are normal. These three pieces make a consumer correct anyway. They don't depend on any particular broker.

1. **Idempotency.** `Dedupe(next)` runs `next` at most once per message `ID`. A redelivered duplicate returns `nil` (ack) without re-running. An ID only counts as seen once `next` **succeeds**, so a failed attempt can retry. It must be safe for concurrent use.
2. **Retry + dead-letter queue.** `WithDeadLetter(next, maxAttempts, dlq)`: if `next` fails and `m.Attempt < maxAttempts`, return the error (nack, and the broker redelivers). On the final attempt, hand the message to `dlq` and return `nil`, so one poison message can't block the partition forever.
3. **Per-key ordering, cross-key parallelism.** `NewDispatcher(partitions, h)` routes each message to `hash(Key) % partitions`, and each partition is one goroutine with its own channel. Messages for one key (say one account) are processed **in order**, while different keys run in parallel. Kafka's partition model works the same way. `Close()` drains everything and waits.

`broker_test.go` contains a fake broker that duplicates and redelivers messages, so the tests exercise real failure modes.

Reference solution: `solution.go`.

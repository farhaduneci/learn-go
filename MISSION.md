# Mission: Go

## Why
The goal is working, senior-level Go proficiency for backend roles that build low-latency, high-throughput systems. The benchmark role is a core-trading engineer at a broker/exchange-style company: Go services, heavy concurrency, message brokers, distributed systems, performance and observability. The learner comes from Python, so every Go concept is taught against its Python equivalent and points out where Python habits cause bugs in Go.

## Success looks like
- Writes idiomatic Go that passes a strict code review: errors, interfaces, packages, naming, testing.
- Designs and debugs concurrent code with confidence: goroutines, channels, `sync`, `context`, `errgroup`. Can explain the scheduler and the memory model, and can find leaks, deadlocks and races with tooling.
- Refactors untestable Go into small, well-bounded, testable packages in safe steps.
- Measures and improves performance with benchmarks, pprof, allocation control and GC tuning, and reports latency as percentiles.
- Builds production services: HTTP/gRPC, graceful shutdown, structured logs, metrics, tracing, idempotent message consumers.
- Has built and can explain the capstone matching engine (single-writer concurrency, price-time priority, p99 latency).
- Answers the questions in lesson 21 fluently and completes its live-coding drills within 25 minutes each.

## Constraints
- Self-paced. No weekly schedule; the order of lessons is the plan.
- Prior experience: Python. Go is new. Docker basics are known.
- Learning style: short, dense lessons with a primary source each, test-driven exercises (`go test -race`), and a retrieval quiz per lesson.
- Local tooling only: Go 1.26+, Docker for the capstone's later stages. No cloud spend.

## Out of scope (owned by sibling tracks)
This track is **Go-specific**. Other skills the benchmark role asks for are left to separate tracks, created when reached, rather than taught shallowly here:

1. **`learn-system-design`**: distributed systems, message-broker internals (Kafka/NATS/Pub/Sub), consistency, databases, scaling, system-design interview practice.
2. **`learn-trading-systems`**: order books and matching in depth, OMS/RMS, the FIX protocol, market data, financial ledgers, exchange connectivity.
3. **`learn-cpp`**: C/C++ fundamentals for low-latency systems.

Docker and Kubernetes are owned by [`learn-kubernetes`](https://github.com/farhaduneci/learn-kubernetes), and GCP/cloud fundamentals by a cloud track. Lesson 20 and the capstone touch brokers and trading only as far as Go code needs them.

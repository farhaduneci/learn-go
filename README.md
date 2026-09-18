# Learn Go

Go learning track: from Python to senior-level Go for low-latency, high-throughput backend roles, with test-driven exercises and a matching-engine capstone.

## Structure

- `MISSION.md`: the "why", success criteria, constraints, and what sibling tracks own
- `RESOURCES.md`: curated primary sources and communities
- `NOTES.md`: environment and build notes
- `content/`: lesson and reference sources (Markdown + front matter), the single source of truth
- `lessons/`, `reference/`: generated HTML pages (`python3 scripts/build.py`, requires pandoc)
- `exercises/`: one test-driven exercise per lesson; reference solutions behind the `solution` build tag
- `capstone/`: a limit order book and single-writer matching engine, built in six stages
- `learning-records/`: session and progress logs
- `assets/`: shared CSS/JS for the lesson site (table of contents, quiz engine, theme)

## Curriculum

1. **Go core**: toolchain, types and pointers, slices/maps/strings, interfaces, errors, generics, packages
2. **Concurrency**: goroutines and the scheduler, channels and select, sync and the memory model, context, patterns, concurrency bugs
3. **Senior Go**: testing, refactoring, design patterns, services (HTTP/gRPC), performance, observability, messaging
4. **Interview**: rapid-fire questions and timed live-coding drills

## Usage

```sh
npm run dev                                   # serve the lesson site at localhost:8000
go test -race ./exercises/08-goroutines/      # work on an exercise (fails until implemented)
go test -race -tags solution ./...            # run every reference solution
./scripts/check.sh                            # gofmt + vet + all solutions under -race
python3 scripts/build.py                      # regenerate HTML after editing content/
```

Requires Go 1.26+.

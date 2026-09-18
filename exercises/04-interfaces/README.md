# Exercise 04 — Methods, interfaces & embedding

Lesson: [04 · Methods, interfaces & composition](../../lessons/0004-interfaces.html)

```sh
go test -race ./exercises/04-interfaces/
```

A broker charges fees in different ways. Model each way as a type that satisfies one small interface.

1. `FlatFee` charges a fixed amount per trade.
2. `BpsFee` charges basis points (1 bps = 0.01%) of the notional amount, rounded down.
3. `Capped` wraps **any** `FeeModel` and caps its result — composition (a decorator), not inheritance.
4. `TotalFee(models, notional)` sums a slice of fee models.
5. `Describe(v any)` uses a **type switch** to return `"flat"`, `"bps"`, `"capped"`, or `"unknown"`.
6. `AuditedAccount` embeds `*Ledger`, so `Post` is promoted. Override `Post` to also append the amount to `Log` and then call the embedded method.

No type declares that it implements `FeeModel`. Go checks this structurally, at compile time.

Reference solution: `solution.go`.

# Exercise 07 — Packages, visibility & API design

Lesson: [07 · Packages, visibility & project layout](../../lessons/0007-packages.html)

```sh
go test -race ./exercises/07-packages/...
```

Build a tiny **double-entry ledger** as a package with a deliberate public API. Every posting moves an amount from one account to another, so the ledger always sums to zero. Financial ledgers are listed as a plus in trading-platform job postings.

The tests live in package `ledger_test` (an *external* test package). They can only use what `ledger` exports, just like a real caller.

1. `New()` returns a ready `*Ledger`. The account map stays unexported.
2. `Post(from, to, amount)` rejects a non-positive amount (`ErrInvalidAmount`) and `from == to` (`ErrSameAccount`). Both are exported sentinels.
3. `Balance(account)` returns the account's balance. Unknown accounts return 0.
4. `Accounts()` returns account names sorted. It returns a **copy**, so callers can't mutate internals.
5. `Sum()` returns the sum of all balances. It should always be 0.

Reference solution: `ledger/solution.go`.

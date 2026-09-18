# Exercise 02 — Types, structs & pointers

Lesson: [02 · Types, zero values, structs & pointers](../../lessons/0002-types-pointers.html)

```sh
go test -race ./exercises/02-types-pointers/
```

Money in trading systems is never a float. These tasks model it as integer cents.

1. `Cents.String()` formats as `$12.34`, `-$0.05`, `$0.00`. Implementing `String()` makes `Cents` a `fmt.Stringer`, so `fmt.Println` uses it automatically.
2. `Account.Deposit` adds to the balance. It needs a **pointer receiver** — with a value receiver it would change a copy.
3. `Account.Withdraw` refuses (returns `false`) if the balance would go negative.
4. `Swap(a, b *int)` swaps two ints through pointers.
5. `NewCounter()` returns a closure; each call returns the next integer, starting at 1.
6. `NewAccount(id)` returns a `*Account` whose zero-value fields are ready to use.

Reference solution: `solution.go` (`go test -tags solution ./exercises/02-types-pointers/`).

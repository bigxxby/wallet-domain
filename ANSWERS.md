# ANSWERS

## Q1: Validation Order

For a FROZEN wallet with balance=100, calling `Withdraw(-50)`:

**Version A returns `ErrWalletFrozen`** — it checks frozen status first.

**Version B returns `ErrInvalidAmount`** — it checks the amount first.

**Version B is correct.** We should always validate the input before checking wallet state. Reasons:

1. A negative amount is always wrong, no matter what state the wallet is in. The user should get the same error for bad input every time.
2. Input errors are the caller's mistake. We should tell them about that first so they can fix it.
3. Checking state first would leak information about the wallet (like whether it's frozen) even when the request is invalid.

---

## Q2: Required Field — 150 or 50?

If balance=100 and someone tries to withdraw 150:

**Option A — Required = 150 (the requested amount):**
> "insufficient balance: required 150 cents, available 100 cents"

**Option B — Required = 50 (the deficit):**
> "insufficient balance: 50 cents short (available 100 cents)"

**Option A is better.** The user already knows what they tried to withdraw, but seeing it echoed back makes it clear. With both numbers (150 and 100) they can immediately understand the situation and decide what to do. With the deficit (50), they'd have to do the math themselves.

---

## Q3: Why not wrap domain errors?

Domain errors like `ErrInsufficientBalance` shouldn't be wrapped with `fmt.Errorf("failed: %w", err)` in the usecase layer because:

1. The handler layer needs to check these errors with `errors.Is` to return the right HTTP status code. Wrapping can make this harder to reason about.
2. Domain errors already carry their meaning. Wrapping just adds noise like "withdraw failed: insufficient balance" — that's redundant.
3. Instead of wrapping, we should log context separately with `slog` and return the domain error as-is.

```go
// don't do this:
return fmt.Errorf("withdraw failed: %w", err)

// do this instead:
logger.Warn("withdraw rejected", "wallet_id", id, "error", err.Error())
return err
```

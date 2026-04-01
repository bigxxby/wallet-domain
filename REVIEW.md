# REVIEW — Bug Analysis

## Why does `errors.Is` return false?

The problem is that `errors.Is` uses `==` to compare structs by default. It compares **all fields**.

The error returned from `Withdraw` looks like this:
```go
WalletError{Code: "E001", Message: "need 500, have 100"}
```

And the sentinel error looks like this:
```go
WalletError{Code: "E001", Message: "insufficient balance"}
```

The `Code` is the same, but the `Message` is different. So `==` returns `false`.

`WalletError` doesn't have an `Is()` method, so `errors.Is` has no other way to compare them — it just fails.

## How to fix it

Add an `Is` method that only compares the `Code` field:

```go
func (e WalletError) Is(target error) bool {
    t, ok := target.(WalletError)
    if !ok {
        return false
    }
    return e.Code == t.Code
}
```

Now `errors.Is` will call this method instead of using `==`, and it will match on `Code: "E001"` regardless of the `Message`.

package main

import (
	"errors"
	"testing"
)

func TestDeposit_Success(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)

	if err := w.Deposit(50); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Balance() != 150 {
		t.Fatalf("expected balance 150, got %d", w.Balance())
	}
}

func TestDeposit_ZeroAmount(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	err := w.Deposit(0)

	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestDeposit_NegativeAmount(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	err := w.Deposit(-10)

	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestDeposit_FrozenWallet(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	_ = w.Freeze()

	err := w.Deposit(50)
	if !errors.Is(err, ErrWalletFrozen) {
		t.Fatalf("expected ErrWalletFrozen, got %v", err)
	}
}

func TestDeposit_InvalidAmountOnFrozenWallet(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	_ = w.Freeze()

	// invalid input should be caught before checking frozen status
	err := w.Deposit(-50)
	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestWithdraw_Success(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)

	if err := w.Withdraw(40); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Balance() != 60 {
		t.Fatalf("expected balance 60, got %d", w.Balance())
	}
}

func TestWithdraw_ExactBalance(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)

	if err := w.Withdraw(100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Balance() != 0 {
		t.Fatalf("expected balance 0, got %d", w.Balance())
	}
}

func TestWithdraw_ZeroAmount(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	err := w.Withdraw(0)

	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestWithdraw_NegativeAmount(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	err := w.Withdraw(-10)

	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestWithdraw_InsufficientBalance_ErrorsIs(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	err := w.Withdraw(150)

	if !errors.Is(err, ErrInsufficientBalance) {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestWithdraw_InsufficientBalance_ErrorsAs(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	err := w.Withdraw(150)

	var balErr *InsufficientBalanceError
	if !errors.As(err, &balErr) {
		t.Fatalf("expected *InsufficientBalanceError, got %T", err)
	}
	if balErr.Required != 150 {
		t.Fatalf("expected Required=150, got %d", balErr.Required)
	}
	if balErr.Available != 100 {
		t.Fatalf("expected Available=100, got %d", balErr.Available)
	}
}

func TestWithdraw_FrozenWallet(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	_ = w.Freeze()

	err := w.Withdraw(50)
	if !errors.Is(err, ErrWalletFrozen) {
		t.Fatalf("expected ErrWalletFrozen, got %v", err)
	}
}

func TestWithdraw_InvalidAmountOnFrozenWallet(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	_ = w.Freeze()

	// invalid input should be caught before checking frozen status
	err := w.Withdraw(-50)
	if !errors.Is(err, ErrInvalidAmount) {
		t.Fatalf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestFreeze_Success(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)

	if err := w.Freeze(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Status() != StatusFrozen {
		t.Fatalf("expected FROZEN status, got %s", w.Status())
	}
}

func TestFreeze_AlreadyFrozen(t *testing.T) {
	w := NewWallet("w1", "owner1", 100)
	_ = w.Freeze()

	err := w.Freeze()
	if !errors.Is(err, ErrWalletFrozen) {
		t.Fatalf("expected ErrWalletFrozen, got %v", err)
	}
}

func TestNewWallet_Defaults(t *testing.T) {
	w := NewWallet("w1", "owner1", 500)

	if w.ID() != "w1" {
		t.Fatalf("expected ID w1, got %s", w.ID())
	}
	if w.OwnerID() != "owner1" {
		t.Fatalf("expected OwnerID owner1, got %s", w.OwnerID())
	}
	if w.Balance() != 500 {
		t.Fatalf("expected balance 500, got %d", w.Balance())
	}
	if w.Status() != StatusActive {
		t.Fatalf("expected ACTIVE status, got %s", w.Status())
	}
}

func TestInsufficientBalanceError_Message(t *testing.T) {
	err := &InsufficientBalanceError{Required: 150, Available: 100}
	expected := "insufficient balance: required 150 cents, available 100 cents"
	if err.Error() != expected {
		t.Fatalf("expected %q, got %q", expected, err.Error())
	}
}

package main

import (
	"errors"
	"fmt"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrWalletFrozen        = errors.New("wallet is frozen")
)

// InsufficientBalanceError holds details about a failed withdrawal.
type InsufficientBalanceError struct {
	Required  int64
	Available int64
}

func (e *InsufficientBalanceError) Error() string {
	return fmt.Sprintf("insufficient balance: required %d cents, available %d cents", e.Required, e.Available)
}

// Is makes errors.Is(err, ErrInsufficientBalance) work.
func (e *InsufficientBalanceError) Is(target error) bool {
	return target == ErrInsufficientBalance
}

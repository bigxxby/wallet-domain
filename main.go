package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	w := NewWallet("w-001", "owner-42", 1000)

	// deposit 500
	if err := w.Deposit(500); err != nil {
		logger.Warn("deposit rejected",
			"wallet_id", w.ID(),
			"amount_cents", 500,
			"reason", err.Error(),
			"error_type", fmt.Sprintf("%T", err),
		)
	} else {
		logger.Info("deposit completed",
			"wallet_id", w.ID(),
			"amount_cents", 500,
			"balance_cents", w.Balance(),
		)
	}

	// withdraw 300
	if err := w.Withdraw(300); err != nil {
		logger.Warn("withdraw rejected",
			"wallet_id", w.ID(),
			"amount_cents", 300,
			"reason", err.Error(),
			"error_type", fmt.Sprintf("%T", err),
		)
	} else {
		logger.Info("withdraw completed",
			"wallet_id", w.ID(),
			"amount_cents", 300,
			"balance_cents", w.Balance(),
		)
	}

	// try to withdraw more than balance
	if err := w.Withdraw(5000); err != nil {
		logger.Warn("withdraw rejected",
			"wallet_id", w.ID(),
			"amount_cents", 5000,
			"reason", err.Error(),
			"error_type", fmt.Sprintf("%T", err),
		)
	}

	// try invalid amount
	if err := w.Deposit(-10); err != nil {
		logger.Warn("deposit rejected",
			"wallet_id", w.ID(),
			"amount_cents", -10,
			"reason", err.Error(),
			"error_type", fmt.Sprintf("%T", err),
		)
	}

	// freeze the wallet
	if err := w.Freeze(); err != nil {
		logger.Warn("freeze rejected",
			"wallet_id", w.ID(),
			"reason", err.Error(),
			"error_type", fmt.Sprintf("%T", err),
		)
	} else {
		logger.Info("wallet frozen", "wallet_id", w.ID())
	}

	// try to deposit on frozen wallet
	if err := w.Deposit(100); err != nil {
		logger.Warn("deposit rejected",
			"wallet_id", w.ID(),
			"amount_cents", 100,
			"reason", err.Error(),
			"error_type", fmt.Sprintf("%T", err),
		)
	}

	// health check
	http.HandleFunc("/healthz", func(hw http.ResponseWriter, r *http.Request) {
		hw.Header().Set("Content-Type", "application/json")
		hw.WriteHeader(http.StatusOK)
		hw.Write([]byte(`{"status":"ok"}`))
	})

	logger.Info("server starting", "addr", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("server failed", "error", err.Error())
		os.Exit(1)
	}
}

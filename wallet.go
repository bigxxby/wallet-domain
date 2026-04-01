package main

type WalletID string
type OwnerID string
type WalletStatus string

const (
	StatusActive WalletStatus = "ACTIVE"
	StatusFrozen WalletStatus = "FROZEN"
)

type Wallet struct {
	id      WalletID
	ownerID OwnerID
	balance int64
	status  WalletStatus
}

func NewWallet(id WalletID, ownerID OwnerID, initialBalance int64) *Wallet {
	return &Wallet{
		id:      id,
		ownerID: ownerID,
		balance: initialBalance,
		status:  StatusActive,
	}
}

func (w *Wallet) ID() WalletID         { return w.id }
func (w *Wallet) OwnerID() OwnerID     { return w.ownerID }
func (w *Wallet) Balance() int64       { return w.balance }
func (w *Wallet) Status() WalletStatus { return w.status }

// Deposit adds money to the wallet.
// Checks: valid amount -> not frozen.
func (w *Wallet) Deposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if w.status == StatusFrozen {
		return ErrWalletFrozen
	}
	w.balance += amount
	return nil
}

// Withdraw takes money from the wallet.
// Checks: valid amount -> not frozen -> enough balance.
func (w *Wallet) Withdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if w.status == StatusFrozen {
		return ErrWalletFrozen
	}
	if w.balance < amount {
		return &InsufficientBalanceError{
			Required:  amount,
			Available: w.balance,
		}
	}
	w.balance -= amount
	return nil
}

// Freeze freezes the wallet. Returns error if already frozen.
func (w *Wallet) Freeze() error {
	if w.status == StatusFrozen {
		return ErrWalletFrozen
	}
	w.status = StatusFrozen
	return nil
}

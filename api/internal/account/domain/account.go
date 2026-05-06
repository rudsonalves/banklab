package domain

import (
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountActive   AccountStatus = "active"
	AccountInactive AccountStatus = "inactive"
	AccountBlocked  AccountStatus = "blocked"
)

type Account struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Number     string
	Branch     string
	Balance    int64
	Status     AccountStatus
	CreatedAt  time.Time
}

// NewAccount creates a new account for the specified customer with the given
// number and branch. It returns an error if the customer ID is invalid.
func NewAccount(customerID uuid.UUID, number, branch string) (*Account, error) {
	if customerID == uuid.Nil {
		return nil, ErrInvalidData
	}

	return &Account{
		ID:         uuid.New(),
		CustomerID: customerID,
		Number:     number,
		Branch:     branch,
		Balance:    0,
		Status:     AccountActive,
		CreatedAt:  time.Now().UTC(),
	}, nil
}

// CanDeposit checks if a deposit transaction can be performed on the account
// with the specified amount. It returns an error if the amount is invalid or
// if the account is not active.
func (a *Account) CanDeposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if a.Status != AccountActive {
		return ErrAccountInactive
	}

	return nil
}

// CanWithdraw checks if a withdrawal transaction can be performed on the account
// with the specified amount. It returns an error if the amount is invalid, if
// the account is not active, or if the account has insufficient balance.
func (a *Account) CanWithdraw(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if a.Status != AccountActive {
		return ErrAccountInactive
	}

	if a.Balance < amount {
		return ErrInsufficientBalance
	}

	return nil
}

// CanTransfer checks if a transfer transaction can be performed from the account
// to the specified destination account with the given amount. It returns an error
// if the amount is invalid, if the account is not active, if the account has
// insufficient balance, or if the destination account is the same as the source
// account.
func (a *Account) CanTransfer(amount int64, destinationID uuid.UUID) error {
	if a.ID == destinationID {
		return ErrSameAccountTransfer
	}

	return a.CanWithdraw(amount)
}

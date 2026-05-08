package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AccountStatus string

const (
	AccountActive   AccountStatus = "active"
	AccountInactive AccountStatus = "inactive"
	AccountBlocked  AccountStatus = "blocked"
)

var (
	ErrInvalidData         = errors.New("invalid data")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrAccountNotFound     = errors.New("account not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrCustomerNotFound    = errors.New("customer not found")
	ErrAccountInactive     = errors.New("account inactive")
	ErrForbidden           = errors.New("forbidden")
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

type TransferRecipient struct {
	AccountID      uuid.UUID
	HolderName     string
	MaskedDocument string
	Branch         string
	AccountNumber  string
	AccountType    string
}

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

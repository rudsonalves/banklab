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

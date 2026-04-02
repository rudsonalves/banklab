package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionDeposit     TransactionType = "deposit"
	TransactionWithdraw    TransactionType = "withdraw"
	TransactionTransferOut TransactionType = "transfer_out"
	TransactionTransferIn  TransactionType = "transfer_in"
)

type Transaction struct {
	ID           uuid.UUID
	AccountID    uuid.UUID
	Type         TransactionType
	Amount       int64
	BalanceAfter int64
	ReferenceID  *uuid.UUID
	CreatedAt    time.Time
}

func NewTransaction(
	accountID uuid.UUID,
	ttype TransactionType,
	amount int64,
	balanceAfter int64,
	referenceID *uuid.UUID,
) *Transaction {
	return &Transaction{
		ID:           uuid.New(),
		AccountID:    accountID,
		Type:         ttype,
		Amount:       amount,
		BalanceAfter: balanceAfter,
		ReferenceID:  referenceID,
		CreatedAt:    time.Now().UTC(),
	}
}

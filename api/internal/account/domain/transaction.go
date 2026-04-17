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
	ID               uuid.UUID
	AccountID        uuid.UUID
	Type             TransactionType
	Amount           int64
	BalanceAfter     int64
	ReferenceID      *uuid.UUID
	RelatedAccountID *uuid.UUID
	IdempotencyKey   *string
	CreatedAt        time.Time
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

func NewTransactionWithIdempotency(
	accountID uuid.UUID,
	ttype TransactionType,
	amount int64,
	balanceAfter int64,
	referenceID *uuid.UUID,
	relatedAccountID *uuid.UUID,
	idempotencyKey string,
) *Transaction {
	key := idempotencyKey
	return &Transaction{
		ID:               uuid.New(),
		AccountID:        accountID,
		Type:             ttype,
		Amount:           amount,
		BalanceAfter:     balanceAfter,
		ReferenceID:      referenceID,
		RelatedAccountID: relatedAccountID,
		IdempotencyKey:   &key,
		CreatedAt:        time.Now().UTC(),
	}
}

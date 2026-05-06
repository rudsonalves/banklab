package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidData         = errors.New("invalid data")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrAccountNotFound     = errors.New("account not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrSameAccountTransfer = errors.New("same account transfer")
	ErrAccountInactive     = errors.New("account inactive")
	ErrForbidden           = errors.New("forbidden")
	ErrTransferDuplicate   = errors.New("transfer already processed")
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
	Balance    int64
	Status     AccountStatus
}

func (a *Account) CanDeposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if a.Status != AccountActive {
		return ErrAccountInactive
	}

	return nil
}

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

func (a *Account) CanTransfer(amount int64, destinationID uuid.UUID) error {
	if a.ID == destinationID {
		return ErrSameAccountTransfer
	}

	return a.CanWithdraw(amount)
}

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

type Repository interface {
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*Account, error)
	IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	CreateTransaction(ctx context.Context, tx *Transaction) error
	GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*Transaction, error)
	GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName TransactionType) (*Transaction, error)
	WithTransaction(ctx context.Context, fn func(tx Tx) error) error
}

type Tx interface {
	Repository
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

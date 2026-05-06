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
	ErrTransactionNotFound = errors.New("transaction not found")
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

// CanDeposit checks if a deposit can be made to the account with the specified
// amount. It validates that the amount is positive and the account is active.
func (a *Account) CanDeposit(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if a.Status != AccountActive {
		return ErrAccountInactive
	}

	return nil
}

// CanWithdraw checks if a withdrawal can be made from the account with the specified
// amount. It validates that the amount is positive, the account is active, and has
// sufficient balance for the withdrawal.
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

// CanTransfer checks if a transfer can be made from the current account to the
// destination account with the specified amount. It validates that the source and
// destination accounts are different, the amount is positive, the source account is
// active, and has sufficient balance for the transfer.
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

// NewTransaction creates a new Transaction with the given parameters.
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

// NewTransactionWithIdempotency creates a new Transaction with the given parameters
// and an idempotency key. The idempotency key is used to ensure that if the same
// transfer is attempted multiple times (e.g., due to retries), only one transaction
// will be created.
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

type TransferReceipt struct {
	OperationType            TransactionType
	Amount                   int64
	Status                   string
	TransactionReference     uuid.UUID
	OperationDate            time.Time
	SourceAccountID          uuid.UUID
	SourceCustomerID         uuid.UUID
	SourceBranch             string
	SourceAccountNumber      string
	DestinationAccountID     uuid.UUID
	DestinationCustomerID    uuid.UUID
	DestinationBranch        string
	DestinationAccountNumber string
	RecipientName            string
	Description              *string
}

type Repository interface {
	GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*Account, error)
	GetByBranchAndNumber(ctx context.Context, branch, number string) (*Account, error)
	GetTransferReceiptByReference(ctx context.Context, referenceID uuid.UUID) (*TransferReceipt, error)
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

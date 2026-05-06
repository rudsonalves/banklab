package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type txRepository struct {
	tx   pgx.Tx
	base baseRepository
}

var _ domain.Tx = (*txRepository)(nil)

// NextAccountNumber generates the next account number using the base repository's
// NextAccountNumber method. It returns the generated account number or an error
// if the generation fails.
func (r *txRepository) NextAccountNumber(ctx context.Context) (string, error) {
	return r.base.NextAccountNumber(ctx)
}

// ListByCustomerID retrieves a list of accounts associated with the specified
// customer ID. It takes a context and a customer ID as input and returns a
// slice of account domain objects or an error if the operation fails.
func (r *txRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return r.base.ListByCustomerID(ctx, customerID)
}

// Create inserts a new account into the database using the base repository's Create
// method. It takes a context and an account domain object as input and returns an
// error if the creation fails.
func (r *txRepository) Create(ctx context.Context, acc *domain.Account) error {
	return r.base.Create(ctx, acc)
}

// CreateTransaction inserts a new transaction into the database using the base
// repository's CreateTransaction method. It takes a context and a transaction domain
// object as input and returns an error if the creation fails.
func (r *txRepository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

// GetTransactionByIdempotencyKey retrieves a transaction based on the provided
// account ID and idempotency key. It takes a context, account ID, and idempotency
// key as input and returns the corresponding transaction domain object or an error
// if the operation fails.
func (r *txRepository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

// GetTransactionByReference retrieves a transaction based on the provided account ID,
// reference ID, and transaction type. It takes a context, account ID, reference ID,
// and transaction type as input and returns the corresponding transaction domain
// object or an error if the operation fails.
func (r *txRepository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

// ExistsByCustomerID checks if any accounts exist for the specified customer ID.
// It takes a context and a customer ID as input and returns a boolean indicating
// whether any accounts exist for that customer, along with an error if the
// operation fails.
func (r *txRepository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return r.base.ExistsByCustomerID(ctx, customerID)
}

// GetByID retrieves an account based on the provided account ID. It takes a context
// and an account ID as input and returns the corresponding account domain object or
// an error if the operation fails. If the account is not found, it returns a
// specific error indicating that the account was not found.
func (r *txRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	account, err := r.base.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}

// GetByIDForUpdate retrieves an account based on the provided account ID and locks
// it for update. It takes a context and an account ID as input and returns the
// corresponding account domain object or an error if the operation fails. If the
// account is not found, it returns a specific error indicating that the account
// was not found.
func (r *txRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

// GetTransactions retrieves a list of transactions for the specified account ID, applying
// pagination and date filtering based on the input parameters. It takes a context,
// account ID, limit, cursor time, cursor ID, from date, and to date as input and
// returns a slice of transaction domain objects or an error if the operation fails. If the account is not found, it returns a specific error indicating that the
// account was not found.
func (r *txRepository) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]domain.Transaction, error) {
	return r.base.GetTransactions(ctx, accountID, limit, cursorTime, cursorID, from, to)
}

// IncreaseBalance increases the balance of the specified account by the given amount.
// It takes a context, account ID, and amount as input and returns the updated balance
// or an error if the operation fails. If the account is not found, it returns a specific error indicating that the
// account was not found.
func (r *txRepository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

// DecreaseBalance decreases the balance of the specified account by the given amount.
// It takes a context, account ID, and amount as input and returns the updated balance
// or an error if the operation fails. If the account is not found, it returns a specific error indicating that the
// account was not found.
func (r *txRepository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

// BeginTx starts a new transaction. Since this repository is already operating
// within a transaction context, it returns an error indicating that nested
// transactions are not supported.
func (r *txRepository) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, fmt.Errorf("nested transactions are not supported")
}

// WithTransaction executes the provided function within a transaction context.
// Since this repository is already operating within a transaction, it returns an
// error indicating that nested transactions are not supported.
func (r *txRepository) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	return fmt.Errorf("nested transactions are not supported")
}

// Commit commits the current transaction. It calls the underlying transaction's
// Commit method and returns an error if the commit operation fails.
func (r *txRepository) Commit(ctx context.Context) error {
	if err := r.tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// Rollback rolls back the current transaction. It calls the underlying transaction's
// Rollback method and returns an error if the rollback operation fails, unless the
// transaction is already closed, in which case it ignores the error.
func (r *txRepository) Rollback(ctx context.Context) error {
	if err := r.tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return fmt.Errorf("rollback transaction: %w", err)
	}
	return nil
}

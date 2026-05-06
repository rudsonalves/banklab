package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type Repository struct {
	db   *pgxpool.Pool
	base baseRepository
}

var _ domain.AccountRepository = (*Repository)(nil)

// New creates a new instance of the Repository with the provided database
// connection pool. It initializes the base repository with the same database
// connection for executing queries.
func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:   db,
		base: baseRepository{exec: db},
	}
}

// NextAccountNumber generates the next account number using the base repository's
// NextAccountNumber method. It returns the generated account number or an error
// if the generation fails.
func (r *Repository) NextAccountNumber(ctx context.Context) (string, error) {
	return r.base.NextAccountNumber(ctx)
}

// Create inserts a new account into the database using the base repository's Create
// method. It takes a context and an account domain object as input and returns an
// error if the creation fails.
func (r *Repository) Create(ctx context.Context, acc *domain.Account) error {
	return r.base.Create(ctx, acc)
}

// ListByCustomerID retrieves a list of accounts associated with the specified
// customer ID. It takes a context and a customer ID as input and returns a
// slice of account domain objects or an error if the operation fails.
func (r *Repository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return r.base.ListByCustomerID(ctx, customerID)
}

// CreateTransaction inserts a new transaction into the database using the base
// repository's CreateTransaction method. It takes a context and a transaction domain
// object as input and returns an error if the creation fails.
func (r *Repository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

// GetTransactionByIdempotencyKey retrieves a transaction based on the provided
// account ID and idempotency key. It takes a context, account ID, and idempotency
// key as input and returns a transaction domain object or an error if the operation
// fails.
func (r *Repository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

// GetTransactionByReference retrieves a transaction based on the provided account ID,
// reference ID, and transaction type. It takes a context, account ID, reference ID,
// and transaction type as input and returns a transaction domain object or an error
// if the operation fails.
func (r *Repository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

// ExistsByCustomerID checks if any accounts exist for the specified customer ID. It
// takes a context and a customer ID as input and returns a boolean indicating whether
// accounts exist for the customer or an error if the operation fails.
func (r *Repository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return r.base.ExistsByCustomerID(ctx, customerID)
}

// GetByID retrieves an account based on the provided account ID. It takes a context
// and an account ID as input and returns an account domain object or an error if the
// operation fails. If the account is not found, it returns a specific error indicating
// that the account was not found.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	account, err := r.base.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, domain.ErrAccountNotFound
	}

	return account, nil
}

// GetByIDForUpdate retrieves an account for update based on the provided account
// ID. It takes a context and an account ID as input and returns an account
// domain object or an error if the operation fails. If the account is not found,
// it returns a specific error indicating that the account was not found. This
// method is typically used within a transaction to lock the account record
// for update, ensuring that concurrent modifications are handled correctly.
func (r *Repository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

// GetTransactions retrieves a list of transactions for the specified account ID, applying
// pagination and date filtering based on the input parameters. It takes a context,
// account ID, limit, cursor time, cursor ID, from time, and to time as input and
// returns a slice of transaction domain objects or an error if the operation fails. If the account is not found, it returns a specific error indicating that the
// account was not found.
func (r *Repository) GetTransactions(
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

// IncreaseBalance increases the balance of the specified account by the
// given amount. It takes a context, account ID, and amount as input and
// returns the updated balance or an error if the operation fails. If the
// account is not found, it returns a specific error indicating that the
// account was not found.
func (r *Repository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

// DecreaseBalance decreases the balance of the specified account by the given
// amount. It takes a context, account ID, and amount as input and returns the
// updated balance or an error if the operation fails.
func (r *Repository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

// BeginTx starts a new database transaction and returns a transaction object
// that can be used to execute queries within the transaction. It takes a context
// as input and returns a transaction domain object or an error if the operation
// fails.
func (r *Repository) BeginTx(ctx context.Context) (domain.Tx, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return &txRepository{
		tx:   tx,
		base: baseRepository{exec: tx},
	}, nil
}

// WithTransaction executes the provided function within a database transaction. It begins a new transaction, executes the function, and commits the transaction if the function
// returns no error. If the function returns an error, it rolls back the transaction.
// It takes a context and a function that accepts a transaction domain object as input
// and returns an error. The method returns an error if starting the transaction,
// executing the function, committing, or rolling back fails.
func (r *Repository) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	tx, err := r.BeginTx(ctx)
	if err != nil {
		return err
	}

	return runInTransaction(ctx, tx, fn)
}

// runInTransaction executes the provided function within the context of a
// transaction. It commits the transaction if the function returns no error,
// or rolls back the transaction if the function returns an error. It takes a
// context, a transaction domain object, and a function that accepts a transaction
// domain object as input and returns an error. The method returns an error if
// committing or rolling back the transaction fails.
func runInTransaction(ctx context.Context, tx domain.Tx, fn func(tx domain.Tx) error) error {
	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("rollback transaction: %w", rollbackErr)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return fmt.Errorf("commit transaction: %w (rollback failed: %v)", err, rollbackErr)
		}
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type executor interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type Repository struct {
	db   *pgxpool.Pool
	base baseRepository
}

var _ transactiondomain.Repository = (*Repository)(nil)

// New creates a Repository backed by the provided connection pool.
func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:   db,
		base: baseRepository{exec: db},
	}
}

// GetByIDForUpdate fetches an account row by ID and locks it for update
// within the current transaction.
func (r *Repository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*transactiondomain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

// GetByBranchAndNumber returns the account identified by the given branch
// and account number.
func (r *Repository) GetByBranchAndNumber(ctx context.Context, branch, number string) (*transactiondomain.Account, error) {
	return r.base.GetByBranchAndNumber(ctx, branch, number)
}

// IncreaseBalance adds amount to the balance of the account identified by id
// and returns the updated balance.
func (r *Repository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

// DecreaseBalance subtracts amount from the balance of the account identified
// by id and returns the updated balance.
func (r *Repository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

// CreateTransaction inserts a ledger transaction record into the database.
func (r *Repository) CreateTransaction(ctx context.Context, tx *transactiondomain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

// GetTransactionByIdempotencyKey looks up the ledger entry for the given
// account that matches the supplied idempotency key. Returns nil when no
// matching entry is found.
func (r *Repository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

// GetTransactionByReference finds the ledger entry for the given account that
// matches the supplied reference ID and transaction type.
func (r *Repository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName transactiondomain.TransactionType) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

// BeginTx starts a new database transaction and returns a Tx that wraps it.
func (r *Repository) BeginTx(ctx context.Context) (transactiondomain.Tx, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return &txRepository{
		tx:   tx,
		base: baseRepository{exec: tx},
	}, nil
}

// WithTransaction begins a transaction, executes fn, and commits on success.
// If fn returns an error the transaction is rolled back.
func (r *Repository) WithTransaction(ctx context.Context, fn func(tx transactiondomain.Tx) error) error {
	tx, err := r.BeginTx(ctx)
	if err != nil {
		return err
	}

	return runInTransaction(ctx, tx, fn)
}

// GetTransferReceiptByReference returns the transfer receipt for the given
// reference ID by delegating to the shared base repository.
func (r *Repository) GetTransferReceiptByReference(ctx context.Context, referenceID uuid.UUID) (*transactiondomain.TransferReceipt, error) {
	return r.base.GetTransferReceiptByReference(ctx, referenceID)
}

// runInTransaction executes fn inside the supplied Tx, committing on success
// and rolling back on any error. A failed rollback wraps both errors.
func runInTransaction(ctx context.Context, tx transactiondomain.Tx, fn func(tx transactiondomain.Tx) error) error {
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

type txRepository struct {
	tx   pgx.Tx
	base baseRepository
}

var _ transactiondomain.Tx = (*txRepository)(nil)

// GetByIDForUpdate fetches an account row by ID and locks it for update
// within the active transaction.
func (r *txRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*transactiondomain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

// GetByBranchAndNumber returns the account identified by the given branch
// and account number, using the active transaction connection.
func (r *txRepository) GetByBranchAndNumber(ctx context.Context, branch, number string) (*transactiondomain.Account, error) {
	return r.base.GetByBranchAndNumber(ctx, branch, number)
}

// IncreaseBalance adds amount to the account balance within the active
// transaction and returns the updated balance.
func (r *txRepository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

// DecreaseBalance subtracts amount from the account balance within the active
// transaction and returns the updated balance.
func (r *txRepository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

// CreateTransaction inserts a ledger transaction record using the active
// database transaction.
func (r *txRepository) CreateTransaction(ctx context.Context, tx *transactiondomain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

// GetTransactionByIdempotencyKey looks up the ledger entry for the given
// account matching the supplied idempotency key within the active transaction.
func (r *txRepository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

// GetTransactionByReference finds the ledger entry for the given account that
// matches the supplied reference ID and transaction type within the active
// transaction.
func (r *txRepository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName transactiondomain.TransactionType) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

// WithTransaction always returns an error because nested transactions are
// not supported.
func (r *txRepository) WithTransaction(ctx context.Context, fn func(tx transactiondomain.Tx) error) error {
	return fmt.Errorf("nested transactions are not supported")
}

// Commit persists all changes made within the current transaction.
func (r *txRepository) Commit(ctx context.Context) error {
	if err := r.tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// Rollback aborts the current transaction. Errors from an already-closed
// transaction are silently ignored.
func (r *txRepository) Rollback(ctx context.Context) error {
	if err := r.tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return fmt.Errorf("rollback transaction: %w", err)
	}
	return nil
}

// GetTransferReceiptByReference returns the transfer receipt for the given
// reference ID using the active transaction connection.
func (r *txRepository) GetTransferReceiptByReference(ctx context.Context, referenceID uuid.UUID) (*transactiondomain.TransferReceipt, error) {
	return r.base.GetTransferReceiptByReference(ctx, referenceID)
}

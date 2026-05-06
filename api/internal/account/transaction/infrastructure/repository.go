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

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:   db,
		base: baseRepository{exec: db},
	}
}

func (r *Repository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*transactiondomain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

func (r *Repository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

func (r *Repository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

func (r *Repository) CreateTransaction(ctx context.Context, tx *transactiondomain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

func (r *Repository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

func (r *Repository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName transactiondomain.TransactionType) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

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

func (r *Repository) WithTransaction(ctx context.Context, fn func(tx transactiondomain.Tx) error) error {
	tx, err := r.BeginTx(ctx)
	if err != nil {
		return err
	}

	return runInTransaction(ctx, tx, fn)
}

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

func (r *txRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*transactiondomain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

func (r *txRepository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

func (r *txRepository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

func (r *txRepository) CreateTransaction(ctx context.Context, tx *transactiondomain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

func (r *txRepository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

func (r *txRepository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName transactiondomain.TransactionType) (*transactiondomain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

func (r *txRepository) WithTransaction(ctx context.Context, fn func(tx transactiondomain.Tx) error) error {
	return fmt.Errorf("nested transactions are not supported")
}

func (r *txRepository) Commit(ctx context.Context) error {
	if err := r.tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (r *txRepository) Rollback(ctx context.Context) error {
	if err := r.tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		return fmt.Errorf("rollback transaction: %w", err)
	}
	return nil
}

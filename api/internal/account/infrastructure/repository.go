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

func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:   db,
		base: baseRepository{exec: db},
	}
}

func (r *Repository) NextAccountNumber(ctx context.Context) (string, error) {
	return r.base.NextAccountNumber(ctx)
}

func (r *Repository) Create(ctx context.Context, acc *domain.Account) error {
	return r.base.Create(ctx, acc)
}

func (r *Repository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return r.base.ListByCustomerID(ctx, customerID)
}

func (r *Repository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

func (r *Repository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

func (r *Repository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

func (r *Repository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return r.base.ExistsByCustomerID(ctx, customerID)
}

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

func (r *Repository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

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

func (r *Repository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

func (r *Repository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

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

func (r *Repository) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	tx, err := r.BeginTx(ctx)
	if err != nil {
		return err
	}

	return runInTransaction(ctx, tx, fn)
}

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

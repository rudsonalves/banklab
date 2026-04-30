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

func (r *txRepository) NextAccountNumber(ctx context.Context) (string, error) {
	return r.base.NextAccountNumber(ctx)
}

func (r *txRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return r.base.ListByCustomerID(ctx, customerID)
}

func (r *txRepository) Create(ctx context.Context, acc *domain.Account) error {
	return r.base.Create(ctx, acc)
}

func (r *txRepository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return r.base.CreateTransaction(ctx, tx)
}

func (r *txRepository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return r.base.GetTransactionByIdempotencyKey(ctx, accountID, key)
}

func (r *txRepository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return r.base.GetTransactionByReference(ctx, accountID, referenceID, typeName)
}

func (r *txRepository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return r.base.ExistsByCustomerID(ctx, customerID)
}

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

func (r *txRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return r.base.GetByIDForUpdate(ctx, id)
}

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

func (r *txRepository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.IncreaseBalance(ctx, id, amount)
}

func (r *txRepository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return r.base.DecreaseBalance(ctx, id, amount)
}

func (r *txRepository) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, fmt.Errorf("nested transactions are not supported")
}

func (r *txRepository) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
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

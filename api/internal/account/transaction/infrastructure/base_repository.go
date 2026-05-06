package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type baseRepository struct {
	exec executor
}

func (r *baseRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*transactiondomain.Account, error) {
	var account transactiondomain.Account

	query := `
		SELECT id, customer_id, balance, status
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	err := r.exec.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Balance,
		&account.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, transactiondomain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account by id for update: %w", err)
	}

	return &account, nil
}

func (r *baseRepository) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance
	`

	err := r.exec.QueryRow(ctx, query, amount, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, transactiondomain.ErrInvalidAmount
		}
		return 0, fmt.Errorf("update balance: %w", err)
	}

	return balance, nil
}

func (r *baseRepository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	query := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		  AND balance >= $1
		RETURNING balance
	`

	err := r.exec.QueryRow(ctx, query, amount, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			exists, existsErr := r.accountExists(ctx, id)
			if existsErr != nil {
				return 0, fmt.Errorf("decrease balance: %w", existsErr)
			}
			if !exists {
				return 0, transactiondomain.ErrAccountNotFound
			}
			return 0, transactiondomain.ErrInsufficientBalance
		}
		return 0, fmt.Errorf("decrease balance: %w", err)
	}

	return balance, nil
}

func (r *baseRepository) accountExists(ctx context.Context, id uuid.UUID) (bool, error) {
	var dummy int
	query := `
		SELECT 1
		FROM accounts
		WHERE id = $1
	`

	err := r.exec.QueryRow(ctx, query, id).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check account exists: %w", err)
	}

	return true, nil
}

func (r *baseRepository) CreateTransaction(ctx context.Context, tx *transactiondomain.Transaction) error {
	query := `
		INSERT INTO transactions (
			id, account_id, type, amount, balance_after, reference_id,
			related_account_id, idempotency_key, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (account_id, idempotency_key)
		WHERE idempotency_key IS NOT NULL
		DO NOTHING
	`

	cmd, err := r.exec.Exec(ctx, query,
		tx.ID,
		tx.AccountID,
		tx.Type,
		tx.Amount,
		tx.BalanceAfter,
		tx.ReferenceID,
		tx.RelatedAccountID,
		tx.IdempotencyKey,
		tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create account transaction: %w", err)
	}

	if tx.IdempotencyKey != nil && cmd.RowsAffected() == 0 {
		return transactiondomain.ErrTransferDuplicate
	}

	return nil
}

func (r *baseRepository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*transactiondomain.Transaction, error) {
	var t transactiondomain.Transaction

	query := `
		SELECT id, account_id, type, amount, balance_after, reference_id,
		       related_account_id, idempotency_key, created_at
		FROM transactions
		WHERE account_id = $1 AND idempotency_key = $2
		LIMIT 1
	`

	err := r.exec.QueryRow(ctx, query, accountID, key).Scan(
		&t.ID,
		&t.AccountID,
		&t.Type,
		&t.Amount,
		&t.BalanceAfter,
		&t.ReferenceID,
		&t.RelatedAccountID,
		&t.IdempotencyKey,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get transaction by idempotency key: %w", err)
	}

	return &t, nil
}

func (r *baseRepository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName transactiondomain.TransactionType) (*transactiondomain.Transaction, error) {
	var t transactiondomain.Transaction

	query := `
		SELECT id, account_id, type, amount, balance_after, reference_id,
		       related_account_id, idempotency_key, created_at
		FROM transactions
		WHERE account_id = $1 AND reference_id = $2 AND type = $3
		LIMIT 1
	`

	err := r.exec.QueryRow(ctx, query, accountID, referenceID, typeName).Scan(
		&t.ID,
		&t.AccountID,
		&t.Type,
		&t.Amount,
		&t.BalanceAfter,
		&t.ReferenceID,
		&t.RelatedAccountID,
		&t.IdempotencyKey,
		&t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get transaction by reference: %w", err)
	}

	return &t, nil
}

package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type executor interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type baseRepository struct {
	exec executor
}

func (r *baseRepository) NextAccountNumber(ctx context.Context) (string, error) {
	var number int64

	err := r.exec.QueryRow(ctx, `
		SELECT nextval('account_number_seq')
	`).Scan(&number)
	if err != nil {
		return "", fmt.Errorf("next account number: %w", err)
	}

	return fmt.Sprintf("%08d", number), nil
}

func (r *baseRepository) Create(ctx context.Context, acc *domain.Account) error {
	query := `
		INSERT INTO accounts (
			id, customer_id, number, branch, balance, status, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.exec.Exec(ctx, query,
		acc.ID,
		acc.CustomerID,
		acc.Number,
		acc.Branch,
		acc.Balance,
		acc.Status,
		acc.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create account: %w", err)
	}

	return nil
}

func (r *baseRepository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	query := `
		SELECT id, customer_id, number, branch, balance, status, created_at
		FROM accounts
		WHERE customer_id = $1
		ORDER BY created_at ASC, id ASC
	`

	rows, err := r.exec.Query(ctx, query, customerID)
	if err != nil {
		return nil, fmt.Errorf("list accounts by customer id: %w", err)
	}
	defer rows.Close()

	accounts := make([]domain.Account, 0)
	for rows.Next() {
		var account domain.Account
		if err := rows.Scan(
			&account.ID,
			&account.CustomerID,
			&account.Number,
			&account.Branch,
			&account.Balance,
			&account.Status,
			&account.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate accounts by customer id: %w", err)
	}

	return accounts, nil
}

func (r *baseRepository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
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
		return domain.ErrTransferDuplicate
	}

	return nil
}

func (r *baseRepository) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	var t domain.Transaction

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

func (r *baseRepository) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	var t domain.Transaction

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

func (r *baseRepository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	query := `
		SELECT 1
		FROM accounts
		WHERE customer_id = $1
	`

	var dummy int
	err := r.exec.QueryRow(ctx, query, customerID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists by customer id: %w", err)
	}

	return true, nil
}

func (r *baseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account

	query := `
		SELECT id, customer_id, number, branch, balance, status, created_at
		FROM accounts
		WHERE id = $1
	`

	err := r.exec.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Number,
		&account.Branch,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	return &account, nil
}

func (r *baseRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account

	query := `
		SELECT id, customer_id, number, branch, balance, status, created_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	err := r.exec.QueryRow(ctx, query, id).Scan(
		&account.ID,
		&account.CustomerID,
		&account.Number,
		&account.Branch,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account by id for update: %w", err)
	}

	return &account, nil
}

func (r *baseRepository) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]domain.Transaction, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	if cursorTime == nil || cursorID == nil {
		cursorTime = nil
		cursorID = nil
	}

	query := `
		SELECT id, account_id, type, amount, balance_after, reference_id,
		       related_account_id, idempotency_key, created_at
		FROM transactions
		WHERE account_id = $1
		  AND ($2::timestamptz IS NULL OR created_at >= $2)
		  AND ($3::timestamptz IS NULL OR created_at <= $3)
		  AND (
			$4::timestamptz IS NULL OR
			(created_at, id) < ($4, $5)
		  )
		ORDER BY created_at DESC, id DESC
		LIMIT $6
	`

	rows, err := r.exec.Query(ctx, query, accountID, from, to, cursorTime, cursorID, limit)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}
	defer rows.Close()

	transactions := make([]domain.Transaction, 0, limit)
	for rows.Next() {
		var transaction domain.Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.BalanceAfter,
			&transaction.ReferenceID,
			&transaction.RelatedAccountID,
			&transaction.IdempotencyKey,
			&transaction.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get transactions: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}

	return transactions, nil
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
			return 0, domain.ErrInvalidAmount
		}
		return 0, fmt.Errorf("update balance: %w", err)
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
				return 0, domain.ErrAccountNotFound
			}
			return 0, domain.ErrInsufficientBalance
		}
		return 0, fmt.Errorf("decrease balance: %w", err)
	}

	return balance, nil
}

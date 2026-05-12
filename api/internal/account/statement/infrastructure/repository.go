package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	bankaccountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	statementdomain "github.com/seu-usuario/bank-api/internal/account/statement/domain"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type executor interface {
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
}

type Repository struct {
	exec executor
}

var _ statementdomain.Repository = (*Repository)(nil)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{exec: db}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*bankaccountdomain.Account, error) {
	var account bankaccountdomain.Account

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
			return nil, bankaccountdomain.ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	return &account, nil
}

func (r *Repository) GetTransactions(
	ctx context.Context,
	accountID uuid.UUID,
	limit int,
	cursorTime *time.Time,
	cursorID *uuid.UUID,
	from *time.Time,
	to *time.Time,
) ([]transactiondomain.Transaction, error) {
	if limit <= 0 || limit > 1000 {
		limit = 100
	}

	if cursorTime == nil || cursorID == nil {
		cursorTime = nil
		cursorID = nil
	}

	query := `
		SELECT id, account_id, type, amount, balance_after, reference_id,
		       related_account_id, idempotency_key, description, created_at
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

	transactions := make([]transactiondomain.Transaction, 0, limit)
	for rows.Next() {
		var transaction transactiondomain.Transaction
		if err := rows.Scan(
			&transaction.ID,
			&transaction.AccountID,
			&transaction.Type,
			&transaction.Amount,
			&transaction.BalanceAfter,
			&transaction.ReferenceID,
			&transaction.RelatedAccountID,
			&transaction.IdempotencyKey,
			&transaction.Description,
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

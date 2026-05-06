package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	bankaccountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
)

type executor interface {
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type Repository struct {
	exec executor
}

var _ bankaccountdomain.Repository = (*Repository)(nil)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{exec: db}
}

func (r *Repository) NextAccountNumber(ctx context.Context) (string, error) {
	var number int64

	err := r.exec.QueryRow(ctx, `
		SELECT nextval('account_number_seq')
	`).Scan(&number)
	if err != nil {
		return "", fmt.Errorf("next account number: %w", err)
	}

	return fmt.Sprintf("%08d", number), nil
}

func (r *Repository) Create(ctx context.Context, acc *bankaccountdomain.Account) error {
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

func (r *Repository) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]bankaccountdomain.Account, error) {
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

	accounts := make([]bankaccountdomain.Account, 0)
	for rows.Next() {
		var account bankaccountdomain.Account
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

func (r *Repository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
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

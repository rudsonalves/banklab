package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

type txRepository struct {
	tx pgx.Tx
}

var _ domain.AccountRepository = (*Repository)(nil)
var _ domain.Tx = (*txRepository)(nil)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) NextAccountNumber(ctx context.Context) (string, error) {
	var number int64

	err := r.db.QueryRow(ctx, `
		SELECT nextval('account_number_seq')
	`).Scan(&number)
	if err != nil {
		return "", fmt.Errorf("next account number: %w", err)
	}

	return fmt.Sprintf("%08d", number), nil
}

func (r *Repository) Create(ctx context.Context, acc *domain.Account) error {
	query := `
		INSERT INTO accounts (
			id, customer_id, number, branch, balance, status, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
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

func (r *Repository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO account_transactions (
			id, account_id, type, amount, balance_after, reference_id, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(ctx, query,
		tx.ID,
		tx.AccountID,
		tx.Type,
		tx.Amount,
		tx.BalanceAfter,
		tx.ReferenceID,
		tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create account transaction: %w", err)
	}

	return nil
}

func (r *Repository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	query := `
		SELECT 1
		FROM accounts
		WHERE customer_id = $1
	`

	var dummy int
	err := r.db.QueryRow(ctx, query, customerID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists by customer id: %w", err)
	}

	return true, nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account

	query := `
		SELECT id, customer_id, number, branch, balance, status, created_at
		FROM accounts
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
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

func (r *Repository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account

	query := `
		SELECT id, customer_id, number, branch, balance, status, created_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
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

func (r *Repository) GetTransactions(
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
		SELECT id, account_id, type, amount, balance_after, reference_id, created_at
		FROM account_transactions
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

	rows, err := r.db.Query(ctx, query, accountID, from, to, cursorTime, cursorID, limit)
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

func (r *Repository) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance
	`

	err := r.db.QueryRow(ctx, query, amount, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrAccountNotFound
		}
		return 0, fmt.Errorf("update balance: %w", err)
	}

	return balance, nil
}

func (r *Repository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	query := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		  AND balance >= $1
	`

	cmdTag, err := r.db.Exec(ctx, query, amount, id)
	if err != nil {
		return fmt.Errorf("decrease balance: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrInsufficientBalance
	}

	return nil
}

func (r *Repository) BeginTx(ctx context.Context) (domain.Tx, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return &txRepository{tx: tx}, nil
}

func (r *txRepository) NextAccountNumber(ctx context.Context) (string, error) {
	var number int64

	err := r.tx.QueryRow(ctx, `
		SELECT nextval('account_number_seq')
	`).Scan(&number)
	if err != nil {
		return "", fmt.Errorf("next account number: %w", err)
	}

	return fmt.Sprintf("%08d", number), nil
}

func (r *txRepository) Create(ctx context.Context, acc *domain.Account) error {
	query := `
		INSERT INTO accounts (
			id, customer_id, number, branch, balance, status, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.tx.Exec(ctx, query,
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

func (r *txRepository) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	query := `
		INSERT INTO account_transactions (
			id, account_id, type, amount, balance_after, reference_id, created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.tx.Exec(ctx, query,
		tx.ID,
		tx.AccountID,
		tx.Type,
		tx.Amount,
		tx.BalanceAfter,
		tx.ReferenceID,
		tx.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create account transaction: %w", err)
	}

	return nil
}

func (r *txRepository) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	query := `
		SELECT 1
		FROM accounts
		WHERE customer_id = $1
	`

	var dummy int
	err := r.tx.QueryRow(ctx, query, customerID).Scan(&dummy)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("exists by customer id: %w", err)
	}

	return true, nil
}

func (r *txRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account

	query := `
		SELECT id, customer_id, number, branch, balance, status, created_at
		FROM accounts
		WHERE id = $1
	`

	err := r.tx.QueryRow(ctx, query, id).Scan(
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

func (r *txRepository) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	var account domain.Account

	query := `
		SELECT id, customer_id, number, branch, balance, status, created_at
		FROM accounts
		WHERE id = $1
		FOR UPDATE
	`

	err := r.tx.QueryRow(ctx, query, id).Scan(
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

func (r *txRepository) GetTransactions(
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
		SELECT id, account_id, type, amount, balance_after, reference_id, created_at
		FROM account_transactions
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

	rows, err := r.tx.Query(ctx, query, accountID, from, to, cursorTime, cursorID, limit)
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

func (r *txRepository) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	query := `
		UPDATE accounts
		SET balance = balance + $1
		WHERE id = $2
		RETURNING balance
	`

	err := r.tx.QueryRow(ctx, query, amount, id).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, domain.ErrAccountNotFound
		}
		return 0, fmt.Errorf("update balance: %w", err)
	}

	return balance, nil
}

func (r *txRepository) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) error {
	query := `
		UPDATE accounts
		SET balance = balance - $1
		WHERE id = $2
		  AND balance >= $1
	`

	cmdTag, err := r.tx.Exec(ctx, query, amount, id)
	if err != nil {
		return fmt.Errorf("decrease balance: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return domain.ErrInsufficientBalance
	}

	return nil
}

func (r *txRepository) BeginTx(ctx context.Context) (domain.Tx, error) {
	return nil, fmt.Errorf("nested transactions are not supported")
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

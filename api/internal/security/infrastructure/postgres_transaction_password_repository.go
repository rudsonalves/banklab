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
	"github.com/seu-usuario/bank-api/internal/database"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type dbExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type PostgresTransactionPasswordRepository struct {
	db *pgxpool.Pool
}

var _ domain.TransactionPasswordRepository = (*PostgresTransactionPasswordRepository)(nil)

func NewPostgresTransactionPasswordRepository(db *pgxpool.Pool) *PostgresTransactionPasswordRepository {
	return &PostgresTransactionPasswordRepository{db: db}
}

// executor returns the current transaction if it exists in the context,
// otherwise it returns the main database connection pool.
func (r *PostgresTransactionPasswordRepository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

// Create inserts a new transaction password record into the database.
// It returns an error if the user already has a transaction password set
// or if the input data is invalid.
func (r *PostgresTransactionPasswordRepository) Create(
	ctx context.Context,
	password *domain.TransactionPassword,
) error {
	query := `
		INSERT INTO transaction_passwords (
			user_id,
			password_hash,
			status,
			failed_attempts,
			locked_until,
			created_at,
			updated_at,
			changed_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.executor(ctx).QueryRow(
		ctx,
		query,
		password.UserID,
		password.PasswordHash,
		string(password.Status),
		password.FailedAttempts,
		password.LockedUntil,
		password.CreatedAt,
		password.UpdatedAt,
		password.ChangedAt,
	).Scan(&password.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return domain.ErrTransactionPasswordAlreadySet
			case "23514", "23503":
				return domain.ErrInvalidTransactionPassword
			}
		}

		return fmt.Errorf("create transaction password: %w", err)
	}

	return nil
}

// FindByUserID retrieves a transaction password record by the associated user ID.
// It returns nil if no record is found, or an error if the query fails.
func (r *PostgresTransactionPasswordRepository) FindByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*domain.TransactionPassword, error) {
	query := `
		SELECT
			id,
			user_id,
			password_hash,
			status,
			failed_attempts,
			locked_until,
			created_at,
			updated_at,
			changed_at
		FROM transaction_passwords
		WHERE user_id = $1
	`

	password, err := scanTransactionPassword(r.executor(ctx).QueryRow(ctx, query, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("find transaction password by user id: %w", err)
	}

	return password, nil
}

// SaveValidationState updates the validation state of a transaction
// password, including its status, failed attempts, and lock duration.
func (r *PostgresTransactionPasswordRepository) SaveValidationState(
	ctx context.Context,
	password *domain.TransactionPassword,
) error {
	query := `
		UPDATE transaction_passwords
		SET
			status = $1,
			failed_attempts = $2,
			locked_until = $3,
			updated_at = $4
		WHERE id = $5
	`

	result, err := r.executor(ctx).Exec(
		ctx,
		query,
		string(password.Status),
		password.FailedAttempts,
		password.LockedUntil,
		password.UpdatedAt,
		password.ID,
	)
	if err != nil {
		return fmt.Errorf("save transaction password validation state: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTransactionPasswordNotSet
	}

	return nil
}

// UpdatePasswordHash updates the password hash and resets the validation
// state of a transaction password.
// It returns an error if the transaction password is not found or if the
// update fails.
func (r *PostgresTransactionPasswordRepository) UpdatePasswordHash(
	ctx context.Context,
	id uuid.UUID,
	passwordHash string,
	changedAt,
	updatedAt time.Time,
) error {
	query := `
		UPDATE transaction_passwords
		SET
			password_hash = $1,
			status = $2,
			failed_attempts = 0,
			locked_until = NULL,
			changed_at = $3,
			updated_at = $4
		WHERE id = $5
	`

	result, err := r.executor(ctx).Exec(
		ctx,
		query,
		passwordHash,
		string(domain.TransactionPasswordActive),
		changedAt,
		updatedAt,
		id,
	)
	if err != nil {
		return fmt.Errorf("update transaction password hash: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrTransactionPasswordNotSet
	}

	return nil
}

// scanTransactionPassword is a helper function that scans a database row into
// a TransactionPassword struct.
func scanTransactionPassword(row pgx.Row) (*domain.TransactionPassword, error) {
	var password domain.TransactionPassword
	var status string

	err := row.Scan(
		&password.ID,
		&password.UserID,
		&password.PasswordHash,
		&status,
		&password.FailedAttempts,
		&password.LockedUntil,
		&password.CreatedAt,
		&password.UpdatedAt,
		&password.ChangedAt,
	)
	if err != nil {
		return nil, err
	}

	password.Status = domain.TransactionPasswordStatus(status)

	return &password, nil
}

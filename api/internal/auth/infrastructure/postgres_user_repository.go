package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/database"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

type dbExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var _ domain.UserRepository = (*PostgresUserRepository)(nil)

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (
			id,
			email,
			password_hash,
			role,
			customer_id,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.executor(ctx).Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		string(user.Role),
		nullableUUIDValue(user.CustomerID),
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			role,
			customer_id,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	row := r.executor(ctx).QueryRow(ctx, query, email)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			role,
			customer_id,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	row := r.executor(ctx).QueryRow(ctx, query, id)
	user, err := scanUser(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *PostgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT 1
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var exists int
	err := r.executor(ctx).QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *PostgresUserRepository) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin user transaction: %w", err)
	}

	txCtx := database.ContextWithTx(ctx, tx)
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return fmt.Errorf("rollback user transaction after callback error: %v (original: %w)", rollbackErr, err)
		}
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return fmt.Errorf("rollback user transaction after commit error: %v (commit: %w)", rollbackErr, err)
		}
		return fmt.Errorf("commit user transaction: %w", err)
	}

	return nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (*domain.User, error) {
	var user domain.User
	var role string
	var customerID *uuid.UUID

	err := s.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&role,
		&customerID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.Role = domain.Role(role)
	user.CustomerID = customerID

	return &user, nil
}

func nullableUUIDValue(value *uuid.UUID) any {
	if value == nil {
		return nil
	}

	return *value
}

func nullableStringValue(value *string) any {
	if value == nil {
		return nil
	}

	return *value
}

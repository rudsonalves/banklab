package infrastructure

import (
	"context"
	"errors"

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
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.executor(ctx).Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		string(user.Role),
		nullableUUIDValue(user.CustomerID),
		string(user.Status),
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) UpdateStatus(ctx context.Context, userID uuid.UUID, status domain.UserStatus) error {
	query := `
		UPDATE users
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.executor(ctx).Exec(ctx, query, string(status), userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
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
			status,
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
			status,
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

type scanner interface {
	Scan(dest ...any) error
}

func scanUser(s scanner) (*domain.User, error) {
	var user domain.User
	var role string
	var customerID *uuid.UUID
	var status string

	err := s.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&role,
		&customerID,
		&status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	user.Role = domain.Role(role)
	user.CustomerID = customerID
	user.Status = domain.UserStatus(status)

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

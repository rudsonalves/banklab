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

// NewPostgresUserRepository creates a new instance of PostgresUserRepository
// with the provided pgxpool.Pool. This repository will be used to manage user
// data in a PostgreSQL database. It provides methods for creating users, updating
// user status, finding users by email or ID, and checking if a user exists by
// email. The repository uses the pgx library to interact with the database and
// supports transactions through the context.
func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

// executor returns the appropriate database executor based on the context. If a
// transaction is present in the context, it returns the transaction; otherwise,
// it returns the main database connection pool. This allows the repository methods
// to work seamlessly with both transactional and non-transactional contexts.
func (r *PostgresUserRepository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

// Create inserts a new user record into the users table with the provided user data.
// It generates a new UUID for the user and executes the SQL query to store the user
// information in the database. If the operation is successful, it returns nil;
// otherwise, it returns an error.
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

// UpdateStatus updates the status of a user identified by the given userID. It
// executes an SQL query to update the user's status in the database. If the user
// is not found, it returns a domain.ErrUserNotFound error. If the operation is
// successful, it returns nil; otherwise, it returns an error.
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

// FindByEmail retrieves a user from the users table based on the provided email. It
// returns a pointer to the domain.User and an error if the operation fails. If no
// user is found with the given email, it returns nil for the user and nil for the
// error.
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

// FindByID retrieves a user from the users table based on the provided user ID. It
// returns a pointer to the domain.User and an error if the operation fails. If no
// user is found with the given ID, it returns nil for the user and nil for the
// error.
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

// FindByIDForUpdate retrieves a user from the users table based on the provided
// user ID and locks the row for update. It returns a pointer to the domain.User
// and an error if the operation fails. If no user is found with the given ID, it
// returns nil for the user and nil for the error.
func (r *PostgresUserRepository) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error) {
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
		FOR UPDATE
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

// ExistsByEmail checks if a user exists in the users table based on the provided email. It
// returns true if a user with the given email exists, false otherwise. If the operation
// fails, it returns an error.
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

// scanUser is a helper function that scans a database row into a domain.User
// struct. It returns a pointer to the domain.User and an error if the operation
// fails.
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

// nullableUUIDValue is a helper function that returns the value of a nullable UUID.
// If the UUID is nil, it returns nil. Otherwise, it returns the UUID value.
func nullableUUIDValue(value *uuid.UUID) any {
	if value == nil {
		return nil
	}

	return *value
}

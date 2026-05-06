package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
	"github.com/seu-usuario/bank-api/internal/database"
)

type Repository struct {
	db *pgxpool.Pool
}

type dbExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var _ domain.CustomerRepository = (*Repository)(nil)

// New creates a new instance of the customer Repository with the provided pgxpool.Pool. This repository will be used to manage customer data in a PostgreSQL database. It
// provides methods for creating customers, checking if a customer exists by ID,
// and retrieving customer information by ID. The repository uses the pgx library
// to interact with the database and supports transactions through the context.
func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// executor returns the appropriate database executor based on the context. If a
// transaction is present in the context, it returns the transaction; otherwise,
// it returns the main database connection pool. This allows the repository methods
// to work seamlessly with both transactional and non-transactional contexts.
func (r *Repository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

// Create inserts a new customer record into the customers table with the provided
// customer data. It generates a new UUID for the customer and executes the SQL
// query to store the customer information in the database. If the operation is
// successful, it returns nil; otherwise, it returns an error. The method also
// handles specific PostgreSQL errors, such as unique constraint violations and
// check constraint violations, and maps them to domain errors.
// If an unknown error occurs during the database operation, it wraps the error
// with additional context before returning it. This helps to provide more
// informative error messages for debugging and troubleshooting.
func (r *Repository) Create(ctx context.Context, c *domain.Customer) error {
	query := `
		INSERT INTO customers (id, name, cpf, created_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.executor(ctx).Exec(ctx, query,
		c.ID,
		c.Name,
		c.CPF,
		c.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "customers_cpf_key" {
					return domain.ErrCPFAlreadyExists
				}

			case "23514": // check_violation
				return domain.ErrInvalidData
			}
		}

		// wrap unknown infra errors
		return fmt.Errorf("repository create: %w", err)
	}

	return nil
}

// Exists checks if a customer with the given ID exists in the customers table. It
// executes a SQL query that uses the EXISTS clause to determine if a record with
// the specified ID exists. If the query executes successfully, it returns a boolean
// indicating whether the customer exists and a nil error. If there is an error
// during the database operation, it wraps the error with additional context before
// returning it. This helps to provide more informative error messages for debugging
// and troubleshooting.
func (r *Repository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM customers WHERE id = $1
		)
	`

	var exists bool
	err := r.executor(ctx).QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("customer repository exists: %w", err)
	}

	return exists, nil
}

// GetByID retrieves a customer record from the customers table based on the provided
// customer ID. It executes a SQL query that joins the customers and users tables to
// fetch the customer information along with the associated email. If a customer with
// the specified ID is found, it returns a pointer to the Customer entity, the email,
// and a nil error. If no customer is found, it returns nil for the customer, an empty
// string for the email, and a domain.ErrNotFound error. If there is an error during
// the database operation, it wraps the error with additional context before returning
// it. This helps to provide more informative error messages for debugging and
// troubleshooting.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, string, error) {
	query := `
		SELECT c.id, c.name, c.cpf, u.email, c.created_at
		FROM customers c
		JOIN users u ON u.customer_id = c.id
		WHERE c.id = $1
	`

	var customer domain.Customer
	var email string
	err := r.executor(ctx).QueryRow(ctx, query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.CPF,
		&email,
		&customer.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "", domain.ErrNotFound
		}

		return nil, "", fmt.Errorf("customer repository get by id: %w", err)
	}

	return &customer, email, nil
}

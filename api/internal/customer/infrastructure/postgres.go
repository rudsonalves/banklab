package infrastructure

import (
	"context"
	"database/sql"
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
var _ domain.CustomerDocumentRepository = (*Repository)(nil)

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
// customer data. The database generates the ID and the repository writes it back
// to the entity. If the operation is successful, it returns nil; otherwise, it
// returns an error.
// The method handles PostgreSQL constraint violations, particularly check constraint
// violations (code 23514), and maps them to domain errors. For any other database
// errors, it wraps them with context information for better error reporting.
func (r *Repository) Create(ctx context.Context, c *domain.Customer) error {
	query := `
		INSERT INTO customers (name, birth_date, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.executor(ctx).QueryRow(ctx, query,
		c.Name,
		c.BirthDate,
		c.CreatedAt,
	).Scan(&c.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
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

// GetByID retrieves a customer profile from the customers table based on the provided
// customer ID. It executes a SQL query that joins the customers and users tables and
// left joins the customer_documents table to fetch the customer information along with
// the associated email and CPF document. If a customer with the specified ID is found,
// it returns a pointer to the CustomerProfile entity and a nil error. If no customer is
// found, it returns nil and a domain.ErrNotFound error. If there is an error during the
// database operation, it wraps the error with additional context before returning it.
// This helps to provide more informative error messages for debugging and troubleshooting.
func (r *Repository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.CustomerProfile, error) {
	query := `
		SELECT 
			c.id,
		 	c.name, 
			c.birth_date,
			c.created_at, 
			u.email, 
			cd.value
		FROM customers c
		JOIN users u ON u.customer_id = c.id
		LEFT JOIN customer_documents cd 
			ON cd.customer_id = c.id 
			AND cd.type = 'cpf'
			AND cd.country = 'BR' 
			AND cd.is_primary = true
		WHERE c.id = $1
	`

	var profile domain.CustomerProfile
	var birthDate sql.NullTime
	var cpf sql.NullString
	err := r.executor(ctx).QueryRow(ctx, query, id).Scan(
		&profile.Customer.ID,
		&profile.Customer.Name,
		&birthDate,
		&profile.Customer.CreatedAt,
		&profile.Email,
		&cpf,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("customer repository get by id: %w", err)
	}

	if birthDate.Valid {
		profile.Customer.BirthDate = birthDate.Time
	}
	if cpf.Valid {
		profile.CPF = cpf.String
	}

	return &profile, nil
}

// CreateDocument inserts a new customer document record into the customer_documents table
// with the provided document data. The database generates the ID and the repository
// writes it back to the entity. If the operation is successful, it returns nil;
// otherwise, it returns an error.
// The method handles specific PostgreSQL errors such as unique constraint violations
// (returning domain.ErrCPFAlreadyExists for duplicate documents) and check constraint violations
// (returning domain.ErrInvalidData for invalid data). Any other database errors are wrapped
// with additional context before returning.
func (r *Repository) CreateDocument(
	ctx context.Context,
	document *domain.CustomerDocument,
) error {
	query := `
		INSERT INTO customer_documents (
			customer_id,
			type,
			value,
			issuer,
			issuer_state,
			country,
			is_primary,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.executor(ctx).QueryRow(ctx, query,
		document.CustomerID,
		document.Type,
		document.Value,
		document.Issuer,
		document.IssuerState,
		document.Country,
		document.IsPrimary,
		document.CreatedAt,
		document.UpdatedAt,
	).Scan(&document.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				if pgErr.ConstraintName == "customer_documents_unique_document" {
					return domain.ErrCPFAlreadyExists
				}

			case "23514": // check_violation
				return domain.ErrInvalidData
			}
		}

		return fmt.Errorf("repository create document: %w", err)
	}

	return nil
}

// ExistsCPF checks whether a CPF document already exists in the customer_documents table.
func (r *Repository) ExistsCPF(ctx context.Context, cpf string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM customer_documents
			WHERE type = 'cpf'
			  AND country = 'BR'
			  AND value = $1
		)
	`

	var exists bool
	err := r.executor(ctx).QueryRow(ctx, query, cpf).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("customer document repository exists cpf: %w", err)
	}

	return exists, nil
}

// GetPrimaryDocumentByCustomerID retrieves the primary document associated with a given customer ID.
// It queries the customer_documents table for the document where is_primary is true.
// Returns the customer document if found, or domain.ErrNotFound if no primary document exists.
func (r *Repository) GetPrimaryDocumentByCustomerID(
	ctx context.Context,
	customerID uuid.UUID,
) (*domain.CustomerDocument, error) {
	query := `
		SELECT
			id,
			customer_id,
			type,
			value,
			issuer,
			issuer_state,
			country,
			is_primary,
			created_at,
			updated_at
		FROM customer_documents
		WHERE customer_id = $1
		  AND is_primary = true
	`

	return r.getDocument(ctx, query, customerID)
}

// GetCPFByCustomerID retrieves the CPF document associated with a given customer ID.
// It queries the customer_documents table for the document where type is 'cpf' and country is 'BR'.
// Returns the customer document if found, or domain.ErrNotFound if no CPF document exists.
func (r *Repository) GetCPFByCustomerID(
	ctx context.Context,
	customerID uuid.UUID,
) (*domain.CustomerDocument, error) {
	query := `
		SELECT
			id,
			customer_id,
			type,
			value,
			issuer,
			issuer_state,
			country,
			is_primary,
			created_at,
			updated_at
		FROM customer_documents
		WHERE customer_id = $1
		  AND type = 'cpf'
		  AND country = 'BR'
	`

	return r.getDocument(ctx, query, customerID)
}

// getDocument executes a query to retrieve a single customer document from the database
// and maps the result to a domain.CustomerDocument struct. It handles null values for
// optional fields like Issuer and IssuerState, and returns domain.ErrNotFound if no
// rows match the query conditions. Any other database errors are wrapped with context.
func (r *Repository) getDocument(
	ctx context.Context,
	query string,
	args ...any,
) (*domain.CustomerDocument, error) {
	var document domain.CustomerDocument
	var issuer sql.NullString
	var issuerState sql.NullString

	err := r.executor(ctx).QueryRow(ctx, query, args...).Scan(
		&document.ID,
		&document.CustomerID,
		&document.Type,
		&document.Value,
		&issuer,
		&issuerState,
		&document.Country,
		&document.IsPrimary,
		&document.CreatedAt,
		&document.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, fmt.Errorf("customer document repository get: %w", err)
	}

	if issuer.Valid {
		document.Issuer = &issuer.String
	}
	if issuerState.Valid {
		document.IssuerState = &issuerState.String
	}

	return &document, nil
}

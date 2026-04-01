package infrastructure

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

var _ domain.Repository = (*Repository)(nil)

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *domain.Customer) error {
	query := `
		INSERT INTO customers (id, name, cpf, email, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		c.ID,
		c.Name,
		c.CPF,
		c.Email,
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
				if pgErr.ConstraintName == "customers_email_key" {
					return domain.ErrEmailAlreadyExists
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

func (r *Repository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM customers WHERE id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("customer repository exists: %w", err)
	}

	return exists, nil
}

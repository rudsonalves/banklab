package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *domain.Customer) error {
	query := `
		INSERT INTO customers (id, name, cpf, email)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(ctx, query,
		c.ID,
		c.Name,
		c.CPF,
		c.Email,
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
			return err
		}
	}

	return nil
}

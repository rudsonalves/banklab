package infrastructure

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/database"
)

type PostgresTransactor struct {
	db *pgxpool.Pool
}

var _ domain.Transactor = (*PostgresTransactor)(nil)

func NewPostgresTransactor(db *pgxpool.Pool) *PostgresTransactor {
	return &PostgresTransactor{db: db}
}

func (t *PostgresTransactor) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // rollback after commit is a no-op

	if err := fn(database.ContextWithTx(ctx, tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

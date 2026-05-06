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

// NewPostgresTransactor creates a new instance of PostgresTransactor with the
// provided database connection pool.
func NewPostgresTransactor(db *pgxpool.Pool) *PostgresTransactor {
	return &PostgresTransactor{db: db}
}

// RunInTx executes the provided function within a database transaction. It begins
// a new transaction, and if the function returns an error, the transaction is
// rolled back.
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

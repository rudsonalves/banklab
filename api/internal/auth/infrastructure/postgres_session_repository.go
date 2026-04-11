package infrastructure

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/database"
)

type PostgresSessionRepository struct {
	db *pgxpool.Pool
}

var _ domain.SessionRepository = (*PostgresSessionRepository)(nil)

func NewPostgresSessionRepository(db *pgxpool.Pool) *PostgresSessionRepository {
	return &PostgresSessionRepository{db: db}
}

func (r *PostgresSessionRepository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

func (r *PostgresSessionRepository) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO user_sessions (
			id,
			user_id,
			token_hash,
			expires_at
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.executor(ctx).Exec(ctx, query, uuid.New(), userID, tokenHash, expiresAt)
	return err
}

func (r *PostgresSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, bool, error) {
	query := `
		SELECT user_id, expires_at, revoked_at
		FROM user_sessions
		WHERE token_hash = $1
	`

	var userID uuid.UUID
	var expiresAt time.Time
	var revokedAt *time.Time

	err := r.executor(ctx).QueryRow(ctx, query, tokenHash).Scan(&userID, &expiresAt, &revokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, time.Time{}, false, nil
		}
		return uuid.Nil, time.Time{}, false, err
	}

	return userID, expiresAt, revokedAt != nil, nil
}

func (r *PostgresSessionRepository) Revoke(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE user_sessions
		SET revoked_at = NOW()
		WHERE token_hash = $1
	`

	_, err := r.executor(ctx).Exec(ctx, query, tokenHash)
	return err
}

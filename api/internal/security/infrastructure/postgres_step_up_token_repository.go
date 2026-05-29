package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/database"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type PostgresStepUpTokenRepository struct {
	db *pgxpool.Pool
}

var _ domain.StepUpTokenRepository = (*PostgresStepUpTokenRepository)(nil)

func NewPostgresStepUpTokenRepository(db *pgxpool.Pool) *PostgresStepUpTokenRepository {
	return &PostgresStepUpTokenRepository{db: db}
}

// executor returns the current transaction if it exists in the context,
// otherwise it returns the main database connection pool.
func (r *PostgresStepUpTokenRepository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

// Create inserts a new step-up token record into the database and populates
// its generated ID. It maps constraint failures to domain errors so callers do
// not need to inspect PostgreSQL details.
func (r *PostgresStepUpTokenRepository) Create(ctx context.Context, token *domain.StepUpToken) error {
	if err := token.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO step_up_tokens (
			jti,
			user_id,
			endpoint_key,
			status,
			expires_at,
			consumed_at,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := r.executor(ctx).QueryRow(
		ctx,
		query,
		token.JTI,
		token.UserID,
		token.EndpointKey,
		string(token.Status),
		token.ExpiresAt,
		token.ConsumedAt,
		token.CreatedAt,
	).Scan(&token.ID)
	if err != nil {
		if isStepUpTokenConstraintError(err) {
			return domain.ErrInvalidStepUpToken
		}

		return fmt.Errorf("create step-up token: %w", err)
	}

	return nil
}

// FindByJTI retrieves a step-up token by its JWT ID. It returns nil when no
// token exists for the provided JTI.
func (r *PostgresStepUpTokenRepository) FindByJTI(
	ctx context.Context,
	jti string,
) (*domain.StepUpToken, error) {
	query := `
		SELECT
			id,
			jti,
			user_id,
			endpoint_key,
			status,
			expires_at,
			consumed_at,
			created_at
		FROM step_up_tokens
		WHERE jti = $1
	`

	token, err := scanStepUpToken(r.executor(ctx).QueryRow(ctx, query, jti))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		if errors.Is(err, domain.ErrInvalidStepUpToken) {
			return nil, err
		}

		return nil, fmt.Errorf("find step-up token by jti: %w", err)
	}

	return token, nil
}

// ConsumeByJTI atomically marks an active, non-expired token as consumed. If
// the update cannot happen, it maps the current persisted state to the matching
// domain error.
func (r *PostgresStepUpTokenRepository) ConsumeByJTI(
	ctx context.Context,
	jti string,
	now time.Time,
) (*domain.StepUpToken, error) {
	now = now.UTC()

	query := `
		UPDATE step_up_tokens
		SET
			status = $1,
			consumed_at = $2
		WHERE jti = $3
			AND status = $4
			AND expires_at >= $2
		RETURNING
			id,
			jti,
			user_id,
			endpoint_key,
			status,
			expires_at,
			consumed_at,
			created_at
	`

	token, err := scanStepUpToken(r.executor(ctx).QueryRow(
		ctx,
		query,
		string(domain.StepUpTokenConsumed),
		now,
		jti,
		string(domain.StepUpTokenActive),
	))
	if err == nil {
		return token, nil
	}

	if !errors.Is(err, pgx.ErrNoRows) {
		if errors.Is(err, domain.ErrInvalidStepUpToken) {
			return nil, err
		}

		return nil, fmt.Errorf("consume step-up token by jti: %w", err)
	}

	return nil, r.consumeFailureError(ctx, jti, now)
}

// consumeFailureError checks the current state of the token to determine the
// appropriate domain error to return when a token consumption attempt fails.
func (r *PostgresStepUpTokenRepository) consumeFailureError(
	ctx context.Context,
	jti string,
	now time.Time,
) error {
	token, err := r.FindByJTI(ctx, jti)
	if err != nil {
		return err
	}
	if token == nil {
		return domain.ErrInvalidStepUpToken
	}

	if token.Status == domain.StepUpTokenConsumed {
		return domain.ErrStepUpTokenConsumed
	}
	if token.IsExpired(now) {
		return domain.ErrStepUpTokenExpired
	}

	return domain.ErrInvalidStepUpToken
}

// scanStepUpToken maps a database row to a StepUpToken domain entity, handling
// any necessary conversions and validations.
func scanStepUpToken(row pgx.Row) (*domain.StepUpToken, error) {
	var id uuid.UUID
	var userID uuid.UUID
	var jti string
	var endpointKey string
	var status string
	var expiresAt time.Time
	var consumedAt *time.Time
	var createdAt time.Time

	err := row.Scan(
		&id,
		&jti,
		&userID,
		&endpointKey,
		&status,
		&expiresAt,
		&consumedAt,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	return domain.RestoreStepUpToken(
		id,
		jti,
		userID,
		endpointKey,
		domain.StepUpTokenStatus(status),
		expiresAt,
		consumedAt,
		createdAt,
	)
}

// isStepUpTokenConstraintError checks if the provided error is a PostgreSQL
// constraint violation related to step-up tokens.
func isStepUpTokenConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23503"
}

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
	"github.com/seu-usuario/bank-api/internal/installation/domain"
)

type PostgresRestrictedAuthorizationRepository struct {
	db *pgxpool.Pool
}

var _ domain.RestrictedAuthorizationRepository = (*PostgresRestrictedAuthorizationRepository)(nil)

func NewPostgresRestrictedAuthorizationRepository(db *pgxpool.Pool) *PostgresRestrictedAuthorizationRepository {
	return &PostgresRestrictedAuthorizationRepository{db: db}
}

func (r *PostgresRestrictedAuthorizationRepository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

func (r *PostgresRestrictedAuthorizationRepository) Create(
	ctx context.Context,
	authorization *domain.RestrictedAuthorization,
) error {
	if err := authorization.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO installation_registration_authorizations (
			jti,
			user_id,
			installation_id,
			scope,
			status,
			expires_at,
			consumed_at,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.executor(ctx).QueryRow(
		ctx,
		query,
		authorization.JTI,
		authorization.UserID,
		authorization.InstallationID.UUID(),
		authorization.Scope,
		string(authorization.Status),
		authorization.ExpiresAt,
		authorization.ConsumedAt,
		authorization.CreatedAt,
	).Scan(&authorization.ID)
	if err != nil {
		return mapRestrictedAuthorizationCreateError(err)
	}

	return nil
}

func (r *PostgresRestrictedAuthorizationRepository) FindByJTI(
	ctx context.Context,
	jti string,
) (*domain.RestrictedAuthorization, error) {
	query := restrictedAuthorizationSelectSQL + `
		WHERE jti = $1
	`

	authorization, err := scanRestrictedAuthorization(r.executor(ctx).QueryRow(ctx, query, jti))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrRestrictedAuthorizationNotFound
		}
		if errors.Is(err, domain.ErrInvalidRestrictedAuthorization) {
			return nil, err
		}
		return nil, fmt.Errorf("find restricted authorization by jti: %w", err)
	}

	return authorization, nil
}

func (r *PostgresRestrictedAuthorizationRepository) ConsumeByJTI(
	ctx context.Context,
	jti string,
	now time.Time,
) (*domain.RestrictedAuthorization, error) {
	now = now.UTC()
	query := `
		UPDATE installation_registration_authorizations
		SET
			status = 'consumed',
			consumed_at = $2
		WHERE jti = $1
			AND status = 'active'
			AND expires_at > $2
		RETURNING
			id,
			jti,
			user_id,
			installation_id,
			scope,
			status,
			expires_at,
			consumed_at,
			created_at
	`

	authorization, err := scanRestrictedAuthorization(r.executor(ctx).QueryRow(ctx, query, jti, now))
	if err == nil {
		return authorization, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		if errors.Is(err, domain.ErrInvalidRestrictedAuthorization) {
			return nil, err
		}
		return nil, fmt.Errorf("consume restricted authorization by jti: %w", err)
	}

	return nil, r.consumeFailureError(ctx, jti, now)
}

func (r *PostgresRestrictedAuthorizationRepository) RevokeByJTI(ctx context.Context, jti string) error {
	commandTag, err := r.executor(ctx).Exec(
		ctx,
		`UPDATE installation_registration_authorizations
		 SET status = 'revoked'
		 WHERE jti = $1
			AND status = 'active'`,
		jti,
	)
	if err != nil {
		return fmt.Errorf("revoke restricted authorization by jti: %w", err)
	}
	if commandTag.RowsAffected() > 0 {
		return nil
	}

	authorization, err := r.FindByJTI(ctx, jti)
	if err != nil {
		return err
	}
	switch authorization.Status {
	case domain.RestrictedAuthorizationStatusConsumed:
		return domain.ErrRestrictedAuthorizationConsumed
	case domain.RestrictedAuthorizationStatusRevoked:
		return domain.ErrRestrictedAuthorizationRevoked
	default:
		return domain.ErrRestrictedAuthorizationInvalid
	}
}

func (r *PostgresRestrictedAuthorizationRepository) RevokeActiveByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID domain.InstallationID,
	scope string,
) error {
	_, err := r.executor(ctx).Exec(
		ctx,
		`UPDATE installation_registration_authorizations
		 SET status = 'revoked'
		 WHERE user_id = $1
			AND installation_id = $2
			AND scope = $3
			AND status = 'active'`,
		userID,
		installationID.UUID(),
		scope,
	)
	if err != nil {
		return fmt.Errorf("revoke active restricted authorization: %w", err)
	}

	return nil
}

func (r *PostgresRestrictedAuthorizationRepository) consumeFailureError(
	ctx context.Context,
	jti string,
	now time.Time,
) error {
	authorization, err := r.FindByJTI(ctx, jti)
	if err != nil {
		return err
	}

	switch authorization.Status {
	case domain.RestrictedAuthorizationStatusConsumed:
		return domain.ErrRestrictedAuthorizationConsumed
	case domain.RestrictedAuthorizationStatusRevoked:
		return domain.ErrRestrictedAuthorizationRevoked
	case domain.RestrictedAuthorizationStatusActive:
		if authorization.IsExpired(now) {
			return domain.ErrRestrictedAuthorizationExpired
		}
	}

	return domain.ErrRestrictedAuthorizationInvalid
}

const restrictedAuthorizationSelectSQL = `
	SELECT
		id,
		jti,
		user_id,
		installation_id,
		scope,
		status,
		expires_at,
		consumed_at,
		created_at
	FROM installation_registration_authorizations
`

func scanRestrictedAuthorization(row pgx.Row) (*domain.RestrictedAuthorization, error) {
	var id uuid.UUID
	var jti string
	var userID uuid.UUID
	var installationUUID uuid.UUID
	var scope string
	var status string
	var expiresAt time.Time
	var consumedAt *time.Time
	var createdAt time.Time

	err := row.Scan(
		&id,
		&jti,
		&userID,
		&installationUUID,
		&scope,
		&status,
		&expiresAt,
		&consumedAt,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	installationID, err := domain.NewInstallationID(installationUUID)
	if err != nil {
		return nil, err
	}

	return domain.RestoreRestrictedAuthorization(
		id,
		jti,
		userID,
		installationID,
		scope,
		domain.RestrictedAuthorizationStatus(status),
		expiresAt,
		consumedAt,
		createdAt,
	)
}

func mapRestrictedAuthorizationCreateError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return fmt.Errorf("create restricted authorization: %w", err)
	}

	if pgErr.Code == "23505" {
		if pgErr.ConstraintName == "ux_installation_registration_authorizations_active" {
			return domain.ErrRestrictedAuthorizationAlreadyActive
		}
		return domain.ErrInvalidRestrictedAuthorization
	}
	if pgErr.Code == "23514" || pgErr.Code == "23503" {
		return domain.ErrInvalidRestrictedAuthorization
	}

	return fmt.Errorf("create restricted authorization: %w", err)
}

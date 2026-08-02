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

type dbExecutor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type PostgresInstallationRepository struct {
	db *pgxpool.Pool
}

var _ domain.InstallationRepository = (*PostgresInstallationRepository)(nil)

func NewPostgresInstallationRepository(db *pgxpool.Pool) *PostgresInstallationRepository {
	return &PostgresInstallationRepository{db: db}
}

func (r *PostgresInstallationRepository) executor(ctx context.Context) dbExecutor {
	if tx, ok := database.TxFromContext(ctx); ok {
		return tx
	}

	return r.db
}

func (r *PostgresInstallationRepository) FindByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID domain.InstallationID,
) (*domain.Installation, error) {
	query := installationSelectSQL + `
		WHERE user_id = $1
			AND installation_id = $2
	`

	installation, err := scanInstallation(r.executor(ctx).QueryRow(ctx, query, userID, installationID.UUID()))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInstallationNotFound
		}
		if errors.Is(err, domain.ErrInvalidInstallation) {
			return nil, err
		}

		return nil, fmt.Errorf("find installation by user and installation id: %w", err)
	}

	return installation, nil
}

func (r *PostgresInstallationRepository) FindByResourceID(
	ctx context.Context,
	userID uuid.UUID,
	resourceID domain.InstallationResourceID,
) (*domain.Installation, error) {
	query := installationSelectSQL + `
		WHERE user_id = $1
			AND resource_id = $2
	`

	installation, err := scanInstallation(r.executor(ctx).QueryRow(ctx, query, userID, resourceID.UUID()))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInstallationNotFound
		}
		if errors.Is(err, domain.ErrInvalidInstallation) {
			return nil, err
		}

		return nil, fmt.Errorf("find installation by resource id: %w", err)
	}

	return installation, nil
}

func (r *PostgresInstallationRepository) CountKnownByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.executor(ctx).QueryRow(
		ctx,
		`SELECT COUNT(*) FROM app_installations WHERE user_id = $1 AND status = 'known'`,
		userID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count known installations by user id: %w", err)
	}

	return count, nil
}

func (r *PostgresInstallationRepository) HasAnyByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.executor(ctx).QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM app_installations WHERE user_id = $1)`,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check installations history by user id: %w", err)
	}

	return exists, nil
}

func (r *PostgresInstallationRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Installation, error) {
	query := installationSelectSQL + `
		WHERE user_id = $1
		ORDER BY created_at DESC, resource_id DESC
	`

	rows, err := r.executor(ctx).Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list installations by user id: %w", err)
	}
	defer rows.Close()

	var installations []*domain.Installation
	for rows.Next() {
		installation, err := scanInstallation(rows)
		if err != nil {
			return nil, err
		}
		installations = append(installations, installation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate installations by user id: %w", err)
	}

	return installations, nil
}

func (r *PostgresInstallationRepository) BootstrapFirstInstallation(
	ctx context.Context,
	userID uuid.UUID,
	resourceID domain.InstallationResourceID,
	installationID domain.InstallationID,
	now time.Time,
) (*domain.Installation, error) {
	var result *domain.Installation

	err := r.inSerializableTx(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		if err := lockUserForUpdate(txCtx, tx, userID); err != nil {
			return err
		}

		hasAny, err := r.hasAnyByUserID(txCtx, tx, userID)
		if err != nil {
			return err
		}
		if hasAny {
			return domain.ErrFirstInstallationAlreadyBootstrapped
		}

		installation, err := insertKnownInstallation(txCtx, tx, userID, resourceID, installationID, 1, now)
		if err != nil {
			return err
		}

		result = installation
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PostgresInstallationRepository) ReserveKnownInstallation(
	ctx context.Context,
	userID uuid.UUID,
	resourceID domain.InstallationResourceID,
	installationID domain.InstallationID,
	maxKnownInstallations int,
	now time.Time,
) (*domain.Installation, error) {
	if maxKnownInstallations <= 0 || maxKnownInstallations > domain.MaxKnownInstallations {
		return nil, domain.ErrInvalidInstallation
	}

	var result *domain.Installation

	err := r.inSerializableTx(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		if err := lockUserForUpdate(txCtx, tx, userID); err != nil {
			return err
		}

		slot, err := findAvailableKnownSlot(txCtx, tx, userID, maxKnownInstallations)
		if err != nil {
			return err
		}

		installation, err := insertKnownInstallation(txCtx, tx, userID, resourceID, installationID, slot, now)
		if err != nil {
			return err
		}

		result = installation
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PostgresInstallationRepository) RevokeByResourceID(
	ctx context.Context,
	userID uuid.UUID,
	resourceID domain.InstallationResourceID,
	now time.Time,
) (*domain.Installation, error) {
	var result *domain.Installation

	err := r.inSerializableTx(ctx, func(txCtx context.Context, tx pgx.Tx) error {
		if err := lockUserForUpdate(txCtx, tx, userID); err != nil {
			return err
		}

		query := installationSelectSQL + `
			WHERE user_id = $1
				AND resource_id = $2
			FOR UPDATE
		`

		current, err := scanInstallation(tx.QueryRow(txCtx, query, userID, resourceID.UUID()))
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrInstallationNotFound
			}
			return err
		}
		if current.Status == domain.InstallationStatusRevoked {
			return domain.ErrInstallationRevoked
		}

		revokedAt := now.UTC()
		updateQuery := `
			UPDATE app_installations
			SET
				status = 'revoked',
				known_slot = NULL,
				revoked_at = $2,
				updated_at = $2
			WHERE id = $1
			RETURNING
				id,
				resource_id,
				user_id,
				installation_id,
				status,
				platform,
				app_version,
				app_build,
				first_seen_at,
				last_seen_at,
				revoked_at,
				created_at,
				updated_at
		`

		installation, err := scanInstallation(tx.QueryRow(txCtx, updateQuery, current.ID, revokedAt))
		if err != nil {
			return err
		}

		result = installation
		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PostgresInstallationRepository) inSerializableTx(
	ctx context.Context,
	fn func(context.Context, pgx.Tx) error,
) error {
	if tx, ok := database.TxFromContext(ctx); ok {
		return fn(ctx, tx)
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("begin installation transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txCtx := database.ContextWithTx(ctx, tx)
	if err := fn(txCtx, tx); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit installation transaction: %w", err)
	}

	return nil
}

func (r *PostgresInstallationRepository) hasAnyByUserID(ctx context.Context, executor dbExecutor, userID uuid.UUID) (bool, error) {
	var exists bool
	err := executor.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM app_installations WHERE user_id = $1)`,
		userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check installations history by user id: %w", err)
	}

	return exists, nil
}

func lockUserForUpdate(ctx context.Context, executor dbExecutor, userID uuid.UUID) error {
	var lockedID uuid.UUID
	err := executor.QueryRow(ctx, `SELECT id FROM users WHERE id = $1 FOR UPDATE`, userID).Scan(&lockedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInvalidInstallation
		}
		return fmt.Errorf("lock user for installation operation: %w", err)
	}

	return nil
}

func findAvailableKnownSlot(
	ctx context.Context,
	executor dbExecutor,
	userID uuid.UUID,
	maxKnownInstallations int,
) (int, error) {
	rows, err := executor.Query(
		ctx,
		`SELECT known_slot
		 FROM app_installations
		 WHERE user_id = $1
			AND status = 'known'
		 ORDER BY known_slot
		 FOR UPDATE`,
		userID,
	)
	if err != nil {
		return 0, fmt.Errorf("query known installation slots: %w", err)
	}
	defer rows.Close()

	used := make(map[int]struct{}, maxKnownInstallations)
	for rows.Next() {
		var slot int
		if err := rows.Scan(&slot); err != nil {
			return 0, fmt.Errorf("scan known installation slot: %w", err)
		}
		used[slot] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("iterate known installation slots: %w", err)
	}

	for slot := 1; slot <= maxKnownInstallations; slot++ {
		if _, ok := used[slot]; !ok {
			return slot, nil
		}
	}

	return 0, domain.ErrInstallationLimitReached
}

func insertKnownInstallation(
	ctx context.Context,
	executor dbExecutor,
	userID uuid.UUID,
	resourceID domain.InstallationResourceID,
	installationID domain.InstallationID,
	knownSlot int,
	now time.Time,
) (*domain.Installation, error) {
	query := `
		INSERT INTO app_installations (
			resource_id,
			user_id,
			installation_id,
			status,
			known_slot,
			first_seen_at,
			last_seen_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, 'known', $4, $5, $5, $5, $5)
		RETURNING
			id,
			resource_id,
			user_id,
			installation_id,
			status,
			platform,
			app_version,
			app_build,
			first_seen_at,
			last_seen_at,
			revoked_at,
			created_at,
			updated_at
	`

	installation, err := scanInstallation(executor.QueryRow(
		ctx,
		query,
		resourceID.UUID(),
		userID,
		installationID.UUID(),
		knownSlot,
		now.UTC(),
	))
	if err != nil {
		if isInstallationUniqueViolation(err) {
			return nil, domain.ErrInvalidInstallation
		}
		if isInstallationConstraintError(err) {
			return nil, domain.ErrInvalidInstallation
		}
		return nil, fmt.Errorf("insert known installation: %w", err)
	}

	return installation, nil
}

const installationSelectSQL = `
	SELECT
		id,
		resource_id,
		user_id,
		installation_id,
		status,
		platform,
		app_version,
		app_build,
		first_seen_at,
		last_seen_at,
		revoked_at,
		created_at,
		updated_at
	FROM app_installations
`

func scanInstallation(row pgx.Row) (*domain.Installation, error) {
	var id uuid.UUID
	var resourceUUID uuid.UUID
	var userID uuid.UUID
	var installationUUID uuid.UUID
	var status string
	var platform *string
	var appVersion *string
	var appBuild *string
	var firstSeenAt time.Time
	var lastSeenAt time.Time
	var revokedAt *time.Time
	var createdAt time.Time
	var updatedAt time.Time

	err := row.Scan(
		&id,
		&resourceUUID,
		&userID,
		&installationUUID,
		&status,
		&platform,
		&appVersion,
		&appBuild,
		&firstSeenAt,
		&lastSeenAt,
		&revokedAt,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	resourceID, err := domain.NewInstallationResourceID(resourceUUID)
	if err != nil {
		return nil, err
	}
	installationID, err := domain.NewInstallationID(installationUUID)
	if err != nil {
		return nil, err
	}

	return domain.RestoreInstallation(
		id,
		resourceID,
		userID,
		installationID,
		domain.InstallationStatus(status),
		stringValue(platform),
		stringValue(appVersion),
		stringValue(appBuild),
		firstSeenAt,
		lastSeenAt,
		revokedAt,
		createdAt,
		updatedAt,
	)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func isInstallationConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" || pgErr.Code == "23514" || pgErr.Code == "23503"
}

func isInstallationUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

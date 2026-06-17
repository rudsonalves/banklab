package infrastructure

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/installation/domain"
)

func TestPostgresInstallationRepository_Integration(t *testing.T) {
	ctx := context.Background()
	pool := newInstallationTestPool(t, ctx)
	defer pool.Close()

	ensureInstallationTestSchema(t, ctx, pool)

	repo := NewPostgresInstallationRepository(pool)

	t.Run("bootstrap first installation", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		resourceID := mustTestResourceID(t)
		installationID := mustTestInstallationID(t)
		now := time.Now().UTC().Truncate(time.Microsecond)

		installation, err := repo.BootstrapFirstInstallation(ctx, userID, resourceID, installationID, now)
		if err != nil {
			t.Fatalf("expected bootstrap to succeed, got %v", err)
		}
		if installation.Status != domain.InstallationStatusKnown {
			t.Fatalf("expected known status, got %q", installation.Status)
		}

		got, err := repo.FindByUserIDAndInstallationID(ctx, userID, installationID)
		if err != nil {
			t.Fatalf("expected find by installation id to succeed, got %v", err)
		}
		if got.ResourceID != resourceID {
			t.Fatalf("expected resource id %q, got %q", resourceID, got.ResourceID)
		}

		hasAny, err := repo.HasAnyByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("expected has any to succeed, got %v", err)
		}
		if !hasAny {
			t.Fatal("expected user to have installation history")
		}

		count, err := repo.CountKnownByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("expected count known to succeed, got %v", err)
		}
		if count != 1 {
			t.Fatalf("expected known count 1, got %d", count)
		}
	})

	t.Run("bootstrap cannot run when revoked history exists", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		first, err := repo.BootstrapFirstInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), time.Now())
		if err != nil {
			t.Fatalf("expected first bootstrap, got %v", err)
		}
		if _, err := repo.RevokeByResourceID(ctx, userID, first.ResourceID, time.Now().Add(time.Second)); err != nil {
			t.Fatalf("expected revoke to succeed, got %v", err)
		}

		_, err = repo.BootstrapFirstInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), time.Now().Add(2*time.Second))
		if !errors.Is(err, domain.ErrFirstInstallationAlreadyBootstrapped) {
			t.Fatalf("expected ErrFirstInstallationAlreadyBootstrapped, got %v", err)
		}
	})

	t.Run("reserve known installation respects limit and revoked slot", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)

		first, err := repo.ReserveKnownInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), domain.MaxKnownInstallations, time.Now())
		if err != nil {
			t.Fatalf("expected first reservation, got %v", err)
		}
		for i := 0; i < 2; i++ {
			_, err := repo.ReserveKnownInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), domain.MaxKnownInstallations, time.Now())
			if err != nil {
				t.Fatalf("expected reservation %d to succeed, got %v", i+2, err)
			}
		}

		_, err = repo.ReserveKnownInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), domain.MaxKnownInstallations, time.Now())
		if !errors.Is(err, domain.ErrInstallationLimitReached) {
			t.Fatalf("expected ErrInstallationLimitReached, got %v", err)
		}

		revoked, err := repo.RevokeByResourceID(ctx, userID, first.ResourceID, time.Now().Add(time.Minute))
		if err != nil {
			t.Fatalf("expected revoke to succeed, got %v", err)
		}
		if revoked.Status != domain.InstallationStatusRevoked {
			t.Fatalf("expected revoked status, got %q", revoked.Status)
		}

		_, err = repo.ReserveKnownInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), domain.MaxKnownInstallations, time.Now())
		if err != nil {
			t.Fatalf("expected reservation after revoked slot to succeed, got %v", err)
		}
	})

	t.Run("list includes revoked history", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		first, err := repo.BootstrapFirstInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), time.Now())
		if err != nil {
			t.Fatalf("expected bootstrap to succeed, got %v", err)
		}
		if _, err := repo.RevokeByResourceID(ctx, userID, first.ResourceID, time.Now().Add(time.Second)); err != nil {
			t.Fatalf("expected revoke to succeed, got %v", err)
		}

		installations, err := repo.ListByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("expected list to succeed, got %v", err)
		}
		if len(installations) != 1 {
			t.Fatalf("expected 1 installation, got %d", len(installations))
		}
		if installations[0].Status != domain.InstallationStatusRevoked {
			t.Fatalf("expected revoked history, got %q", installations[0].Status)
		}
	})

	t.Run("concurrent bootstrap creates only one first installation", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)

		var wg sync.WaitGroup
		errs := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := repo.BootstrapFirstInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), time.Now())
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)

		successes := 0
		conflicts := 0
		for err := range errs {
			switch {
			case err == nil:
				successes++
			case errors.Is(err, domain.ErrFirstInstallationAlreadyBootstrapped):
				conflicts++
			default:
				t.Fatalf("unexpected bootstrap error: %v", err)
			}
		}
		if successes != 1 || conflicts != 1 {
			t.Fatalf("expected 1 success and 1 conflict, got successes=%d conflicts=%d", successes, conflicts)
		}
	})

	t.Run("concurrent reserve does not exceed known limit", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)

		var wg sync.WaitGroup
		errs := make(chan error, 5)
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := repo.ReserveKnownInstallation(ctx, userID, mustTestResourceID(t), mustTestInstallationID(t), domain.MaxKnownInstallations, time.Now())
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)

		successes := 0
		limitReached := 0
		for err := range errs {
			switch {
			case err == nil:
				successes++
			case errors.Is(err, domain.ErrInstallationLimitReached):
				limitReached++
			default:
				t.Fatalf("unexpected reserve error: %v", err)
			}
		}
		if successes != domain.MaxKnownInstallations || limitReached != 2 {
			t.Fatalf("expected %d successes and 2 limit errors, got successes=%d limit=%d", domain.MaxKnownInstallations, successes, limitReached)
		}
	})
}

func TestPostgresRestrictedAuthorizationRepository_Integration(t *testing.T) {
	ctx := context.Background()
	pool := newInstallationTestPool(t, ctx)
	defer pool.Close()

	ensureInstallationTestSchema(t, ctx, pool)

	repo := NewPostgresRestrictedAuthorizationRepository(pool)

	t.Run("create and find by jti", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		authorization := mustRestrictedAuthorization(t, userID, mustTestInstallationID(t), time.Now())

		if err := repo.Create(ctx, authorization); err != nil {
			t.Fatalf("expected create to succeed, got %v", err)
		}
		if authorization.ID == uuid.Nil {
			t.Fatal("expected repository to populate id")
		}

		got, err := repo.FindByJTI(ctx, authorization.JTI)
		if err != nil {
			t.Fatalf("expected find to succeed, got %v", err)
		}
		if got.JTI != authorization.JTI {
			t.Fatalf("expected jti %q, got %q", authorization.JTI, got.JTI)
		}
	})

	t.Run("active authorization is unique per user installation scope", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		installationID := mustTestInstallationID(t)

		first := mustRestrictedAuthorization(t, userID, installationID, time.Now())
		second := mustRestrictedAuthorization(t, userID, installationID, time.Now())
		if err := repo.Create(ctx, first); err != nil {
			t.Fatalf("expected first create, got %v", err)
		}

		err := repo.Create(ctx, second)
		if !errors.Is(err, domain.ErrRestrictedAuthorizationAlreadyActive) {
			t.Fatalf("expected ErrRestrictedAuthorizationAlreadyActive, got %v", err)
		}
	})

	t.Run("consume is single use", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		now := time.Now().UTC().Truncate(time.Microsecond)
		authorization := mustRestrictedAuthorization(t, userID, mustTestInstallationID(t), now)
		if err := repo.Create(ctx, authorization); err != nil {
			t.Fatalf("expected create, got %v", err)
		}

		consumed, err := repo.ConsumeByJTI(ctx, authorization.JTI, now.Add(time.Minute))
		if err != nil {
			t.Fatalf("expected consume, got %v", err)
		}
		if consumed.Status != domain.RestrictedAuthorizationStatusConsumed {
			t.Fatalf("expected consumed status, got %q", consumed.Status)
		}

		_, err = repo.ConsumeByJTI(ctx, authorization.JTI, now.Add(2*time.Minute))
		if !errors.Is(err, domain.ErrRestrictedAuthorizationConsumed) {
			t.Fatalf("expected ErrRestrictedAuthorizationConsumed, got %v", err)
		}
	})

	t.Run("consume expired authorization returns expired", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		now := time.Now().UTC().Truncate(time.Microsecond)
		authorization := mustRestrictedAuthorization(t, userID, mustTestInstallationID(t), now)
		if err := repo.Create(ctx, authorization); err != nil {
			t.Fatalf("expected create, got %v", err)
		}

		_, err := repo.ConsumeByJTI(ctx, authorization.JTI, now.Add(domain.RestrictedAuthorizationDefaultDuration+time.Second))
		if !errors.Is(err, domain.ErrRestrictedAuthorizationExpired) {
			t.Fatalf("expected ErrRestrictedAuthorizationExpired, got %v", err)
		}
	})

	t.Run("concurrent consume succeeds once", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		now := time.Now().UTC().Truncate(time.Microsecond)
		authorization := mustRestrictedAuthorization(t, userID, mustTestInstallationID(t), now)
		if err := repo.Create(ctx, authorization); err != nil {
			t.Fatalf("expected create, got %v", err)
		}

		var wg sync.WaitGroup
		errs := make(chan error, 2)
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := repo.ConsumeByJTI(ctx, authorization.JTI, now.Add(time.Minute))
				errs <- err
			}()
		}
		wg.Wait()
		close(errs)

		successes := 0
		consumedErrors := 0
		for err := range errs {
			switch {
			case err == nil:
				successes++
			case errors.Is(err, domain.ErrRestrictedAuthorizationConsumed):
				consumedErrors++
			default:
				t.Fatalf("unexpected consume error: %v", err)
			}
		}
		if successes != 1 || consumedErrors != 1 {
			t.Fatalf("expected 1 success and 1 consumed error, got successes=%d consumed=%d", successes, consumedErrors)
		}
	})

	t.Run("cleanup removes only authorizations outside retention", func(t *testing.T) {
		userID := createInstallationTestUser(t, ctx, pool)
		defer cleanupInstallationTestUser(t, ctx, pool, userID)
		now := time.Now().UTC().Truncate(time.Microsecond)
		retention := 24 * time.Hour

		oldActive := mustRestrictedAuthorizationAt(
			t,
			userID,
			mustTestInstallationID(t),
			domain.RestrictedAuthorizationStatusActive,
			now.Add(-48*time.Hour),
			nil,
		)
		oldConsumedAt := now.Add(-25 * time.Hour)
		oldConsumed := mustRestrictedAuthorizationAt(
			t,
			userID,
			mustTestInstallationID(t),
			domain.RestrictedAuthorizationStatusConsumed,
			now.Add(-48*time.Hour),
			&oldConsumedAt,
		)
		oldRevoked := mustRestrictedAuthorizationAt(
			t,
			userID,
			mustTestInstallationID(t),
			domain.RestrictedAuthorizationStatusRevoked,
			now.Add(-48*time.Hour),
			nil,
		)
		recentActive := mustRestrictedAuthorizationAt(
			t,
			userID,
			mustTestInstallationID(t),
			domain.RestrictedAuthorizationStatusActive,
			now.Add(-time.Hour),
			nil,
		)

		for _, authorization := range []*domain.RestrictedAuthorization{
			oldActive,
			oldConsumed,
			oldRevoked,
			recentActive,
		} {
			if err := repo.Create(ctx, authorization); err != nil {
				t.Fatalf("expected create %q, got %v", authorization.JTI, err)
			}
		}

		deleted, err := repo.CleanupExpired(ctx, now, retention)
		if err != nil {
			t.Fatalf("expected cleanup to succeed, got %v", err)
		}
		if deleted != 3 {
			t.Fatalf("expected 3 deleted authorizations, got %d", deleted)
		}

		for _, authorization := range []*domain.RestrictedAuthorization{oldActive, oldConsumed, oldRevoked} {
			_, err := repo.FindByJTI(ctx, authorization.JTI)
			if !errors.Is(err, domain.ErrRestrictedAuthorizationNotFound) {
				t.Fatalf("expected authorization %q to be removed, got %v", authorization.JTI, err)
			}
		}

		got, err := repo.FindByJTI(ctx, recentActive.JTI)
		if err != nil {
			t.Fatalf("expected recent authorization to remain, got %v", err)
		}
		if got.Status != domain.RestrictedAuthorizationStatusActive {
			t.Fatalf("expected recent authorization to stay active, got %q", got.Status)
		}
	})
}

func newInstallationTestPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	connString := os.Getenv("BANK_TEST_DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:postgres@localhost:5432/bank_test?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Skipf("skipping integration test: cannot create pool: %v", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		t.Skipf("skipping integration test: database unavailable: %v", err)
	}

	return pool
}

func ensureInstallationTestSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS pgcrypto`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(120) NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS app_installations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			resource_id UUID NOT NULL DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			installation_id UUID NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'known',
			known_slot SMALLINT,
			platform VARCHAR(40),
			app_version VARCHAR(40),
			app_build VARCHAR(40),
			first_seen_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			last_seen_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			revoked_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_app_installations_status CHECK (status IN ('known', 'revoked')),
			CONSTRAINT chk_app_installations_revoked_at_consistency CHECK (
				(status = 'known' AND revoked_at IS NULL)
				OR
				(status = 'revoked' AND revoked_at IS NOT NULL)
			),
			CONSTRAINT chk_app_installations_known_slot_consistency CHECK (
				(status = 'known' AND known_slot BETWEEN 1 AND 3)
				OR
				(status = 'revoked' AND known_slot IS NULL)
			)
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_app_installations_resource_id ON app_installations (resource_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_app_installations_user_installation_id ON app_installations (user_id, installation_id)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_app_installations_known_slot ON app_installations (user_id, known_slot) WHERE status = 'known'`,
		`CREATE TABLE IF NOT EXISTS installation_registration_authorizations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			jti VARCHAR(120) NOT NULL,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			installation_id UUID NOT NULL,
			scope VARCHAR(80) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			consumed_at TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_installation_registration_authorizations_scope CHECK (scope = 'installation.register'),
			CONSTRAINT chk_installation_registration_authorizations_status CHECK (status IN ('active', 'consumed', 'revoked')),
			CONSTRAINT chk_installation_registration_authorizations_consumed_at_consistency CHECK (
				(status = 'active' AND consumed_at IS NULL)
				OR
				(status = 'consumed' AND consumed_at IS NOT NULL)
				OR
				(status = 'revoked')
			)
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_installation_registration_authorizations_jti ON installation_registration_authorizations (jti)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_installation_registration_authorizations_active ON installation_registration_authorizations (user_id, installation_id, scope) WHERE status = 'active'`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("failed to ensure installation repository test schema: %v", err)
		}
	}
}

func createInstallationTestUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	var userID uuid.UUID
	email := strings.ToLower(uuid.NewString()) + "@example.com"
	err := pool.QueryRow(
		ctx,
		`INSERT INTO users (email, password_hash, role, status)
		 VALUES ($1, 'hash', 'admin', 'active')
		 RETURNING id`,
		email,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create installation test user: %v", err)
	}

	return userID
}

func cleanupInstallationTestUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) {
	t.Helper()

	if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
		t.Logf("cleanup warning: failed to delete user %q: %v", userID, err)
	}
}

func mustTestInstallationID(t *testing.T) domain.InstallationID {
	t.Helper()

	id, err := domain.NewInstallationID(uuid.New())
	if err != nil {
		t.Fatalf("expected installation id, got %v", err)
	}
	return id
}

func mustTestResourceID(t *testing.T) domain.InstallationResourceID {
	t.Helper()

	id, err := domain.NewInstallationResourceID(uuid.New())
	if err != nil {
		t.Fatalf("expected resource id, got %v", err)
	}
	return id
}

func mustRestrictedAuthorization(
	t *testing.T,
	userID uuid.UUID,
	installationID domain.InstallationID,
	now time.Time,
) *domain.RestrictedAuthorization {
	t.Helper()

	authorization, err := domain.NewRestrictedAuthorization(uuid.NewString(), userID, installationID, now)
	if err != nil {
		t.Fatalf("expected restricted authorization, got %v", err)
	}
	return authorization
}

func mustRestrictedAuthorizationAt(
	t *testing.T,
	userID uuid.UUID,
	installationID domain.InstallationID,
	status domain.RestrictedAuthorizationStatus,
	createdAt time.Time,
	consumedAt *time.Time,
) *domain.RestrictedAuthorization {
	t.Helper()

	authorization, err := domain.RestoreRestrictedAuthorization(
		uuid.Nil,
		uuid.NewString(),
		userID,
		installationID,
		domain.RestrictedAuthorizationScopeInstallationRegister,
		status,
		createdAt.Add(domain.RestrictedAuthorizationDefaultDuration),
		consumedAt,
		createdAt,
	)
	if err != nil {
		t.Fatalf("expected restricted authorization, got %v", err)
	}
	return authorization
}

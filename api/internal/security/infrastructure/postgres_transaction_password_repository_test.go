package infrastructure

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

func TestPostgresTransactionPasswordRepository_Integration(t *testing.T) {
	ctx := context.Background()
	pool := newSecurityTestPool(t, ctx)
	defer pool.Close()

	ensureSecurityRepoTestSchema(t, ctx, pool)

	repo := NewPostgresTransactionPasswordRepository(pool)

	t.Run("create and find by user id", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		password, err := domain.NewTransactionPassword(userID, "hashed-pin", now)
		if err != nil {
			t.Fatalf("expected no error creating domain password, got %v", err)
		}

		if err := repo.Create(ctx, password); err != nil {
			t.Fatalf("expected no error creating transaction password, got %v", err)
		}
		if password.ID == uuid.Nil {
			t.Fatal("expected repository to populate generated transaction password id")
		}

		got, err := repo.FindByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("expected no error finding transaction password, got %v", err)
		}
		if got == nil {
			t.Fatal("expected transaction password, got nil")
		}
		if got.ID != password.ID {
			t.Fatalf("expected id %q, got %q", password.ID, got.ID)
		}
		if got.UserID != userID {
			t.Fatalf("expected user id %q, got %q", userID, got.UserID)
		}
		if got.PasswordHash != "hashed-pin" {
			t.Fatalf("expected password hash %q, got %q", "hashed-pin", got.PasswordHash)
		}
		if got.Status != domain.TransactionPasswordActive {
			t.Fatalf("expected status %q, got %q", domain.TransactionPasswordActive, got.Status)
		}
		if got.FailedAttempts != 0 {
			t.Fatalf("expected failed attempts 0, got %d", got.FailedAttempts)
		}
		if got.LockedUntil != nil {
			t.Fatalf("expected nil locked_until, got %v", got.LockedUntil)
		}
	})

	t.Run("find missing returns nil", func(t *testing.T) {
		got, err := repo.FindByUserID(ctx, uuid.New())
		if err != nil {
			t.Fatalf("expected no error finding missing transaction password, got %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil transaction password, got %+v", got)
		}
	})

	t.Run("duplicate user returns domain error", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		first, err := domain.NewTransactionPassword(userID, "first-hash", now)
		if err != nil {
			t.Fatalf("expected no error creating first domain password, got %v", err)
		}
		second, err := domain.NewTransactionPassword(userID, "second-hash", now)
		if err != nil {
			t.Fatalf("expected no error creating second domain password, got %v", err)
		}

		if err := repo.Create(ctx, first); err != nil {
			t.Fatalf("expected first create to succeed, got %v", err)
		}

		err = repo.Create(ctx, second)
		if !errors.Is(err, domain.ErrTransactionPasswordAlreadySet) {
			t.Fatalf("expected ErrTransactionPasswordAlreadySet, got %v", err)
		}
	})

	t.Run("save validation state", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		password, err := domain.NewTransactionPassword(userID, "hashed-pin", now)
		if err != nil {
			t.Fatalf("expected no error creating domain password, got %v", err)
		}
		if err := repo.Create(ctx, password); err != nil {
			t.Fatalf("expected create to succeed, got %v", err)
		}

		lockTime := now.Add(time.Minute)
		_ = password.RegisterFailure(lockTime)
		_ = password.RegisterFailure(lockTime.Add(time.Minute))
		err = password.RegisterFailure(lockTime.Add(2 * time.Minute))
		if !errors.Is(err, domain.ErrTransactionPasswordLocked) {
			t.Fatalf("expected lock after third failure, got %v", err)
		}

		if err := repo.SaveValidationState(ctx, password); err != nil {
			t.Fatalf("expected save validation state to succeed, got %v", err)
		}

		got, err := repo.FindByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("expected find after save to succeed, got %v", err)
		}
		if got == nil {
			t.Fatal("expected transaction password after save, got nil")
		}
		if got.Status != domain.TransactionPasswordBlocked {
			t.Fatalf("expected blocked status, got %q", got.Status)
		}
		if got.FailedAttempts != domain.TransactionPasswordMaxFailures {
			t.Fatalf("expected %d failed attempts, got %d", domain.TransactionPasswordMaxFailures, got.FailedAttempts)
		}
		if got.LockedUntil == nil {
			t.Fatal("expected locked_until to be set")
		}
	})

	t.Run("save validation state missing returns domain error", func(t *testing.T) {
		now := time.Now().UTC().Truncate(time.Microsecond)
		password := &domain.TransactionPassword{
			ID:             uuid.New(),
			UserID:         uuid.New(),
			PasswordHash:   "hash",
			Status:         domain.TransactionPasswordActive,
			FailedAttempts: 0,
			UpdatedAt:      now,
		}

		err := repo.SaveValidationState(ctx, password)
		if !errors.Is(err, domain.ErrTransactionPasswordNotSet) {
			t.Fatalf("expected ErrTransactionPasswordNotSet, got %v", err)
		}
	})

	t.Run("update password hash resets validation state", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		password, err := domain.NewTransactionPassword(userID, "old-hash", now)
		if err != nil {
			t.Fatalf("expected no error creating domain password, got %v", err)
		}
		if err := repo.Create(ctx, password); err != nil {
			t.Fatalf("expected create to succeed, got %v", err)
		}

		lockTime := now.Add(time.Minute)
		_ = password.RegisterFailure(lockTime)
		_ = password.RegisterFailure(lockTime.Add(time.Minute))
		_ = password.RegisterFailure(lockTime.Add(2 * time.Minute))
		if err := repo.SaveValidationState(ctx, password); err != nil {
			t.Fatalf("expected save validation state to succeed, got %v", err)
		}

		changedAt := now.Add(time.Hour)
		if err := repo.UpdatePasswordHash(ctx, password.ID, "new-hash", changedAt, changedAt); err != nil {
			t.Fatalf("expected update password hash to succeed, got %v", err)
		}

		got, err := repo.FindByUserID(ctx, userID)
		if err != nil {
			t.Fatalf("expected find after update hash to succeed, got %v", err)
		}
		if got == nil {
			t.Fatal("expected transaction password after update hash, got nil")
		}
		if got.PasswordHash != "new-hash" {
			t.Fatalf("expected new hash, got %q", got.PasswordHash)
		}
		if got.Status != domain.TransactionPasswordActive {
			t.Fatalf("expected active status, got %q", got.Status)
		}
		if got.FailedAttempts != 0 {
			t.Fatalf("expected failed attempts reset, got %d", got.FailedAttempts)
		}
		if got.LockedUntil != nil {
			t.Fatalf("expected locked_until reset, got %v", got.LockedUntil)
		}
		if got.ChangedAt == nil || !got.ChangedAt.Equal(changedAt) {
			t.Fatalf("expected changed_at %v, got %v", changedAt, got.ChangedAt)
		}
	})
}

func newSecurityTestPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
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

func ensureSecurityRepoTestSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	statements := []string{
		`CREATE EXTENSION IF NOT EXISTS pgcrypto`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email VARCHAR(120) NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role VARCHAR(20) NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS transaction_passwords (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			password_hash TEXT NOT NULL,
			status VARCHAR(20) NOT NULL DEFAULT 'active',
			failed_attempts INTEGER NOT NULL DEFAULT 0,
			locked_until TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			changed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			CONSTRAINT chk_transaction_passwords_status CHECK (status IN ('active', 'blocked')),
			CONSTRAINT chk_transaction_passwords_failed_attempts CHECK (failed_attempts >= 0),
			CONSTRAINT chk_transaction_passwords_blocked_locked_until CHECK (status <> 'blocked' OR locked_until IS NOT NULL)
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS ux_transaction_passwords_user_id
			ON transaction_passwords(user_id)`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("failed to ensure security repository test schema: %v", err)
		}
	}

	if _, err := pool.Exec(ctx, `ALTER TABLE users ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active'`); err != nil {
		t.Fatalf("failed to ensure users.status column: %v", err)
	}
	if _, err := pool.Exec(ctx, `ALTER TABLE transaction_passwords ALTER COLUMN id SET DEFAULT gen_random_uuid()`); err != nil {
		t.Fatalf("failed to ensure transaction_passwords.id default: %v", err)
	}
}

func createSecurityTestUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	var userID uuid.UUID
	email := strings.ToLower(uuid.NewString()) + "@example.com"
	now := time.Now().UTC()

	err := pool.QueryRow(ctx, `
		INSERT INTO users (
			email,
			password_hash,
			role,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, email, "hashed-password", "admin", "active", now, now).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create security test user: %v", err)
	}

	return userID
}

func cleanupSecurityTestUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) {
	t.Helper()

	if _, err := pool.Exec(ctx, `DELETE FROM transaction_passwords WHERE user_id = $1`, userID); err != nil {
		t.Logf("cleanup warning: failed to delete transaction password for user %q: %v", userID, err)
	}

	if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
		t.Logf("cleanup warning: failed to delete user %q: %v", userID, err)
	}
}

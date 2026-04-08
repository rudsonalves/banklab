package infrastructure

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestPostgresUserRepository_Integration(t *testing.T) {
	ctx := context.Background()
	pool := newAuthTestPool(t, ctx)
	defer pool.Close()

	ensureAuthRepoTestSchema(t, ctx, pool)

	repo := NewPostgresUserRepository(pool)

	t.Run("create and query user with null customer id", func(t *testing.T) {
		user := testUser(time.Now().UTC())
		defer cleanupUserByID(t, ctx, pool, user.ID)

		err := repo.Create(ctx, user)
		if err != nil {
			t.Fatalf("expected no error creating user, got %v", err)
		}

		gotByEmail, err := repo.FindByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("expected no error finding by email, got %v", err)
		}
		if gotByEmail == nil {
			t.Fatal("expected user from FindByEmail, got nil")
		}
		if gotByEmail.Email != user.Email {
			t.Fatalf("expected email %q, got %q", user.Email, gotByEmail.Email)
		}
		if gotByEmail.CustomerID != nil {
			t.Fatalf("expected nil customer id, got %v", *gotByEmail.CustomerID)
		}

		gotByID, err := repo.FindByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("expected no error finding by id, got %v", err)
		}
		if gotByID == nil {
			t.Fatal("expected user from FindByID, got nil")
		}
		if gotByID.ID != user.ID {
			t.Fatalf("expected id %q, got %q", user.ID, gotByID.ID)
		}

		exists, err := repo.ExistsByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("expected no error checking exists by email, got %v", err)
		}
		if !exists {
			t.Fatal("expected ExistsByEmail to return true")
		}
	})

	t.Run("duplicate email returns error", func(t *testing.T) {
		now := time.Now().UTC()
		email := strings.ToLower(uuid.NewString()) + "@example.com"

		userA := &domain.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: "hash-a",
			Role:         domain.RoleCustomer,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		defer cleanupUserByID(t, ctx, pool, userA.ID)

		userB := &domain.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: "hash-b",
			Role:         domain.RoleCustomer,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		defer cleanupUserByID(t, ctx, pool, userB.ID)

		if err := repo.Create(ctx, userA); err != nil {
			t.Fatalf("expected first create to succeed, got %v", err)
		}

		err := repo.Create(ctx, userB)
		if err == nil {
			t.Fatal("expected duplicate email error, got nil")
		}
	})

	t.Run("not found returns nil", func(t *testing.T) {
		nonExistingEmail := strings.ToLower(uuid.NewString()) + "@example.com"
		nonExistingID := uuid.New()

		foundByEmail, err := repo.FindByEmail(ctx, nonExistingEmail)
		if err != nil {
			t.Fatalf("expected no error for missing email, got %v", err)
		}
		if foundByEmail != nil {
			t.Fatalf("expected nil user for missing email, got %#v", foundByEmail)
		}

		foundByID, err := repo.FindByID(ctx, nonExistingID)
		if err != nil {
			t.Fatalf("expected no error for missing id, got %v", err)
		}
		if foundByID != nil {
			t.Fatalf("expected nil user for missing id, got %#v", foundByID)
		}

		exists, err := repo.ExistsByEmail(ctx, nonExistingEmail)
		if err != nil {
			t.Fatalf("expected no error for missing email exists check, got %v", err)
		}
		if exists {
			t.Fatal("expected ExistsByEmail to return false")
		}
	})
}

func newAuthTestPool(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	connString := os.Getenv("BANK_TEST_DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:postgres@localhost:5432/bank?sslmode=disable"
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

func ensureAuthRepoTestSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	statements := []string{
		`CREATE TABLE IF NOT EXISTS customers (
			id UUID PRIMARY KEY,
			name VARCHAR(120) NOT NULL,
			cpf VARCHAR(11) NOT NULL UNIQUE,
			email VARCHAR(120) NOT NULL UNIQUE,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			email VARCHAR(120) NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			role VARCHAR(20) NOT NULL,
			customer_id UUID UNIQUE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_users_customer_id FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
		)`,
	}

	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			t.Fatalf("failed to ensure auth repository test schema: %v", err)
		}
	}
}

func cleanupUserByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID uuid.UUID) {
	t.Helper()

	if _, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID); err != nil {
		t.Logf("cleanup warning: failed to delete user %q: %v", userID, err)
	}
}

func testUser(now time.Time) *domain.User {
	return &domain.User{
		ID:           uuid.New(),
		Email:        strings.ToLower(uuid.NewString()) + "@example.com",
		PasswordHash: "hashed-password",
		Role:         domain.RoleCustomer,
		CustomerID:   nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

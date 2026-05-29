package infrastructure

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

func TestPostgresStepUpTokenRepository_Integration(t *testing.T) {
	ctx := context.Background()
	pool := newSecurityTestPool(t, ctx)
	defer pool.Close()

	ensureSecurityRepoTestSchema(t, ctx, pool)

	repo := NewPostgresStepUpTokenRepository(pool)

	t.Run("create and find by jti", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		token, err := domain.NewStepUpToken(
			uuid.NewString(),
			userID,
			domain.StepUpEndpointInternalTransferCreate,
			now,
		)
		if err != nil {
			t.Fatalf("expected no error creating domain token, got %v", err)
		}

		if err := repo.Create(ctx, token); err != nil {
			t.Fatalf("expected no error creating step-up token, got %v", err)
		}
		if token.ID == uuid.Nil {
			t.Fatal("expected repository to populate generated step-up token id")
		}

		got, err := repo.FindByJTI(ctx, token.JTI)
		if err != nil {
			t.Fatalf("expected no error finding step-up token, got %v", err)
		}
		if got == nil {
			t.Fatal("expected step-up token, got nil")
		}
		if got.ID != token.ID {
			t.Fatalf("expected id %q, got %q", token.ID, got.ID)
		}
		if got.JTI != token.JTI {
			t.Fatalf("expected jti %q, got %q", token.JTI, got.JTI)
		}
		if got.UserID != userID {
			t.Fatalf("expected user id %q, got %q", userID, got.UserID)
		}
		if got.EndpointKey != domain.StepUpEndpointInternalTransferCreate {
			t.Fatalf("expected endpoint key %q, got %q", domain.StepUpEndpointInternalTransferCreate, got.EndpointKey)
		}
		if got.Status != domain.StepUpTokenActive {
			t.Fatalf("expected active status, got %q", got.Status)
		}
		if got.ConsumedAt != nil {
			t.Fatalf("expected nil consumed_at, got %v", got.ConsumedAt)
		}
	})

	t.Run("find missing returns nil", func(t *testing.T) {
		got, err := repo.FindByJTI(ctx, uuid.NewString())
		if err != nil {
			t.Fatalf("expected no error finding missing step-up token, got %v", err)
		}
		if got != nil {
			t.Fatalf("expected nil step-up token, got %+v", got)
		}
	})

	t.Run("duplicate jti returns domain error", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		jti := uuid.NewString()
		first, err := domain.NewStepUpToken(jti, userID, domain.StepUpEndpointInternalTransferCreate, now)
		if err != nil {
			t.Fatalf("expected no error creating first domain token, got %v", err)
		}
		second, err := domain.NewStepUpToken(jti, userID, domain.StepUpEndpointInternalTransferCreate, now)
		if err != nil {
			t.Fatalf("expected no error creating second domain token, got %v", err)
		}

		if err := repo.Create(ctx, first); err != nil {
			t.Fatalf("expected first create to succeed, got %v", err)
		}

		err = repo.Create(ctx, second)
		if !errors.Is(err, domain.ErrInvalidStepUpToken) {
			t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
		}
	})

	t.Run("consume active token", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		token, err := domain.NewStepUpToken(
			uuid.NewString(),
			userID,
			domain.StepUpEndpointInternalTransferCreate,
			now,
		)
		if err != nil {
			t.Fatalf("expected no error creating domain token, got %v", err)
		}
		if err := repo.Create(ctx, token); err != nil {
			t.Fatalf("expected create to succeed, got %v", err)
		}

		consumedAt := now.Add(time.Minute)
		consumed, err := repo.ConsumeByJTI(ctx, token.JTI, consumedAt)
		if err != nil {
			t.Fatalf("expected consume to succeed, got %v", err)
		}
		if consumed.Status != domain.StepUpTokenConsumed {
			t.Fatalf("expected consumed status, got %q", consumed.Status)
		}
		if consumed.ConsumedAt == nil {
			t.Fatal("expected consumed_at to be set")
		}
		if !consumed.ConsumedAt.Equal(consumedAt) {
			t.Fatalf("expected consumed_at %v, got %v", consumedAt, *consumed.ConsumedAt)
		}
	})

	t.Run("consume same token twice returns consumed error", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		token, err := domain.NewStepUpToken(
			uuid.NewString(),
			userID,
			domain.StepUpEndpointInternalTransferCreate,
			now,
		)
		if err != nil {
			t.Fatalf("expected no error creating domain token, got %v", err)
		}
		if err := repo.Create(ctx, token); err != nil {
			t.Fatalf("expected create to succeed, got %v", err)
		}
		if _, err := repo.ConsumeByJTI(ctx, token.JTI, now.Add(time.Minute)); err != nil {
			t.Fatalf("expected first consume to succeed, got %v", err)
		}

		_, err = repo.ConsumeByJTI(ctx, token.JTI, now.Add(90*time.Second))
		if !errors.Is(err, domain.ErrStepUpTokenConsumed) {
			t.Fatalf("expected ErrStepUpTokenConsumed, got %v", err)
		}
	})

	t.Run("consume expired token returns expired error without changing status", func(t *testing.T) {
		userID := createSecurityTestUser(t, ctx, pool)
		defer cleanupSecurityTestUser(t, ctx, pool, userID)

		now := time.Now().UTC().Truncate(time.Microsecond)
		token, err := domain.RestoreStepUpToken(
			uuid.Nil,
			uuid.NewString(),
			userID,
			domain.StepUpEndpointInternalTransferCreate,
			domain.StepUpTokenActive,
			now.Add(time.Minute),
			nil,
			now,
		)
		if err != nil {
			t.Fatalf("expected no error creating domain token, got %v", err)
		}
		if err := repo.Create(ctx, token); err != nil {
			t.Fatalf("expected create to succeed, got %v", err)
		}

		_, err = repo.ConsumeByJTI(ctx, token.JTI, now.Add(2*time.Minute))
		if !errors.Is(err, domain.ErrStepUpTokenExpired) {
			t.Fatalf("expected ErrStepUpTokenExpired, got %v", err)
		}

		got, err := repo.FindByJTI(ctx, token.JTI)
		if err != nil {
			t.Fatalf("expected find after expired consume to succeed, got %v", err)
		}
		if got == nil {
			t.Fatal("expected step-up token after expired consume, got nil")
		}
		if got.Status != domain.StepUpTokenActive {
			t.Fatalf("expected active status after expired consume, got %q", got.Status)
		}
		if got.ConsumedAt != nil {
			t.Fatalf("expected consumed_at to remain nil, got %v", got.ConsumedAt)
		}
	})

	t.Run("consume missing token returns invalid error", func(t *testing.T) {
		_, err := repo.ConsumeByJTI(ctx, uuid.NewString(), time.Now().UTC())
		if !errors.Is(err, domain.ErrInvalidStepUpToken) {
			t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
		}
	})
}

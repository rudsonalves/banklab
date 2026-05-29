package infrastructure

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	securityapp "github.com/seu-usuario/bank-api/internal/security/application"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

func TestStepUpEnforcement_Integration_ConsumesPersistedJTIOnce(t *testing.T) {
	ctx := context.Background()
	pool := newSecurityTestPool(t, ctx)
	defer pool.Close()

	ensureSecurityRepoTestSchema(t, ctx, pool)

	userID := createSecurityTestUser(t, ctx, pool)
	defer cleanupSecurityTestUser(t, ctx, pool, userID)

	secret := "step-up-enforcement-integration-secret"
	repo := NewPostgresStepUpTokenRepository(pool)
	signer := NewJWTStepUpTokenSigner(secret)
	verifier := NewJWTStepUpTokenVerifier(secret)
	enforce := securityapp.NewEnforceStepUpUseCase(verifier, repo)

	now := time.Now().UTC().Truncate(time.Second)
	stepUpToken, err := domain.NewStepUpToken(
		uuid.NewString(),
		userID,
		domain.StepUpEndpointInternalTransferCreate,
		now,
	)
	if err != nil {
		t.Fatalf("expected no error creating step-up token, got %v", err)
	}
	if err := repo.Create(ctx, stepUpToken); err != nil {
		t.Fatalf("expected no error persisting step-up token, got %v", err)
	}

	signedToken, err := signer.Sign(stepUpToken)
	if err != nil {
		t.Fatalf("expected no error signing step-up token, got %v", err)
	}

	input := securityapp.EnforceStepUpInput{
		User: &authdomain.AuthenticatedUser{
			UserID: userID,
			Role:   authdomain.RoleCustomer,
		},
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       signedToken,
		Now:         now.Add(time.Second),
	}

	if err := enforce.Execute(ctx, input); err != nil {
		t.Fatalf("expected first enforcement to succeed, got %v", err)
	}

	consumed, err := repo.FindByJTI(ctx, stepUpToken.JTI)
	if err != nil {
		t.Fatalf("expected no error finding consumed token, got %v", err)
	}
	if consumed == nil {
		t.Fatal("expected persisted step-up token, got nil")
	}
	if consumed.Status != domain.StepUpTokenConsumed {
		t.Fatalf("expected consumed status, got %q", consumed.Status)
	}
	if consumed.ConsumedAt == nil {
		t.Fatal("expected consumed_at to be set")
	}

	err = enforce.Execute(ctx, input)
	if !errors.Is(err, domain.ErrStepUpTokenConsumed) {
		t.Fatalf("expected ErrStepUpTokenConsumed on retry with same token, got %v", err)
	}
}

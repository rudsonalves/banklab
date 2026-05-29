package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type stepUpTokenVerifierMock struct {
	verifyCalls int
	token       string
	claims      *domain.VerifiedStepUpTokenClaims
	err         error
	events      *[]string
}

func (m *stepUpTokenVerifierMock) Verify(token string) (*domain.VerifiedStepUpTokenClaims, error) {
	m.verifyCalls++
	m.token = token
	if m.events != nil {
		*m.events = append(*m.events, "verify-step-up-token")
	}
	return m.claims, m.err
}

func TestEnforceStepUpUseCase_Execute_Success(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	claims := verifiedStepUpClaims(t, userID, domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now)
	consumedToken := consumedStepUpToken(t, userID, domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now)
	events := []string{}

	verifier := &stepUpTokenVerifierMock{claims: claims, events: &events}
	repo := &stepUpTokenRepositoryMock{consume: consumedToken, events: &events}
	uc := NewEnforceStepUpUseCase(verifier, repo)

	err := uc.Execute(context.Background(), EnforceStepUpInput{
		User:        authenticatedUser(userID),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       "signed-step-up-token",
		Now:         now,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if verifier.verifyCalls != 1 {
		t.Fatalf("expected verifier once, got %d", verifier.verifyCalls)
	}
	if verifier.token != "signed-step-up-token" {
		t.Fatalf("expected verifier token %q, got %q", "signed-step-up-token", verifier.token)
	}
	if repo.consumeCalls != 1 {
		t.Fatalf("expected ConsumeByJTI once, got %d", repo.consumeCalls)
	}
	if repo.consumeJTI != "step-up-jti" {
		t.Fatalf("expected consumed jti %q, got %q", "step-up-jti", repo.consumeJTI)
	}
	if !repo.consumeNow.Equal(now) {
		t.Fatalf("expected consume time %v, got %v", now, repo.consumeNow)
	}
	if got := events; len(got) != 2 || got[0] != "verify-step-up-token" || got[1] != "consume-step-up-token" {
		t.Fatalf("expected verify before consume, got events %v", got)
	}
}

func TestEnforceStepUpUseCase_Execute_TokenRequired(t *testing.T) {
	verifier := &stepUpTokenVerifierMock{}
	repo := &stepUpTokenRepositoryMock{}
	uc := NewEnforceStepUpUseCase(verifier, repo)

	err := uc.Execute(context.Background(), EnforceStepUpInput{
		User:        authenticatedUser(uuid.New()),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       " ",
		Now:         time.Now().UTC(),
	})

	if !errors.Is(err, domain.ErrStepUpTokenRequired) {
		t.Fatalf("expected ErrStepUpTokenRequired, got %v", err)
	}
	if verifier.verifyCalls != 0 {
		t.Fatalf("expected verifier not to be called, got %d", verifier.verifyCalls)
	}
	if repo.consumeCalls != 0 {
		t.Fatalf("expected ConsumeByJTI not to be called, got %d", repo.consumeCalls)
	}
}

func TestEnforceStepUpUseCase_Execute_InvalidAuthenticatedUser(t *testing.T) {
	uc := NewEnforceStepUpUseCase(&stepUpTokenVerifierMock{}, &stepUpTokenRepositoryMock{})

	err := uc.Execute(context.Background(), EnforceStepUpInput{
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       "signed-step-up-token",
		Now:         time.Now().UTC(),
	})

	if !errors.Is(err, authdomain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestEnforceStepUpUseCase_Execute_VerifierErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "invalid token", err: domain.ErrInvalidStepUpToken},
		{name: "expired token", err: domain.ErrStepUpTokenExpired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verifier := &stepUpTokenVerifierMock{err: tt.err}
			repo := &stepUpTokenRepositoryMock{}
			uc := NewEnforceStepUpUseCase(verifier, repo)

			err := uc.Execute(context.Background(), EnforceStepUpInput{
				User:        authenticatedUser(uuid.New()),
				EndpointKey: domain.StepUpEndpointInternalTransferCreate,
				Token:       "signed-step-up-token",
				Now:         time.Now().UTC(),
			})

			if !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
			if repo.consumeCalls != 0 {
				t.Fatalf("expected ConsumeByJTI not to be called, got %d", repo.consumeCalls)
			}
		})
	}
}

func TestEnforceStepUpUseCase_Execute_UserMismatch(t *testing.T) {
	now := time.Now().UTC()
	claims := verifiedStepUpClaims(t, uuid.New(), domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now)
	repo := &stepUpTokenRepositoryMock{}
	uc := NewEnforceStepUpUseCase(&stepUpTokenVerifierMock{claims: claims}, repo)

	err := uc.Execute(context.Background(), EnforceStepUpInput{
		User:        authenticatedUser(uuid.New()),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       "signed-step-up-token",
		Now:         now,
	})

	if !errors.Is(err, domain.ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
	}
	if repo.consumeCalls != 0 {
		t.Fatalf("expected ConsumeByJTI not to be called, got %d", repo.consumeCalls)
	}
}

func TestEnforceStepUpUseCase_Execute_EndpointMismatch(t *testing.T) {
	userID := uuid.New()
	now := time.Now().UTC()
	claims := verifiedStepUpClaims(t, userID, "pix.create", "step-up-jti", now)
	repo := &stepUpTokenRepositoryMock{}
	uc := NewEnforceStepUpUseCase(&stepUpTokenVerifierMock{claims: claims}, repo)

	err := uc.Execute(context.Background(), EnforceStepUpInput{
		User:        authenticatedUser(userID),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       "signed-step-up-token",
		Now:         now,
	})

	if !errors.Is(err, domain.ErrStepUpEndpointMismatch) {
		t.Fatalf("expected ErrStepUpEndpointMismatch, got %v", err)
	}
	if repo.consumeCalls != 0 {
		t.Fatalf("expected ConsumeByJTI not to be called, got %d", repo.consumeCalls)
	}
}

func TestEnforceStepUpUseCase_Execute_ConsumeErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "missing jti", err: domain.ErrInvalidStepUpToken},
		{name: "consumed jti", err: domain.ErrStepUpTokenConsumed},
		{name: "expired persisted token", err: domain.ErrStepUpTokenExpired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := uuid.New()
			now := time.Now().UTC()
			claims := verifiedStepUpClaims(t, userID, domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now)
			repo := &stepUpTokenRepositoryMock{consumeErr: tt.err}
			uc := NewEnforceStepUpUseCase(&stepUpTokenVerifierMock{claims: claims}, repo)

			err := uc.Execute(context.Background(), EnforceStepUpInput{
				User:        authenticatedUser(userID),
				EndpointKey: domain.StepUpEndpointInternalTransferCreate,
				Token:       "signed-step-up-token",
				Now:         now,
			})

			if !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
			if repo.consumeCalls != 1 {
				t.Fatalf("expected ConsumeByJTI once, got %d", repo.consumeCalls)
			}
		})
	}
}

func TestEnforceStepUpUseCase_Execute_PersistedRecordDivergence(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name          string
		consumedToken *domain.StepUpToken
		wantErr       error
	}{
		{
			name:          "user mismatch",
			consumedToken: consumedStepUpToken(t, uuid.New(), domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now),
			wantErr:       domain.ErrInvalidStepUpToken,
		},
		{
			name:          "endpoint mismatch",
			consumedToken: consumedStepUpToken(t, uuid.New(), "pix.create", "step-up-jti", now),
			wantErr:       domain.ErrInvalidStepUpToken,
		},
		{
			name:          "jti mismatch",
			consumedToken: consumedStepUpToken(t, uuid.New(), domain.StepUpEndpointInternalTransferCreate, "other-jti", now),
			wantErr:       domain.ErrInvalidStepUpToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := uuid.New()
			claims := verifiedStepUpClaims(t, userID, domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now)
			consumed := *tt.consumedToken
			if tt.name != "user mismatch" {
				consumed.UserID = userID
			}

			uc := NewEnforceStepUpUseCase(
				&stepUpTokenVerifierMock{claims: claims},
				&stepUpTokenRepositoryMock{consume: &consumed},
			)

			err := uc.Execute(context.Background(), EnforceStepUpInput{
				User:        authenticatedUser(userID),
				EndpointKey: domain.StepUpEndpointInternalTransferCreate,
				Token:       "signed-step-up-token",
				Now:         now,
			})

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestEnforceStepUpUseCase_Execute_ExpiredPersistedRecord(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	claims := verifiedStepUpClaims(t, userID, domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now)
	consumedAt := now.Add(-time.Minute)
	consumedToken, err := domain.RestoreStepUpToken(
		uuid.New(),
		"step-up-jti",
		userID,
		domain.StepUpEndpointInternalTransferCreate,
		domain.StepUpTokenConsumed,
		now.Add(-time.Nanosecond),
		&consumedAt,
		now.Add(-domain.StepUpTokenDefaultDuration),
	)
	if err != nil {
		t.Fatalf("expected valid consumed token, got %v", err)
	}

	uc := NewEnforceStepUpUseCase(
		&stepUpTokenVerifierMock{claims: claims},
		&stepUpTokenRepositoryMock{consume: consumedToken},
	)

	err = uc.Execute(context.Background(), EnforceStepUpInput{
		User:        authenticatedUser(userID),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       "signed-step-up-token",
		Now:         now,
	})

	if !errors.Is(err, domain.ErrStepUpTokenExpired) {
		t.Fatalf("expected ErrStepUpTokenExpired, got %v", err)
	}
}

func TestEnforceStepUpUseCase_Execute_ExpiredJWTClaim(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	claims := verifiedStepUpClaims(t, userID, domain.StepUpEndpointInternalTransferCreate, "step-up-jti", now.Add(-2*time.Minute))
	repo := &stepUpTokenRepositoryMock{}
	uc := NewEnforceStepUpUseCase(&stepUpTokenVerifierMock{claims: claims}, repo)

	err := uc.Execute(context.Background(), EnforceStepUpInput{
		User:        authenticatedUser(userID),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Token:       "signed-step-up-token",
		Now:         now,
	})

	if !errors.Is(err, domain.ErrStepUpTokenExpired) {
		t.Fatalf("expected ErrStepUpTokenExpired, got %v", err)
	}
	if repo.consumeCalls != 0 {
		t.Fatalf("expected ConsumeByJTI not to be called, got %d", repo.consumeCalls)
	}
}

func verifiedStepUpClaims(
	t *testing.T,
	userID uuid.UUID,
	endpointKey string,
	jti string,
	issuedAt time.Time,
) *domain.VerifiedStepUpTokenClaims {
	t.Helper()

	claims, err := domain.NewVerifiedStepUpTokenClaims(
		userID,
		endpointKey,
		domain.StepUpTokenScope,
		jti,
		issuedAt.Add(domain.StepUpTokenDefaultDuration),
		issuedAt,
	)
	if err != nil {
		t.Fatalf("expected valid verified claims, got %v", err)
	}

	return claims
}

func consumedStepUpToken(
	t *testing.T,
	userID uuid.UUID,
	endpointKey string,
	jti string,
	now time.Time,
) *domain.StepUpToken {
	t.Helper()

	consumedAt := now
	token, err := domain.RestoreStepUpToken(
		uuid.New(),
		jti,
		userID,
		endpointKey,
		domain.StepUpTokenConsumed,
		now.Add(domain.StepUpTokenDefaultDuration),
		&consumedAt,
		now,
	)
	if err != nil {
		t.Fatalf("expected valid consumed token, got %v", err)
	}

	return token
}

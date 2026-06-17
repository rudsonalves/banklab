package infrastructure

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/installation/domain"
)

type restrictedAuthorizationRepositoryMock struct {
	findByJTIJTI string
	findByJTIOut *domain.RestrictedAuthorization
	findByJTIErr error
}

func (m *restrictedAuthorizationRepositoryMock) Create(ctx context.Context, authorization *domain.RestrictedAuthorization) error {
	return nil
}

func (m *restrictedAuthorizationRepositoryMock) FindByJTI(ctx context.Context, jti string) (*domain.RestrictedAuthorization, error) {
	m.findByJTIJTI = jti
	if m.findByJTIErr != nil {
		return nil, m.findByJTIErr
	}
	return m.findByJTIOut, nil
}

func (m *restrictedAuthorizationRepositoryMock) ConsumeByJTI(ctx context.Context, jti string, now time.Time) (*domain.RestrictedAuthorization, error) {
	return nil, nil
}

func (m *restrictedAuthorizationRepositoryMock) RevokeByJTI(ctx context.Context, jti string) error {
	return nil
}

func (m *restrictedAuthorizationRepositoryMock) RevokeActiveByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID domain.InstallationID,
	scope string,
) error {
	return nil
}

func TestJWTRestrictedAccessTokenService_SignAndVerify(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	userID := uuid.New()
	installationID := mustRestrictedTokenInstallationID(t)
	authorization, err := domain.NewRestrictedAuthorization("jti-1", userID, installationID, now)
	if err != nil {
		t.Fatalf("expected authorization, got %v", err)
	}
	repo := &restrictedAuthorizationRepositoryMock{findByJTIOut: authorization}
	service := NewJWTRestrictedAccessTokenService("test-secret", repo)
	service.now = func() time.Time { return now.Add(time.Minute) }

	token, err := service.SignRestrictedAccessToken(&domain.RestrictedAccessTokenClaims{
		UserID:         userID,
		InstallationID: installationID,
		JTI:            authorization.JTI,
		TokenType:      domain.RestrictedAccessTokenType,
		Scope:          domain.RestrictedAuthorizationScopeInstallationRegister,
		IssuedAt:       authorization.CreatedAt,
		ExpiresAt:      authorization.ExpiresAt,
	})
	if err != nil {
		t.Fatalf("expected sign to succeed, got %v", err)
	}

	claims, err := service.VerifyRestrictedAccessToken(context.Background(), token)
	if err != nil {
		t.Fatalf("expected verify to succeed, got %v", err)
	}
	if claims.UserID != userID {
		t.Fatalf("expected user id %q, got %q", userID, claims.UserID)
	}
	if claims.InstallationID.UUID() != installationID.UUID() {
		t.Fatalf("expected installation id %q, got %q", installationID, claims.InstallationID)
	}
	if repo.findByJTIJTI != authorization.JTI {
		t.Fatalf("expected repository lookup by jti %q, got %q", authorization.JTI, repo.findByJTIJTI)
	}
}

func TestJWTRestrictedAccessTokenService_VerifyRejectsMismatchedAuthorization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	userID := uuid.New()
	installationID := mustRestrictedTokenInstallationID(t)
	authorization, err := domain.NewRestrictedAuthorization("jti-1", userID, installationID, now)
	if err != nil {
		t.Fatalf("expected authorization, got %v", err)
	}
	otherInstallationID := mustRestrictedTokenInstallationID(t)
	otherAuthorization, err := domain.NewRestrictedAuthorization("jti-1", userID, otherInstallationID, now)
	if err != nil {
		t.Fatalf("expected authorization, got %v", err)
	}
	repo := &restrictedAuthorizationRepositoryMock{findByJTIOut: otherAuthorization}
	service := NewJWTRestrictedAccessTokenService("test-secret", repo)
	service.now = func() time.Time { return now.Add(time.Minute) }

	token, err := service.SignRestrictedAccessToken(&domain.RestrictedAccessTokenClaims{
		UserID:         userID,
		InstallationID: installationID,
		JTI:            authorization.JTI,
		TokenType:      domain.RestrictedAccessTokenType,
		Scope:          domain.RestrictedAuthorizationScopeInstallationRegister,
		IssuedAt:       authorization.CreatedAt,
		ExpiresAt:      authorization.ExpiresAt,
	})
	if err != nil {
		t.Fatalf("expected sign to succeed, got %v", err)
	}

	_, err = service.VerifyRestrictedAccessToken(context.Background(), token)
	if !errors.Is(err, domain.ErrRestrictedAuthorizationInvalid) {
		t.Fatalf("expected ErrRestrictedAuthorizationInvalid, got %v", err)
	}
}

func TestJWTRestrictedAccessTokenService_VerifyRejectsConsumedAuthorization(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	userID := uuid.New()
	installationID := mustRestrictedTokenInstallationID(t)
	authorization, err := domain.NewRestrictedAuthorization("jti-1", userID, installationID, now)
	if err != nil {
		t.Fatalf("expected authorization, got %v", err)
	}
	if err := authorization.Consume(now.Add(time.Minute)); err != nil {
		t.Fatalf("expected consume, got %v", err)
	}
	repo := &restrictedAuthorizationRepositoryMock{findByJTIOut: authorization}
	service := NewJWTRestrictedAccessTokenService("test-secret", repo)

	token, err := service.SignRestrictedAccessToken(&domain.RestrictedAccessTokenClaims{
		UserID:         userID,
		InstallationID: installationID,
		JTI:            authorization.JTI,
		TokenType:      domain.RestrictedAccessTokenType,
		Scope:          domain.RestrictedAuthorizationScopeInstallationRegister,
		IssuedAt:       authorization.CreatedAt,
		ExpiresAt:      authorization.ExpiresAt,
	})
	if err != nil {
		t.Fatalf("expected sign to succeed, got %v", err)
	}

	_, err = service.VerifyRestrictedAccessToken(context.Background(), token)
	if !errors.Is(err, domain.ErrRestrictedAuthorizationConsumed) {
		t.Fatalf("expected ErrRestrictedAuthorizationConsumed, got %v", err)
	}
}

func mustRestrictedTokenInstallationID(t *testing.T) domain.InstallationID {
	t.Helper()

	installationID, err := domain.NewInstallationID(uuid.New())
	if err != nil {
		t.Fatalf("expected installation id, got %v", err)
	}
	return installationID
}

package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
	securityapplication "github.com/seu-usuario/bank-api/internal/security/application"
	securitydomain "github.com/seu-usuario/bank-api/internal/security/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type userReaderMock struct {
	findByIDCalls int
	findByIDValue uuid.UUID
	user          *authdomain.User
	err           error
}

func (m *userReaderMock) FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	m.findByIDCalls++
	m.findByIDValue = id
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

type installationRepositoryMock struct {
	findByResourceCalls int
	findByResourceUser  uuid.UUID
	findByResourceID    installationdomain.InstallationResourceID
	findByResourceOut   *installationdomain.Installation
	findByResourceErr   error

	listCalls int
	listUser  uuid.UUID
	listOut   []*installationdomain.Installation
	listErr   error

	reserveCalls          int
	reserveUser           uuid.UUID
	reserveResourceID     installationdomain.InstallationResourceID
	reserveInstallationID installationdomain.InstallationID
	reserveMax            int
	reserveNow            time.Time
	reserveOut            *installationdomain.Installation
	reserveErr            error

	revokeCalls      int
	revokeUser       uuid.UUID
	revokeResourceID installationdomain.InstallationResourceID
	revokeNow        time.Time
	revokeOut        *installationdomain.Installation
	revokeErr        error
}

func (m *installationRepositoryMock) FindByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID installationdomain.InstallationID,
) (*installationdomain.Installation, error) {
	return nil, installationdomain.ErrInstallationNotFound
}

func (m *installationRepositoryMock) FindByResourceID(
	ctx context.Context,
	userID uuid.UUID,
	resourceID installationdomain.InstallationResourceID,
) (*installationdomain.Installation, error) {
	m.findByResourceCalls++
	m.findByResourceUser = userID
	m.findByResourceID = resourceID
	if m.findByResourceErr != nil {
		return nil, m.findByResourceErr
	}
	return m.findByResourceOut, nil
}

func (m *installationRepositoryMock) CountKnownByUserID(ctx context.Context, userID uuid.UUID) (int, error) {
	return 0, nil
}

func (m *installationRepositoryMock) HasAnyByUserID(ctx context.Context, userID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *installationRepositoryMock) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]*installationdomain.Installation, error) {
	m.listCalls++
	m.listUser = userID
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.listOut, nil
}

func (m *installationRepositoryMock) BootstrapFirstInstallation(
	ctx context.Context,
	userID uuid.UUID,
	resourceID installationdomain.InstallationResourceID,
	installationID installationdomain.InstallationID,
	now time.Time,
) (*installationdomain.Installation, error) {
	return nil, nil
}

func (m *installationRepositoryMock) ReserveKnownInstallation(
	ctx context.Context,
	userID uuid.UUID,
	resourceID installationdomain.InstallationResourceID,
	installationID installationdomain.InstallationID,
	maxKnownInstallations int,
	now time.Time,
) (*installationdomain.Installation, error) {
	m.reserveCalls++
	m.reserveUser = userID
	m.reserveResourceID = resourceID
	m.reserveInstallationID = installationID
	m.reserveMax = maxKnownInstallations
	m.reserveNow = now
	if m.reserveErr != nil {
		return nil, m.reserveErr
	}
	return m.reserveOut, nil
}

func (m *installationRepositoryMock) RevokeByResourceID(
	ctx context.Context,
	userID uuid.UUID,
	resourceID installationdomain.InstallationResourceID,
	now time.Time,
) (*installationdomain.Installation, error) {
	m.revokeCalls++
	m.revokeUser = userID
	m.revokeResourceID = resourceID
	m.revokeNow = now
	if m.revokeErr != nil {
		return nil, m.revokeErr
	}
	return m.revokeOut, nil
}

type restrictedAuthorizationRepositoryMock struct {
	consumeCalls int
	consumeJTI   string
	consumeNow   time.Time
	consumeOut   *installationdomain.RestrictedAuthorization
	consumeErr   error
}

func (m *restrictedAuthorizationRepositoryMock) Create(
	ctx context.Context,
	authorization *installationdomain.RestrictedAuthorization,
) error {
	return nil
}

func (m *restrictedAuthorizationRepositoryMock) FindByJTI(
	ctx context.Context,
	jti string,
) (*installationdomain.RestrictedAuthorization, error) {
	return nil, nil
}

func (m *restrictedAuthorizationRepositoryMock) ConsumeByJTI(
	ctx context.Context,
	jti string,
	now time.Time,
) (*installationdomain.RestrictedAuthorization, error) {
	m.consumeCalls++
	m.consumeJTI = jti
	m.consumeNow = now
	if m.consumeErr != nil {
		return nil, m.consumeErr
	}
	return m.consumeOut, nil
}

func (m *restrictedAuthorizationRepositoryMock) RevokeByJTI(ctx context.Context, jti string) error {
	return nil
}

func (m *restrictedAuthorizationRepositoryMock) RevokeActiveByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID installationdomain.InstallationID,
	scope string,
) error {
	return nil
}

type tokenServiceMock struct {
	accessCalls  int
	accessClaims authdomain.TokenClaims
	accessToken  string
	accessErr    error

	refreshCalls int
	refreshUser  uuid.UUID
	refreshToken string
	refreshErr   error
}

func (m *tokenServiceMock) GenerateAccessToken(claims authdomain.TokenClaims) (string, error) {
	m.accessCalls++
	m.accessClaims = claims
	if m.accessErr != nil {
		return "", m.accessErr
	}
	return m.accessToken, nil
}

func (m *tokenServiceMock) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	m.refreshCalls++
	m.refreshUser = userID
	if m.refreshErr != nil {
		return "", m.refreshErr
	}
	return m.refreshToken, nil
}

func (m *tokenServiceMock) ParseAccessToken(token string) (*authdomain.TokenClaims, error) {
	return nil, nil
}

func (m *tokenServiceMock) ParseRefreshToken(token string) (uuid.UUID, error) {
	return uuid.Nil, nil
}

type sessionRepositoryMock struct {
	createCalls         int
	createInput         authdomain.CreateSessionInput
	createErr           error
	invalidateCalls     int
	invalidateUser      uuid.UUID
	invalidateInstallID uuid.UUID
	invalidateRevokedAt time.Time
}

func (m *sessionRepositoryMock) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	return m.CreateWithInstallation(ctx, authdomain.CreateSessionInput{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
}

func (m *sessionRepositoryMock) CreateWithInstallation(ctx context.Context, input authdomain.CreateSessionInput) error {
	m.createCalls++
	m.createInput = input
	return m.createErr
}

func (m *sessionRepositoryMock) FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, bool, error) {
	return uuid.Nil, time.Time{}, false, nil
}

func (m *sessionRepositoryMock) FindByTokenHashWithInstallation(
	ctx context.Context,
	tokenHash string,
) (*authdomain.SessionRecord, error) {
	return nil, nil
}

func (m *sessionRepositoryMock) Revoke(ctx context.Context, tokenHash string) error {
	return nil
}

func (m *sessionRepositoryMock) RevokeByUserIDAndInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID uuid.UUID,
	revokedAt time.Time,
) error {
	return nil
}

func (m *sessionRepositoryMock) InvalidateByInstallationID(
	ctx context.Context,
	userID uuid.UUID,
	installationID installationdomain.InstallationID,
	now time.Time,
) error {
	m.invalidateCalls++
	m.invalidateUser = userID
	m.invalidateInstallID = installationID.UUID()
	m.invalidateRevokedAt = now
	return nil
}

type transactorMock struct {
	calls int
	err   error
}

func (m *transactorMock) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	m.calls++
	if m.err != nil {
		return m.err
	}
	return fn(ctx)
}

type stepUpEnforcerMock struct {
	calls int
	input securityapplication.EnforceStepUpInput
	err   error
}

func (m *stepUpEnforcerMock) Execute(ctx context.Context, input securityapplication.EnforceStepUpInput) error {
	m.calls++
	m.input = input
	return m.err
}

func TestRegisterInstallationUseCase_Execute_Success(t *testing.T) {
	userID := uuid.New()
	customerID := uuid.New()
	installationUUID := uuid.New()
	resourceUUID := uuid.New()
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	installationID := mustInstallationID(t, installationUUID)
	resourceID := mustResourceID(t, resourceUUID)
	authorization := mustConsumedAuthorization(t, userID, installationID, now)
	registered := mustKnownInstallation(t, userID, resourceID, installationID, now)

	users := &userReaderMock{user: &authdomain.User{
		ID:         userID,
		Email:      "user@example.com",
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	}}
	installations := &installationRepositoryMock{reserveOut: registered}
	authorizations := &restrictedAuthorizationRepositoryMock{consumeOut: authorization}
	tokens := &tokenServiceMock{accessToken: "access-token", refreshToken: "refresh-token"}
	sessions := &sessionRepositoryMock{}
	tx := &transactorMock{}
	stepUp := &stepUpEnforcerMock{}

	uc := NewRegisterInstallationUseCase(users, installations, authorizations, tokens, sessions, tx, stepUp).
		WithRefreshSessionTTL(time.Hour)
	uc.installationResourceID = func() uuid.UUID { return resourceUUID }

	ctx := sharedauthctx.WithRestrictedSession(context.Background(), sharedauthctx.RestrictedSession{
		UserID:         userID,
		InstallationID: installationUUID,
		JTI:            authorization.JTI,
		Scope:          installationdomain.RestrictedAuthorizationScopeInstallationRegister,
	})
	output, err := uc.Execute(ctx, RegisterInstallationInput{
		PresentedInstallationID: installationUUID,
		StepUpToken:             "step-up-token",
		Now:                     now,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output.AccessToken != "access-token" || output.RefreshToken != "refresh-token" {
		t.Fatalf("unexpected tokens: %#v", output)
	}
	if output.InstallationResourceID != resourceUUID || output.InstallationStatus != string(installationdomain.InstallationStatusKnown) {
		t.Fatalf("unexpected installation output: %#v", output)
	}
	if stepUp.calls != 1 {
		t.Fatalf("expected step-up once, got %d", stepUp.calls)
	}
	if stepUp.input.EndpointKey != securitydomain.StepUpEndpointInstallationRegisterCreate {
		t.Fatalf("expected endpoint %q, got %q", securitydomain.StepUpEndpointInstallationRegisterCreate, stepUp.input.EndpointKey)
	}
	if stepUp.input.Token != "step-up-token" || stepUp.input.User == nil || stepUp.input.User.UserID != userID {
		t.Fatalf("unexpected step-up input: %#v", stepUp.input)
	}
	if tx.calls != 1 {
		t.Fatalf("expected transaction once, got %d", tx.calls)
	}
	if authorizations.consumeCalls != 1 || authorizations.consumeJTI != authorization.JTI {
		t.Fatalf("expected authorization consumed by jti, got calls=%d jti=%q", authorizations.consumeCalls, authorizations.consumeJTI)
	}
	if installations.reserveCalls != 1 {
		t.Fatalf("expected reserve once, got %d", installations.reserveCalls)
	}
	if installations.reserveMax != installationdomain.MaxKnownInstallations {
		t.Fatalf("expected max installations %d, got %d", installationdomain.MaxKnownInstallations, installations.reserveMax)
	}
	if sessions.createCalls != 1 {
		t.Fatalf("expected session create once, got %d", sessions.createCalls)
	}
	if sessions.createInput.InstallationID == nil || *sessions.createInput.InstallationID != installationUUID {
		t.Fatalf("expected session installation %q, got %#v", installationUUID, sessions.createInput.InstallationID)
	}
	expectedHash := sha256.Sum256([]byte("refresh-token"))
	if sessions.createInput.TokenHash != hex.EncodeToString(expectedHash[:]) {
		t.Fatalf("unexpected refresh hash %q", sessions.createInput.TokenHash)
	}
	if tokens.accessClaims.InstallationID == nil || *tokens.accessClaims.InstallationID != installationUUID {
		t.Fatalf("expected access token installation %q, got %#v", installationUUID, tokens.accessClaims.InstallationID)
	}
}

func TestRegisterInstallationUseCase_Execute_MismatchedInstallationDoesNotConsumeAuthorization(t *testing.T) {
	userID := uuid.New()
	ctxInstallationID := uuid.New()
	presentedInstallationID := uuid.New()
	authorizations := &restrictedAuthorizationRepositoryMock{}
	stepUp := &stepUpEnforcerMock{}
	uc := NewRegisterInstallationUseCase(
		&userReaderMock{},
		&installationRepositoryMock{},
		authorizations,
		&tokenServiceMock{},
		&sessionRepositoryMock{},
		&transactorMock{},
		stepUp,
	)

	ctx := sharedauthctx.WithRestrictedSession(context.Background(), sharedauthctx.RestrictedSession{
		UserID:         userID,
		InstallationID: ctxInstallationID,
		JTI:            "jti-1",
		Scope:          installationdomain.RestrictedAuthorizationScopeInstallationRegister,
	})
	output, err := uc.Execute(ctx, RegisterInstallationInput{
		PresentedInstallationID: presentedInstallationID,
		StepUpToken:             "step-up-token",
	})
	if !errors.Is(err, installationdomain.ErrInstallationMismatch) {
		t.Fatalf("expected ErrInstallationMismatch, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %#v", output)
	}
	if authorizations.consumeCalls != 0 {
		t.Fatalf("expected authorization not to be consumed, got %d", authorizations.consumeCalls)
	}
	if stepUp.calls != 0 {
		t.Fatalf("expected step-up not to be consumed, got %d", stepUp.calls)
	}
}

func TestRegisterInstallationUseCase_Execute_ConsumedAuthorizationMustMatchContext(t *testing.T) {
	userID := uuid.New()
	installationUUID := uuid.New()
	otherInstallationID := mustInstallationID(t, uuid.New())
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	authorization := mustConsumedAuthorization(t, userID, otherInstallationID, now)
	authorizations := &restrictedAuthorizationRepositoryMock{consumeOut: authorization}
	installations := &installationRepositoryMock{}
	sessions := &sessionRepositoryMock{}
	uc := NewRegisterInstallationUseCase(
		&userReaderMock{user: &authdomain.User{ID: userID, Role: authdomain.RoleCustomer}},
		installations,
		authorizations,
		&tokenServiceMock{accessToken: "access-token", refreshToken: "refresh-token"},
		sessions,
		&transactorMock{},
		&stepUpEnforcerMock{},
	)

	ctx := sharedauthctx.WithRestrictedSession(context.Background(), sharedauthctx.RestrictedSession{
		UserID:         userID,
		InstallationID: installationUUID,
		JTI:            authorization.JTI,
		Scope:          installationdomain.RestrictedAuthorizationScopeInstallationRegister,
	})
	output, err := uc.Execute(ctx, RegisterInstallationInput{
		PresentedInstallationID: installationUUID,
		StepUpToken:             "step-up-token",
		Now:                     now,
	})
	if !errors.Is(err, installationdomain.ErrRestrictedAuthorizationInvalid) {
		t.Fatalf("expected ErrRestrictedAuthorizationInvalid, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %#v", output)
	}
	if installations.reserveCalls != 0 {
		t.Fatalf("expected no reserve, got %d", installations.reserveCalls)
	}
	if sessions.createCalls != 0 {
		t.Fatalf("expected no session, got %d", sessions.createCalls)
	}
}

func mustInstallationID(t *testing.T, value uuid.UUID) installationdomain.InstallationID {
	t.Helper()
	id, err := installationdomain.NewInstallationID(value)
	if err != nil {
		t.Fatalf("expected installation id, got %v", err)
	}
	return id
}

func mustResourceID(t *testing.T, value uuid.UUID) installationdomain.InstallationResourceID {
	t.Helper()
	id, err := installationdomain.NewInstallationResourceID(value)
	if err != nil {
		t.Fatalf("expected resource id, got %v", err)
	}
	return id
}

func mustKnownInstallation(
	t *testing.T,
	userID uuid.UUID,
	resourceID installationdomain.InstallationResourceID,
	installationID installationdomain.InstallationID,
	now time.Time,
) *installationdomain.Installation {
	t.Helper()
	installation, err := installationdomain.NewKnownInstallation(userID, resourceID, installationID, now)
	if err != nil {
		t.Fatalf("expected known installation, got %v", err)
	}
	return installation
}

func mustConsumedAuthorization(
	t *testing.T,
	userID uuid.UUID,
	installationID installationdomain.InstallationID,
	now time.Time,
) *installationdomain.RestrictedAuthorization {
	t.Helper()
	authorization, err := installationdomain.NewRestrictedAuthorization("jti-1", userID, installationID, now)
	if err != nil {
		t.Fatalf("expected restricted authorization, got %v", err)
	}
	if err := authorization.Consume(now.Add(time.Second)); err != nil {
		t.Fatalf("expected consumed authorization, got %v", err)
	}
	return authorization
}

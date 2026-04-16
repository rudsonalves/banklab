package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type refreshUserRepositoryMock struct {
	findByIDCalls int
	findByIDID    uuid.UUID
	findByIDUser  *domain.User
	findByIDErr   error
}

func (m *refreshUserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *refreshUserRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status domain.UserStatus) error {
	return nil
}

func (m *refreshUserRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *refreshUserRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	m.findByIDCalls++
	m.findByIDID = id
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	return m.findByIDUser, nil
}

func (m *refreshUserRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

type refreshTokenServiceMock struct {
	parseCalls           int
	parseToken           string
	parseUserID          uuid.UUID
	parseErr             error
	generateCalls        int
	generateIn           domain.TokenClaims
	generateToken        string
	generateErr          error
	generateRefreshCalls int
	generateRefreshToken string
	generateRefreshErr   error
}

func (m *refreshTokenServiceMock) GenerateAccessToken(claims domain.TokenClaims) (string, error) {
	m.generateCalls++
	m.generateIn = claims
	if m.generateErr != nil {
		return "", m.generateErr
	}
	return m.generateToken, nil
}

func (m *refreshTokenServiceMock) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	m.generateRefreshCalls++
	if m.generateRefreshErr != nil {
		return "", m.generateRefreshErr
	}
	if m.generateRefreshToken != "" {
		return m.generateRefreshToken, nil
	}
	return "new-refresh-token", nil
}

func (m *refreshTokenServiceMock) ParseAccessToken(token string) (*domain.TokenClaims, error) {
	return nil, nil
}

func (m *refreshTokenServiceMock) ParseRefreshToken(token string) (uuid.UUID, error) {
	m.parseCalls++
	m.parseToken = token
	if m.parseErr != nil {
		return uuid.Nil, m.parseErr
	}
	return m.parseUserID, nil
}

type refreshSessionRepositoryMock struct {
	findCalls       int
	findTokenHash   string
	findUserID      uuid.UUID
	findExpiresAt   time.Time
	findRevoked     bool
	findErr         error
	revokeCalls     int
	revokeTokenHash string
	revokeErr       error
	createCalls     int
	createTokenHash string
	createErr       error
}

func (m *refreshSessionRepositoryMock) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	m.createCalls++
	m.createTokenHash = tokenHash
	return m.createErr
}

func (m *refreshSessionRepositoryMock) FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, bool, error) {
	m.findCalls++
	m.findTokenHash = tokenHash
	if m.findErr != nil {
		return uuid.Nil, time.Time{}, false, m.findErr
	}
	return m.findUserID, m.findExpiresAt, m.findRevoked, nil
}

func (m *refreshSessionRepositoryMock) Revoke(ctx context.Context, tokenHash string) error {
	m.revokeCalls++
	m.revokeTokenHash = tokenHash
	return m.revokeErr
}

type transactorMock struct{}

func (m *transactorMock) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func TestRefreshAccessTokenUseCase_Execute_Success(t *testing.T) {
	customerID := uuid.New()
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000111")
	repo := &refreshUserRepositoryMock{findByIDUser: &domain.User{
		ID:         userID,
		Email:      "user@example.com",
		Role:       domain.RoleCustomer,
		CustomerID: &customerID,
	}}
	tokens := &refreshTokenServiceMock{
		parseUserID:   userID,
		generateToken: "new-access-token",
	}
	sessions := &refreshSessionRepositoryMock{
		findUserID:    userID,
		findExpiresAt: time.Now().UTC().Add(10 * time.Minute),
	}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if out == nil || out.AccessToken != "new-access-token" {
		t.Fatalf("expected access token %q, got %#v", "new-access-token", out)
	}
	if tokens.parseCalls != 1 || tokens.parseToken != "refresh-token" {
		t.Fatalf("expected ParseRefreshToken called once with token, got calls=%d token=%q", tokens.parseCalls, tokens.parseToken)
	}

	expectedHash := sha256.Sum256([]byte("refresh-token"))
	if sessions.findCalls != 1 || sessions.findTokenHash != hex.EncodeToString(expectedHash[:]) {
		t.Fatalf("expected FindByTokenHash called with hash %q, got calls=%d hash=%q", hex.EncodeToString(expectedHash[:]), sessions.findCalls, sessions.findTokenHash)
	}

	if repo.findByIDCalls != 1 || repo.findByIDID != userID {
		t.Fatalf("expected FindByID called once with %q, got calls=%d id=%q", userID, repo.findByIDCalls, repo.findByIDID)
	}
	if tokens.generateCalls != 1 {
		t.Fatalf("expected GenerateAccessToken called once, got %d", tokens.generateCalls)
	}
	if tokens.generateIn.UserID != userID {
		t.Fatalf("expected token claims user id %q, got %q", userID, tokens.generateIn.UserID)
	}
	if tokens.generateIn.Role != domain.RoleCustomer {
		t.Fatalf("expected token claims role %q, got %q", domain.RoleCustomer, tokens.generateIn.Role)
	}
	if tokens.generateIn.CustomerID == nil || *tokens.generateIn.CustomerID != customerID {
		t.Fatalf("expected token claims customer_id %q, got %#v", customerID, tokens.generateIn.CustomerID)
	}

	// rotation: old token revoked
	if sessions.revokeCalls != 1 || sessions.revokeTokenHash != hex.EncodeToString(expectedHash[:]) {
		t.Fatalf("expected Revoke called once with old hash, got calls=%d hash=%q", sessions.revokeCalls, sessions.revokeTokenHash)
	}

	// rotation: new refresh token generated and persisted
	if tokens.generateRefreshCalls != 1 {
		t.Fatalf("expected GenerateRefreshToken called once, got %d", tokens.generateRefreshCalls)
	}
	if out.RefreshToken == "" {
		t.Fatal("expected non-empty refresh token in output")
	}
	if sessions.createCalls != 1 {
		t.Fatalf("expected Create called once for new session, got %d", sessions.createCalls)
	}
	newExpectedHash := sha256.Sum256([]byte(out.RefreshToken))
	if sessions.createTokenHash != hex.EncodeToString(newExpectedHash[:]) {
		t.Fatalf("expected Create called with new hash %q, got %q", hex.EncodeToString(newExpectedHash[:]), sessions.createTokenHash)
	}
}

func TestRefreshAccessTokenUseCase_Execute_InvalidRefreshToken(t *testing.T) {
	repo := &refreshUserRepositoryMock{}
	tokens := &refreshTokenServiceMock{parseErr: errors.New("bad token")}
	sessions := &refreshSessionRepositoryMock{}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "invalid-token"})
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidToken, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
	if sessions.findCalls != 0 {
		t.Fatalf("expected FindByTokenHash not to be called, got %d", sessions.findCalls)
	}
	if repo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d", repo.findByIDCalls)
	}
	if tokens.generateCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d", tokens.generateCalls)
	}
}

func TestRefreshAccessTokenUseCase_Execute_InvalidSession_NotFound(t *testing.T) {
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{}
	tokens := &refreshTokenServiceMock{parseUserID: userID}
	sessions := &refreshSessionRepositoryMock{}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidToken, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
	if repo.findByIDCalls != 0 {
		t.Fatalf("expected FindByID not to be called, got %d", repo.findByIDCalls)
	}
}

func TestRefreshAccessTokenUseCase_Execute_InvalidSession_Revoked(t *testing.T) {
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{}
	tokens := &refreshTokenServiceMock{parseUserID: userID}
	sessions := &refreshSessionRepositoryMock{findUserID: userID, findExpiresAt: time.Now().UTC().Add(time.Minute), findRevoked: true}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidToken, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}

func TestRefreshAccessTokenUseCase_Execute_InvalidSession_Expired(t *testing.T) {
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{}
	tokens := &refreshTokenServiceMock{parseUserID: userID}
	sessions := &refreshSessionRepositoryMock{findUserID: userID, findExpiresAt: time.Now().UTC().Add(-1 * time.Minute)}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidToken, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}

func TestRefreshAccessTokenUseCase_Execute_InvalidSession_UserMismatch(t *testing.T) {
	repo := &refreshUserRepositoryMock{}
	tokens := &refreshTokenServiceMock{parseUserID: uuid.New()}
	sessions := &refreshSessionRepositoryMock{findUserID: uuid.New(), findExpiresAt: time.Now().UTC().Add(time.Minute)}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidToken, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}

func TestRefreshAccessTokenUseCase_Execute_FindSessionFailure(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	repo := &refreshUserRepositoryMock{}
	tokens := &refreshTokenServiceMock{parseUserID: uuid.New()}
	sessions := &refreshSessionRepositoryMock{findErr: expectedErr}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}

func TestRefreshAccessTokenUseCase_Execute_UserNotFound(t *testing.T) {
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{findByIDUser: nil}
	tokens := &refreshTokenServiceMock{parseUserID: userID}
	sessions := &refreshSessionRepositoryMock{findUserID: userID, findExpiresAt: time.Now().UTC().Add(time.Minute)}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, domain.ErrUnauthorized) {
		t.Fatalf("expected error %v, got %v", domain.ErrUnauthorized, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
	if tokens.generateCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d", tokens.generateCalls)
	}
}

func TestRefreshAccessTokenUseCase_Execute_FindUserFailure(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{findByIDErr: expectedErr}
	tokens := &refreshTokenServiceMock{parseUserID: userID}
	sessions := &refreshSessionRepositoryMock{findUserID: userID, findExpiresAt: time.Now().UTC().Add(time.Minute)}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}

func TestRefreshAccessTokenUseCase_Execute_AccessTokenGenerationFailure(t *testing.T) {
	expectedErr := errors.New("token unavailable")
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{findByIDUser: &domain.User{ID: userID, Role: domain.RoleAdmin}}
	tokens := &refreshTokenServiceMock{parseUserID: userID, generateErr: expectedErr}
	sessions := &refreshSessionRepositoryMock{findUserID: userID, findExpiresAt: time.Now().UTC().Add(time.Minute)}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}

func TestRefreshAccessTokenUseCase_Execute_RevokeFailure(t *testing.T) {
	expectedErr := errors.New("revoke failed")
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{findByIDUser: &domain.User{ID: userID, Role: domain.RoleAdmin}}
	tokens := &refreshTokenServiceMock{parseUserID: userID, generateToken: "new-access-token"}
	sessions := &refreshSessionRepositoryMock{
		findUserID:    userID,
		findExpiresAt: time.Now().UTC().Add(time.Minute),
		revokeErr:     expectedErr,
	}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
	if sessions.createCalls != 0 {
		t.Fatalf("expected Create not to be called after revoke failure, got %d", sessions.createCalls)
	}
}

func TestRefreshAccessTokenUseCase_Execute_CreateNewSessionFailure(t *testing.T) {
	expectedErr := errors.New("create session failed")
	userID := uuid.New()
	repo := &refreshUserRepositoryMock{findByIDUser: &domain.User{ID: userID, Role: domain.RoleAdmin}}
	tokens := &refreshTokenServiceMock{parseUserID: userID, generateToken: "new-access-token"}
	sessions := &refreshSessionRepositoryMock{
		findUserID:    userID,
		findExpiresAt: time.Now().UTC().Add(time.Minute),
		createErr:     expectedErr,
	}
	uc := NewRefreshAccessTokenUseCase(repo, tokens, sessions, &transactorMock{})

	out, err := uc.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: "refresh-token"})
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected wrapped error %v, got %v", expectedErr, err)
	}
	if out != nil {
		t.Fatalf("expected nil output, got %#v", out)
	}
}

// statefulSessionMock is a session repository that tracks actual state across
// calls, used to verify rotation integrity end-to-end.
type statefulSessionEntry struct {
	userID    uuid.UUID
	expiresAt time.Time
	revoked   bool
}

type statefulSessionMock struct {
	store map[string]*statefulSessionEntry
}

func newStatefulSessionMock() *statefulSessionMock {
	return &statefulSessionMock{store: map[string]*statefulSessionEntry{}}
}

func (m *statefulSessionMock) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	m.store[tokenHash] = &statefulSessionEntry{userID: userID, expiresAt: expiresAt}
	return nil
}

func (m *statefulSessionMock) FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, bool, error) {
	e, ok := m.store[tokenHash]
	if !ok {
		return uuid.Nil, time.Time{}, false, nil
	}
	return e.userID, e.expiresAt, e.revoked, nil
}

func (m *statefulSessionMock) Revoke(ctx context.Context, tokenHash string) error {
	e, ok := m.store[tokenHash]
	if !ok {
		return domain.ErrSessionNotFound
	}
	e.revoked = true
	return nil
}

// TestRefreshAccessTokenUseCase_Execute_RotationIntegrity verifies Case 5:
// after a successful rotation, the old token is unusable and the new token works.
func TestRefreshAccessTokenUseCase_Execute_RotationIntegrity(t *testing.T) {
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000333")
	user := &domain.User{ID: userID, Email: "user@example.com", Role: domain.RoleCustomer}

	sessions := newStatefulSessionMock()

	// Seed the initial session for "old-refresh-token".
	oldToken := "old-refresh-token"
	oldHash := sha256.Sum256([]byte(oldToken))
	oldTokenHash := hex.EncodeToString(oldHash[:])
	sessions.store[oldTokenHash] = &statefulSessionEntry{
		userID:    userID,
		expiresAt: time.Now().UTC().Add(30 * 24 * time.Hour),
	}

	newToken := "new-refresh-token"

	// --- Call 1: use old token → must succeed and rotate.
	tokens1 := &refreshTokenServiceMock{
		parseUserID:          userID,
		generateToken:        "access-token-1",
		generateRefreshToken: newToken,
	}
	uc1 := NewRefreshAccessTokenUseCase(
		&refreshUserRepositoryMock{findByIDUser: user},
		tokens1,
		sessions,
		&transactorMock{},
	)

	out1, err1 := uc1.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: oldToken})
	if err1 != nil {
		t.Fatalf("call 1: expected no error, got %v", err1)
	}
	if out1 == nil || out1.RefreshToken != newToken {
		t.Fatalf("call 1: expected refresh token %q, got %#v", newToken, out1)
	}
	if out1.AccessToken != "access-token-1" {
		t.Fatalf("call 1: expected access token %q, got %q", "access-token-1", out1.AccessToken)
	}

	// --- Call 2: reuse old token → must fail (revoked).
	tokens2 := &refreshTokenServiceMock{parseUserID: userID}
	uc2 := NewRefreshAccessTokenUseCase(
		&refreshUserRepositoryMock{findByIDUser: user},
		tokens2,
		sessions,
		&transactorMock{},
	)

	out2, err2 := uc2.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: oldToken})
	if !errors.Is(err2, domain.ErrInvalidToken) {
		t.Fatalf("call 2 (old token reuse): expected %v, got %v", domain.ErrInvalidToken, err2)
	}
	if out2 != nil {
		t.Fatalf("call 2 (old token reuse): expected nil output, got %#v", out2)
	}

	// --- Call 3: use new token → must succeed.
	tokens3 := &refreshTokenServiceMock{
		parseUserID:   userID,
		generateToken: "access-token-3",
	}
	uc3 := NewRefreshAccessTokenUseCase(
		&refreshUserRepositoryMock{findByIDUser: user},
		tokens3,
		sessions,
		&transactorMock{},
	)

	out3, err3 := uc3.Execute(context.Background(), RefreshAccessTokenInput{RefreshToken: newToken})
	if err3 != nil {
		t.Fatalf("call 3 (new token): expected no error, got %v", err3)
	}
	if out3 == nil || out3.AccessToken != "access-token-3" {
		t.Fatalf("call 3 (new token): expected access token %q, got %#v", "access-token-3", out3)
	}
}

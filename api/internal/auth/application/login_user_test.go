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

type loginUserRepositoryMock struct {
	findByEmailCalls int
	findByEmailUser  *domain.User
	findByEmailErr   error
	findByEmailValue string
}

func (m *loginUserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *loginUserRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status domain.UserStatus) error {
	return nil
}

func (m *loginUserRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	m.findByEmailCalls++
	m.findByEmailValue = email
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	return m.findByEmailUser, nil
}

func (m *loginUserRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (m *loginUserRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}

type loginPasswordHasherMock struct {
	compareCalls    int
	compareHash     string
	comparePassword string
	compareErr      error
}

func (m *loginPasswordHasherMock) Hash(password string) (string, error) {
	return "", nil
}

func (m *loginPasswordHasherMock) Compare(hash string, password string) error {
	m.compareCalls++
	m.compareHash = hash
	m.comparePassword = password
	return m.compareErr
}

type tokenServiceMock struct {
	generateAccessCalls  int
	generateAccessClaims domain.TokenClaims
	accessToken          string
	generateAccessErr    error

	generateRefreshCalls int
	generateRefreshUser  uuid.UUID
	refreshToken         string
	generateRefreshErr   error
}

type sessionRepositoryMock struct {
	createCalls   int
	createUserID  uuid.UUID
	createHash    string
	createExpires time.Time
	createErr     error
}

func (m *sessionRepositoryMock) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	m.createCalls++
	m.createUserID = userID
	m.createHash = tokenHash
	m.createExpires = expiresAt
	return m.createErr
}

func (m *sessionRepositoryMock) FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, bool, error) {
	return uuid.Nil, time.Time{}, false, nil
}

func (m *sessionRepositoryMock) Revoke(ctx context.Context, tokenHash string) error {
	return nil
}

func (m *tokenServiceMock) GenerateAccessToken(claims domain.TokenClaims) (string, error) {
	m.generateAccessCalls++
	m.generateAccessClaims = claims
	if m.generateAccessErr != nil {
		return "", m.generateAccessErr
	}
	return m.accessToken, nil

}

func (m *tokenServiceMock) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	m.generateRefreshCalls++
	m.generateRefreshUser = userID
	if m.generateRefreshErr != nil {
		return "", m.generateRefreshErr
	}
	return m.refreshToken, nil
}

func (m *tokenServiceMock) ParseAccessToken(token string) (*domain.TokenClaims, error) {
	return nil, nil
}

func (m *tokenServiceMock) ParseRefreshToken(token string) (uuid.UUID, error) {
	return uuid.Nil, nil
}

func TestLoginUserUseCase_Execute_Success(t *testing.T) {
	customerID := uuid.New()
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           userID,
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleCustomer,
			CustomerID:   &customerID,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "jwt-token", refreshToken: "refresh-token"}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    " USER@EXAMPLE.COM ",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("expected output to be non-nil")
	}

	if output.AccessToken != "jwt-token" {
		t.Fatalf("expected access token %q, got %q", "jwt-token", output.AccessToken)
	}

	if output.RefreshToken != "refresh-token" {
		t.Fatalf("expected refresh token %q, got %q", "refresh-token", output.RefreshToken)
	}

	if output.UserID != userID {
		t.Fatalf("expected user ID %q, got %q", userID, output.UserID)
	}

	if output.Email != "user@example.com" {
		t.Fatalf("expected email %q, got %q", "user@example.com", output.Email)
	}

	if output.Role != string(domain.RoleCustomer) {
		t.Fatalf("expected role %q, got %q", domain.RoleCustomer, output.Role)
	}

	if output.CustomerID == nil || *output.CustomerID != customerID {
		t.Fatalf("expected customer ID %q, got %v", customerID, output.CustomerID)
	}

	if userRepo.findByEmailCalls != 1 {
		t.Fatalf("expected FindByEmail to be called once, got %d", userRepo.findByEmailCalls)
	}

	if userRepo.findByEmailValue != "user@example.com" {
		t.Fatalf("expected normalized lookup email, got %q", userRepo.findByEmailValue)
	}

	if hasher.compareCalls != 1 {
		t.Fatalf("expected Compare to be called once, got %d", hasher.compareCalls)
	}

	if hasher.compareHash != "stored-hash" {
		t.Fatalf("expected Compare hash %q, got %q", "stored-hash", hasher.compareHash)
	}

	if hasher.comparePassword != "password123" {
		t.Fatalf("expected Compare password %q, got %q", "password123", hasher.comparePassword)
	}

	if tokenService.generateAccessCalls != 1 {
		t.Fatalf("expected GenerateAccessToken to be called once, got %d", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 1 {
		t.Fatalf("expected GenerateRefreshToken to be called once, got %d", tokenService.generateRefreshCalls)
	}

	if tokenService.generateAccessClaims.UserID != userID {
		t.Fatalf("expected token user ID %q, got %q", userID, tokenService.generateAccessClaims.UserID)
	}

	if tokenService.generateAccessClaims.Role != domain.RoleCustomer {
		t.Fatalf("expected token role %q, got %q", domain.RoleCustomer, tokenService.generateAccessClaims.Role)
	}

	if tokenService.generateAccessClaims.CustomerID == nil || *tokenService.generateAccessClaims.CustomerID != customerID {
		t.Fatalf("expected token customer ID %q, got %v", customerID, tokenService.generateAccessClaims.CustomerID)
	}

	if tokenService.generateRefreshUser != userID {
		t.Fatalf("expected refresh token user ID %q, got %q", userID, tokenService.generateRefreshUser)
	}

	if sessionRepo.createCalls != 1 {
		t.Fatalf("expected session Create to be called once, got %d", sessionRepo.createCalls)
	}

	if sessionRepo.createUserID != userID {
		t.Fatalf("expected session user ID %q, got %q", userID, sessionRepo.createUserID)
	}

	expectedHash := sha256.Sum256([]byte("refresh-token"))
	if sessionRepo.createHash != hex.EncodeToString(expectedHash[:]) {
		t.Fatalf("expected session token hash %q, got %q", hex.EncodeToString(expectedHash[:]), sessionRepo.createHash)
	}

	if sessionRepo.createExpires.IsZero() {
		t.Fatal("expected session expires_at to be set")
	}
}

func TestLoginUserUseCase_Execute_UserNotFound(t *testing.T) {
	userRepo := &loginUserRepositoryMock{}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidCredentials, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if hasher.compareCalls != 0 {
		t.Fatalf("expected Compare not to be called, got %d calls", hasher.compareCalls)
	}

	if tokenService.generateAccessCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 0 {
		t.Fatalf("expected GenerateRefreshToken not to be called, got %d calls", tokenService.generateRefreshCalls)
	}
}

func TestLoginUserUseCase_Execute_WrongPassword(t *testing.T) {
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           uuid.New(),
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleCustomer,
		},
	}
	hasher := &loginPasswordHasherMock{compareErr: errors.New("wrong password")}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "bad-password",
	})

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidCredentials, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if tokenService.generateAccessCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 0 {
		t.Fatalf("expected GenerateRefreshToken not to be called, got %d calls", tokenService.generateRefreshCalls)
	}
}

func TestLoginUserUseCase_Execute_TokenGenerationFailure(t *testing.T) {
	expectedErr := errors.New("token unavailable")
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           uuid.New(),
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleAdmin,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{generateAccessErr: expectedErr}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if tokenService.generateAccessCalls != 1 {
		t.Fatalf("expected GenerateAccessToken to be called once, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 0 {
		t.Fatalf("expected GenerateRefreshToken not to be called, got %d calls", tokenService.generateRefreshCalls)
	}
}

func TestLoginUserUseCase_Execute_RefreshTokenGenerationFailure(t *testing.T) {
	expectedErr := errors.New("refresh token unavailable")
	userID := uuid.New()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           userID,
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleAdmin,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{
		accessToken:        "jwt-token",
		generateRefreshErr: expectedErr,
	}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if tokenService.generateAccessCalls != 1 {
		t.Fatalf("expected GenerateAccessToken to be called once, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 1 {
		t.Fatalf("expected GenerateRefreshToken to be called once, got %d calls", tokenService.generateRefreshCalls)
	}

	if tokenService.generateRefreshUser != userID {
		t.Fatalf("expected refresh token user ID %q, got %q", userID, tokenService.generateRefreshUser)
	}
}

func TestLoginUserUseCase_Execute_SessionPersistenceFailure(t *testing.T) {
	expectedErr := errors.New("session unavailable")
	userID := uuid.New()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           userID,
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleAdmin,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "jwt-token", refreshToken: "refresh-token"}
	sessionRepo := &sessionRepositoryMock{createErr: expectedErr}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if sessionRepo.createCalls != 1 {
		t.Fatalf("expected session Create to be called once, got %d", sessionRepo.createCalls)
	}
}

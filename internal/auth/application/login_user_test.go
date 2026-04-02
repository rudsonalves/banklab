package application

import (
	"context"
	"errors"
	"testing"

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

func (m *loginUserRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	m.findByEmailCalls++
	m.findByEmailValue = email
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	return m.findByEmailUser, nil
}

func (m *loginUserRepositoryMock) FindByID(ctx context.Context, id string) (*domain.User, error) {
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
	generateCalls  int
	generateClaims domain.TokenClaims
	generateToken  string
	generateErr    error
}

func (m *tokenServiceMock) GenerateAccessToken(claims domain.TokenClaims) (string, error) {
	m.generateCalls++
	m.generateClaims = claims
	if m.generateErr != nil {
		return "", m.generateErr
	}
	return m.generateToken, nil
}

func (m *tokenServiceMock) ParseAccessToken(token string) (*domain.TokenClaims, error) {
	return nil, nil
}

func TestLoginUserUseCase_Execute_Success(t *testing.T) {
	customerID := "customer-1"
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           "user-1",
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleCustomer,
			CustomerID:   &customerID,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{generateToken: "jwt-token"}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService)

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

	if output.UserID != "user-1" {
		t.Fatalf("expected user ID %q, got %q", "user-1", output.UserID)
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

	if tokenService.generateCalls != 1 {
		t.Fatalf("expected GenerateAccessToken to be called once, got %d", tokenService.generateCalls)
	}

	if tokenService.generateClaims.UserID != "user-1" {
		t.Fatalf("expected token user ID %q, got %q", "user-1", tokenService.generateClaims.UserID)
	}

	if tokenService.generateClaims.Role != domain.RoleCustomer {
		t.Fatalf("expected token role %q, got %q", domain.RoleCustomer, tokenService.generateClaims.Role)
	}

	if tokenService.generateClaims.CustomerID == nil || *tokenService.generateClaims.CustomerID != customerID {
		t.Fatalf("expected token customer ID %q, got %v", customerID, tokenService.generateClaims.CustomerID)
	}
}

func TestLoginUserUseCase_Execute_UserNotFound(t *testing.T) {
	userRepo := &loginUserRepositoryMock{}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected error %v, got %v", ErrInvalidCredentials, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if hasher.compareCalls != 0 {
		t.Fatalf("expected Compare not to be called, got %d calls", hasher.compareCalls)
	}

	if tokenService.generateCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateCalls)
	}
}

func TestLoginUserUseCase_Execute_WrongPassword(t *testing.T) {
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           "user-1",
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleCustomer,
		},
	}
	hasher := &loginPasswordHasherMock{compareErr: errors.New("wrong password")}
	tokenService := &tokenServiceMock{}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "bad-password",
	})

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected error %v, got %v", ErrInvalidCredentials, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if tokenService.generateCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateCalls)
	}
}

func TestLoginUserUseCase_Execute_TokenGenerationFailure(t *testing.T) {
	expectedErr := errors.New("token unavailable")
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:           "user-1",
			Email:        "user@example.com",
			PasswordHash: "stored-hash",
			Role:         domain.RoleAdmin,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{generateErr: expectedErr}
	useCase := NewLoginUserUseCase(userRepo, hasher, tokenService)

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

	if tokenService.generateCalls != 1 {
		t.Fatalf("expected GenerateAccessToken to be called once, got %d calls", tokenService.generateCalls)
	}
}

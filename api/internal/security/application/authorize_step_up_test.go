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

type stepUpTokenRepositoryMock struct {
	createCalls  int
	createErr    error
	created      *domain.StepUpToken
	findByJTI    *domain.StepUpToken
	consume      *domain.StepUpToken
	consumeCalls int
	consumeErr   error
	consumeJTI   string
	consumeNow   time.Time
	events       *[]string
}

func (m *stepUpTokenRepositoryMock) Create(ctx context.Context, token *domain.StepUpToken) error {
	m.createCalls++
	m.created = token
	if m.events != nil {
		*m.events = append(*m.events, "create-step-up-token")
	}
	if m.createErr == nil && token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	return m.createErr
}

func (m *stepUpTokenRepositoryMock) FindByJTI(ctx context.Context, jti string) (*domain.StepUpToken, error) {
	return m.findByJTI, nil
}

func (m *stepUpTokenRepositoryMock) ConsumeByJTI(ctx context.Context, jti string, now time.Time) (*domain.StepUpToken, error) {
	m.consumeCalls++
	m.consumeJTI = jti
	m.consumeNow = now
	if m.events != nil {
		*m.events = append(*m.events, "consume-step-up-token")
	}
	return m.consume, m.consumeErr
}

type stepUpTokenSignerMock struct {
	signCalls int
	signErr   error
	signed    string
	token     *domain.StepUpToken
	events    *[]string
}

func (m *stepUpTokenSignerMock) Sign(token *domain.StepUpToken) (string, error) {
	m.signCalls++
	m.token = token
	if m.events != nil {
		*m.events = append(*m.events, "sign-step-up-token")
	}
	return m.signed, m.signErr
}

func TestAuthorizeStepUpUseCase_Execute_Success(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	password := activeTransactionPassword(t, userID, "hashed-pin", now)
	password.FailedAttempts = 2
	events := []string{}

	passwordRepo := &transactionPasswordRepositoryMock{findByUserID: password}
	userRepo := &transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)}
	hasher := &transactionPasswordHasherMock{
		compareSet:     true,
		compareMatches: true,
	}
	tokenRepo := &stepUpTokenRepositoryMock{events: &events}
	signer := &stepUpTokenSignerMock{signed: "signed-step-up-token", events: &events}
	resolver := domain.NewDefaultStepUpPublicOperationResolver()
	uc := NewAuthorizeStepUpUseCase(passwordRepo, userRepo, hasher, tokenRepo, signer, resolver)
	uc.now = func() time.Time { return now }
	uc.newJTI = func() string { return "deterministic-jti" }

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if output.StepUpToken != "signed-step-up-token" {
		t.Fatalf("expected signed token, got %q", output.StepUpToken)
	}
	if output.ExpiresIn != 120 {
		t.Fatalf("expected expires_in 120, got %d", output.ExpiresIn)
	}
	if hasher.compareCalls != 1 {
		t.Fatalf("expected Compare once, got %d", hasher.compareCalls)
	}
	if hasher.compareHash != "hashed-pin" || hasher.comparePIN != "123456" {
		t.Fatalf("expected Compare with stored hash and PIN, got hash=%q pin=%q", hasher.compareHash, hasher.comparePIN)
	}
	if passwordRepo.saveCalls != 1 {
		t.Fatalf("expected SaveValidationState once, got %d", passwordRepo.saveCalls)
	}
	if passwordRepo.savedPassword.FailedAttempts != 0 {
		t.Fatalf("expected failed attempts reset, got %d", passwordRepo.savedPassword.FailedAttempts)
	}
	if tokenRepo.createCalls != 1 {
		t.Fatalf("expected step-up token Create once, got %d", tokenRepo.createCalls)
	}
	if tokenRepo.created == nil {
		t.Fatal("expected created step-up token to be captured")
	}
	if tokenRepo.created.JTI != "deterministic-jti" {
		t.Fatalf("expected jti %q, got %q", "deterministic-jti", tokenRepo.created.JTI)
	}
	if tokenRepo.created.UserID != userID {
		t.Fatalf("expected user id %q, got %q", userID, tokenRepo.created.UserID)
	}
	if tokenRepo.created.EndpointKey != domain.StepUpEndpointInternalTransferCreate {
		t.Fatalf("expected endpoint key %q, got %q", domain.StepUpEndpointInternalTransferCreate, tokenRepo.created.EndpointKey)
	}
	if !tokenRepo.created.ExpiresAt.Equal(now.Add(domain.StepUpTokenDefaultDuration)) {
		t.Fatalf("expected expires_at %v, got %v", now.Add(domain.StepUpTokenDefaultDuration), tokenRepo.created.ExpiresAt)
	}
	if signer.signCalls != 1 {
		t.Fatalf("expected signer once, got %d", signer.signCalls)
	}
	if signer.token != tokenRepo.created {
		t.Fatal("expected signer to receive the persisted step-up token instance")
	}
	if got := events; len(got) != 2 || got[0] != "create-step-up-token" || got[1] != "sign-step-up-token" {
		t.Fatalf("expected token persistence before signing, got events %v", got)
	}
}

func TestAuthorizeStepUpUseCase_Execute_MissingAuthenticatedUser(t *testing.T) {
	uc := NewAuthorizeStepUpUseCase(
		&transactionPasswordRepositoryMock{},
		&transactionPasswordUserRepositoryMock{},
		&transactionPasswordHasherMock{},
		&stepUpTokenRepositoryMock{},
		&stepUpTokenSignerMock{},
		domain.NewDefaultStepUpPublicOperationResolver(),
	)

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "123456",
	})

	if !errors.Is(err, authdomain.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
}

func TestAuthorizeStepUpUseCase_Execute_EndpointNotAllowed(t *testing.T) {
	userID := uuid.New()
	passwordRepo := &transactionPasswordRepositoryMock{}
	userRepo := &transactionPasswordUserRepositoryMock{}
	hasher := &transactionPasswordHasherMock{}
	tokenRepo := &stepUpTokenRepositoryMock{}
	signer := &stepUpTokenSignerMock{}
	uc := NewAuthorizeStepUpUseCase(
		passwordRepo,
		userRepo,
		hasher,
		tokenRepo,
		signer,
		domain.NewDefaultStepUpPublicOperationResolver(),
	)

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/pix-transfers",
		TransactionPassword: "123456",
	})

	if !errors.Is(err, domain.ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected user lookup not to be called, got %d", userRepo.findByIDCalls)
	}
	if hasher.compareCalls != 0 {
		t.Fatalf("expected Compare not to be called, got %d", hasher.compareCalls)
	}
	if tokenRepo.createCalls != 0 {
		t.Fatalf("expected token Create not to be called, got %d", tokenRepo.createCalls)
	}
	if signer.signCalls != 0 {
		t.Fatalf("expected signer not to be called, got %d", signer.signCalls)
	}
}

func TestAuthorizeStepUpUseCase_Execute_InvalidMethod(t *testing.T) {
	userID := uuid.New()
	passwordRepo := &transactionPasswordRepositoryMock{}
	userRepo := &transactionPasswordUserRepositoryMock{}
	hasher := &transactionPasswordHasherMock{}
	tokenRepo := &stepUpTokenRepositoryMock{}
	signer := &stepUpTokenSignerMock{}
	uc := NewAuthorizeStepUpUseCase(
		passwordRepo,
		userRepo,
		hasher,
		tokenRepo,
		signer,
		domain.NewDefaultStepUpPublicOperationResolver(),
	)

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "123456",
	})

	if !errors.Is(err, domain.ErrInvalidStepUpPublicOperationMethod) {
		t.Fatalf("expected ErrInvalidStepUpPublicOperationMethod, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if userRepo.findByIDCalls != 0 {
		t.Fatalf("expected user lookup not to be called, got %d", userRepo.findByIDCalls)
	}
	if hasher.compareCalls != 0 {
		t.Fatalf("expected Compare not to be called, got %d", hasher.compareCalls)
	}
	if tokenRepo.createCalls != 0 {
		t.Fatalf("expected token Create not to be called, got %d", tokenRepo.createCalls)
	}
	if signer.signCalls != 0 {
		t.Fatalf("expected signer not to be called, got %d", signer.signCalls)
	}
	if passwordRepo.saveCalls != 0 {
		t.Fatalf("expected SaveValidationState not to be called, got %d", passwordRepo.saveCalls)
	}
}

func TestAuthorizeStepUpUseCase_Execute_InvalidPathFormats(t *testing.T) {
	testCases := []string{
		"http://api.banklab.local/accounts/internal-transfers",
		"/accounts/internal-transfers?tenant=banklab",
		"/accounts/internal-transfers#details",
	}

	for _, path := range testCases {
		t.Run(path, func(t *testing.T) {
			userID := uuid.New()
			passwordRepo := &transactionPasswordRepositoryMock{}
			userRepo := &transactionPasswordUserRepositoryMock{}
			hasher := &transactionPasswordHasherMock{}
			tokenRepo := &stepUpTokenRepositoryMock{}
			signer := &stepUpTokenSignerMock{}
			uc := NewAuthorizeStepUpUseCase(
				passwordRepo,
				userRepo,
				hasher,
				tokenRepo,
				signer,
				domain.NewDefaultStepUpPublicOperationResolver(),
			)

			output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
				User:                authenticatedUser(userID),
				Method:              "POST",
				Path:                path,
				TransactionPassword: "123456",
			})

			if !errors.Is(err, domain.ErrInvalidStepUpPublicOperationPath) {
				t.Fatalf("expected ErrInvalidStepUpPublicOperationPath, got %v", err)
			}
			if output != nil {
				t.Fatalf("expected nil output, got %+v", output)
			}
			if userRepo.findByIDCalls != 0 {
				t.Fatalf("expected user lookup not to be called, got %d", userRepo.findByIDCalls)
			}
			if hasher.compareCalls != 0 {
				t.Fatalf("expected Compare not to be called, got %d", hasher.compareCalls)
			}
			if tokenRepo.createCalls != 0 {
				t.Fatalf("expected token Create not to be called, got %d", tokenRepo.createCalls)
			}
			if signer.signCalls != 0 {
				t.Fatalf("expected signer not to be called, got %d", signer.signCalls)
			}
			if passwordRepo.saveCalls != 0 {
				t.Fatalf("expected SaveValidationState not to be called, got %d", passwordRepo.saveCalls)
			}
		})
	}
}

func TestAuthorizeStepUpUseCase_Execute_TransactionPasswordNotSet(t *testing.T) {
	userID := uuid.New()
	uc := NewAuthorizeStepUpUseCase(
		&transactionPasswordRepositoryMock{},
		&transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)},
		&transactionPasswordHasherMock{},
		&stepUpTokenRepositoryMock{},
		&stepUpTokenSignerMock{},
		domain.NewDefaultStepUpPublicOperationResolver(),
	)

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "123456",
	})

	if !errors.Is(err, domain.ErrTransactionPasswordNotSet) {
		t.Fatalf("expected ErrTransactionPasswordNotSet, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
}

func TestAuthorizeStepUpUseCase_Execute_InvalidTransactionPassword(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	password := activeTransactionPassword(t, userID, "hashed-pin", now)
	password.FailedAttempts = 1
	passwordRepo := &transactionPasswordRepositoryMock{findByUserID: password}
	hasher := &transactionPasswordHasherMock{
		compareSet:     true,
		compareMatches: false,
	}
	tokenRepo := &stepUpTokenRepositoryMock{}
	signer := &stepUpTokenSignerMock{}
	uc := NewAuthorizeStepUpUseCase(
		passwordRepo,
		&transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)},
		hasher,
		tokenRepo,
		signer,
		domain.NewDefaultStepUpPublicOperationResolver(),
	)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "000000",
	})

	if !errors.Is(err, domain.ErrTransactionPasswordInvalid) {
		t.Fatalf("expected ErrTransactionPasswordInvalid, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if passwordRepo.saveCalls != 1 {
		t.Fatalf("expected SaveValidationState once, got %d", passwordRepo.saveCalls)
	}
	if passwordRepo.savedPassword.FailedAttempts != 2 {
		t.Fatalf("expected failed attempts 2, got %d", passwordRepo.savedPassword.FailedAttempts)
	}
	if tokenRepo.createCalls != 0 {
		t.Fatalf("expected token Create not to be called, got %d", tokenRepo.createCalls)
	}
	if signer.signCalls != 0 {
		t.Fatalf("expected signer not to be called, got %d", signer.signCalls)
	}
}

func TestAuthorizeStepUpUseCase_Execute_ThirdInvalidTransactionPasswordLocks(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	password := activeTransactionPassword(t, userID, "hashed-pin", now)
	password.FailedAttempts = 2
	passwordRepo := &transactionPasswordRepositoryMock{findByUserID: password}
	uc := NewAuthorizeStepUpUseCase(
		passwordRepo,
		&transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)},
		&transactionPasswordHasherMock{compareSet: true, compareMatches: false},
		&stepUpTokenRepositoryMock{},
		&stepUpTokenSignerMock{},
		domain.NewDefaultStepUpPublicOperationResolver(),
	)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "000000",
	})

	if !errors.Is(err, domain.ErrTransactionPasswordLocked) {
		t.Fatalf("expected ErrTransactionPasswordLocked, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if passwordRepo.saveCalls != 1 {
		t.Fatalf("expected SaveValidationState once, got %d", passwordRepo.saveCalls)
	}
	if passwordRepo.savedPassword.Status != domain.TransactionPasswordBlocked {
		t.Fatalf("expected blocked status, got %q", passwordRepo.savedPassword.Status)
	}
	if passwordRepo.savedPassword.LockedUntil == nil {
		t.Fatal("expected locked_until to be set")
	}
}

func TestAuthorizeStepUpUseCase_Execute_BlockedTransactionPassword(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	lockedUntil := now.Add(10 * time.Minute)
	password := &domain.TransactionPassword{
		ID:             uuid.New(),
		UserID:         userID,
		PasswordHash:   "hashed-pin",
		Status:         domain.TransactionPasswordBlocked,
		FailedAttempts: domain.TransactionPasswordMaxFailures,
		LockedUntil:    &lockedUntil,
		CreatedAt:      now.Add(-time.Hour),
		UpdatedAt:      now.Add(-time.Minute),
	}
	passwordRepo := &transactionPasswordRepositoryMock{findByUserID: password}
	hasher := &transactionPasswordHasherMock{}
	tokenRepo := &stepUpTokenRepositoryMock{}
	signer := &stepUpTokenSignerMock{}
	uc := NewAuthorizeStepUpUseCase(
		passwordRepo,
		&transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)},
		hasher,
		tokenRepo,
		signer,
		domain.NewDefaultStepUpPublicOperationResolver(),
	)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "123456",
	})

	if !errors.Is(err, domain.ErrTransactionPasswordLocked) {
		t.Fatalf("expected ErrTransactionPasswordLocked, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if hasher.compareCalls != 0 {
		t.Fatalf("expected Compare not to be called, got %d", hasher.compareCalls)
	}
	if passwordRepo.saveCalls != 0 {
		t.Fatalf("expected SaveValidationState not to be called, got %d", passwordRepo.saveCalls)
	}
	if tokenRepo.createCalls != 0 {
		t.Fatalf("expected token Create not to be called, got %d", tokenRepo.createCalls)
	}
	if signer.signCalls != 0 {
		t.Fatalf("expected signer not to be called, got %d", signer.signCalls)
	}
}

func TestAuthorizeStepUpUseCase_Execute_ExpiredLockIsNormalizedBeforeValidation(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	lockedUntil := now.Add(-time.Minute)
	password := &domain.TransactionPassword{
		ID:             uuid.New(),
		UserID:         userID,
		PasswordHash:   "hashed-pin",
		Status:         domain.TransactionPasswordBlocked,
		FailedAttempts: domain.TransactionPasswordMaxFailures,
		LockedUntil:    &lockedUntil,
		CreatedAt:      now.Add(-time.Hour),
		UpdatedAt:      now.Add(-time.Minute),
	}
	passwordRepo := &transactionPasswordRepositoryMock{findByUserID: password}
	uc := NewAuthorizeStepUpUseCase(
		passwordRepo,
		&transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)},
		&transactionPasswordHasherMock{compareSet: true, compareMatches: true},
		&stepUpTokenRepositoryMock{},
		&stepUpTokenSignerMock{signed: "signed-step-up-token"},
		domain.NewDefaultStepUpPublicOperationResolver(),
	)
	uc.now = func() time.Time { return now }
	uc.newJTI = func() string { return "deterministic-jti" }

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "123456",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if passwordRepo.saveCalls != 2 {
		t.Fatalf("expected SaveValidationState for normalize and success, got %d", passwordRepo.saveCalls)
	}
	if password.Status != domain.TransactionPasswordActive {
		t.Fatalf("expected normalized active status, got %q", password.Status)
	}
	if password.LockedUntil != nil {
		t.Fatalf("expected locked_until cleared, got %v", password.LockedUntil)
	}
}

func TestAuthorizeStepUpUseCase_Execute_TokenPersistenceFailureDoesNotSign(t *testing.T) {
	userID := uuid.New()
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	expectedErr := errors.New("insert step-up token failed")
	password := activeTransactionPassword(t, userID, "hashed-pin", now)
	tokenRepo := &stepUpTokenRepositoryMock{createErr: expectedErr}
	signer := &stepUpTokenSignerMock{signed: "signed-step-up-token"}
	uc := NewAuthorizeStepUpUseCase(
		&transactionPasswordRepositoryMock{findByUserID: password},
		&transactionPasswordUserRepositoryMock{findByIDValue: activeUser(userID)},
		&transactionPasswordHasherMock{compareSet: true, compareMatches: true},
		tokenRepo,
		signer,
		domain.NewDefaultStepUpPublicOperationResolver(),
	)
	uc.now = func() time.Time { return now }

	output, err := uc.Execute(context.Background(), AuthorizeStepUpInput{
		User:                authenticatedUser(userID),
		Method:              "POST",
		Path:                "/accounts/internal-transfers",
		TransactionPassword: "123456",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if tokenRepo.createCalls != 1 {
		t.Fatalf("expected token Create once, got %d", tokenRepo.createCalls)
	}
	if signer.signCalls != 0 {
		t.Fatalf("expected signer not to be called, got %d", signer.signCalls)
	}
}

func activeTransactionPassword(
	t *testing.T,
	userID uuid.UUID,
	hash string,
	now time.Time,
) *domain.TransactionPassword {
	t.Helper()

	password, err := domain.NewTransactionPassword(userID, hash, now)
	if err != nil {
		t.Fatalf("expected no error creating transaction password, got %v", err)
	}
	password.ID = uuid.New()

	return password
}

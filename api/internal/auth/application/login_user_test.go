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
	findByEmailCalls       int
	findByEmailUser        *domain.User
	findByEmailErr         error
	findByEmailValue       string
	findByIDForUpdateCalls int
	findByIDForUpdateUser  *domain.User
	findByIDForUpdateErr   error
	findByIDForUpdateValue uuid.UUID
}

func (m *loginUserRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	return nil
}

func (m *loginUserRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status domain.UserStatus) error {
	return nil
}

func (m *loginUserRepositoryMock) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	m.findByIDForUpdateCalls++
	m.findByIDForUpdateValue = id
	if m.findByIDForUpdateErr != nil {
		return nil, m.findByIDForUpdateErr
	}
	return m.findByIDForUpdateUser, nil
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

func (m *loginUserRepositoryMock) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
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
	createCalls          int
	createUserID         uuid.UUID
	createHash           string
	createExpires        time.Time
	createInstallationID *uuid.UUID
	createErr            error
}

type accountProvisioningCheckerMock struct {
	existsByCustomerIDCalls int
	existsByCustomerIDValue uuid.UUID
	existsByCustomerIDOK    bool
	existsByCustomerIDErr   error
}

type installationLoginClassifierMock struct {
	calls          int
	userID         uuid.UUID
	installationID uuid.UUID
	decision       *InstallationLoginDecision
	err            error
	decisions      []*InstallationLoginDecision
	errs           []error
}

type loginTransactorMock struct {
	runInTxCalls int
	runInTxErr   error
}

type firstInstallationBootstrapperMock struct {
	calls          int
	userID         uuid.UUID
	installationID uuid.UUID
	now            time.Time
	err            error
}

func (m *accountProvisioningCheckerMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	m.existsByCustomerIDCalls++
	m.existsByCustomerIDValue = customerID
	if m.existsByCustomerIDErr != nil {
		return false, m.existsByCustomerIDErr
	}
	return m.existsByCustomerIDOK, nil
}

func (m *installationLoginClassifierMock) Classify(
	ctx context.Context,
	userID uuid.UUID,
	installationID uuid.UUID,
) (*InstallationLoginDecision, error) {
	m.calls++
	m.userID = userID
	m.installationID = installationID
	if len(m.decisions) > 0 || len(m.errs) > 0 {
		var decision *InstallationLoginDecision
		var err error
		if len(m.decisions) > 0 {
			decision = m.decisions[0]
			m.decisions = m.decisions[1:]
		}
		if len(m.errs) > 0 {
			err = m.errs[0]
			m.errs = m.errs[1:]
		}
		return decision, err
	}
	return m.decision, m.err
}

func (m *loginTransactorMock) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	m.runInTxCalls++
	if m.runInTxErr != nil {
		return m.runInTxErr
	}
	return fn(ctx)
}

func (m *firstInstallationBootstrapperMock) BootstrapFirstInstallation(
	ctx context.Context,
	userID uuid.UUID,
	installationID uuid.UUID,
	now time.Time,
) error {
	m.calls++
	m.userID = userID
	m.installationID = installationID
	m.now = now
	return m.err
}

func (m *sessionRepositoryMock) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	return m.CreateWithInstallation(ctx, domain.CreateSessionInput{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
}

func (m *sessionRepositoryMock) CreateWithInstallation(ctx context.Context, input domain.CreateSessionInput) error {
	m.createCalls++
	m.createUserID = input.UserID
	m.createHash = input.TokenHash
	m.createExpires = input.ExpiresAt
	m.createInstallationID = input.InstallationID
	return m.createErr
}

func (m *sessionRepositoryMock) FindByTokenHash(ctx context.Context, tokenHash string) (uuid.UUID, time.Time, bool, error) {
	return uuid.Nil, time.Time{}, false, nil
}

func (m *sessionRepositoryMock) FindByTokenHashWithInstallation(ctx context.Context, tokenHash string) (*domain.SessionRecord, error) {
	return nil, nil
}

func (m *sessionRepositoryMock) Revoke(ctx context.Context, tokenHash string) error {
	return nil
}

func (m *sessionRepositoryMock) RevokeByUserIDAndInstallationID(ctx context.Context, userID uuid.UUID, installationID uuid.UUID, revokedAt time.Time) error {
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
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              userID,
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			CustomerID:      &customerID,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "jwt-token", refreshToken: "refresh-token"}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

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

	if accountProvisioning.existsByCustomerIDCalls != 1 {
		t.Fatalf("expected ExistsByCustomerID to be called once, got %d", accountProvisioning.existsByCustomerIDCalls)
	}

	if accountProvisioning.existsByCustomerIDValue != customerID {
		t.Fatalf("expected ExistsByCustomerID customer ID %q, got %q", customerID, accountProvisioning.existsByCustomerIDValue)
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
	if sessionRepo.createInstallationID != nil {
		t.Fatalf("expected no installation id for legacy login, got %q", *sessionRepo.createInstallationID)
	}
}

func TestLoginUserUseCase_Execute_CallsInstallationClassifierWhenConfigured(t *testing.T) {
	customerID := uuid.New()
	userID := uuid.New()
	installationID := uuid.New()
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              userID,
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			CustomerID:      &customerID,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "jwt-token", refreshToken: "refresh-token"}
	sessionRepo := &sessionRepositoryMock{}
	classifier := &installationLoginClassifierMock{
		decision: &InstallationLoginDecision{Classification: InstallationLoginKnown},
	}

	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo).
		WithInstallationClassifier(classifier)

	_, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:          "user@example.com",
		Password:       "password123",
		InstallationID: installationID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if classifier.calls != 1 {
		t.Fatalf("expected classifier to be called once, got %d", classifier.calls)
	}
	if classifier.userID != userID {
		t.Fatalf("expected classifier userID %q, got %q", userID, classifier.userID)
	}
	if classifier.installationID != installationID {
		t.Fatalf("expected classifier installationID %q, got %q", installationID, classifier.installationID)
	}
	if tokenService.generateAccessClaims.InstallationID == nil || *tokenService.generateAccessClaims.InstallationID != installationID {
		t.Fatalf("expected token installation ID %q, got %#v", installationID, tokenService.generateAccessClaims.InstallationID)
	}
	if sessionRepo.createInstallationID == nil || *sessionRepo.createInstallationID != installationID {
		t.Fatalf("expected session installation ID %q, got %#v", installationID, sessionRepo.createInstallationID)
	}
}

func TestLoginUserUseCase_Execute_FirstInstallationBootstrapsAtomically(t *testing.T) {
	customerID := uuid.New()
	userID := uuid.New()
	installationID := uuid.New()
	verifiedAt := time.Now().UTC()
	user := &domain.User{
		ID:              userID,
		Email:           "user@example.com",
		PasswordHash:    "stored-hash",
		Role:            domain.RoleCustomer,
		CustomerID:      &customerID,
		Status:          domain.UserStatusActive,
		EmailVerifiedAt: &verifiedAt,
		PhoneVerifiedAt: &verifiedAt,
	}
	userRepo := &loginUserRepositoryMock{
		findByEmailUser:       user,
		findByIDForUpdateUser: user,
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "jwt-token", refreshToken: "refresh-token"}
	sessionRepo := &sessionRepositoryMock{}
	classifier := &installationLoginClassifierMock{
		decisions: []*InstallationLoginDecision{
			{Classification: InstallationLoginFirst},
			{Classification: InstallationLoginFirst},
		},
	}
	transactor := &loginTransactorMock{}
	bootstrapper := &firstInstallationBootstrapperMock{}

	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo).
		WithInstallationClassifier(classifier).
		WithTransactor(transactor).
		WithFirstInstallationBootstrapper(bootstrapper)

	_, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:          "user@example.com",
		Password:       "password123",
		InstallationID: installationID,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transactor.runInTxCalls != 1 {
		t.Fatalf("expected RunInTx once, got %d", transactor.runInTxCalls)
	}
	if userRepo.findByIDForUpdateCalls != 1 {
		t.Fatalf("expected FindByIDForUpdate once, got %d", userRepo.findByIDForUpdateCalls)
	}
	if userRepo.findByIDForUpdateValue != userID {
		t.Fatalf("expected lock userID %q, got %q", userID, userRepo.findByIDForUpdateValue)
	}
	if classifier.calls != 2 {
		t.Fatalf("expected classifier twice, got %d", classifier.calls)
	}
	if bootstrapper.calls != 1 {
		t.Fatalf("expected bootstrapper once, got %d", bootstrapper.calls)
	}
	if bootstrapper.userID != userID {
		t.Fatalf("expected bootstrap userID %q, got %q", userID, bootstrapper.userID)
	}
	if bootstrapper.installationID != installationID {
		t.Fatalf("expected bootstrap installationID %q, got %q", installationID, bootstrapper.installationID)
	}
}

func TestLoginUserUseCase_Execute_FirstInstallationBootstrapLostRace(t *testing.T) {
	customerID := uuid.New()
	userID := uuid.New()
	installationID := uuid.New()
	verifiedAt := time.Now().UTC()
	user := &domain.User{
		ID:              userID,
		Email:           "user@example.com",
		PasswordHash:    "stored-hash",
		Role:            domain.RoleCustomer,
		CustomerID:      &customerID,
		Status:          domain.UserStatusActive,
		EmailVerifiedAt: &verifiedAt,
		PhoneVerifiedAt: &verifiedAt,
	}
	userRepo := &loginUserRepositoryMock{
		findByEmailUser:       user,
		findByIDForUpdateUser: user,
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "jwt-token", refreshToken: "refresh-token"}
	sessionRepo := &sessionRepositoryMock{}
	classifier := &installationLoginClassifierMock{
		decisions: []*InstallationLoginDecision{
			{Classification: InstallationLoginFirst},
			{Classification: InstallationLoginNew},
		},
	}
	transactor := &loginTransactorMock{}
	bootstrapper := &firstInstallationBootstrapperMock{}

	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo).
		WithInstallationClassifier(classifier).
		WithTransactor(transactor).
		WithFirstInstallationBootstrapper(bootstrapper)

	_, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:          "user@example.com",
		Password:       "password123",
		InstallationID: installationID,
	})
	if !errors.Is(err, errFirstInstallationBootstrapLostRace) {
		t.Fatalf("expected errFirstInstallationBootstrapLostRace, got %v", err)
	}
	if bootstrapper.calls != 0 {
		t.Fatalf("expected bootstrapper not to be called, got %d", bootstrapper.calls)
	}
}

func TestLoginUserUseCase_Execute_PendingCustomerRequiresApproval(t *testing.T) {
	customerID := uuid.New()
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              uuid.New(),
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			CustomerID:      &customerID,
			Status:          domain.UserStatusPending,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, domain.ErrAccountApprovalRequired) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountApprovalRequired, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if accountProvisioning.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected account provisioning not to be checked, got %d calls", accountProvisioning.existsByCustomerIDCalls)
	}

	if tokenService.generateAccessCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 0 {
		t.Fatalf("expected GenerateRefreshToken not to be called, got %d calls", tokenService.generateRefreshCalls)
	}

	if sessionRepo.createCalls != 0 {
		t.Fatalf("expected session Create not to be called, got %d calls", sessionRepo.createCalls)
	}
}

func TestLoginUserUseCase_Execute_ActiveCustomerWithoutAccountRequiresApproval(t *testing.T) {
	customerID := uuid.New()
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              uuid.New(),
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			CustomerID:      &customerID,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: false}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, domain.ErrAccountApprovalRequired) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountApprovalRequired, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if accountProvisioning.existsByCustomerIDCalls != 1 {
		t.Fatalf("expected ExistsByCustomerID to be called once, got %d", accountProvisioning.existsByCustomerIDCalls)
	}

	if accountProvisioning.existsByCustomerIDValue != customerID {
		t.Fatalf("expected ExistsByCustomerID customer ID %q, got %q", customerID, accountProvisioning.existsByCustomerIDValue)
	}

	if tokenService.generateAccessCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 0 {
		t.Fatalf("expected GenerateRefreshToken not to be called, got %d calls", tokenService.generateRefreshCalls)
	}

	if sessionRepo.createCalls != 0 {
		t.Fatalf("expected session Create not to be called, got %d calls", sessionRepo.createCalls)
	}
}

func TestLoginUserUseCase_Execute_ActiveCustomerWithoutCustomerIDRequiresApproval(t *testing.T) {
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              uuid.New(),
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, domain.ErrAccountApprovalRequired) {
		t.Fatalf("expected error %v, got %v", domain.ErrAccountApprovalRequired, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if accountProvisioning.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected account provisioning not to be checked, got %d calls", accountProvisioning.existsByCustomerIDCalls)
	}

	if tokenService.generateAccessCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateAccessCalls)
	}

	if sessionRepo.createCalls != 0 {
		t.Fatalf("expected session Create not to be called, got %d calls", sessionRepo.createCalls)
	}
}

func TestLoginUserUseCase_Execute_AdminWithoutAccountCanLogin(t *testing.T) {
	userID := uuid.New()
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              userID,
			Email:           "admin@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleAdmin,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "admin-jwt", refreshToken: "admin-refresh"}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "admin@example.com",
		Password: "password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("expected output to be non-nil")
	}

	if output.Role != string(domain.RoleAdmin) {
		t.Fatalf("expected role %q, got %q", domain.RoleAdmin, output.Role)
	}

	if output.CustomerID != nil {
		t.Fatalf("expected customer ID to be nil, got %v", output.CustomerID)
	}

	if accountProvisioning.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected account provisioning not to be checked for admin, got %d calls", accountProvisioning.existsByCustomerIDCalls)
	}

	if tokenService.generateAccessCalls != 1 {
		t.Fatalf("expected GenerateAccessToken to be called once, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 1 {
		t.Fatalf("expected GenerateRefreshToken to be called once, got %d calls", tokenService.generateRefreshCalls)
	}

	if sessionRepo.createCalls != 1 {
		t.Fatalf("expected session Create to be called once, got %d calls", sessionRepo.createCalls)
	}
}

func TestLoginUserUseCase_Execute_AccountProvisioningErrorIsWrapped(t *testing.T) {
	customerID := uuid.New()
	expectedErr := errors.New("database unavailable")
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              uuid.New(),
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			CustomerID:      &customerID,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDErr: expectedErr}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if errors.Is(err, domain.ErrAccountApprovalRequired) {
		t.Fatalf("expected repository error not to map to %v", domain.ErrAccountApprovalRequired)
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

	if sessionRepo.createCalls != 0 {
		t.Fatalf("expected session Create not to be called, got %d calls", sessionRepo.createCalls)
	}
}

func TestLoginUserUseCase_Execute_UserNotFound(t *testing.T) {
	userRepo := &loginUserRepositoryMock{}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, nil, hasher, tokenService, sessionRepo)

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
	useCase := NewLoginUserUseCase(userRepo, nil, hasher, tokenService, sessionRepo)

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
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              uuid.New(),
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleAdmin,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{generateAccessErr: expectedErr}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, nil, hasher, tokenService, sessionRepo)

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
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              userID,
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleAdmin,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{
		accessToken:        "jwt-token",
		generateRefreshErr: expectedErr,
	}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, nil, hasher, tokenService, sessionRepo)

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
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              userID,
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleAdmin,
			EmailVerifiedAt: &verifiedAt,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{accessToken: "jwt-token", refreshToken: "refresh-token"}
	sessionRepo := &sessionRepositoryMock{createErr: expectedErr}
	useCase := NewLoginUserUseCase(userRepo, nil, hasher, tokenService, sessionRepo)

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

func TestLoginUserUseCase_Execute_EmailNotVerified(t *testing.T) {
	customerID := uuid.New()
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              uuid.New(),
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			CustomerID:      &customerID,
			Status:          domain.UserStatusActive,
			PhoneVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, domain.ErrContactNotVerified) {
		t.Fatalf("expected error %v, got %v", domain.ErrContactNotVerified, err)
	}

	var contactErr *domain.ContactNotVerifiedError
	if !errors.As(err, &contactErr) {
		t.Fatal("expected *domain.ContactNotVerifiedError")
	}

	if contactErr.EmailVerified {
		t.Fatal("expected EmailVerified to be false")
	}

	if !contactErr.PhoneVerified {
		t.Fatal("expected PhoneVerified to be true")
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if accountProvisioning.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected account provisioning not to be checked, got %d calls", accountProvisioning.existsByCustomerIDCalls)
	}

	if tokenService.generateAccessCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 0 {
		t.Fatalf("expected GenerateRefreshToken not to be called, got %d calls", tokenService.generateRefreshCalls)
	}

	if sessionRepo.createCalls != 0 {
		t.Fatalf("expected session Create not to be called, got %d calls", sessionRepo.createCalls)
	}
}

func TestLoginUserUseCase_Execute_PhoneNotVerified(t *testing.T) {
	customerID := uuid.New()
	verifiedAt := time.Now().UTC()
	userRepo := &loginUserRepositoryMock{
		findByEmailUser: &domain.User{
			ID:              uuid.New(),
			Email:           "user@example.com",
			PasswordHash:    "stored-hash",
			Role:            domain.RoleCustomer,
			CustomerID:      &customerID,
			Status:          domain.UserStatusActive,
			EmailVerifiedAt: &verifiedAt,
		},
	}
	accountProvisioning := &accountProvisioningCheckerMock{existsByCustomerIDOK: true}
	hasher := &loginPasswordHasherMock{}
	tokenService := &tokenServiceMock{}
	sessionRepo := &sessionRepositoryMock{}
	useCase := NewLoginUserUseCase(userRepo, accountProvisioning, hasher, tokenService, sessionRepo)

	output, err := useCase.Execute(context.Background(), LoginUserInput{
		Email:    "user@example.com",
		Password: "password123",
	})

	if !errors.Is(err, domain.ErrContactNotVerified) {
		t.Fatalf("expected error %v, got %v", domain.ErrContactNotVerified, err)
	}

	var contactErr *domain.ContactNotVerifiedError
	if !errors.As(err, &contactErr) {
		t.Fatal("expected *domain.ContactNotVerifiedError")
	}

	if !contactErr.EmailVerified {
		t.Fatal("expected EmailVerified to be true")
	}

	if contactErr.PhoneVerified {
		t.Fatal("expected PhoneVerified to be false")
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if accountProvisioning.existsByCustomerIDCalls != 0 {
		t.Fatalf("expected account provisioning not to be checked, got %d calls", accountProvisioning.existsByCustomerIDCalls)
	}

	if tokenService.generateAccessCalls != 0 {
		t.Fatalf("expected GenerateAccessToken not to be called, got %d calls", tokenService.generateAccessCalls)
	}

	if tokenService.generateRefreshCalls != 0 {
		t.Fatalf("expected GenerateRefreshToken not to be called, got %d calls", tokenService.generateRefreshCalls)
	}

	if sessionRepo.createCalls != 0 {
		t.Fatalf("expected session Create not to be called, got %d calls", sessionRepo.createCalls)
	}
}

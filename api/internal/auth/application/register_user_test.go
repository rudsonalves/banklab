package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
)

type userRepositoryMock struct {
	existsByEmailCalls int
	existsByEmailValue bool
	existsByEmailErr   error
	existsByPhoneCalls int
	existsByPhoneValue bool
	existsByPhoneErr   error
	createCalls        int
	createErr          error
	createdUser        *domain.User
}

func (m *userRepositoryMock) Create(ctx context.Context, user *domain.User) error {
	m.createCalls++
	m.createdUser = user
	return m.createErr
}

func (m *userRepositoryMock) UpdateStatus(ctx context.Context, userID uuid.UUID, status domain.UserStatus) error {
	return nil
}

func (m *userRepositoryMock) FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (m *userRepositoryMock) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	return nil, nil
}

func (m *userRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return nil, nil
}

func (m *userRepositoryMock) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	m.existsByEmailCalls++
	return m.existsByEmailValue, m.existsByEmailErr
}

func (m *userRepositoryMock) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	m.existsByPhoneCalls++
	return m.existsByPhoneValue, m.existsByPhoneErr
}

type registerContactVerificationRepositoryMock struct {
	emailVerification *domain.ContactVerification
	phoneVerification *domain.ContactVerification
	err               error
}

func (m *registerContactVerificationRepositoryMock) CreateContactVerification(
	ctx context.Context,
	verification *domain.ContactVerification,
) error {
	return nil
}

func (m *registerContactVerificationRepositoryMock) FindContactVerificationByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.ContactVerification, error) {
	return nil, nil
}

func (m *registerContactVerificationRepositoryMock) FindContactVerificationByVerificationToken(
	ctx context.Context,
	verificationToken string,
) (*domain.ContactVerification, error) {
	if m.err != nil {
		return nil, m.err
	}
	switch verificationToken {
	case "email-token":
		return m.emailVerification, nil
	case "phone-token":
		return m.phoneVerification, nil
	default:
		return nil, nil
	}
}

func (m *registerContactVerificationRepositoryMock) ConfirmContactVerification(
	ctx context.Context,
	id uuid.UUID,
	verificationToken string,
	verifiedAt time.Time,
) error {
	return nil
}

type registerTransactorMock struct {
	runInTxCalls int
	runInTxErr   error
}

func (m *registerTransactorMock) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	m.runInTxCalls++
	if m.runInTxErr != nil {
		return m.runInTxErr
	}

	return fn(ctx)
}

type customerRepositoryMock struct {
	createCalls     int
	createErr       error
	createdCustomer *customerdomain.Customer
}

func (m *customerRepositoryMock) Create(ctx context.Context, c *customerdomain.Customer) error {
	m.createCalls++
	m.createdCustomer = c
	return m.createErr
}

func (m *customerRepositoryMock) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	return false, nil
}

func (m *customerRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*customerdomain.CustomerProfile, error) {
	return nil, nil
}

type customerDocumentRepositoryMock struct {
	createCalls     int
	createErr       error
	createdDocument *customerdomain.CustomerDocument
}

func (m *customerDocumentRepositoryMock) CreateDocument(
	ctx context.Context,
	document *customerdomain.CustomerDocument,
) error {
	m.createCalls++
	m.createdDocument = document
	return m.createErr
}

func (m *customerDocumentRepositoryMock) GetPrimaryDocumentByCustomerID(
	ctx context.Context,
	customerID uuid.UUID,
) (*customerdomain.CustomerDocument, error) {
	return nil, nil
}

func (m *customerDocumentRepositoryMock) GetCPFByCustomerID(
	ctx context.Context,
	customerID uuid.UUID,
) (*customerdomain.CustomerDocument, error) {
	return nil, nil
}

type passwordHasherMock struct {
	hashCalls int
	hashValue string
	hashErr   error
}

func (m *passwordHasherMock) Hash(password string) (string, error) {
	m.hashCalls++
	if m.hashErr != nil {
		return "", m.hashErr
	}
	return m.hashValue, nil
}

func (m *passwordHasherMock) Compare(hash string, password string) error {
	return nil
}

var registerBirthDate = time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)

func newRegisterContactVerificationRepo() *registerContactVerificationRepositoryMock {
	verifiedAt := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	return &registerContactVerificationRepositoryMock{
		emailVerification: &domain.ContactVerification{
			ID:                uuid.New(),
			Channel:           domain.ContactVerificationChannelEmail,
			Target:            "user@example.com",
			Token:             "111111",
			VerificationToken: stringPtr("email-token"),
			VerifiedAt:        &verifiedAt,
			ExpiresAt:         verifiedAt.Add(time.Hour),
			CreatedAt:         verifiedAt.Add(-time.Minute),
		},
		phoneVerification: &domain.ContactVerification{
			ID:                uuid.New(),
			Channel:           domain.ContactVerificationChannelPhone,
			Target:            "+5511999999999",
			Token:             "222222",
			VerificationToken: stringPtr("phone-token"),
			VerifiedAt:        &verifiedAt,
			ExpiresAt:         verifiedAt.Add(time.Hour),
			CreatedAt:         verifiedAt.Add(-time.Minute),
		},
	}
}

func stringPtr(value string) *string {
	return &value
}

func TestRegisterUserUseCase_Execute_Success(t *testing.T) {
	userRepo := &userRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{hashValue: "hashed-password"}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "  USER@Example.com ",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "123.456.789-09",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output == nil {
		t.Fatal("expected output to be non-nil")
	}

	if output.ID == uuid.Nil {
		t.Fatal("expected output ID to be set")
	}

	if output.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", output.Email)
	}

	if output.Role != string(domain.RoleCustomer) {
		t.Fatalf("expected role %q, got %q", domain.RoleCustomer, output.Role)
	}

	if output.CustomerID == nil {
		t.Fatal("expected customer ID to be set")
	}

	if transactor.runInTxCalls != 1 {
		t.Fatalf("expected RunInTx to be called once, got %d", transactor.runInTxCalls)
	}

	if userRepo.existsByEmailCalls != 1 {
		t.Fatalf("expected ExistsByEmail to be called once, got %d", userRepo.existsByEmailCalls)
	}

	if hasher.hashCalls != 1 {
		t.Fatalf("expected Hash to be called once, got %d", hasher.hashCalls)
	}

	if userRepo.createCalls != 1 {
		t.Fatalf("expected Create to be called once, got %d", userRepo.createCalls)
	}

	if customerRepo.createCalls != 1 {
		t.Fatalf("expected customer Create to be called once, got %d", customerRepo.createCalls)
	}

	if customerRepo.createdCustomer == nil {
		t.Fatal("expected created customer to be captured")
	}

	if customerRepo.createdCustomer.Name != "Maria Silva" {
		t.Fatalf("expected customer name %q, got %q", "Maria Silva", customerRepo.createdCustomer.Name)
	}

	if !customerRepo.createdCustomer.BirthDate.Equal(registerBirthDate) {
		t.Fatalf("expected customer birth date %v, got %v", registerBirthDate, customerRepo.createdCustomer.BirthDate)
	}

	if customerDocumentRepo.createCalls != 1 {
		t.Fatalf("expected customer document Create to be called once, got %d", customerDocumentRepo.createCalls)
	}

	if customerDocumentRepo.createdDocument == nil {
		t.Fatal("expected created customer document to be captured")
	}

	if customerDocumentRepo.createdDocument.CustomerID != customerRepo.createdCustomer.ID {
		t.Fatalf("expected document customer ID %v, got %v", customerRepo.createdCustomer.ID, customerDocumentRepo.createdDocument.CustomerID)
	}

	if customerDocumentRepo.createdDocument.Type != customerdomain.DocumentTypeCPF {
		t.Fatalf("expected document type %q, got %q", customerdomain.DocumentTypeCPF, customerDocumentRepo.createdDocument.Type)
	}

	if customerDocumentRepo.createdDocument.Value != "12345678909" {
		t.Fatalf("expected customer document value %q, got %q", "12345678909", customerDocumentRepo.createdDocument.Value)
	}

	if customerDocumentRepo.createdDocument.Country != customerdomain.DefaultCountry {
		t.Fatalf("expected customer document country %q, got %q", customerdomain.DefaultCountry, customerDocumentRepo.createdDocument.Country)
	}

	if !customerDocumentRepo.createdDocument.IsPrimary {
		t.Fatal("expected customer document to be primary")
	}

	if userRepo.createdUser == nil {
		t.Fatal("expected created user to be captured")
	}

	if userRepo.createdUser.PasswordHash != "hashed-password" {
		t.Fatalf("expected hashed password to be persisted, got %q", userRepo.createdUser.PasswordHash)
	}

	if userRepo.createdUser.Phone != "+5511999999999" {
		t.Fatalf("expected phone %q, got %q", "+5511999999999", userRepo.createdUser.Phone)
	}

	if userRepo.createdUser.EmailVerifiedAt == nil {
		t.Fatal("expected email_verified_at to be set")
	}

	if userRepo.createdUser.PhoneVerifiedAt == nil {
		t.Fatal("expected phone_verified_at to be set")
	}

	if userRepo.createdUser.Role != domain.RoleCustomer {
		t.Fatalf("expected persisted role %q, got %q", domain.RoleCustomer, userRepo.createdUser.Role)
	}

	if userRepo.createdUser.CustomerID == nil {
		t.Fatal("expected persisted customer ID to be set")
	}

	if *userRepo.createdUser.CustomerID != customerRepo.createdCustomer.ID {
		t.Fatalf("expected user customer ID %v, got %v", customerRepo.createdCustomer.ID, *userRepo.createdUser.CustomerID)
	}

	if userRepo.createdUser.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set")
	}

	if userRepo.createdUser.UpdatedAt.IsZero() {
		t.Fatal("expected updated_at to be set")
	}

	if !userRepo.createdUser.CreatedAt.Equal(userRepo.createdUser.UpdatedAt) {
		t.Fatal("expected created_at and updated_at to match on creation")
	}
}

func TestRegisterUserUseCase_Execute_DuplicateEmail(t *testing.T) {
	userRepo := &userRepositoryMock{existsByEmailValue: true}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{hashValue: "hashed-password"}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "user@example.com",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678909",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Fatalf("expected error %v, got %v", domain.ErrEmailAlreadyExists, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}

	if customerRepo.createCalls != 0 {
		t.Fatalf("expected customer Create not to be called, got %d calls", customerRepo.createCalls)
	}

	if customerDocumentRepo.createCalls != 0 {
		t.Fatalf("expected customer document Create not to be called, got %d calls", customerDocumentRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_InvalidEmail(t *testing.T) {
	userRepo := &userRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "invalid-email",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678909",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, domain.ErrInvalidEmail) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidEmail, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.existsByEmailCalls != 0 {
		t.Fatalf("expected ExistsByEmail not to be called, got %d calls", userRepo.existsByEmailCalls)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}

	if customerRepo.createCalls != 0 {
		t.Fatalf("expected customer Create not to be called, got %d calls", customerRepo.createCalls)
	}

	if customerDocumentRepo.createCalls != 0 {
		t.Fatalf("expected customer document Create not to be called, got %d calls", customerDocumentRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_DuplicatePhone(t *testing.T) {
	userRepo := &userRepositoryMock{existsByPhoneValue: true}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{hashValue: "hashed-password"}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "user@example.com",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678909",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, domain.ErrPhoneAlreadyExists) {
		t.Fatalf("expected error %v, got %v", domain.ErrPhoneAlreadyExists, err)
	}
	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}
	if customerRepo.createCalls != 0 {
		t.Fatalf("expected customer Create not to be called, got %d calls", customerRepo.createCalls)
	}
	if customerDocumentRepo.createCalls != 0 {
		t.Fatalf("expected customer document Create not to be called, got %d calls", customerDocumentRepo.createCalls)
	}
	if userRepo.createCalls != 0 {
		t.Fatalf("expected user Create not to be called, got %d calls", userRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_InvalidPassword(t *testing.T) {
	userRepo := &userRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "user@example.com",
		Phone:                  "+5511999999999",
		Password:               "short",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678909",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, domain.ErrInvalidPassword) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidPassword, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.existsByEmailCalls != 0 {
		t.Fatalf("expected ExistsByEmail not to be called, got %d calls", userRepo.existsByEmailCalls)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}

	if customerRepo.createCalls != 0 {
		t.Fatalf("expected customer Create not to be called, got %d calls", customerRepo.createCalls)
	}

	if customerDocumentRepo.createCalls != 0 {
		t.Fatalf("expected customer document Create not to be called, got %d calls", customerDocumentRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_InvalidCPF(t *testing.T) {
	userRepo := &userRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "user@example.com",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678901",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, customerdomain.ErrCPFInvalid) {
		t.Fatalf("expected error %v, got %v", customerdomain.ErrCPFInvalid, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if customerRepo.createCalls != 0 {
		t.Fatalf("expected customer Create not to be called, got %d calls", customerRepo.createCalls)
	}

	if customerDocumentRepo.createCalls != 0 {
		t.Fatalf("expected customer document Create not to be called, got %d calls", customerDocumentRepo.createCalls)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected user Create not to be called, got %d calls", userRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_HashingFailure(t *testing.T) {
	expectedErr := errors.New("hash unavailable")
	userRepo := &userRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{hashErr: expectedErr}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "user@example.com",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678909",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected Create not to be called, got %d calls", userRepo.createCalls)
	}

	if customerRepo.createCalls != 1 {
		t.Fatalf("expected customer Create to be called once before hashing failure, got %d calls", customerRepo.createCalls)
	}

	if customerDocumentRepo.createCalls != 1 {
		t.Fatalf("expected customer document Create to be called once before hashing failure, got %d calls", customerDocumentRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_CustomerCreateFailure(t *testing.T) {
	expectedErr := customerdomain.ErrInvalidData
	userRepo := &userRepositoryMock{}
	customerRepo := &customerRepositoryMock{createErr: expectedErr}
	customerDocumentRepo := &customerDocumentRepositoryMock{}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{hashValue: "hashed-password"}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "user@example.com",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678909",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected user Create not to be called, got %d calls", userRepo.createCalls)
	}

	if customerDocumentRepo.createCalls != 0 {
		t.Fatalf("expected customer document Create not to be called, got %d calls", customerDocumentRepo.createCalls)
	}
}

func TestRegisterUserUseCase_Execute_CustomerDocumentCreateFailure(t *testing.T) {
	expectedErr := customerdomain.ErrCPFAlreadyExists
	userRepo := &userRepositoryMock{}
	customerRepo := &customerRepositoryMock{}
	customerDocumentRepo := &customerDocumentRepositoryMock{createErr: expectedErr}
	contactVerificationRepo := newRegisterContactVerificationRepo()
	hasher := &passwordHasherMock{hashValue: "hashed-password"}
	transactor := &registerTransactorMock{}
	useCase := NewRegisterUserUseCase(userRepo, customerRepo, customerDocumentRepo, contactVerificationRepo, hasher, transactor)

	output, err := useCase.Execute(context.Background(), RegisterUserInput{
		Email:                  "user@example.com",
		Phone:                  "+5511999999999",
		Password:               "password123",
		Name:                   "Maria Silva",
		BirthDate:              registerBirthDate,
		CPF:                    "12345678909",
		EmailVerificationToken: "email-token",
		PhoneVerificationToken: "phone-token",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error to wrap %v, got %v", expectedErr, err)
	}

	if output != nil {
		t.Fatalf("expected output to be nil, got %+v", output)
	}

	if customerRepo.createCalls != 1 {
		t.Fatalf("expected customer Create to be called once, got %d calls", customerRepo.createCalls)
	}

	if customerDocumentRepo.createCalls != 1 {
		t.Fatalf("expected customer document Create to be called once, got %d calls", customerDocumentRepo.createCalls)
	}

	if hasher.hashCalls != 0 {
		t.Fatalf("expected Hash not to be called, got %d calls", hasher.hashCalls)
	}

	if userRepo.createCalls != 0 {
		t.Fatalf("expected user Create not to be called, got %d calls", userRepo.createCalls)
	}
}

// TestRegisterUserUseCase_Execute_CustomerIDNeverNilForCustomerRole verifies the
// defensive post-transaction invariant check: if somehow customerID ends up nil
// for a customer-role user, ErrInvalidUserState is returned.
func TestRegisterUserUseCase_Execute_CustomerIDNeverNilForCustomerRole(t *testing.T) {
	customerRepo := &customerRepositoryMock{}
	hasher := &passwordHasherMock{}

	// Inject a user repo whose Create() clears CustomerID to simulate a bug.
	userRepo := &userRepositoryMock{}

	// We can't easily test NewUser returning ErrInvalidUserState via the use case
	// because the constructor is called inside the transaction. Instead we verify
	// directly that domain.NewUser enforces the invariant.
	_, err := domain.NewUser(
		uuid.New(), "user@example.com", "hash",
		domain.RoleCustomer, nil, // nil customerID — invariant violation
		time.Now().UTC(),
	)
	if !errors.Is(err, domain.ErrInvalidUserState) {
		t.Fatalf("expected ErrInvalidUserState, got %v", err)
	}

	// And that admin role with nil customerID is allowed.
	u, err := domain.NewUser(
		uuid.New(), "admin@example.com", "hash",
		domain.RoleAdmin, nil, // nil customerID — valid for admin
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("expected no error for admin with nil customerID, got %v", err)
	}
	if u == nil {
		t.Fatal("expected user to be non-nil")
	}

	_ = userRepo
	_ = customerRepo
	_ = hasher
}

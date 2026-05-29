package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
)

type RegisterUserUseCase struct {
	userRepo                domain.UserRepository
	transactor              domain.Transactor
	customerRepo            customerdomain.CustomerRepository
	customerDocumentRepo    customerdomain.CustomerDocumentRepository
	contactVerificationRepo domain.ContactVerificationRepository
	hasher                  domain.PasswordHasher
}

// NewRegisterUserUseCase cria uma nova instância do RegisterUserUseCase com as
// dependências fornecidas. Requer um repositório de usuários para gerenciar dados
// de usuários, um repositório de clientes para gerenciar dados de clientes, um
// hasher de senha para criptografar senhas de usuários com segurança, um repositório
// de documentos de cliente para gerenciar documentos do cliente, um repositório de
// verificação de contato para validar tokens de e-mail e telefone já verificados, e
// um transactor para executar operações de banco de dados dentro de uma transação.
// Este caso de uso é responsável por lidar com o registro de novos usuários,
// incluindo validação de dados de entrada, criação de registros de clientes
// associados, e garantindo que toda a operação seja executada atomicamente para
// manter a integridade dos dados.
func NewRegisterUserUseCase(
	userRepo domain.UserRepository,
	customerRepo customerdomain.CustomerRepository,
	customerDocumentRepo customerdomain.CustomerDocumentRepository,
	contactVerificationRepo domain.ContactVerificationRepository,
	hasher domain.PasswordHasher,
	transactor domain.Transactor,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:                userRepo,
		transactor:              transactor,
		customerRepo:            customerRepo,
		customerDocumentRepo:    customerDocumentRepo,
		contactVerificationRepo: contactVerificationRepo,
		hasher:                  hasher,
	}
}

type RegisterUserInput struct {
	Email                  string
	Phone                  string
	Password               string
	Name                   string
	BirthDate              time.Time
	CPF                    string
	EmailVerificationToken string
	PhoneVerificationToken string
}

type RegisterUserOutput struct {
	ID         uuid.UUID
	Email      string
	Role       string
	CustomerID *uuid.UUID
}

// Execute performs the user registration process. It normalizes and validates e-mail,
// phone and password, validates e-mail and phone verification tokens, creates a new
// customer record with associated CPF document, hashes the password, and creates a
// new user record. The entire operation is executed within a database transaction to
// ensure atomicity and data integrity.
func (uc *RegisterUserUseCase) Execute(
	ctx context.Context,
	input RegisterUserInput,
) (*RegisterUserOutput, error) {
	email := normalizeEmail(input.Email)
	phone := normalizePhone(input.Phone)
	if !isValidEmail(email) {
		return nil, domain.ErrInvalidEmail
	}

	if phone == "" {
		return nil, domain.ErrInvalidData
	}

	if !isValidPassword(input.Password) {
		return nil, domain.ErrInvalidPassword
	}

	emailVerification, err := uc.validateContactVerification(
		ctx,
		input.EmailVerificationToken,
		domain.ContactVerificationChannelEmail,
		email,
	)
	if err != nil {
		return nil, err
	}

	phoneVerification, err := uc.validateContactVerification(
		ctx,
		input.PhoneVerificationToken,
		domain.ContactVerificationChannelPhone,
		phone,
	)
	if err != nil {
		return nil, err
	}

	var user *domain.User

	err = uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		exists, err := uc.userRepo.ExistsByEmail(txCtx, email)
		if err != nil {
			return fmt.Errorf("check email uniqueness: %w", err)
		}
		if exists {
			return domain.ErrEmailAlreadyExists
		}

		exists, err = uc.userRepo.ExistsByPhone(txCtx, phone)
		if err != nil {
			return fmt.Errorf("check phone uniqueness: %w", err)
		}
		if exists {
			return domain.ErrPhoneAlreadyExists
		}

		customer, err := customerdomain.NewCustomer(input.Name, input.BirthDate)
		if err != nil {
			return err
		}

		normalizedCPF := customerdomain.NormalizeCPF(input.CPF)
		if normalizedCPF == "" {
			return customerdomain.ErrCPFRequired
		}
		if !customerdomain.ValidateCPF(normalizedCPF) {
			return customerdomain.ErrCPFInvalid
		}

		if err := uc.customerRepo.Create(txCtx, customer); err != nil {
			return fmt.Errorf("create customer: %w", err)
		}

		cpfDocument, err := customerdomain.NewCPFDocument(customer.ID, normalizedCPF, true)
		if err != nil {
			return err
		}

		if err := uc.customerDocumentRepo.CreateDocument(txCtx, cpfDocument); err != nil {
			return fmt.Errorf("create customer cpf document: %w", err)
		}

		hash, err := uc.hasher.Hash(input.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		customerID := customer.ID
		now := time.Now().UTC()
		var newUserErr error
		user, newUserErr = domain.NewUser(email, hash, domain.RoleCustomer, &customerID, now)
		if newUserErr != nil {
			return newUserErr
		}
		user.Phone = phone
		user.EmailVerifiedAt = emailVerification.VerifiedAt
		user.PhoneVerifiedAt = phoneVerification.VerifiedAt

		if err := uc.userRepo.Create(txCtx, user); err != nil {
			return fmt.Errorf("create user: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if user == nil || (user.Role == domain.RoleCustomer && user.CustomerID == nil) {
		return nil, domain.ErrInvalidUserState
	}

	return &RegisterUserOutput{
		ID:         user.ID,
		Email:      user.Email,
		Role:       string(user.Role),
		CustomerID: user.CustomerID,
	}, nil
}

// normalizeEmail trims whitespace and converts the email to lowercase to ensure
// a consistent format for storage and comparison.
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizePhone(phone string) string {
	return strings.TrimSpace(phone)
}

// validateContactVerification checks whether a verification token points to a
// verified contact for the expected channel and target value.
func (uc *RegisterUserUseCase) validateContactVerification(
	ctx context.Context,
	verificationToken string,
	channel domain.ContactVerificationChannel,
	target string,
) (*domain.ContactVerification, error) {
	verificationToken = strings.TrimSpace(verificationToken)
	if verificationToken == "" || target == "" {
		return nil, domain.ErrInvalidData
	}

	verification, err := uc.contactVerificationRepo.FindContactVerificationByVerificationToken(ctx, verificationToken)
	if err != nil {
		return nil, err
	}
	if verification == nil ||
		verification.VerifiedAt == nil ||
		verification.Channel != channel ||
		verification.Target != target {
		return nil, domain.ErrInvalidVerificationToken
	}

	return verification, nil
}

// isValidEmail performs basic validation to check if the email has a valid
// format. It checks for the presence of an "@" symbol, ensures that there are
// local and domain parts, and that the domain part contains a dot.
func isValidEmail(email string) bool {
	if email == "" {
		return false
	}

	if strings.Count(email, "@") != 1 {
		return false
	}

	parts := strings.Split(email, "@")
	localPart := parts[0]
	domainPart := parts[1]
	if localPart == "" || domainPart == "" {
		return false
	}

	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") {
		return false
	}

	return strings.Contains(domainPart, ".")
}

// isValidPassword checks if the provided password meets the minimum requirements.
// In this case, it ensures that the password is not empty and has a minimum
// length of 8 characters.
func isValidPassword(password string) bool {
	if strings.TrimSpace(password) == "" {
		return false
	}

	return len(password) >= 8
}

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
	userRepo             domain.UserRepository
	transactor           domain.Transactor
	customerRepo         customerdomain.CustomerRepository
	customerDocumentRepo customerdomain.CustomerDocumentRepository
	hasher               domain.PasswordHasher
}

// NewRegisterUserUseCase cria uma nova instância do RegisterUserUseCase com as
// dependências fornecidas. Requer um repositório de usuários para gerenciar dados
// de usuários, um repositório de clientes para gerenciar dados de clientes, um
// hasher de senha para criptografar senhas de usuários com segurança, um repositório
// de documentos de cliente para gerenciar documentos do cliente, e um transactor
// para executar operações de banco de dados dentro de uma transação. Este caso de
// uso é responsável por lidar com o registro de novos usuários, incluindo validação
// de dados de entrada, criação de registros de clientes associados, e garantindo que
// toda a operação seja executada atomicamente para manter a integridade dos dados.
func NewRegisterUserUseCase(
	userRepo domain.UserRepository,
	customerRepo customerdomain.CustomerRepository,
	customerDocumentRepo customerdomain.CustomerDocumentRepository,
	hasher domain.PasswordHasher,
	transactor domain.Transactor,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:             userRepo,
		transactor:           transactor,
		customerRepo:         customerRepo,
		customerDocumentRepo: customerDocumentRepo,
		hasher:               hasher,
	}
}

type RegisterUserInput struct {
	Email     string
	Password  string
	Name      string
	BirthDate time.Time
	CPF       string
}

type RegisterUserOutput struct {
	ID         uuid.UUID
	Email      string
	Role       string
	CustomerID *uuid.UUID
}

// Execute performs the user registration process. It validates the provided email
// and password, creates a new customer record with associated CPF document, hashes
// the password, and creates a new user record. The entire operation is executed
// within a database transaction to ensure atomicity and data integrity.
func (uc *RegisterUserUseCase) Execute(
	ctx context.Context,
	input RegisterUserInput,
) (*RegisterUserOutput, error) {
	email := normalizeEmail(input.Email)
	if !isValidEmail(email) {
		return nil, domain.ErrInvalidEmail
	}

	if !isValidPassword(input.Password) {
		return nil, domain.ErrInvalidPassword
	}

	var user *domain.User

	err := uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		exists, err := uc.userRepo.ExistsByEmail(txCtx, email)
		if err != nil {
			return fmt.Errorf("check email uniqueness: %w", err)
		}
		if exists {
			return domain.ErrEmailAlreadyExists
		}

		customer, err := customerdomain.NewCustomer(input.Name, input.BirthDate)
		if err != nil {
			return err
		}

		cpfDocument, err := customerdomain.NewCPFDocument(customer.ID, input.CPF, true)
		if err != nil {
			return err
		}

		if err := uc.customerRepo.Create(txCtx, customer); err != nil {
			return fmt.Errorf("create customer: %w", err)
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
		user, newUserErr = domain.NewUser(uuid.New(), email, hash, domain.RoleCustomer, &customerID, now)
		if newUserErr != nil {
			return newUserErr
		}

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

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
	userRepo     domain.UserRepository
	transactor   domain.Transactor
	customerRepo customerdomain.CustomerRepository
	hasher       domain.PasswordHasher
}

// NewRegisterUserUseCase creates a new instance of the RegisterUserUseCase with the
// provided dependencies. It requires a user repository for managing user data, a
// customer repository for managing customer data, a password hasher for securely
// hashing user passwords, and a transactor for executing database operations within
// a transaction. This use case is responsible for handling the registration of new
// users, including validating input data, creating associated customer records, and
// ensuring that the entire operation is performed atomically to maintain data
// integrity.
func NewRegisterUserUseCase(
	userRepo domain.UserRepository,
	customerRepo customerdomain.CustomerRepository,
	hasher domain.PasswordHasher,
	transactor domain.Transactor,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:     userRepo,
		transactor:   transactor,
		customerRepo: customerRepo,
		hasher:       hasher,
	}
}

type RegisterUserInput struct {
	Email    string
	Password string
	Name     string
	CPF      string
}

type RegisterUserOutput struct {
	ID         uuid.UUID
	Email      string
	Role       string
	CustomerID *uuid.UUID
}

// Execute performs the user registration process. It validates the provided email
// and password, creates a new customer record, hashes the password, and creates
// a new user record.
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

		now := time.Now().UTC()
		customer := &customerdomain.Customer{
			ID:        uuid.New(),
			Name:      strings.TrimSpace(input.Name),
			CPF:       strings.TrimSpace(input.CPF),
			CreatedAt: now,
		}

		if err := uc.customerRepo.Create(txCtx, customer); err != nil {
			return fmt.Errorf("create customer: %w", err)
		}

		hash, err := uc.hasher.Hash(input.Password)
		if err != nil {
			return fmt.Errorf("hash password: %w", err)
		}

		customerID := customer.ID
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

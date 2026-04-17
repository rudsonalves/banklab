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
			Email:     email,
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

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

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

func isValidPassword(password string) bool {
	if strings.TrimSpace(password) == "" {
		return false
	}

	return len(password) >= 8
}

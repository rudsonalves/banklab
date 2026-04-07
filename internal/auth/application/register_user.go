package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type RegisterUserUseCase struct {
	userRepo domain.UserRepository
	hasher   domain.PasswordHasher
}

func NewRegisterUserUseCase(
	userRepo domain.UserRepository,
	hasher domain.PasswordHasher,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

type RegisterUserInput struct {
	Email    string
	Password string
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

	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("check email uniqueness: %w", err)
	}
	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	hash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	user := &domain.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hash,
		Role:         domain.RoleCustomer,
		CustomerID:   nil,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
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

package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type LoginUserUseCase struct {
	userRepo     domain.UserRepository
	hasher       domain.PasswordHasher
	tokenService domain.TokenService
}

func NewLoginUserUseCase(
	userRepo domain.UserRepository,
	hasher domain.PasswordHasher,
	tokenService domain.TokenService,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenService: tokenService,
	}
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	AccessToken string
	UserID      uuid.UUID
	Email       string
	Role        string
	CustomerID  *uuid.UUID
}

func (uc *LoginUserUseCase) Execute(
	ctx context.Context,
	input LoginUserInput,
) (*LoginUserOutput, error) {
	email := normalizeEmail(input.Email)
	if email == "" {
		return nil, domain.ErrInvalidEmail
	}

	if strings.TrimSpace(input.Password) == "" {
		return nil, domain.ErrInvalidPassword
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if user == nil {
		return nil, domain.ErrInvalidCredentials
	}

	if err := uc.hasher.Compare(user.PasswordHash, input.Password); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	token, err := uc.tokenService.GenerateAccessToken(domain.TokenClaims{
		UserID:     user.ID,
		Role:       user.Role,
		CustomerID: user.CustomerID,
	})
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &LoginUserOutput{
		AccessToken: token,
		UserID:      user.ID,
		Email:       user.Email,
		Role:        string(user.Role),
		CustomerID:  user.CustomerID,
	}, nil
}

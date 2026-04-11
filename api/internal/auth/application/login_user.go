package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type LoginUserUseCase struct {
	userRepo     domain.UserRepository
	hasher       domain.PasswordHasher
	tokenService domain.TokenService
	sessionRepo  domain.SessionRepository
}

const refreshSessionTTL = 30 * 24 * time.Hour

func NewLoginUserUseCase(
	userRepo domain.UserRepository,
	hasher domain.PasswordHasher,
	tokenService domain.TokenService,
	sessionRepo ...domain.SessionRepository,
) *LoginUserUseCase {
	var sr domain.SessionRepository
	if len(sessionRepo) > 0 {
		sr = sessionRepo[0]
	}

	return &LoginUserUseCase{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenService: tokenService,
		sessionRepo:  sr,
	}
}

type LoginUserInput struct {
	Email    string
	Password string
}

type LoginUserOutput struct {
	AccessToken  string
	RefreshToken string
	UserID       uuid.UUID
	Email        string
	Role         string
	CustomerID   *uuid.UUID
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

	accessToken, err := uc.tokenService.GenerateAccessToken(domain.TokenClaims{
		UserID:     user.ID,
		Role:       user.Role,
		CustomerID: user.CustomerID,
	})
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if uc.sessionRepo != nil {
		hash := sha256.Sum256([]byte(refreshToken))
		tokenHash := hex.EncodeToString(hash[:])

		err = uc.sessionRepo.Create(ctx, user.ID, tokenHash, time.Now().UTC().Add(refreshSessionTTL))
		if err != nil {
			return nil, fmt.Errorf("create session: %w", err)
		}
	}

	return &LoginUserOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       user.ID,
		Email:        user.Email,
		Role:         string(user.Role),
		CustomerID:   user.CustomerID,
	}, nil
}

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

// NewLoginUserUseCase creates a new instance of the LoginUserUseCase with the
// provided dependencies. It requires a user repository for fetching user data, a
// password hasher for verifying passwords, a token service for generating access
// and refresh tokens, and a session repository for managing user sessions. This
// use case is responsible for handling the login process, including validating
// credentials, generating tokens, and creating sessions.
func NewLoginUserUseCase(
	userRepo domain.UserRepository,
	hasher domain.PasswordHasher,
	tokenService domain.TokenService,
	sessionRepo domain.SessionRepository,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:     userRepo,
		hasher:       hasher,
		tokenService: tokenService,
		sessionRepo:  sessionRepo,
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

// Execute performs the login operation for a user. It validates the provided email
// and password, checks the user's credentials against the database, and if valid,
// generates an access token and a refresh token. It also creates a session for the
// user with the refresh token's hash. The output includes the access token, refresh
// token, and user information. If any step in the process fails, it returns an
// appropriate error.
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

	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	err = uc.sessionRepo.Create(ctx, user.ID, tokenHash, time.Now().UTC().Add(refreshSessionTTL))
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
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

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

type RefreshAccessTokenUseCase struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
	sessionRepo  domain.SessionRepository
	transactor   domain.Transactor
}

// NewRefreshAccessTokenUseCase creates a new instance of the RefreshAccessTokenUseCase with the
// provided dependencies. It requires a user repository for fetching user data, a token
// service for parsing and generating tokens, a session repository for managing user
// sessions, and a transactor for executing database operations within a transaction. This
// use case is responsible for handling the process of refreshing an access token using a valid refresh token, including validating the refresh token, checking the associated session,
// generating a new access token and refresh token, and updating the session accordingly.
func NewRefreshAccessTokenUseCase(
	userRepo domain.UserRepository,
	tokenService domain.TokenService,
	sessionRepo domain.SessionRepository,
	transactor domain.Transactor,
) *RefreshAccessTokenUseCase {
	return &RefreshAccessTokenUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
		sessionRepo:  sessionRepo,
		transactor:   transactor,
	}
}

type RefreshAccessTokenInput struct {
	RefreshToken string
}

type RefreshAccessTokenOutput struct {
	AccessToken  string
	RefreshToken string
}

// Execute performs the operation of refreshing an access token using a provided
// refresh token. It validates the refresh token, checks the associated session
// for validity, retrieves the user information, generates a new access token and
// refresh token, and updates the session with the new refresh token. If any step
// in the process fails, it returns an appropriate error indicating the reason for
// failure.
func (uc *RefreshAccessTokenUseCase) Execute(
	ctx context.Context,
	input RefreshAccessTokenInput,
) (*RefreshAccessTokenOutput, error) {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil, domain.ErrInvalidToken
	}

	userID, err := uc.tokenService.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	hash := sha256.Sum256([]byte(refreshToken))
	tokenHash := hex.EncodeToString(hash[:])

	storedUserID, expiresAt, revoked, err := uc.sessionRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("find session by token hash: %w", err)
	}

	if storedUserID == uuid.Nil {
		return nil, domain.ErrInvalidToken
	}

	if revoked {
		return nil, domain.ErrInvalidToken
	}

	if time.Now().UTC().After(expiresAt.UTC()) {
		return nil, domain.ErrInvalidToken
	}

	if storedUserID != userID {
		return nil, domain.ErrInvalidToken
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUnauthorized
	}

	accessToken, err := uc.tokenService.GenerateAccessToken(domain.TokenClaims{
		UserID:     user.ID,
		Role:       user.Role,
		CustomerID: user.CustomerID,
	})
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	newRefreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	newHash := sha256.Sum256([]byte(newRefreshToken))
	newTokenHash := hex.EncodeToString(newHash[:])

	if err := uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		if err := uc.sessionRepo.Revoke(txCtx, tokenHash); err != nil {
			return fmt.Errorf("revoke old session: %w", err)
		}
		if err := uc.sessionRepo.Create(txCtx, user.ID, newTokenHash, time.Now().UTC().Add(refreshSessionTTL)); err != nil {
			return fmt.Errorf("create new session: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &RefreshAccessTokenOutput{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

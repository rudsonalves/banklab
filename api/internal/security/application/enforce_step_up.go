package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type EnforceStepUpUseCase struct {
	tokenVerifier domain.StepUpTokenVerifier
	tokenRepo     domain.StepUpTokenRepository
	now           func() time.Time
}

func NewEnforceStepUpUseCase(
	tokenVerifier domain.StepUpTokenVerifier,
	tokenRepo domain.StepUpTokenRepository,
) *EnforceStepUpUseCase {
	return &EnforceStepUpUseCase{
		tokenVerifier: tokenVerifier,
		tokenRepo:     tokenRepo,
		now:           time.Now,
	}
}

type EnforceStepUpInput struct {
	User        *authdomain.AuthenticatedUser
	EndpointKey string
	Token       string
	Now         time.Time
}

func (uc *EnforceStepUpUseCase) Execute(ctx context.Context, input EnforceStepUpInput) error {
	if input.User == nil || input.User.UserID == uuid.Nil {
		return authdomain.ErrUnauthorized
	}

	if strings.TrimSpace(input.Token) == "" {
		return domain.ErrStepUpTokenRequired
	}

	if uc.tokenVerifier == nil || uc.tokenRepo == nil {
		return domain.ErrInvalidStepUpToken
	}

	now := input.Now.UTC()
	if now.IsZero() {
		now = uc.now().UTC()
	}

	claims, err := uc.tokenVerifier.Verify(input.Token)
	if err != nil {
		return err
	}
	if claims == nil {
		return domain.ErrInvalidStepUpToken
	}

	expectedEndpointKey := strings.TrimSpace(input.EndpointKey)
	if claims.UserID != input.User.UserID {
		return domain.ErrInvalidStepUpToken
	}
	if claims.EndpointKey != expectedEndpointKey {
		return domain.ErrStepUpEndpointMismatch
	}
	if !claims.ExpiresAt.After(now) {
		return domain.ErrStepUpTokenExpired
	}

	consumedToken, err := uc.tokenRepo.ConsumeByJTI(ctx, claims.JTI, now)
	if err != nil {
		return err
	}
	if consumedToken == nil {
		return domain.ErrInvalidStepUpToken
	}

	if consumedToken.UserID != claims.UserID ||
		consumedToken.EndpointKey != claims.EndpointKey ||
		consumedToken.JTI != claims.JTI {
		return domain.ErrInvalidStepUpToken
	}
	if consumedToken.EndpointKey != expectedEndpointKey {
		return domain.ErrStepUpEndpointMismatch
	}
	if consumedToken.ExpiresAt.Before(now) {
		return domain.ErrStepUpTokenExpired
	}

	return nil
}

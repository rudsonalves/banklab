package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

type AuthenticatedUser = domain.AuthenticatedUser

type GetCurrentUserUseCase struct {
	userRepo domain.UserRepository
}

func NewGetCurrentUserUseCase(userRepo domain.UserRepository) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{userRepo: userRepo}
}

type GetCurrentUserOutput struct {
	ID         uuid.UUID
	Email      string
	Role       string
	CustomerID *uuid.UUID
}

func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return sharedauthctx.WithAuthenticatedUser(ctx, user)
}

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	return sharedauthctx.GetAuthenticatedUser(ctx)
}

func (uc *GetCurrentUserUseCase) Execute(ctx context.Context) (*GetCurrentUserOutput, error) {
	principal, ok := GetAuthenticatedUser(ctx)
	if !ok || principal.UserID == uuid.Nil {
		return nil, domain.ErrUnauthorized
	}

	if uc.userRepo == nil {
		return &GetCurrentUserOutput{
			ID:         principal.UserID,
			Role:       string(principal.Role),
			CustomerID: principal.CustomerID,
		}, nil
	}

	user, err := uc.userRepo.FindByID(ctx, principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUnauthorized
	}

	return &GetCurrentUserOutput{
		ID:         user.ID,
		Email:      user.Email,
		Role:       string(user.Role),
		CustomerID: user.CustomerID,
	}, nil
}

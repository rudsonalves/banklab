package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

var ErrUnauthorized = errors.New("unauthorized")

type AuthenticatedUser struct {
	UserID     string
	Role       domain.Role
	CustomerID *string
}

type authenticatedUserContextKey struct{}

type GetCurrentUserUseCase struct {
	userRepo domain.UserRepository
}

func NewGetCurrentUserUseCase(userRepo domain.UserRepository) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{userRepo: userRepo}
}

type GetCurrentUserOutput struct {
	ID         string
	Email      string
	Role       string
	CustomerID *string
}

func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return context.WithValue(ctx, authenticatedUserContextKey{}, user)
}

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(authenticatedUserContextKey{}).(AuthenticatedUser)
	if ok {
		return &user, true
	}

	userPtr, ok := ctx.Value(authenticatedUserContextKey{}).(*AuthenticatedUser)
	if !ok || userPtr == nil {
		return nil, false
	}

	return userPtr, true
}

func (uc *GetCurrentUserUseCase) Execute(ctx context.Context) (*GetCurrentUserOutput, error) {
	principal, ok := GetAuthenticatedUser(ctx)
	if !ok || principal.UserID == "" {
		return nil, ErrUnauthorized
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
		return nil, ErrUnauthorized
	}

	return &GetCurrentUserOutput{
		ID:         user.ID,
		Email:      user.Email,
		Role:       string(user.Role),
		CustomerID: user.CustomerID,
	}, nil
}
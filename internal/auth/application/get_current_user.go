package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type AuthenticatedUser = domain.AuthenticatedUser

type contextKey string

const authenticatedUserKey contextKey = "authenticatedUser"

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
	return context.WithValue(ctx, authenticatedUserKey, user)
}

func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(authenticatedUserKey).(AuthenticatedUser)
	if ok {
		return &user, true
	}

	userPtr, ok := ctx.Value(authenticatedUserKey).(*AuthenticatedUser)
	if !ok || userPtr == nil {
		return nil, false
	}

	return userPtr, true
}

func (uc *GetCurrentUserUseCase) Execute(ctx context.Context) (*GetCurrentUserOutput, error) {
	principal, ok := GetAuthenticatedUser(ctx)
	if !ok || principal.UserID == "" {
		return nil, domain.ErrUnauthorized
	}

	if uc.userRepo == nil {
		return &GetCurrentUserOutput{
			ID:         principal.UserID,
			Role:       string(principal.Role),
			CustomerID: nullableUUIDToString(principal.CustomerID),
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

func nullableUUIDToString(value *uuid.UUID) *string {
	if value == nil {
		return nil
	}

	s := value.String()
	return &s
}

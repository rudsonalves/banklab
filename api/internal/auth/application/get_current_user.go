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

// NewGetCurrentUserUseCase creates a new instance of the GetCurrentUserUseCase with
// the provided user repository. This use case is responsible for retrieving the
// currently authenticated user's information based on the context. If the user
// repository is not provided, it will return basic information from the context
// without fetching additional details from the database.
func NewGetCurrentUserUseCase(userRepo domain.UserRepository) *GetCurrentUserUseCase {
	return &GetCurrentUserUseCase{userRepo: userRepo}
}

type GetCurrentUserOutput struct {
	ID         uuid.UUID
	Email      string
	Role       string
	CustomerID *uuid.UUID
}

// WithAuthenticatedUser adds the authenticated user information to the context. This
// function is typically used in middleware to set the authenticated user in the
// context for subsequent handlers to access.
func WithAuthenticatedUser(ctx context.Context, user AuthenticatedUser) context.Context {
	return sharedauthctx.WithAuthenticatedUser(ctx, user)
}

// GetAuthenticatedUser retrieves the authenticated user information from the context.
// It returns the authenticated user and a boolean indicating whether the user was
// successfully retrieved from the context. If the user is not found in the context,
// it returns nil and false.
func GetAuthenticatedUser(ctx context.Context) (*AuthenticatedUser, bool) {
	return sharedauthctx.GetAuthenticatedUser(ctx)
}

// Execute retrieves the current authenticated user's information. It first attempts
// to get the authenticated user from the context. If the user is not found or if the user ID is invalid, it returns an unauthorized error. If the user repository is
// not provided, it returns basic information from the context. If the user repository
// is available, it fetches the user's details from the database and returns them in
// the output. If any error occurs during this process, it returns an appropriate error.
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

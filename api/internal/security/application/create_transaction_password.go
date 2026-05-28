package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type CreateTransactionPasswordUseCase struct {
	passwordRepo domain.TransactionPasswordRepository
	userRepo     authdomain.UserRepository
	hasher       domain.TransactionPasswordHasher
	now          func() time.Time
}

func NewCreateTransactionPasswordUseCase(
	passwordRepo domain.TransactionPasswordRepository,
	userRepo authdomain.UserRepository,
	hasher domain.TransactionPasswordHasher,
) *CreateTransactionPasswordUseCase {
	return &CreateTransactionPasswordUseCase{
		passwordRepo: passwordRepo,
		userRepo:     userRepo,
		hasher:       hasher,
		now:          time.Now,
	}
}

type CreateTransactionPasswordInput struct {
	User                            *authdomain.AuthenticatedUser
	TransactionPassword             string
	TransactionPasswordConfirmation string
}

type CreateTransactionPasswordOutput struct {
	UserID    string
	Status    string
	CreatedAt time.Time
}

func (uc *CreateTransactionPasswordUseCase) Execute(
	ctx context.Context,
	input CreateTransactionPasswordInput,
) (*CreateTransactionPasswordOutput, error) {
	if input.User == nil || input.User.UserID == uuid.Nil {
		return nil, authdomain.ErrUnauthorized
	}

	if err := domain.ValidateTransactionPasswordPIN(input.TransactionPassword); err != nil {
		return nil, err
	}

	if input.TransactionPassword != input.TransactionPasswordConfirmation {
		return nil, domain.ErrInvalidTransactionPasswordPIN
	}

	user, err := uc.userRepo.FindByID(ctx, input.User.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if user == nil {
		return nil, authdomain.ErrUnauthorized
	}
	if user.Status != authdomain.UserStatusActive {
		return nil, authdomain.ErrForbidden
	}

	current, err := uc.passwordRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("find transaction password by user id: %w", err)
	}
	if current != nil {
		return nil, domain.ErrTransactionPasswordAlreadySet
	}

	hash, err := uc.hasher.Hash(input.TransactionPassword)
	if err != nil {
		return nil, fmt.Errorf("hash transaction password: %w", err)
	}

	now := uc.now().UTC()
	password, err := domain.NewTransactionPassword(user.ID, hash, now)
	if err != nil {
		return nil, err
	}

	if err := uc.passwordRepo.Create(ctx, password); err != nil {
		return nil, err
	}

	return &CreateTransactionPasswordOutput{
		UserID:    user.ID.String(),
		Status:    string(password.Status),
		CreatedAt: password.CreatedAt,
	}, nil
}

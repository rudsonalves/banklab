package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	customerdomain "github.com/seu-usuario/bank-api/internal/customer/domain"
	securitydomain "github.com/seu-usuario/bank-api/internal/security/domain"
)

const (
	TransactionPasswordSessionStatusActive  = "active"
	TransactionPasswordSessionStatusNotSet  = "not_set"
	TransactionPasswordSessionStatusLocked  = "locked"
	TransactionPasswordSessionStatusUnknown = "unknown"
)

type GetSessionUseCase struct {
	userRepo                authdomain.UserRepository
	customerRepo            customerdomain.CustomerRepository
	accountRepo             accountdomain.AccountRepository
	transactionPasswordRepo securitydomain.TransactionPasswordRepository
}

func NewGetSessionUseCase(
	userRepo authdomain.UserRepository,
	customerRepo customerdomain.CustomerRepository,
	accountRepo accountdomain.AccountRepository,
	transactionPasswordRepo securitydomain.TransactionPasswordRepository,
) *GetSessionUseCase {
	return &GetSessionUseCase{
		userRepo:                userRepo,
		customerRepo:            customerRepo,
		accountRepo:             accountRepo,
		transactionPasswordRepo: transactionPasswordRepo,
	}
}

type GetSessionOutput struct {
	User      GetSessionUserOutput
	Customer  GetSessionCustomerOutput
	Readiness GetSessionReadinessOutput
}

type GetSessionUserOutput struct {
	ID    uuid.UUID
	Email string
	Phone string
	Role  string
}

type GetSessionCustomerOutput struct {
	ID        uuid.UUID
	Name      string
	CPF       string
	BirthDate time.Time
	CreatedAt time.Time
}

type GetSessionReadinessOutput struct {
	OnboardingCompleted       bool
	Approved                  bool
	HasOperationalAccount     bool
	TransactionPasswordStatus string
	CanAccessHome             bool
}

func (uc *GetSessionUseCase) Execute(ctx context.Context) (*GetSessionOutput, error) {
	principal, ok := GetAuthenticatedUser(ctx)
	if !ok || principal.UserID == uuid.Nil {
		return nil, authdomain.ErrUnauthorized
	}

	if uc.userRepo == nil ||
		uc.customerRepo == nil ||
		uc.accountRepo == nil ||
		uc.transactionPasswordRepo == nil {
		return nil, fmt.Errorf("session use case dependencies not configured")
	}

	user, err := uc.userRepo.FindByID(ctx, principal.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if user == nil {
		return nil, authdomain.ErrUnauthorized
	}

	if user.CustomerID == nil {
		return nil, authdomain.ErrInvalidUserState
	}

	customerProfile, err := uc.customerRepo.GetByID(ctx, *user.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("get customer by id: %w", err)
	}
	if customerProfile == nil {
		return nil, customerdomain.ErrNotFound
	}

	approved := user.Status == authdomain.UserStatusActive

	accounts, err := uc.accountRepo.ListByCustomerID(ctx, customerProfile.Customer.ID)
	if err != nil {
		return nil, fmt.Errorf("list accounts by customer id: %w", err)
	}

	hasOperationalAccount := hasActiveAccount(accounts)

	transactionPasswordStatus, err := uc.transactionPasswordStatus(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	onboardingCompleted := true
	canAccessHome := onboardingCompleted &&
		approved &&
		hasOperationalAccount &&
		transactionPasswordStatus == TransactionPasswordSessionStatusActive

	return &GetSessionOutput{
		User: GetSessionUserOutput{
			ID:    user.ID,
			Email: user.Email,
			Phone: user.Phone,
			Role:  string(user.Role),
		},
		Customer: GetSessionCustomerOutput{
			ID:        customerProfile.Customer.ID,
			Name:      customerProfile.Customer.Name,
			CPF:       customerProfile.CPF,
			BirthDate: customerProfile.Customer.BirthDate,
			CreatedAt: customerProfile.Customer.CreatedAt,
		},
		Readiness: GetSessionReadinessOutput{
			OnboardingCompleted:       true,
			Approved:                  approved,
			HasOperationalAccount:     hasOperationalAccount,
			TransactionPasswordStatus: transactionPasswordStatus,
			CanAccessHome:             canAccessHome,
		},
	}, nil
}

func (uc *GetSessionUseCase) transactionPasswordStatus(
	ctx context.Context,
	userID uuid.UUID,
) (string, error) {
	password, err := uc.transactionPasswordRepo.FindByUserID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("find transaction password by user id: %w", err)
	}
	if password == nil {
		return TransactionPasswordSessionStatusNotSet, nil
	}

	switch password.Status {
	case securitydomain.TransactionPasswordActive:
		return TransactionPasswordSessionStatusActive, nil
	case securitydomain.TransactionPasswordBlocked:
		return TransactionPasswordSessionStatusLocked, nil
	default:
		return TransactionPasswordSessionStatusUnknown, nil
	}
}

func hasActiveAccount(accounts []accountdomain.Account) bool {
	for _, account := range accounts {
		if account.Status == accountdomain.AccountActive {
			return true
		}
	}

	return false
}

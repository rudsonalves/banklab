package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type AuthorizeStepUpUseCase struct {
	passwordRepo            domain.TransactionPasswordRepository
	userRepo                authdomain.UserRepository
	hasher                  domain.TransactionPasswordHasher
	tokenRepo               domain.StepUpTokenRepository
	tokenSigner             domain.StepUpTokenSigner
	publicOperationResolver domain.StepUpPublicOperationResolver
	now                     func() time.Time
	newJTI                  func() string
}

func NewAuthorizeStepUpUseCase(
	passwordRepo domain.TransactionPasswordRepository,
	userRepo authdomain.UserRepository,
	hasher domain.TransactionPasswordHasher,
	tokenRepo domain.StepUpTokenRepository,
	tokenSigner domain.StepUpTokenSigner,
	publicOperationResolver domain.StepUpPublicOperationResolver,
) *AuthorizeStepUpUseCase {
	return &AuthorizeStepUpUseCase{
		passwordRepo:            passwordRepo,
		userRepo:                userRepo,
		hasher:                  hasher,
		tokenRepo:               tokenRepo,
		tokenSigner:             tokenSigner,
		publicOperationResolver: publicOperationResolver,
		now:                     time.Now,
		newJTI:                  uuid.NewString,
	}
}

type AuthorizeStepUpInput struct {
	User                *authdomain.AuthenticatedUser
	Method              string
	Path                string
	TransactionPassword string
}

type AuthorizeStepUpOutput struct {
	StepUpToken string `json:"step_up_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (uc *AuthorizeStepUpUseCase) Execute(
	ctx context.Context,
	input AuthorizeStepUpInput,
) (*AuthorizeStepUpOutput, error) {
	if input.User == nil || input.User.UserID == uuid.Nil {
		return nil, authdomain.ErrUnauthorized
	}

	if uc.publicOperationResolver == nil {
		return nil, domain.ErrStepUpEndpointNotAllowed
	}

	publicOperation, err := domain.NewPublicHTTPOperation(input.Method, input.Path)
	if err != nil {
		return nil, err
	}

	endpointKey, err := uc.publicOperationResolver.Resolve(publicOperation)
	if err != nil {
		return nil, err
	}

	if err := domain.ValidateTransactionPasswordPIN(input.TransactionPassword); err != nil {
		return nil, err
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

	password, err := uc.passwordRepo.FindByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("find transaction password by user id: %w", err)
	}
	if password == nil {
		return nil, domain.ErrTransactionPasswordNotSet
	}

	now := uc.now().UTC()
	if err := uc.normalizeExpiredLock(ctx, password, now); err != nil {
		return nil, err
	}

	if err := password.CanValidate(now); err != nil {
		return nil, err
	}

	if !uc.hasher.Compare(password.PasswordHash, input.TransactionPassword) {
		err := password.RegisterFailure(now)
		if saveErr := uc.passwordRepo.SaveValidationState(ctx, password); saveErr != nil {
			return nil, fmt.Errorf("save transaction password validation state: %w", saveErr)
		}

		return nil, err
	}

	password.RegisterSuccess(now)
	if err := uc.passwordRepo.SaveValidationState(ctx, password); err != nil {
		return nil, fmt.Errorf("save transaction password validation state: %w", err)
	}

	stepUpToken, err := domain.NewStepUpToken(uc.newJTI(), user.ID, endpointKey, now)
	if err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.Create(ctx, stepUpToken); err != nil {
		return nil, err
	}

	signedToken, err := uc.tokenSigner.Sign(stepUpToken)
	if err != nil {
		return nil, err
	}

	return &AuthorizeStepUpOutput{
		StepUpToken: signedToken,
		ExpiresIn:   int(domain.StepUpTokenDefaultDuration / time.Second),
	}, nil
}

func (uc *AuthorizeStepUpUseCase) normalizeExpiredLock(
	ctx context.Context,
	password *domain.TransactionPassword,
	now time.Time,
) error {
	statusBefore := password.Status
	failedAttemptsBefore := password.FailedAttempts
	lockedUntilBefore := password.LockedUntil

	password.NormalizeLock(now)

	if password.Status == statusBefore &&
		password.FailedAttempts == failedAttemptsBefore &&
		sameTimePtr(password.LockedUntil, lockedUntilBefore) {
		return nil
	}

	if err := uc.passwordRepo.SaveValidationState(ctx, password); err != nil {
		return fmt.Errorf("save transaction password validation state: %w", err)
	}

	return nil
}

func sameTimePtr(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.Equal(*right)
}

package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type LoginUserUseCase struct {
	userRepo                   domain.UserRepository
	accountProvisioningChecker AccountProvisioningChecker
	hasher                     domain.PasswordHasher
	tokenService               domain.TokenService
	sessionRepo                domain.SessionRepository
	installationClassifier     InstallationLoginClassifier
	firstInstallationBootstrap FirstInstallationBootstrapper
	transactor                 domain.Transactor
	refreshSessionTTL          time.Duration
}

const defaultRefreshSessionTTL = 30 * 24 * time.Hour

var errFirstInstallationBootstrapLostRace = errors.New("first installation bootstrap lost race")

type AccountProvisioningChecker interface {
	ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error)
}

// NewLoginUserUseCase creates a new instance of the LoginUserUseCase with the
// provided dependencies. It requires a user repository for fetching user data,
// a provisioning checker for account approval checks, a password hasher for
// verifying passwords, a token service for generating access and refresh tokens,
// and a session repository for managing user sessions. This use case is
// responsible for handling the login process, including validating credentials,
// generating tokens, and creating sessions.
func NewLoginUserUseCase(
	userRepo domain.UserRepository,
	accountProvisioningChecker AccountProvisioningChecker,
	hasher domain.PasswordHasher,
	tokenService domain.TokenService,
	sessionRepo domain.SessionRepository,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:                   userRepo,
		accountProvisioningChecker: accountProvisioningChecker,
		hasher:                     hasher,
		tokenService:               tokenService,
		sessionRepo:                sessionRepo,
		refreshSessionTTL:          defaultRefreshSessionTTL,
	}
}

func (uc *LoginUserUseCase) WithRefreshSessionTTL(ttl time.Duration) *LoginUserUseCase {
	if ttl > 0 {
		uc.refreshSessionTTL = ttl
	}

	return uc
}

func (uc *LoginUserUseCase) WithInstallationClassifier(classifier InstallationLoginClassifier) *LoginUserUseCase {
	if classifier != nil {
		uc.installationClassifier = classifier
	}

	return uc
}

func (uc *LoginUserUseCase) WithFirstInstallationBootstrapper(
	bootstrapper FirstInstallationBootstrapper,
) *LoginUserUseCase {
	if bootstrapper != nil {
		uc.firstInstallationBootstrap = bootstrapper
	}

	return uc
}

func (uc *LoginUserUseCase) WithTransactor(transactor domain.Transactor) *LoginUserUseCase {
	if transactor != nil {
		uc.transactor = transactor
	}

	return uc
}

type LoginUserInput struct {
	Email          string
	Password       string
	InstallationID uuid.UUID
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

	if err := validateContactVerification(user); err != nil {
		return nil, err
	}

	if err := uc.validateLoginEligibility(ctx, user); err != nil {
		return nil, err
	}

	var installationDecision *InstallationLoginDecision
	if uc.installationClassifier != nil {
		installationDecision, err = uc.installationClassifier.Classify(ctx, user.ID, input.InstallationID)
		if err != nil {
			return nil, fmt.Errorf("classify login installation: %w", err)
		}
	}

	if installationDecision != nil && installationDecision.Classification == InstallationLoginFirst {
		if err := uc.bootstrapFirstInstallation(ctx, user.ID, input.InstallationID); err != nil {
			return nil, err
		}
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

	err = uc.sessionRepo.Create(ctx, user.ID, tokenHash, time.Now().UTC().Add(uc.refreshSessionTTL))
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

func (uc *LoginUserUseCase) validateLoginEligibility(ctx context.Context, user *domain.User) error {
	if user.Role != domain.RoleCustomer {
		return nil
	}

	switch user.Status {
	case domain.UserStatusPending:
		return domain.ErrAccountApprovalRequired
	case domain.UserStatusBlocked:
		return domain.ErrForbidden
	case domain.UserStatusActive:
	default:
		return domain.ErrForbidden
	}

	if user.CustomerID == nil {
		return domain.ErrAccountApprovalRequired
	}
	if uc.accountProvisioningChecker == nil {
		return fmt.Errorf("account provisioning checker not configured")
	}

	exists, err := uc.accountProvisioningChecker.ExistsByCustomerID(ctx, *user.CustomerID)
	if err != nil {
		return fmt.Errorf("check account provisioning: %w", err)
	}
	if !exists {
		return domain.ErrAccountApprovalRequired
	}

	return nil
}

func validateContactVerification(user *domain.User) error {
	emailVerified := user.EmailVerifiedAt != nil
	phoneVerified := user.PhoneVerifiedAt != nil
	if emailVerified && phoneVerified {
		return nil
	}

	return domain.NewContactNotVerifiedError(emailVerified, phoneVerified)
}

type FirstInstallationBootstrapper interface {
	BootstrapFirstInstallation(ctx context.Context, userID uuid.UUID, installationID uuid.UUID, now time.Time) error
}

func (uc *LoginUserUseCase) bootstrapFirstInstallation(
	ctx context.Context,
	userID uuid.UUID,
	installationID uuid.UUID,
) error {
	if uc.transactor == nil {
		return fmt.Errorf("first installation transactor not configured")
	}
	if uc.firstInstallationBootstrap == nil {
		return fmt.Errorf("first installation bootstrapper not configured")
	}
	if uc.installationClassifier == nil {
		return fmt.Errorf("installation classifier not configured")
	}

	return uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		lockedUser, err := uc.userRepo.FindByIDForUpdate(txCtx, userID)
		if err != nil {
			return fmt.Errorf("lock user for first installation bootstrap: %w", err)
		}
		if lockedUser == nil {
			return domain.ErrUserNotFound
		}

		decision, err := uc.installationClassifier.Classify(txCtx, userID, installationID)
		if err != nil {
			return fmt.Errorf("reclassify login installation in bootstrap: %w", err)
		}
		if decision == nil || decision.Classification != InstallationLoginFirst {
			return errFirstInstallationBootstrapLostRace
		}

		if err := uc.firstInstallationBootstrap.BootstrapFirstInstallation(
			txCtx,
			userID,
			installationID,
			time.Now().UTC(),
		); err != nil {
			return fmt.Errorf("bootstrap first installation: %w", err)
		}

		return nil
	})
}

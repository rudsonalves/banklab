package application

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
	securityapplication "github.com/seu-usuario/bank-api/internal/security/application"
	securitydomain "github.com/seu-usuario/bank-api/internal/security/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

const defaultRefreshSessionTTL = 30 * 24 * time.Hour

type UserReader interface {
	FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error)
}

type StepUpEnforcer interface {
	Execute(ctx context.Context, input securityapplication.EnforceStepUpInput) error
}

type RegisterInstallationUseCase struct {
	users                  UserReader
	installations          installationdomain.InstallationRepository
	authorizations         installationdomain.RestrictedAuthorizationRepository
	tokenService           authdomain.TokenService
	sessionRepo            authdomain.SessionRepository
	transactor             authdomain.Transactor
	stepUpEnforcer         StepUpEnforcer
	refreshSessionTTL      time.Duration
	installationResourceID func() uuid.UUID
	now                    func() time.Time
}

func NewRegisterInstallationUseCase(
	users UserReader,
	installations installationdomain.InstallationRepository,
	authorizations installationdomain.RestrictedAuthorizationRepository,
	tokenService authdomain.TokenService,
	sessionRepo authdomain.SessionRepository,
	transactor authdomain.Transactor,
	stepUpEnforcer StepUpEnforcer,
) *RegisterInstallationUseCase {
	return &RegisterInstallationUseCase{
		users:                  users,
		installations:          installations,
		authorizations:         authorizations,
		tokenService:           tokenService,
		sessionRepo:            sessionRepo,
		transactor:             transactor,
		stepUpEnforcer:         stepUpEnforcer,
		refreshSessionTTL:      defaultRefreshSessionTTL,
		installationResourceID: uuid.New,
		now:                    func() time.Time { return time.Now().UTC() },
	}
}

func (uc *RegisterInstallationUseCase) WithRefreshSessionTTL(ttl time.Duration) *RegisterInstallationUseCase {
	if ttl > 0 {
		uc.refreshSessionTTL = ttl
	}

	return uc
}

type RegisterInstallationInput struct {
	PresentedInstallationID uuid.UUID
	StepUpToken             string
	Now                     time.Time
}

type RegisterInstallationOutput struct {
	AccessToken            string
	RefreshToken           string
	InstallationResourceID uuid.UUID
	InstallationStatus     string
}

func (uc *RegisterInstallationUseCase) Execute(
	ctx context.Context,
	input RegisterInstallationInput,
) (*RegisterInstallationOutput, error) {
	restricted, err := sharedauthctx.RequireRestrictedSession(ctx)
	if err != nil {
		return nil, err
	}
	if err := validateRestrictedSession(restricted); err != nil {
		return nil, err
	}

	presentedInstallationID, err := installationdomain.NewInstallationID(input.PresentedInstallationID)
	if err != nil {
		return nil, err
	}
	restrictedInstallationID, err := installationdomain.NewInstallationID(restricted.InstallationID)
	if err != nil {
		return nil, err
	}
	if presentedInstallationID.UUID() != restrictedInstallationID.UUID() {
		return nil, installationdomain.ErrInstallationMismatch
	}

	if uc == nil ||
		uc.users == nil ||
		uc.installations == nil ||
		uc.authorizations == nil ||
		uc.tokenService == nil ||
		uc.sessionRepo == nil ||
		uc.transactor == nil ||
		uc.stepUpEnforcer == nil {
		return nil, authdomain.ErrInvalidData
	}

	user, err := uc.users.FindByID(ctx, restricted.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if user == nil {
		return nil, authdomain.ErrUnauthorized
	}

	now := input.Now.UTC()
	if now.IsZero() {
		now = uc.now().UTC()
	}

	if err := uc.stepUpEnforcer.Execute(ctx, securityapplication.EnforceStepUpInput{
		User: &authdomain.AuthenticatedUser{
			UserID:     user.ID,
			Role:       user.Role,
			CustomerID: user.CustomerID,
		},
		EndpointKey: securitydomain.StepUpEndpointInstallationRegisterCreate,
		Token:       input.StepUpToken,
		Now:         now,
	}); err != nil {
		return nil, err
	}

	resourceID, err := installationdomain.NewInstallationResourceID(uc.installationResourceID())
	if err != nil {
		return nil, err
	}

	accessToken, err := uc.tokenService.GenerateAccessToken(authdomain.TokenClaims{
		UserID:         user.ID,
		Role:           user.Role,
		CustomerID:     user.CustomerID,
		InstallationID: optionalUUID(presentedInstallationID.UUID()),
	})
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	refreshHash := hashToken(refreshToken)

	var registered *installationdomain.Installation
	if err := uc.transactor.RunInTx(ctx, func(txCtx context.Context) error {
		authorization, err := uc.authorizations.ConsumeByJTI(txCtx, restricted.JTI, now)
		if err != nil {
			return fmt.Errorf("consume restricted authorization: %w", err)
		}
		if err := validateConsumedAuthorization(authorization, restricted.UserID, restrictedInstallationID); err != nil {
			return err
		}

		installation, err := uc.installations.ReserveKnownInstallation(
			txCtx,
			restricted.UserID,
			resourceID,
			restrictedInstallationID,
			installationdomain.MaxKnownInstallations,
			now,
		)
		if err != nil {
			return fmt.Errorf("reserve known installation: %w", err)
		}
		registered = installation

		if err := uc.sessionRepo.CreateWithInstallation(txCtx, authdomain.CreateSessionInput{
			UserID:         user.ID,
			TokenHash:      refreshHash,
			ExpiresAt:      now.Add(uc.refreshSessionTTL),
			InstallationID: optionalUUID(presentedInstallationID.UUID()),
		}); err != nil {
			return fmt.Errorf("create installation session: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}
	if registered == nil {
		return nil, installationdomain.ErrInvalidInstallation
	}

	return &RegisterInstallationOutput{
		AccessToken:            accessToken,
		RefreshToken:           refreshToken,
		InstallationResourceID: registered.ResourceID.UUID(),
		InstallationStatus:     string(registered.Status),
	}, nil
}

func validateRestrictedSession(session *sharedauthctx.RestrictedSession) error {
	if session == nil ||
		session.UserID == uuid.Nil ||
		session.InstallationID == uuid.Nil ||
		strings.TrimSpace(session.JTI) == "" ||
		strings.TrimSpace(session.Scope) != installationdomain.RestrictedAuthorizationScopeInstallationRegister {
		return installationdomain.ErrRestrictedAuthorizationInvalid
	}

	return nil
}

func validateConsumedAuthorization(
	authorization *installationdomain.RestrictedAuthorization,
	userID uuid.UUID,
	installationID installationdomain.InstallationID,
) error {
	if authorization == nil {
		return installationdomain.ErrRestrictedAuthorizationNotFound
	}
	if authorization.UserID != userID ||
		authorization.InstallationID.UUID() != installationID.UUID() ||
		authorization.Scope != installationdomain.RestrictedAuthorizationScopeInstallationRegister ||
		authorization.Status != installationdomain.RestrictedAuthorizationStatusConsumed {
		return installationdomain.ErrRestrictedAuthorizationInvalid
	}

	return nil
}

func optionalUUID(value uuid.UUID) *uuid.UUID {
	if value == uuid.Nil {
		return nil
	}

	return &value
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

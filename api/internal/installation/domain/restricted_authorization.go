package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	RestrictedAuthorizationScopeInstallationRegister = "installation.register"
	RestrictedAuthorizationDefaultDuration           = 5 * time.Minute
)

type RestrictedAuthorizationStatus string

const (
	RestrictedAuthorizationStatusActive   RestrictedAuthorizationStatus = "active"
	RestrictedAuthorizationStatusConsumed RestrictedAuthorizationStatus = "consumed"
	RestrictedAuthorizationStatusRevoked  RestrictedAuthorizationStatus = "revoked"
)

type RestrictedAuthorization struct {
	ID             uuid.UUID
	JTI            string
	UserID         uuid.UUID
	InstallationID InstallationID
	Scope          string
	Status         RestrictedAuthorizationStatus
	ExpiresAt      time.Time
	ConsumedAt     *time.Time
	CreatedAt      time.Time
}

func NewRestrictedAuthorization(
	jti string,
	userID uuid.UUID,
	installationID InstallationID,
	now time.Time,
) (*RestrictedAuthorization, error) {
	createdAt := now.UTC()

	return RestoreRestrictedAuthorization(
		uuid.Nil,
		jti,
		userID,
		installationID,
		RestrictedAuthorizationScopeInstallationRegister,
		RestrictedAuthorizationStatusActive,
		createdAt.Add(RestrictedAuthorizationDefaultDuration),
		nil,
		createdAt,
	)
}

func RestoreRestrictedAuthorization(
	id uuid.UUID,
	jti string,
	userID uuid.UUID,
	installationID InstallationID,
	scope string,
	status RestrictedAuthorizationStatus,
	expiresAt time.Time,
	consumedAt *time.Time,
	createdAt time.Time,
) (*RestrictedAuthorization, error) {
	authorization := &RestrictedAuthorization{
		ID:             id,
		JTI:            strings.TrimSpace(jti),
		UserID:         userID,
		InstallationID: installationID,
		Scope:          strings.TrimSpace(scope),
		Status:         status,
		ExpiresAt:      expiresAt.UTC(),
		ConsumedAt:     utcTimePtr(consumedAt),
		CreatedAt:      createdAt.UTC(),
	}

	if err := authorization.Validate(); err != nil {
		return nil, err
	}

	return authorization, nil
}

func (a *RestrictedAuthorization) Validate() error {
	if a == nil ||
		a.JTI == "" ||
		a.UserID == uuid.Nil ||
		a.InstallationID.IsZero() ||
		a.Scope != RestrictedAuthorizationScopeInstallationRegister ||
		a.ExpiresAt.IsZero() ||
		a.CreatedAt.IsZero() ||
		!a.ExpiresAt.After(a.CreatedAt) {
		return ErrInvalidRestrictedAuthorization
	}

	switch a.Status {
	case RestrictedAuthorizationStatusActive:
		if a.ConsumedAt != nil {
			return ErrInvalidRestrictedAuthorization
		}
	case RestrictedAuthorizationStatusConsumed:
		if a.ConsumedAt == nil {
			return ErrInvalidRestrictedAuthorization
		}
	case RestrictedAuthorizationStatusRevoked:
	default:
		return ErrInvalidRestrictedAuthorization
	}

	return nil
}

func (a *RestrictedAuthorization) IsExpired(now time.Time) bool {
	if a == nil {
		return true
	}

	return a.Status == RestrictedAuthorizationStatusActive && !a.ExpiresAt.After(now.UTC())
}

func (a *RestrictedAuthorization) Consume(now time.Time) error {
	if err := a.Validate(); err != nil {
		return err
	}
	switch a.Status {
	case RestrictedAuthorizationStatusConsumed:
		return ErrRestrictedAuthorizationConsumed
	case RestrictedAuthorizationStatusRevoked:
		return ErrRestrictedAuthorizationRevoked
	}
	if a.IsExpired(now) {
		return ErrRestrictedAuthorizationExpired
	}

	consumedAt := now.UTC()
	a.Status = RestrictedAuthorizationStatusConsumed
	a.ConsumedAt = &consumedAt

	return nil
}

func (a *RestrictedAuthorization) Revoke() error {
	if err := a.Validate(); err != nil {
		return err
	}
	if a.Status == RestrictedAuthorizationStatusConsumed {
		return ErrRestrictedAuthorizationConsumed
	}
	if a.Status == RestrictedAuthorizationStatusRevoked {
		return ErrRestrictedAuthorizationRevoked
	}

	a.Status = RestrictedAuthorizationStatusRevoked

	return nil
}

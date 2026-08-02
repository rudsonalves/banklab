package domain

import (
	"time"

	"github.com/google/uuid"
)

const MaxKnownInstallations = 3

type InstallationStatus string

const (
	InstallationStatusKnown   InstallationStatus = "known"
	InstallationStatusRevoked InstallationStatus = "revoked"
)

type Installation struct {
	ID             uuid.UUID
	ResourceID     InstallationResourceID
	UserID         uuid.UUID
	InstallationID InstallationID
	Status         InstallationStatus
	Platform       string
	AppVersion     string
	AppBuild       string
	FirstSeenAt    time.Time
	LastSeenAt     time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewKnownInstallation(
	userID uuid.UUID,
	resourceID InstallationResourceID,
	installationID InstallationID,
	now time.Time,
) (*Installation, error) {
	timestamp := now.UTC()

	return RestoreInstallation(
		uuid.Nil,
		resourceID,
		userID,
		installationID,
		InstallationStatusKnown,
		"",
		"",
		"",
		timestamp,
		timestamp,
		nil,
		timestamp,
		timestamp,
	)
}

func RestoreInstallation(
	id uuid.UUID,
	resourceID InstallationResourceID,
	userID uuid.UUID,
	installationID InstallationID,
	status InstallationStatus,
	platform string,
	appVersion string,
	appBuild string,
	firstSeenAt time.Time,
	lastSeenAt time.Time,
	revokedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (*Installation, error) {
	installation := &Installation{
		ID:             id,
		ResourceID:     resourceID,
		UserID:         userID,
		InstallationID: installationID,
		Status:         status,
		Platform:       platform,
		AppVersion:     appVersion,
		AppBuild:       appBuild,
		FirstSeenAt:    firstSeenAt.UTC(),
		LastSeenAt:     lastSeenAt.UTC(),
		RevokedAt:      utcTimePtr(revokedAt),
		CreatedAt:      createdAt.UTC(),
		UpdatedAt:      updatedAt.UTC(),
	}

	if err := installation.Validate(); err != nil {
		return nil, err
	}

	return installation, nil
}

func (i *Installation) Validate() error {
	if i == nil ||
		i.ResourceID.IsZero() ||
		i.UserID == uuid.Nil ||
		i.InstallationID.IsZero() ||
		i.FirstSeenAt.IsZero() ||
		i.LastSeenAt.IsZero() ||
		i.CreatedAt.IsZero() ||
		i.UpdatedAt.IsZero() ||
		i.LastSeenAt.Before(i.FirstSeenAt) {
		return ErrInvalidInstallation
	}

	switch i.Status {
	case InstallationStatusKnown:
		if i.RevokedAt != nil {
			return ErrInvalidInstallation
		}
	case InstallationStatusRevoked:
		if i.RevokedAt == nil {
			return ErrInvalidInstallation
		}
	default:
		return ErrInvalidInstallation
	}

	return nil
}

func (i *Installation) Revoke(now time.Time) error {
	if err := i.Validate(); err != nil {
		return err
	}
	if i.Status == InstallationStatusRevoked {
		return ErrInstallationRevoked
	}

	revokedAt := now.UTC()
	i.Status = InstallationStatusRevoked
	i.RevokedAt = &revokedAt
	i.UpdatedAt = revokedAt

	return nil
}

type LoginClassification string

const (
	LoginClassificationKnown        LoginClassification = "known"
	LoginClassificationFirst        LoginClassification = "first"
	LoginClassificationNew          LoginClassification = "new"
	LoginClassificationRevoked      LoginClassification = "revoked"
	LoginClassificationLimitReached LoginClassification = "limit_reached"
)

type LoginDecision struct {
	Classification             LoginClassification
	KnownInstallationsCount    int
	MaxKnownInstallations      int
	HasAssociatedInstallations bool
}

func NewLoginDecision(
	classification LoginClassification,
	knownInstallationsCount int,
	hasAssociatedInstallations bool,
) (*LoginDecision, error) {
	decision := &LoginDecision{
		Classification:             classification,
		KnownInstallationsCount:    knownInstallationsCount,
		MaxKnownInstallations:      MaxKnownInstallations,
		HasAssociatedInstallations: hasAssociatedInstallations,
	}

	if err := decision.Validate(); err != nil {
		return nil, err
	}

	return decision, nil
}

func (d *LoginDecision) Validate() error {
	if d == nil || d.KnownInstallationsCount < 0 || d.MaxKnownInstallations <= 0 {
		return ErrInvalidInstallation
	}

	switch d.Classification {
	case LoginClassificationKnown, LoginClassificationRevoked:
		return nil
	case LoginClassificationFirst:
		if d.HasAssociatedInstallations || d.KnownInstallationsCount != 0 {
			return ErrInvalidInstallation
		}
	case LoginClassificationNew:
		if !d.HasAssociatedInstallations || d.KnownInstallationsCount >= d.MaxKnownInstallations {
			return ErrInvalidInstallation
		}
	case LoginClassificationLimitReached:
		if !d.HasAssociatedInstallations || d.KnownInstallationsCount < d.MaxKnownInstallations {
			return ErrInvalidInstallation
		}
	default:
		return ErrInvalidInstallation
	}

	return nil
}

func utcTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}

	utc := value.UTC()
	return &utc
}

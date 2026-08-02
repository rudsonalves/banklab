package application

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type InstallationLoginClassification string

const (
	InstallationLoginKnown        InstallationLoginClassification = "known"
	InstallationLoginFirst        InstallationLoginClassification = "first"
	InstallationLoginNew          InstallationLoginClassification = "new"
	InstallationLoginRevoked      InstallationLoginClassification = "revoked"
	InstallationLoginLimitReached InstallationLoginClassification = "limit_reached"
)

const defaultMaxKnownInstallations = 3

type InstallationLoginRecordStatus string

const (
	InstallationLoginRecordKnown   InstallationLoginRecordStatus = "known"
	InstallationLoginRecordRevoked InstallationLoginRecordStatus = "revoked"
)

type InstallationLoginRecord struct {
	UserID         uuid.UUID
	InstallationID uuid.UUID
	Status         InstallationLoginRecordStatus
}

type InstallationLoginDecision struct {
	Classification             InstallationLoginClassification
	KnownInstallationsCount    int
	HasAssociatedInstallations bool
}

type InstallationLoginRepository interface {
	FindByUserIDAndInstallationID(
		ctx context.Context,
		userID uuid.UUID,
		installationID uuid.UUID,
	) (*InstallationLoginRecord, error)
	CountKnownByUserID(ctx context.Context, userID uuid.UUID) (int, error)
	HasAnyByUserID(ctx context.Context, userID uuid.UUID) (bool, error)
}

type InstallationLoginClassifier interface {
	Classify(
		ctx context.Context,
		userID uuid.UUID,
		installationID uuid.UUID,
	) (*InstallationLoginDecision, error)
}

type DefaultInstallationLoginClassifier struct {
	repo                  InstallationLoginRepository
	maxKnownInstallations int
}

func NewDefaultInstallationLoginClassifier(repo InstallationLoginRepository) *DefaultInstallationLoginClassifier {
	return &DefaultInstallationLoginClassifier{
		repo:                  repo,
		maxKnownInstallations: defaultMaxKnownInstallations,
	}
}

func (c *DefaultInstallationLoginClassifier) Classify(
	ctx context.Context,
	userID uuid.UUID,
	installationID uuid.UUID,
) (*InstallationLoginDecision, error) {
	if userID == uuid.Nil || installationID == uuid.Nil {
		return nil, authdomain.ErrInvalidData
	}
	if c == nil || c.repo == nil {
		return nil, fmt.Errorf("installation login repository not configured")
	}

	record, err := c.repo.FindByUserIDAndInstallationID(ctx, userID, installationID)
	if err != nil {
		return nil, fmt.Errorf("find installation by user and installation id: %w", err)
	}
	if record != nil {
		switch record.Status {
		case InstallationLoginRecordKnown:
			return &InstallationLoginDecision{
				Classification: InstallationLoginKnown,
			}, nil
		case InstallationLoginRecordRevoked:
			return &InstallationLoginDecision{
				Classification: InstallationLoginRevoked,
			}, nil
		default:
			return nil, authdomain.ErrInvalidData
		}
	}

	hasAny, err := c.repo.HasAnyByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("check if user has any installation: %w", err)
	}
	if !hasAny {
		return &InstallationLoginDecision{
			Classification:             InstallationLoginFirst,
			HasAssociatedInstallations: false,
		}, nil
	}

	knownCount, err := c.repo.CountKnownByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("count known installations by user: %w", err)
	}
	if knownCount >= c.maxKnownInstallations {
		return &InstallationLoginDecision{
			Classification:             InstallationLoginLimitReached,
			KnownInstallationsCount:    knownCount,
			HasAssociatedInstallations: true,
		}, nil
	}

	return &InstallationLoginDecision{
		Classification:             InstallationLoginNew,
		KnownInstallationsCount:    knownCount,
		HasAssociatedInstallations: true,
	}, nil
}

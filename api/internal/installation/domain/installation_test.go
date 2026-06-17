package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewKnownInstallation(t *testing.T) {
	userID := uuid.New()
	resourceID := mustResourceID(t)
	installationID := mustInstallationID(t)
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)

	installation, err := NewKnownInstallation(userID, resourceID, installationID, now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if installation.UserID != userID {
		t.Fatalf("expected userID %q, got %q", userID, installation.UserID)
	}
	if installation.ResourceID != resourceID {
		t.Fatalf("expected resourceID %q, got %q", resourceID, installation.ResourceID)
	}
	if installation.InstallationID != installationID {
		t.Fatalf("expected installationID %q, got %q", installationID, installation.InstallationID)
	}
	if installation.Status != InstallationStatusKnown {
		t.Fatalf("expected status %q, got %q", InstallationStatusKnown, installation.Status)
	}
	if installation.RevokedAt != nil {
		t.Fatalf("expected revoked_at nil, got %v", installation.RevokedAt)
	}
}

func TestInstallation_Revoke(t *testing.T) {
	installation, err := NewKnownInstallation(uuid.New(), mustResourceID(t), mustInstallationID(t), time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	revokedAt := time.Now().Add(time.Minute)

	if err := installation.Revoke(revokedAt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if installation.Status != InstallationStatusRevoked {
		t.Fatalf("expected revoked status, got %q", installation.Status)
	}
	if installation.RevokedAt == nil || !installation.RevokedAt.Equal(revokedAt.UTC()) {
		t.Fatalf("expected revoked_at %v, got %v", revokedAt.UTC(), installation.RevokedAt)
	}

	if err := installation.Revoke(revokedAt.Add(time.Minute)); !errors.Is(err, ErrInstallationRevoked) {
		t.Fatalf("expected ErrInstallationRevoked, got %v", err)
	}
}

func TestRestoreInstallationValidation(t *testing.T) {
	now := time.Now().UTC()
	revokedAt := now.Add(time.Minute)

	tests := []struct {
		name      string
		status    InstallationStatus
		revokedAt *time.Time
		wantErr   bool
	}{
		{name: "known", status: InstallationStatusKnown},
		{name: "revoked", status: InstallationStatusRevoked, revokedAt: &revokedAt},
		{name: "known with revoked_at", status: InstallationStatusKnown, revokedAt: &revokedAt, wantErr: true},
		{name: "revoked without revoked_at", status: InstallationStatusRevoked, wantErr: true},
		{name: "trusted is not valid", status: InstallationStatus("trusted"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RestoreInstallation(
				uuid.Nil,
				mustResourceID(t),
				uuid.New(),
				mustInstallationID(t),
				tt.status,
				"",
				"",
				"",
				now,
				now,
				tt.revokedAt,
				now,
				now,
			)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidInstallation) {
					t.Fatalf("expected ErrInvalidInstallation, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewLoginDecision(t *testing.T) {
	tests := []struct {
		name    string
		class   LoginClassification
		count   int
		hasAny  bool
		wantErr bool
	}{
		{name: "known", class: LoginClassificationKnown},
		{name: "revoked", class: LoginClassificationRevoked},
		{name: "first", class: LoginClassificationFirst},
		{name: "new", class: LoginClassificationNew, count: 2, hasAny: true},
		{name: "limit reached", class: LoginClassificationLimitReached, count: MaxKnownInstallations, hasAny: true},
		{name: "first cannot have history", class: LoginClassificationFirst, hasAny: true, wantErr: true},
		{name: "new requires history", class: LoginClassificationNew, count: 1, wantErr: true},
		{name: "new requires vacancy", class: LoginClassificationNew, count: MaxKnownInstallations, hasAny: true, wantErr: true},
		{name: "limit reached requires max count", class: LoginClassificationLimitReached, count: MaxKnownInstallations - 1, hasAny: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := NewLoginDecision(tt.class, tt.count, tt.hasAny)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidInstallation) {
					t.Fatalf("expected ErrInvalidInstallation, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if decision.MaxKnownInstallations != MaxKnownInstallations {
				t.Fatalf("expected max %d, got %d", MaxKnownInstallations, decision.MaxKnownInstallations)
			}
		})
	}
}

func mustInstallationID(t *testing.T) InstallationID {
	t.Helper()

	id, err := NewInstallationID(uuid.New())
	if err != nil {
		t.Fatalf("expected installation id, got %v", err)
	}
	return id
}

func mustResourceID(t *testing.T) InstallationResourceID {
	t.Helper()

	id, err := NewInstallationResourceID(uuid.New())
	if err != nil {
		t.Fatalf("expected resource id, got %v", err)
	}
	return id
}

package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewRestrictedAuthorization(t *testing.T) {
	now := time.Date(2026, 6, 17, 10, 0, 0, 0, time.UTC)
	installationID := mustInstallationID(t)

	authorization, err := NewRestrictedAuthorization(" jti-1 ", uuid.New(), installationID, now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if authorization.JTI != "jti-1" {
		t.Fatalf("expected trimmed jti, got %q", authorization.JTI)
	}
	if authorization.Scope != RestrictedAuthorizationScopeInstallationRegister {
		t.Fatalf("expected scope %q, got %q", RestrictedAuthorizationScopeInstallationRegister, authorization.Scope)
	}
	if authorization.Status != RestrictedAuthorizationStatusActive {
		t.Fatalf("expected active status, got %q", authorization.Status)
	}
	if !authorization.ExpiresAt.Equal(now.Add(RestrictedAuthorizationDefaultDuration)) {
		t.Fatalf("expected default expiration, got %v", authorization.ExpiresAt)
	}
}

func TestRestrictedAuthorization_Consume(t *testing.T) {
	now := time.Now().UTC()
	authorization, err := NewRestrictedAuthorization("jti-1", uuid.New(), mustInstallationID(t), now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	consumedAt := now.Add(time.Minute)

	if err := authorization.Consume(consumedAt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if authorization.Status != RestrictedAuthorizationStatusConsumed {
		t.Fatalf("expected consumed status, got %q", authorization.Status)
	}
	if authorization.ConsumedAt == nil || !authorization.ConsumedAt.Equal(consumedAt.UTC()) {
		t.Fatalf("expected consumed_at %v, got %v", consumedAt.UTC(), authorization.ConsumedAt)
	}

	if err := authorization.Consume(consumedAt.Add(time.Minute)); !errors.Is(err, ErrRestrictedAuthorizationConsumed) {
		t.Fatalf("expected ErrRestrictedAuthorizationConsumed, got %v", err)
	}
}

func TestRestrictedAuthorization_ConsumeExpired(t *testing.T) {
	now := time.Now().UTC()
	authorization, err := NewRestrictedAuthorization("jti-1", uuid.New(), mustInstallationID(t), now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = authorization.Consume(now.Add(RestrictedAuthorizationDefaultDuration + time.Second))
	if !errors.Is(err, ErrRestrictedAuthorizationExpired) {
		t.Fatalf("expected ErrRestrictedAuthorizationExpired, got %v", err)
	}
}

func TestRestrictedAuthorization_Revoke(t *testing.T) {
	authorization, err := NewRestrictedAuthorization("jti-1", uuid.New(), mustInstallationID(t), time.Now())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if err := authorization.Revoke(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if authorization.Status != RestrictedAuthorizationStatusRevoked {
		t.Fatalf("expected revoked status, got %q", authorization.Status)
	}
	if err := authorization.Revoke(); !errors.Is(err, ErrRestrictedAuthorizationRevoked) {
		t.Fatalf("expected ErrRestrictedAuthorizationRevoked, got %v", err)
	}
}

func TestRestoreRestrictedAuthorizationValidation(t *testing.T) {
	now := time.Now().UTC()
	consumedAt := now.Add(time.Minute)

	tests := []struct {
		name       string
		jti        string
		scope      string
		status     RestrictedAuthorizationStatus
		consumedAt *time.Time
		wantErr    bool
	}{
		{
			name:   "active",
			jti:    "jti-1",
			scope:  RestrictedAuthorizationScopeInstallationRegister,
			status: RestrictedAuthorizationStatusActive,
		},
		{
			name:       "consumed",
			jti:        "jti-1",
			scope:      RestrictedAuthorizationScopeInstallationRegister,
			status:     RestrictedAuthorizationStatusConsumed,
			consumedAt: &consumedAt,
		},
		{
			name:   "revoked",
			jti:    "jti-1",
			scope:  RestrictedAuthorizationScopeInstallationRegister,
			status: RestrictedAuthorizationStatusRevoked,
		},
		{
			name:    "blank jti",
			scope:   RestrictedAuthorizationScopeInstallationRegister,
			status:  RestrictedAuthorizationStatusActive,
			wantErr: true,
		},
		{
			name:    "invalid scope",
			jti:     "jti-1",
			scope:   "other.scope",
			status:  RestrictedAuthorizationStatusActive,
			wantErr: true,
		},
		{
			name:       "active with consumed_at",
			jti:        "jti-1",
			scope:      RestrictedAuthorizationScopeInstallationRegister,
			status:     RestrictedAuthorizationStatusActive,
			consumedAt: &consumedAt,
			wantErr:    true,
		},
		{
			name:    "consumed without consumed_at",
			jti:     "jti-1",
			scope:   RestrictedAuthorizationScopeInstallationRegister,
			status:  RestrictedAuthorizationStatusConsumed,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RestoreRestrictedAuthorization(
				uuid.Nil,
				tt.jti,
				uuid.New(),
				mustInstallationID(t),
				tt.scope,
				tt.status,
				now.Add(time.Minute),
				tt.consumedAt,
				now,
			)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidRestrictedAuthorization) {
					t.Fatalf("expected ErrInvalidRestrictedAuthorization, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeStepUpTokenVerifier struct {
	claims *VerifiedStepUpTokenClaims
	err    error
}

func (f *fakeStepUpTokenVerifier) Verify(token string) (*VerifiedStepUpTokenClaims, error) {
	if token == "" {
		return nil, ErrStepUpTokenRequired
	}
	if f.err != nil {
		return nil, f.err
	}

	return f.claims, nil
}

func TestNewVerifiedStepUpTokenClaims_Valid(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	claims, err := NewVerifiedStepUpTokenClaims(
		uuid.MustParse("00000000-0000-0000-0000-000000000123"),
		" internal_transfer.create ",
		" step_up ",
		" step-up-jti ",
		now.Add(2*time.Minute),
		now,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if claims.Scope != StepUpTokenScope {
		t.Fatalf("expected scope %q, got %q", StepUpTokenScope, claims.Scope)
	}
	if claims.EndpointKey != "internal_transfer.create" {
		t.Fatalf("expected normalized endpoint key, got %q", claims.EndpointKey)
	}
	if claims.JTI != "step-up-jti" {
		t.Fatalf("expected normalized jti, got %q", claims.JTI)
	}
}

func TestNewVerifiedStepUpTokenClaims_Invalid(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		userID      uuid.UUID
		endpointKey string
		scope       string
		jti         string
		expiresAt   time.Time
		issuedAt    time.Time
	}{
		{
			name:        "missing user id",
			userID:      uuid.Nil,
			endpointKey: "internal_transfer.create",
			scope:       StepUpTokenScope,
			jti:         "step-up-jti",
			expiresAt:   now.Add(2 * time.Minute),
			issuedAt:    now,
		},
		{
			name:        "missing endpoint key",
			userID:      uuid.New(),
			endpointKey: " ",
			scope:       StepUpTokenScope,
			jti:         "step-up-jti",
			expiresAt:   now.Add(2 * time.Minute),
			issuedAt:    now,
		},
		{
			name:        "invalid scope",
			userID:      uuid.New(),
			endpointKey: "internal_transfer.create",
			scope:       "signin",
			jti:         "step-up-jti",
			expiresAt:   now.Add(2 * time.Minute),
			issuedAt:    now,
		},
		{
			name:        "missing jti",
			userID:      uuid.New(),
			endpointKey: "internal_transfer.create",
			scope:       StepUpTokenScope,
			jti:         "",
			expiresAt:   now.Add(2 * time.Minute),
			issuedAt:    now,
		},
		{
			name:        "expires before issued",
			userID:      uuid.New(),
			endpointKey: "internal_transfer.create",
			scope:       StepUpTokenScope,
			jti:         "step-up-jti",
			expiresAt:   now,
			issuedAt:    now,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := NewVerifiedStepUpTokenClaims(
				tt.userID,
				tt.endpointKey,
				tt.scope,
				tt.jti,
				tt.expiresAt,
				tt.issuedAt,
			)
			if !errors.Is(err, ErrInvalidStepUpToken) {
				t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
			}
			if claims != nil {
				t.Fatalf("expected nil claims, got %+v", claims)
			}
		})
	}
}

func TestStepUpTokenVerifierContract_AllowsFakeScenariosWithoutJWT(t *testing.T) {
	validClaims, err := NewVerifiedStepUpTokenClaims(
		uuid.New(),
		"internal_transfer.create",
		StepUpTokenScope,
		"step-up-jti",
		time.Now().UTC().Add(time.Minute),
		time.Now().UTC(),
	)
	if err != nil {
		t.Fatalf("expected valid claims, got %v", err)
	}

	verifier := &fakeStepUpTokenVerifier{claims: validClaims}

	claims, err := verifier.Verify("signed-step-up-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if claims == nil {
		t.Fatal("expected claims, got nil")
	}

	_, err = verifier.Verify("")
	if !errors.Is(err, ErrStepUpTokenRequired) {
		t.Fatalf("expected ErrStepUpTokenRequired, got %v", err)
	}

	verifier.err = ErrInvalidStepUpToken
	_, err = verifier.Verify("invalid-token")
	if !errors.Is(err, ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
	}

	verifier.err = ErrStepUpTokenExpired
	_, err = verifier.Verify("expired-token")
	if !errors.Is(err, ErrStepUpTokenExpired) {
		t.Fatalf("expected ErrStepUpTokenExpired, got %v", err)
	}

	verifier.err = ErrStepUpTokenConsumed
	_, err = verifier.Verify("consumed-token")
	if !errors.Is(err, ErrStepUpTokenConsumed) {
		t.Fatalf("expected ErrStepUpTokenConsumed, got %v", err)
	}
}

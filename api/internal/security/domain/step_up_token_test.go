package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewStepUpToken(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	userID := uuid.New()

	token, err := NewStepUpToken(" jti-value ", userID, " internal_transfer.create ", now)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token.ID != uuid.Nil {
		t.Fatalf("expected empty ID before persistence, got %v", token.ID)
	}

	if token.JTI != "jti-value" {
		t.Fatalf("expected trimmed jti, got %q", token.JTI)
	}

	if token.UserID != userID {
		t.Fatalf("expected user id %v, got %v", userID, token.UserID)
	}

	if token.EndpointKey != "internal_transfer.create" {
		t.Fatalf("expected trimmed endpoint key, got %q", token.EndpointKey)
	}

	if token.Status != StepUpTokenActive {
		t.Fatalf("expected active status, got %s", token.Status)
	}

	if !token.CreatedAt.Equal(now) {
		t.Fatalf("expected created_at %v, got %v", now, token.CreatedAt)
	}

	expectedExpiresAt := now.Add(StepUpTokenDefaultDuration)
	if !token.ExpiresAt.Equal(expectedExpiresAt) {
		t.Fatalf("expected expires_at %v, got %v", expectedExpiresAt, token.ExpiresAt)
	}

	if token.ConsumedAt != nil {
		t.Fatalf("expected nil consumed_at, got %v", token.ConsumedAt)
	}
}

func TestNewStepUpToken_InvalidData(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		jti         string
		userID      uuid.UUID
		endpointKey string
		now         time.Time
	}{
		{name: "empty jti", jti: "", userID: uuid.New(), endpointKey: "internal_transfer.create", now: now},
		{name: "blank jti", jti: "   ", userID: uuid.New(), endpointKey: "internal_transfer.create", now: now},
		{name: "empty user id", jti: "jti-value", userID: uuid.Nil, endpointKey: "internal_transfer.create", now: now},
		{name: "empty endpoint key", jti: "jti-value", userID: uuid.New(), endpointKey: "", now: now},
		{name: "blank endpoint key", jti: "jti-value", userID: uuid.New(), endpointKey: "   ", now: now},
		{name: "empty time", jti: "jti-value", userID: uuid.New(), endpointKey: "internal_transfer.create", now: time.Time{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewStepUpToken(tt.jti, tt.userID, tt.endpointKey, tt.now)

			if !errors.Is(err, ErrInvalidStepUpToken) {
				t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
			}

			if token != nil {
				t.Fatalf("expected nil token, got %+v", token)
			}
		})
	}
}

func TestRestoreStepUpToken_ValidatesPersistedState(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	consumedAt := now.Add(time.Minute)

	tests := []struct {
		name       string
		status     StepUpTokenStatus
		expiresAt  time.Time
		consumedAt *time.Time
	}{
		{
			name:      "active",
			status:    StepUpTokenActive,
			expiresAt: now.Add(StepUpTokenDefaultDuration),
		},
		{
			name:       "consumed",
			status:     StepUpTokenConsumed,
			expiresAt:  now.Add(StepUpTokenDefaultDuration),
			consumedAt: &consumedAt,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := RestoreStepUpToken(
				uuid.New(),
				"jti-value",
				uuid.New(),
				"internal_transfer.create",
				tt.status,
				tt.expiresAt,
				tt.consumedAt,
				now,
			)

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if token.Status != tt.status {
				t.Fatalf("expected status %s, got %s", tt.status, token.Status)
			}
		})
	}
}

func TestRestoreStepUpToken_InvalidPersistedState(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	consumedAt := now.Add(time.Minute)

	tests := []struct {
		name       string
		status     StepUpTokenStatus
		expiresAt  time.Time
		consumedAt *time.Time
	}{
		{
			name:      "invalid status",
			status:    StepUpTokenStatus("expired"),
			expiresAt: now.Add(StepUpTokenDefaultDuration),
		},
		{
			name:       "active with consumed_at",
			status:     StepUpTokenActive,
			expiresAt:  now.Add(StepUpTokenDefaultDuration),
			consumedAt: &consumedAt,
		},
		{
			name:      "consumed without consumed_at",
			status:    StepUpTokenConsumed,
			expiresAt: now.Add(StepUpTokenDefaultDuration),
		},
		{
			name:      "expires before created",
			status:    StepUpTokenActive,
			expiresAt: now.Add(-time.Minute),
		},
		{
			name:      "expires equals created",
			status:    StepUpTokenActive,
			expiresAt: now,
		},
		{
			name:      "empty expires_at",
			status:    StepUpTokenActive,
			expiresAt: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := RestoreStepUpToken(
				uuid.New(),
				"jti-value",
				uuid.New(),
				"internal_transfer.create",
				tt.status,
				tt.expiresAt,
				tt.consumedAt,
				now,
			)

			if !errors.Is(err, ErrInvalidStepUpToken) {
				t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
			}

			if token != nil {
				t.Fatalf("expected nil token, got %+v", token)
			}
		})
	}
}

func TestStepUpToken_IsExpiredRequiresActiveStatusAndExpiresAtBeforeNow(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	consumedAt := now.Add(-time.Minute)

	tests := []struct {
		name       string
		status     StepUpTokenStatus
		expiresAt  time.Time
		consumedAt *time.Time
		want       bool
	}{
		{name: "active before now", status: StepUpTokenActive, expiresAt: now.Add(-time.Nanosecond), want: true},
		{name: "active equal now", status: StepUpTokenActive, expiresAt: now, want: false},
		{name: "active after now", status: StepUpTokenActive, expiresAt: now.Add(time.Nanosecond), want: false},
		{
			name:       "consumed before now",
			status:     StepUpTokenConsumed,
			expiresAt:  now.Add(-time.Nanosecond),
			consumedAt: &consumedAt,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := RestoreStepUpToken(
				uuid.New(),
				"jti-value",
				uuid.New(),
				"internal_transfer.create",
				tt.status,
				tt.expiresAt,
				tt.consumedAt,
				now.Add(-StepUpTokenDefaultDuration),
			)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if got := token.IsExpired(now); got != tt.want {
				t.Fatalf("expected expired=%v, got %v", tt.want, got)
			}
		})
	}
}

func TestStepUpToken_Consume(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	token, err := NewStepUpToken("jti-value", uuid.New(), "internal_transfer.create", now)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	consumedAt := now.Add(time.Minute)
	err = token.Consume(consumedAt)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token.Status != StepUpTokenConsumed {
		t.Fatalf("expected consumed status, got %s", token.Status)
	}

	if token.ConsumedAt == nil {
		t.Fatal("expected consumed_at to be set")
	}

	if !token.ConsumedAt.Equal(consumedAt) {
		t.Fatalf("expected consumed_at %v, got %v", consumedAt, *token.ConsumedAt)
	}
}

func TestStepUpToken_ConsumeConsumedTokenFails(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	consumedAt := now.Add(time.Minute)
	token, err := RestoreStepUpToken(
		uuid.New(),
		"jti-value",
		uuid.New(),
		"internal_transfer.create",
		StepUpTokenConsumed,
		now.Add(StepUpTokenDefaultDuration),
		&consumedAt,
		now,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = token.Consume(now.Add(90 * time.Second))

	if !errors.Is(err, ErrStepUpTokenConsumed) {
		t.Fatalf("expected ErrStepUpTokenConsumed, got %v", err)
	}
}

func TestStepUpToken_ConsumeExpiredTokenFails(t *testing.T) {
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
	token, err := RestoreStepUpToken(
		uuid.New(),
		"jti-value",
		uuid.New(),
		"internal_transfer.create",
		StepUpTokenActive,
		now.Add(time.Minute),
		nil,
		now,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = token.Consume(now.Add(2 * time.Minute))

	if !errors.Is(err, ErrStepUpTokenExpired) {
		t.Fatalf("expected ErrStepUpTokenExpired, got %v", err)
	}

	if token.Status != StepUpTokenActive {
		t.Fatalf("expected token to remain active, got %s", token.Status)
	}

	if token.ConsumedAt != nil {
		t.Fatalf("expected consumed_at to remain nil, got %v", token.ConsumedAt)
	}
}

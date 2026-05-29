package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type StepUpTokenStatus string

const (
	StepUpTokenActive   StepUpTokenStatus = "active"
	StepUpTokenConsumed StepUpTokenStatus = "consumed"
	StepUpTokenScope                      = "step_up"

	StepUpTokenDefaultDuration = 2 * time.Minute
)

type StepUpToken struct {
	ID          uuid.UUID
	JTI         string
	UserID      uuid.UUID
	EndpointKey string
	Status      StepUpTokenStatus
	ExpiresAt   time.Time
	ConsumedAt  *time.Time
	CreatedAt   time.Time
}

// NewStepUpToken creates a new StepUpToken instance with the provided parameters. It
// trims whitespace from the JTI and endpoint key, sets the status to active, and
// calculates the expiration time based on the current time and the default duration.
// If any validation rule is violated, it returns an error indicating that the token
// is invalid.
func NewStepUpToken(
	jti string,
	userID uuid.UUID,
	endpointKey string,
	now time.Time,
) (*StepUpToken, error) {
	createdAt := now.UTC()

	return RestoreStepUpToken(
		uuid.Nil,
		jti,
		userID,
		endpointKey,
		StepUpTokenActive,
		createdAt.Add(StepUpTokenDefaultDuration),
		nil,
		createdAt,
	)
}

// RestoreStepUpToken creates a StepUpToken instance from the provided parameters. It
// trims whitespace from the JTI and endpoint key, converts timestamps to UTC, and
// validates the token's integrity. If any validation rule is violated, it returns
// an error indicating that the token is invalid.
func RestoreStepUpToken(
	id uuid.UUID,
	jti string,
	userID uuid.UUID,
	endpointKey string,
	status StepUpTokenStatus,
	expiresAt time.Time,
	consumedAt *time.Time,
	createdAt time.Time,
) (*StepUpToken, error) {
	token := &StepUpToken{
		ID:          id,
		JTI:         strings.TrimSpace(jti),
		UserID:      userID,
		EndpointKey: strings.TrimSpace(endpointKey),
		Status:      status,
		ExpiresAt:   expiresAt.UTC(),
		ConsumedAt:  utcTimePtr(consumedAt),
		CreatedAt:   createdAt.UTC(),
	}

	if err := token.Validate(); err != nil {
		return nil, err
	}

	return token, nil
}

// Validate checks the integrity of the StepUpToken fields. It ensures that all
// required fields are present and valid, and that the status and timestamps are
// consistent. If any validation rule is violated, it returns an error indicating
// that the token is invalid.
func (t *StepUpToken) Validate() error {
	if t == nil {
		return ErrInvalidStepUpToken
	}

	if t.JTI == "" ||
		t.UserID == uuid.Nil ||
		t.EndpointKey == "" ||
		t.ExpiresAt.IsZero() ||
		t.CreatedAt.IsZero() ||
		!t.ExpiresAt.After(t.CreatedAt) {
		return ErrInvalidStepUpToken
	}

	switch t.Status {
	case StepUpTokenActive:
		if t.ConsumedAt != nil {
			return ErrInvalidStepUpToken
		}
	case StepUpTokenConsumed:
		if t.ConsumedAt == nil {
			return ErrInvalidStepUpToken
		}
	default:
		return ErrInvalidStepUpToken
	}

	return nil
}

// IsExpired checks if the token is expired based on the current time. It returns
// true if the token is active and the current time is after the expiration time,
// or if the token is nil. Otherwise, it returns false.
func (t *StepUpToken) IsExpired(now time.Time) bool {
	if t == nil {
		return true
	}

	return t.Status == StepUpTokenActive && t.ExpiresAt.Before(now.UTC())
}

// Consume marks the token as consumed if it is valid, active, and
// not expired. It returns an error if the token is invalid, already consumed,
// or expired.
func (t *StepUpToken) Consume(now time.Time) error {
	if err := t.Validate(); err != nil {
		return err
	}

	if t.Status == StepUpTokenConsumed {
		return ErrStepUpTokenConsumed
	}

	if t.IsExpired(now) {
		return ErrStepUpTokenExpired
	}

	consumedAt := now.UTC()
	t.Status = StepUpTokenConsumed
	t.ConsumedAt = &consumedAt

	return nil
}

// utcTimePtr converts a time.Time pointer to UTC. If the input
// is nil, it returns nil.
func utcTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}

	utc := value.UTC()
	return &utc
}

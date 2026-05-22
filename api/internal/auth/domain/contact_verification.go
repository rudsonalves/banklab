package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type ContactVerificationChannel string

const (
	ContactVerificationChannelEmail ContactVerificationChannel = "email"
	ContactVerificationChannelPhone ContactVerificationChannel = "phone"
)

type ContactVerification struct {
	ID                uuid.UUID
	Channel           ContactVerificationChannel
	Target            string
	Token             string
	VerificationToken *string
	VerifiedAt        *time.Time
	ExpiresAt         time.Time
	CreatedAt         time.Time
}

func NewContactVerification(
	channel ContactVerificationChannel,
	target string,
	token string,
	now time.Time,
	ttl time.Duration,
) (*ContactVerification, error) {
	channel = NormalizeContactVerificationChannel(string(channel))
	target = strings.TrimSpace(target)
	token = strings.TrimSpace(token)

	if !channel.IsValid() || target == "" || token == "" || ttl <= 0 {
		return nil, ErrInvalidData
	}

	return &ContactVerification{
		ID:        uuid.New(),
		Channel:   channel,
		Target:    target,
		Token:     token,
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}, nil
}

func NormalizeContactVerificationChannel(channel string) ContactVerificationChannel {
	return ContactVerificationChannel(strings.ToLower(strings.TrimSpace(channel)))
}

func (c ContactVerificationChannel) IsValid() bool {
	switch c {
	case ContactVerificationChannelEmail, ContactVerificationChannelPhone:
		return true
	default:
		return false
	}
}

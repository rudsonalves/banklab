package application

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

const contactVerificationTTL = 10 * time.Minute

type RequestContactVerificationUseCase struct {
	repo     domain.ContactVerificationRepository
	userRepo domain.UserRepository
	now      func() time.Time
}

func NewRequestContactVerificationUseCase(
	repo domain.ContactVerificationRepository,
	userRepo domain.UserRepository,
) *RequestContactVerificationUseCase {
	return &RequestContactVerificationUseCase{
		repo:     repo,
		userRepo: userRepo,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

type RequestContactVerificationInput struct {
	Channel string
	Target  string
}

type RequestContactVerificationOutput struct {
	VerificationID uuid.UUID `json:"verification_id"`
	Channel        string    `json:"channel"`
	Target         string    `json:"target"`
	ExpiresAt      time.Time `json:"expires_at"`

	DebugToken *string `json:"debug_token,omitempty"`
}

func (uc *RequestContactVerificationUseCase) Execute(
	ctx context.Context,
	input RequestContactVerificationInput,
) (*RequestContactVerificationOutput, error) {
	channel := domain.NormalizeContactVerificationChannel(input.Channel)
	target := normalizeContactVerificationTarget(channel, input.Target)
	if !channel.IsValid() || target == "" {
		return nil, domain.ErrInvalidData
	}

	if uc.userRepo != nil {
		switch channel {
		case domain.ContactVerificationChannelEmail:
			exists, err := uc.userRepo.ExistsByEmail(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("check email uniqueness: %w", err)
			}
			if exists {
				return nil, domain.ErrEmailAlreadyExists
			}
		case domain.ContactVerificationChannelPhone:
			exists, err := uc.userRepo.ExistsByPhone(ctx, target)
			if err != nil {
				return nil, fmt.Errorf("check phone uniqueness: %w", err)
			}
			if exists {
				return nil, domain.ErrPhoneAlreadyExists
			}
		}
	}

	token, err := generateNumericToken(6)
	if err != nil {
		return nil, fmt.Errorf("generate verification token: %w", err)
	}

	verification, err := domain.NewContactVerification(
		channel,
		target,
		token,
		uc.now(),
		contactVerificationTTL,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.CreateContactVerification(ctx, verification); err != nil {
		return nil, fmt.Errorf("create contact verification: %w", err)
	}

	return &RequestContactVerificationOutput{
		VerificationID: verification.ID,
		Channel:        string(verification.Channel),
		Target:         verification.Target,
		DebugToken:     &verification.Token,
		ExpiresAt:      verification.ExpiresAt,
	}, nil
}

func normalizeContactVerificationTarget(
	channel domain.ContactVerificationChannel,
	target string,
) string {
	target = strings.TrimSpace(target)
	if channel == domain.ContactVerificationChannelEmail {
		return strings.ToLower(target)
	}

	return target
}

type ConfirmContactVerificationUseCase struct {
	repo domain.ContactVerificationRepository
	now  func() time.Time
}

func NewConfirmContactVerificationUseCase(
	repo domain.ContactVerificationRepository,
) *ConfirmContactVerificationUseCase {
	return &ConfirmContactVerificationUseCase{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

type ConfirmContactVerificationInput struct {
	VerificationID uuid.UUID
	Token          string
}

type ConfirmContactVerificationOutput struct {
	VerificationToken string    `json:"verification_token"`
	Channel           string    `json:"channel"`
	Target            string    `json:"target"`
	VerifiedAt        time.Time `json:"verified_at"`
}

func (uc *ConfirmContactVerificationUseCase) Execute(
	ctx context.Context,
	input ConfirmContactVerificationInput,
) (*ConfirmContactVerificationOutput, error) {
	token := strings.TrimSpace(input.Token)
	if input.VerificationID == uuid.Nil || token == "" {
		return nil, domain.ErrInvalidData
	}

	verification, err := uc.repo.FindContactVerificationByID(ctx, input.VerificationID)
	if err != nil {
		return nil, err
	}
	if verification == nil {
		return nil, domain.ErrContactVerificationNotFound
	}
	if verification.Token != token {
		return nil, domain.ErrInvalidVerificationToken
	}

	now := uc.now()
	if now.After(verification.ExpiresAt) {
		return nil, domain.ErrContactVerificationExpired
	}

	verificationToken := uuid.NewString()
	if err := uc.repo.ConfirmContactVerification(ctx, verification.ID, verificationToken, now); err != nil {
		return nil, fmt.Errorf("confirm contact verification: %w", err)
	}

	return &ConfirmContactVerificationOutput{
		VerificationToken: verificationToken,
		Channel:           string(verification.Channel),
		Target:            verification.Target,
		VerifiedAt:        now,
	}, nil
}

func generateNumericToken(length int) (string, error) {
	if length <= 0 {
		return "", domain.ErrInvalidData
	}

	token := make([]byte, length)
	for i := range token {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		token[i] = byte('0' + n.Int64())
	}

	return string(token), nil
}

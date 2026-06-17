package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/installation/domain"
)

type JWTRestrictedAccessTokenService struct {
	secret []byte
	repo   domain.RestrictedAuthorizationRepository
	now    func() time.Time
}

type restrictedAccessJWTClaims struct {
	TokenType      string `json:"token_type"`
	Scope          string `json:"scope"`
	InstallationID string `json:"installation_id"`
	jwt.RegisteredClaims
}

var _ domain.RestrictedAccessTokenSigner = (*JWTRestrictedAccessTokenService)(nil)
var _ domain.RestrictedAccessTokenVerifier = (*JWTRestrictedAccessTokenService)(nil)

func NewJWTRestrictedAccessTokenService(
	secret string,
	repo domain.RestrictedAuthorizationRepository,
) *JWTRestrictedAccessTokenService {
	return &JWTRestrictedAccessTokenService{
		secret: []byte(secret),
		repo:   repo,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

func (s *JWTRestrictedAccessTokenService) SignRestrictedAccessToken(
	claims *domain.RestrictedAccessTokenClaims,
) (string, error) {
	if err := claims.Validate(); err != nil {
		return "", err
	}

	payload := restrictedAccessJWTClaims{
		TokenType:      claims.TokenType,
		Scope:          claims.Scope,
		InstallationID: claims.InstallationID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UserID.String(),
			ID:        claims.JTI,
			IssuedAt:  jwt.NewNumericDate(claims.IssuedAt.UTC()),
			ExpiresAt: jwt.NewNumericDate(claims.ExpiresAt.UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(s.secret)
}

func (s *JWTRestrictedAccessTokenService) VerifyRestrictedAccessToken(
	ctx context.Context,
	token string,
) (*domain.RestrictedAccessTokenClaims, error) {
	claims, err := s.parse(token)
	if err != nil {
		return nil, err
	}
	if s.repo == nil {
		return nil, domain.ErrRestrictedAuthorizationInvalid
	}

	authorization, err := s.repo.FindByJTI(ctx, claims.JTI)
	if err != nil {
		return nil, err
	}
	if err := validateRestrictedTokenAuthorization(claims, authorization, s.now()); err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *JWTRestrictedAccessTokenService) parse(token string) (*domain.RestrictedAccessTokenClaims, error) {
	parsedClaims := &restrictedAccessJWTClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, parsedClaims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid token signing method")
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, domain.ErrRestrictedAuthorizationInvalid
	}
	if !parsedToken.Valid {
		return nil, domain.ErrRestrictedAuthorizationInvalid
	}
	if parsedClaims.Subject == "" || parsedClaims.ID == "" || parsedClaims.IssuedAt == nil || parsedClaims.ExpiresAt == nil {
		return nil, domain.ErrRestrictedAuthorizationInvalid
	}

	userID, err := uuid.Parse(parsedClaims.Subject)
	if err != nil {
		return nil, domain.ErrRestrictedAuthorizationInvalid
	}
	installationID, err := domain.ParseInstallationID(parsedClaims.InstallationID)
	if err != nil {
		return nil, domain.ErrRestrictedAuthorizationInvalid
	}

	claims := &domain.RestrictedAccessTokenClaims{
		UserID:         userID,
		InstallationID: installationID,
		JTI:            parsedClaims.ID,
		TokenType:      parsedClaims.TokenType,
		Scope:          parsedClaims.Scope,
		IssuedAt:       parsedClaims.IssuedAt.Time.UTC(),
		ExpiresAt:      parsedClaims.ExpiresAt.Time.UTC(),
	}
	if err := claims.Validate(); err != nil {
		return nil, err
	}

	return claims, nil
}

func validateRestrictedTokenAuthorization(
	claims *domain.RestrictedAccessTokenClaims,
	authorization *domain.RestrictedAuthorization,
	now time.Time,
) error {
	if authorization == nil {
		return domain.ErrRestrictedAuthorizationNotFound
	}
	if authorization.Status == domain.RestrictedAuthorizationStatusConsumed {
		return domain.ErrRestrictedAuthorizationConsumed
	}
	if authorization.Status == domain.RestrictedAuthorizationStatusRevoked {
		return domain.ErrRestrictedAuthorizationRevoked
	}
	if authorization.Status != domain.RestrictedAuthorizationStatusActive {
		return domain.ErrRestrictedAuthorizationInvalid
	}
	if authorization.IsExpired(now) {
		return domain.ErrRestrictedAuthorizationExpired
	}
	if authorization.UserID != claims.UserID ||
		authorization.InstallationID.UUID() != claims.InstallationID.UUID() ||
		authorization.JTI != claims.JTI ||
		authorization.Scope != claims.Scope {
		return domain.ErrRestrictedAuthorizationInvalid
	}
	if authorization.ExpiresAt.Unix() != claims.ExpiresAt.Unix() {
		return fmt.Errorf("%w: token expiration does not match authorization", domain.ErrRestrictedAuthorizationInvalid)
	}

	return nil
}

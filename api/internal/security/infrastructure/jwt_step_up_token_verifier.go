package infrastructure

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type JWTStepUpTokenVerifier struct {
	secret []byte
}

var _ domain.StepUpTokenVerifier = (*JWTStepUpTokenVerifier)(nil)

func NewJWTStepUpTokenVerifier(secret string) *JWTStepUpTokenVerifier {
	return &JWTStepUpTokenVerifier{
		secret: []byte(secret),
	}
}

func (v *JWTStepUpTokenVerifier) Verify(rawToken string) (*domain.VerifiedStepUpTokenClaims, error) {
	if strings.TrimSpace(rawToken) == "" {
		return nil, domain.ErrInvalidStepUpToken
	}

	claims := &stepUpTokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, domain.ErrInvalidStepUpToken
		}

		return v.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrStepUpTokenExpired
		}

		return nil, domain.ErrInvalidStepUpToken
	}

	if parsedToken == nil || !parsedToken.Valid {
		return nil, domain.ErrInvalidStepUpToken
	}

	if claims.ExpiresAt == nil || claims.IssuedAt == nil {
		return nil, domain.ErrInvalidStepUpToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidStepUpToken
	}

	verifiedClaims, err := domain.NewVerifiedStepUpTokenClaims(
		userID,
		claims.EndpointKey,
		claims.Scope,
		claims.ID,
		claims.ExpiresAt.Time,
		claims.IssuedAt.Time,
	)
	if err != nil {
		return nil, domain.ErrInvalidStepUpToken
	}

	return verifiedClaims, nil
}

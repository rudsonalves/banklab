package infrastructure

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

type JWTStepUpTokenSigner struct {
	secret []byte
}

type stepUpTokenClaims struct {
	UserID      string `json:"user_id"`
	EndpointKey string `json:"endpoint_key"`
	Scope       string `json:"scope"`
	jwt.RegisteredClaims
}

var _ domain.StepUpTokenSigner = (*JWTStepUpTokenSigner)(nil)

func NewJWTStepUpTokenSigner(secret string) *JWTStepUpTokenSigner {
	return &JWTStepUpTokenSigner{
		secret: []byte(secret),
	}
}

// Sign creates a short-lived JWT for a persisted step-up token. The token only
// carries the minimum claims required by the step-up contract and never includes
// the transaction password, password hash, or operation payload.
func (s *JWTStepUpTokenSigner) Sign(token *domain.StepUpToken) (string, error) {
	if err := token.Validate(); err != nil {
		return "", err
	}
	if token.Status != domain.StepUpTokenActive {
		return "", domain.ErrInvalidStepUpToken
	}

	payload := stepUpTokenClaims{
		UserID:      token.UserID.String(),
		EndpointKey: token.EndpointKey,
		Scope:       domain.StepUpTokenScope,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        token.JTI,
			IssuedAt:  jwt.NewNumericDate(token.CreatedAt.UTC()),
			ExpiresAt: jwt.NewNumericDate(token.ExpiresAt.UTC()),
		},
	}

	signedToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

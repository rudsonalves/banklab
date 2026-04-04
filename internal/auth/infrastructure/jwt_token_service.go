package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type JWTTokenService struct {
	secret []byte
	ttl    time.Duration
}

type jwtClaims struct {
	Role string  `json:"role"`
	CID  *string `json:"cid,omitempty"`
	jwt.RegisteredClaims
}

var _ domain.TokenService = (*JWTTokenService)(nil)

func NewJWTTokenService(secret string, ttl time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (s *JWTTokenService) GenerateAccessToken(claims domain.TokenClaims) (string, error) {
	now := time.Now().UTC()

	payload := jwtClaims{
		Role: string(claims.Role),
		CID:  claims.CustomerID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UserID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedToken, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *JWTTokenService) ParseAccessToken(token string) (*domain.TokenClaims, error) {
	parsedClaims := &jwtClaims{}

	parsedToken, err := jwt.ParseWithClaims(token, parsedClaims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("invalid token signing method")
		}

		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	if parsedClaims.Subject == "" {
		return nil, errors.New("missing subject claim")
	}

	if parsedClaims.Role == "" {
		return nil, errors.New("missing role claim")
	}

	return &domain.TokenClaims{
		UserID:     parsedClaims.Subject,
		Role:       domain.Role(parsedClaims.Role),
		CustomerID: parsedClaims.CID,
	}, nil
}

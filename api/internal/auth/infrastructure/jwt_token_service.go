package infrastructure

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type JWTTokenService struct {
	secret []byte
	ttl    time.Duration
}

type jwtClaims struct {
	Role       string  `json:"role"`
	CustomerID *string `json:"customer_id,omitempty"`
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

	var cidStr *string
	if claims.CustomerID != nil {
		s := claims.CustomerID.String()
		cidStr = &s
	}

	payload := jwtClaims{
		Role:       string(claims.Role),
		CustomerID: cidStr,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   claims.UserID.String(),
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

	userID, err := uuid.Parse(parsedClaims.Subject)
	if err != nil {
		return nil, errors.New("invalid subject claim: not a valid uuid")
	}

	var customerID *uuid.UUID
	if parsedClaims.CustomerID != nil {
		cid, err := uuid.Parse(*parsedClaims.CustomerID)
		if err != nil {
			return nil, errors.New("invalid customer_id claim: not a valid uuid")
		}
		customerID = &cid
	}

	return &domain.TokenClaims{
		UserID:     userID,
		Role:       domain.Role(parsedClaims.Role),
		CustomerID: customerID,
	}, nil
}

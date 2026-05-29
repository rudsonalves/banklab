package infrastructure

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

func TestJWTStepUpTokenVerifier_VerifiesTokenIssuedBySigner(t *testing.T) {
	secret := "step-up-secret"
	signer := NewJWTStepUpTokenSigner(secret)
	verifier := NewJWTStepUpTokenVerifier(secret)
	now := time.Now().UTC().Truncate(time.Second)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000123")
	jti := "step-up-jti"

	stepUpToken, err := domain.NewStepUpToken(
		jti,
		userID,
		domain.StepUpEndpointInternalTransferCreate,
		now,
	)
	if err != nil {
		t.Fatalf("expected no error creating step-up token, got %v", err)
	}

	signedToken, err := signer.Sign(stepUpToken)
	if err != nil {
		t.Fatalf("expected no error signing step-up token, got %v", err)
	}

	claims, err := verifier.Verify(signedToken)
	if err != nil {
		t.Fatalf("expected verifier to accept signed token, got %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("expected user_id %q, got %q", userID, claims.UserID)
	}
	if claims.EndpointKey != domain.StepUpEndpointInternalTransferCreate {
		t.Fatalf("expected endpoint_key %q, got %q", domain.StepUpEndpointInternalTransferCreate, claims.EndpointKey)
	}
	if claims.Scope != domain.StepUpTokenScope {
		t.Fatalf("expected scope %q, got %q", domain.StepUpTokenScope, claims.Scope)
	}
	if claims.JTI != jti {
		t.Fatalf("expected jti %q, got %q", jti, claims.JTI)
	}
	if !claims.IssuedAt.Equal(stepUpToken.CreatedAt) {
		t.Fatalf("expected issued_at %v, got %v", stepUpToken.CreatedAt, claims.IssuedAt)
	}
	if !claims.ExpiresAt.Equal(stepUpToken.ExpiresAt) {
		t.Fatalf("expected expires_at %v, got %v", stepUpToken.ExpiresAt, claims.ExpiresAt)
	}
}

func TestJWTStepUpTokenVerifier_RejectsTokenSignedWithDifferentSecret(t *testing.T) {
	signer := NewJWTStepUpTokenSigner("issuer-secret")
	verifier := NewJWTStepUpTokenVerifier("validator-secret")
	signedToken := signValidStepUpJWT(t, signer, time.Now().UTC())

	claims, err := verifier.Verify(signedToken)
	if !errors.Is(err, domain.ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
	}
	if claims != nil {
		t.Fatalf("expected nil claims, got %+v", claims)
	}
}

func TestJWTStepUpTokenVerifier_RejectsInvalidScope(t *testing.T) {
	secret := "step-up-secret"
	now := time.Now().UTC()
	signedToken := signStepUpJWT(t, secret, jwt.SigningMethodHS256, stepUpTokenClaims{
		UserID:      uuid.New().String(),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Scope:       "signin",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "step-up-jti",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
		},
	})

	claims, err := NewJWTStepUpTokenVerifier(secret).Verify(signedToken)
	if !errors.Is(err, domain.ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
	}
	if claims != nil {
		t.Fatalf("expected nil claims, got %+v", claims)
	}
}

func TestJWTStepUpTokenVerifier_RejectsMissingRequiredClaims(t *testing.T) {
	secret := "step-up-secret"
	now := time.Now().UTC()

	tests := []struct {
		name   string
		claims stepUpTokenClaims
	}{
		{
			name: "missing jti",
			claims: stepUpTokenClaims{
				UserID:      uuid.New().String(),
				EndpointKey: domain.StepUpEndpointInternalTransferCreate,
				Scope:       domain.StepUpTokenScope,
				RegisteredClaims: jwt.RegisteredClaims{
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				},
			},
		},
		{
			name: "missing user id",
			claims: stepUpTokenClaims{
				EndpointKey: domain.StepUpEndpointInternalTransferCreate,
				Scope:       domain.StepUpTokenScope,
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        "step-up-jti",
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				},
			},
		},
		{
			name: "missing endpoint key",
			claims: stepUpTokenClaims{
				UserID: uuid.New().String(),
				Scope:  domain.StepUpTokenScope,
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        "step-up-jti",
					IssuedAt:  jwt.NewNumericDate(now),
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				},
			},
		},
		{
			name: "missing exp",
			claims: stepUpTokenClaims{
				UserID:      uuid.New().String(),
				EndpointKey: domain.StepUpEndpointInternalTransferCreate,
				Scope:       domain.StepUpTokenScope,
				RegisteredClaims: jwt.RegisteredClaims{
					ID:       "step-up-jti",
					IssuedAt: jwt.NewNumericDate(now),
				},
			},
		},
		{
			name: "missing iat",
			claims: stepUpTokenClaims{
				UserID:      uuid.New().String(),
				EndpointKey: domain.StepUpEndpointInternalTransferCreate,
				Scope:       domain.StepUpTokenScope,
				RegisteredClaims: jwt.RegisteredClaims{
					ID:        "step-up-jti",
					ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signedToken := signStepUpJWT(t, secret, jwt.SigningMethodHS256, tt.claims)

			claims, err := NewJWTStepUpTokenVerifier(secret).Verify(signedToken)
			if !errors.Is(err, domain.ErrInvalidStepUpToken) {
				t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
			}
			if claims != nil {
				t.Fatalf("expected nil claims, got %+v", claims)
			}
		})
	}
}

func TestJWTStepUpTokenVerifier_RejectsExpiredToken(t *testing.T) {
	secret := "step-up-secret"
	now := time.Now().UTC()
	signedToken := signStepUpJWT(t, secret, jwt.SigningMethodHS256, stepUpTokenClaims{
		UserID:      uuid.New().String(),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Scope:       domain.StepUpTokenScope,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "step-up-jti",
			IssuedAt:  jwt.NewNumericDate(now.Add(-2 * time.Minute)),
			ExpiresAt: jwt.NewNumericDate(now.Add(-time.Minute)),
		},
	})

	claims, err := NewJWTStepUpTokenVerifier(secret).Verify(signedToken)
	if !errors.Is(err, domain.ErrStepUpTokenExpired) {
		t.Fatalf("expected ErrStepUpTokenExpired, got %v", err)
	}
	if claims != nil {
		t.Fatalf("expected nil claims, got %+v", claims)
	}
}

func TestJWTStepUpTokenVerifier_RejectsInvalidAlgorithm(t *testing.T) {
	secret := "step-up-secret"
	now := time.Now().UTC()
	signedToken := signStepUpJWT(t, secret, jwt.SigningMethodHS384, stepUpTokenClaims{
		UserID:      uuid.New().String(),
		EndpointKey: domain.StepUpEndpointInternalTransferCreate,
		Scope:       domain.StepUpTokenScope,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        "step-up-jti",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
		},
	})

	claims, err := NewJWTStepUpTokenVerifier(secret).Verify(signedToken)
	if !errors.Is(err, domain.ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
	}
	if claims != nil {
		t.Fatalf("expected nil claims, got %+v", claims)
	}
}

func TestJWTStepUpTokenVerifier_RejectsMalformedToken(t *testing.T) {
	claims, err := NewJWTStepUpTokenVerifier("step-up-secret").Verify("not-a-jwt")
	if !errors.Is(err, domain.ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken, got %v", err)
	}
	if claims != nil {
		t.Fatalf("expected nil claims, got %+v", claims)
	}
}

func signValidStepUpJWT(t *testing.T, signer *JWTStepUpTokenSigner, now time.Time) string {
	t.Helper()

	stepUpToken, err := domain.NewStepUpToken(
		"step-up-jti",
		uuid.New(),
		domain.StepUpEndpointInternalTransferCreate,
		now,
	)
	if err != nil {
		t.Fatalf("expected no error creating step-up token, got %v", err)
	}

	signedToken, err := signer.Sign(stepUpToken)
	if err != nil {
		t.Fatalf("expected no error signing step-up token, got %v", err)
	}

	return signedToken
}

func signStepUpJWT(
	t *testing.T,
	secret string,
	method jwt.SigningMethod,
	claims stepUpTokenClaims,
) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("expected no error signing test token, got %v", err)
	}

	return signedToken
}

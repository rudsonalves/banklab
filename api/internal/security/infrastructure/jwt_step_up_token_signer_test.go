package infrastructure

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/security/domain"
)

func TestJWTStepUpTokenSigner_SignContainsRequiredClaims(t *testing.T) {
	secret := "step-up-secret"
	signer := NewJWTStepUpTokenSigner(secret)
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
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

	claims := parseSignedStepUpToken(t, signedToken, secret, now)

	if claims.UserID != userID.String() {
		t.Fatalf("expected user_id %q, got %q", userID.String(), claims.UserID)
	}
	if claims.EndpointKey != domain.StepUpEndpointInternalTransferCreate {
		t.Fatalf("expected endpoint_key %q, got %q", domain.StepUpEndpointInternalTransferCreate, claims.EndpointKey)
	}
	if claims.Scope != stepUpTokenScope {
		t.Fatalf("expected scope %q, got %q", stepUpTokenScope, claims.Scope)
	}
	if claims.ID != jti {
		t.Fatalf("expected jti %q, got %q", jti, claims.ID)
	}
	if claims.IssuedAt == nil {
		t.Fatal("expected iat claim to be present")
	}
	if claims.ExpiresAt == nil {
		t.Fatal("expected exp claim to be present")
	}
	if !claims.IssuedAt.Time.Equal(now) {
		t.Fatalf("expected iat %v, got %v", now, claims.IssuedAt.Time)
	}
	if !claims.ExpiresAt.Time.Equal(now.Add(domain.StepUpTokenDefaultDuration)) {
		t.Fatalf("expected exp %v, got %v", now.Add(domain.StepUpTokenDefaultDuration), claims.ExpiresAt.Time)
	}
	if got := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time); got != domain.StepUpTokenDefaultDuration {
		t.Fatalf("expected token duration %v, got %v", domain.StepUpTokenDefaultDuration, got)
	}
}

func TestJWTStepUpTokenSigner_SignedTokenValidatesWithSameSecret(t *testing.T) {
	secret := "step-up-secret"
	signer := NewJWTStepUpTokenSigner(secret)
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
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

	parsed := parseRawStepUpToken(t, signedToken, secret, now)
	if !parsed.Valid {
		t.Fatal("expected signed step-up token to be valid")
	}
}

func TestJWTStepUpTokenSigner_InvalidSecretFailsValidation(t *testing.T) {
	signer := NewJWTStepUpTokenSigner("issuer-secret")
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
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

	claims := &stepUpTokenClaims{}
	parsed, err := jwt.ParseWithClaims(signedToken, claims, func(t *jwt.Token) (any, error) {
		return []byte("validator-secret"), nil
	}, jwt.WithTimeFunc(func() time.Time {
		return now
	}))
	if err == nil {
		t.Fatal("expected validation with wrong secret to fail")
	}
	if parsed != nil && parsed.Valid {
		t.Fatal("expected token validated with wrong secret to be invalid")
	}
}

func TestJWTStepUpTokenSigner_RejectsInvalidOrConsumedToken(t *testing.T) {
	signer := NewJWTStepUpTokenSigner("step-up-secret")
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)

	_, err := signer.Sign(nil)
	if !errors.Is(err, domain.ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken for nil token, got %v", err)
	}

	consumedAt := now.Add(time.Minute)
	consumed, err := domain.RestoreStepUpToken(
		uuid.New(),
		"step-up-jti",
		uuid.New(),
		domain.StepUpEndpointInternalTransferCreate,
		domain.StepUpTokenConsumed,
		now.Add(domain.StepUpTokenDefaultDuration),
		&consumedAt,
		now,
	)
	if err != nil {
		t.Fatalf("expected no error creating consumed token, got %v", err)
	}

	_, err = signer.Sign(consumed)
	if !errors.Is(err, domain.ErrInvalidStepUpToken) {
		t.Fatalf("expected ErrInvalidStepUpToken for consumed token, got %v", err)
	}
}

func TestJWTStepUpTokenSigner_DoesNotIncludeSensitiveOrOperationPayload(t *testing.T) {
	secret := "step-up-secret"
	signer := NewJWTStepUpTokenSigner(secret)
	now := time.Date(2026, 5, 29, 10, 0, 0, 0, time.UTC)
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

	parsed := parseRawStepUpToken(t, signedToken, secret, now)
	payload, err := json.Marshal(parsed.Claims)
	if err != nil {
		t.Fatalf("expected claims to marshal, got %v", err)
	}

	var claims map[string]any
	if err := json.Unmarshal(payload, &claims); err != nil {
		t.Fatalf("expected claims to unmarshal into map, got %v", err)
	}

	for _, forbidden := range []string{
		"transaction_password",
		"password_hash",
		"operation_payload",
		"payload",
		"amount",
	} {
		if _, ok := claims[forbidden]; ok {
			t.Fatalf("expected claim %q to be absent, got claims %v", forbidden, claims)
		}
	}
}

func parseSignedStepUpToken(
	t *testing.T,
	token string,
	secret string,
	now time.Time,
) *stepUpTokenClaims {
	t.Helper()

	parsed := parseRawStepUpToken(t, token, secret, now)

	claims, ok := parsed.Claims.(*stepUpTokenClaims)
	if !ok {
		t.Fatalf("expected stepUpTokenClaims, got %T", parsed.Claims)
	}

	return claims
}

func parseRawStepUpToken(t *testing.T, token string, secret string, now time.Time) *jwt.Token {
	t.Helper()

	claims := &stepUpTokenClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(parsedToken *jwt.Token) (any, error) {
		if parsedToken.Method != jwt.SigningMethodHS256 {
			t.Fatalf("expected HS256 signing method, got %v", parsedToken.Method.Alg())
		}

		return []byte(secret), nil
	}, jwt.WithTimeFunc(func() time.Time {
		return now
	}))
	if err != nil {
		t.Fatalf("expected no error parsing signed step-up token, got %v", err)
	}

	return parsed
}

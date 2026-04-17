package infrastructure

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

func TestJWTTokenService_GenerateAndParseAccessToken_Success(t *testing.T) {
	service := NewJWTTokenService("test-secret", 2*time.Minute)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000123")
	cid := uuid.MustParse("de305d54-75b4-431b-adb2-eb6b9e546014")

	token, err := service.GenerateAccessToken(domain.TokenClaims{
		UserID:     userID,
		Role:       domain.RoleCustomer,
		CustomerID: &cid,
	})
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	claims, err := service.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("expected no error parsing token, got %v", err)
	}

	if claims.UserID != userID {
		t.Fatalf("expected user id %q, got %q", userID, claims.UserID)
	}

	if claims.Role != domain.RoleCustomer {
		t.Fatalf("expected role %q, got %q", domain.RoleCustomer, claims.Role)
	}

	if claims.CustomerID == nil || *claims.CustomerID != cid {
		t.Fatalf("expected customer id %q, got %#v", cid, claims.CustomerID)
	}
}

func TestJWTTokenService_GenerateAccessToken_SetsExpirationClaim(t *testing.T) {
	ttl := 15 * time.Minute
	service := NewJWTTokenService("test-secret", ttl)

	before := time.Now().UTC()
	token, err := service.GenerateAccessToken(domain.TokenClaims{
		UserID: uuid.MustParse("00000000-0000-0000-0000-000000000123"),
		Role:   domain.RoleCustomer,
	})
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}
	after := time.Now().UTC()

	parsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte("test-secret"), nil
	})
	if err != nil {
		t.Fatalf("expected no error parsing generated token, got %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected jwt.MapClaims, got %T", parsed.Claims)
	}

	exp, err := claims.GetExpirationTime()
	if err != nil {
		t.Fatalf("expected exp claim to be present, got %v", err)
	}

	minExp := before.Add(ttl).Add(-1 * time.Second)
	maxExp := after.Add(ttl).Add(1 * time.Second)
	if exp.Time.UTC().Before(minExp) || exp.Time.UTC().After(maxExp) {
		t.Fatalf("expected exp between %v and %v, got %v", minExp, maxExp, exp.Time)
	}
}

func TestJWTTokenService_ParseAccessToken_InvalidSignature(t *testing.T) {
	issuer := NewJWTTokenService("issuer-secret", time.Minute)
	validator := NewJWTTokenService("validator-secret", time.Minute)

	token, err := issuer.GenerateAccessToken(domain.TokenClaims{
		UserID: uuid.MustParse("00000000-0000-0000-0000-000000000123"),
		Role:   domain.RoleCustomer,
	})
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	_, err = validator.ParseAccessToken(token)
	if err == nil {
		t.Fatal("expected parsing to fail due to invalid signature")
	}
}

func TestJWTTokenService_ParseAccessToken_ExpiredToken(t *testing.T) {
	service := NewJWTTokenService("test-secret", -1*time.Second)

	token, err := service.GenerateAccessToken(domain.TokenClaims{
		UserID: uuid.MustParse("00000000-0000-0000-0000-000000000123"),
		Role:   domain.RoleCustomer,
	})
	if err != nil {
		t.Fatalf("expected no error generating token, got %v", err)
	}

	_, err = service.ParseAccessToken(token)
	if err == nil {
		t.Fatal("expected parsing expired token to fail")
	}
}

func TestJWTTokenService_ParseAccessToken_MalformedToken(t *testing.T) {
	service := NewJWTTokenService("test-secret", time.Minute)

	_, err := service.ParseAccessToken("not-a-jwt")
	if err == nil {
		t.Fatal("expected parsing malformed token to fail")
	}
}

func TestJWTTokenService_ParseAccessToken_InvalidSigningMethod(t *testing.T) {
	secret := "test-secret"
	service := NewJWTTokenService(secret, time.Minute)

	payload := jwt.MapClaims{
		"sub":  "user-123",
		"role": string(domain.RoleCustomer),
		"iat":  time.Now().UTC().Unix(),
		"exp":  time.Now().UTC().Add(time.Minute).Unix(),
	}
	unsigned := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := unsigned.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("expected no error creating token with none algorithm, got %v", err)
	}

	_, err = service.ParseAccessToken(token)
	if err == nil {
		t.Fatal("expected parsing token with invalid signing method to fail")
	}

	if !strings.Contains(err.Error(), "invalid token signing method") {
		t.Fatalf("expected signing method error, got %v", err)
	}
}

func TestJWTTokenService_ParseAccessToken_MissingRequiredClaims(t *testing.T) {
	secret := "test-secret"
	service := NewJWTTokenService(secret, time.Minute)

	tests := []struct {
		name    string
		claims  jwt.MapClaims
		errText string
	}{
		{
			name: "missing sub",
			claims: jwt.MapClaims{
				"role": string(domain.RoleCustomer),
				"iat":  time.Now().UTC().Unix(),
				"exp":  time.Now().UTC().Add(time.Minute).Unix(),
			},
			errText: "missing subject claim",
		},
		{
			name: "missing role",
			claims: jwt.MapClaims{
				"sub": "user-123",
				"iat": time.Now().UTC().Unix(),
				"exp": time.Now().UTC().Add(time.Minute).Unix(),
			},
			errText: "missing role claim",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tc.claims)
			signedToken, err := token.SignedString([]byte(secret))
			if err != nil {
				t.Fatalf("expected no error signing token, got %v", err)
			}

			_, err = service.ParseAccessToken(signedToken)
			if err == nil {
				t.Fatal("expected parsing token to fail")
			}

			if !strings.Contains(err.Error(), tc.errText) {
				t.Fatalf("expected error containing %q, got %v", tc.errText, err)
			}
		})
	}
}

func TestJWTTokenService_GenerateAndParseRefreshToken_Success(t *testing.T) {
	service := NewJWTTokenService("test-secret", time.Minute)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000123")

	token, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("expected no error generating refresh token, got %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty refresh token")
	}

	parsedUserID, err := service.ParseRefreshToken(token)
	if err != nil {
		t.Fatalf("expected no error parsing refresh token, got %v", err)
	}

	if parsedUserID != userID {
		t.Fatalf("expected user id %q, got %q", userID, parsedUserID)
	}
}

func TestJWTTokenService_GenerateRefreshToken_HighEntropy(t *testing.T) {
	service := NewJWTTokenService("test-secret", time.Minute)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000123")

	first, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("expected no error generating first refresh token, got %v", err)
	}

	second, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("expected no error generating second refresh token, got %v", err)
	}

	if first == second {
		t.Fatal("expected different refresh tokens for same user")
	}
}

func TestJWTTokenService_ParseRefreshToken_TamperedSignature(t *testing.T) {
	service := NewJWTTokenService("test-secret", time.Minute)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000123")

	token, err := service.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("expected no error generating refresh token, got %v", err)
	}

	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		t.Fatalf("expected token with 2 parts, got %d", len(parts))
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("expected valid signature encoding, got %v", err)
	}
	sig[0] ^= 0x01
	tamperedSig := base64.RawURLEncoding.EncodeToString(sig)
	tampered := parts[0] + "." + tamperedSig

	_, err = service.ParseRefreshToken(tampered)
	if err == nil {
		t.Fatal("expected parsing tampered refresh token to fail")
	}

	if !strings.Contains(err.Error(), "signature") {
		t.Fatalf("expected signature error, got %v", err)
	}
}

func TestJWTTokenService_ParseRefreshToken_Malformed(t *testing.T) {
	service := NewJWTTokenService("test-secret", time.Minute)

	_, err := service.ParseRefreshToken("not-a-refresh-token")
	if err == nil {
		t.Fatal("expected parsing malformed refresh token to fail")
	}
}

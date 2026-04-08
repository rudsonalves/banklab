package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/domain"
)

type tokenServiceMock struct {
	claims *domain.TokenClaims
	err    error
	token  string
	called bool
}

func (m *tokenServiceMock) GenerateAccessToken(claims domain.TokenClaims) (string, error) {
	return "", nil
}

func (m *tokenServiceMock) ParseAccessToken(token string) (*domain.TokenClaims, error) {
	m.called = true
	m.token = token
	if m.err != nil {
		return nil, m.err
	}
	return m.claims, nil
}

func TestJWTMiddleware_RequireAuth_ValidToken(t *testing.T) {
	customerID := uuid.New()
	userUUID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	tokenService := &tokenServiceMock{
		claims: &domain.TokenClaims{
			UserID:     userUUID,
			Role:       domain.RoleCustomer,
			CustomerID: &customerID,
		},
	}
	middleware := NewJWTMiddleware(tokenService)
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()

	var principal *AuthenticatedUser
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		principal, ok = GetAuthenticatedUser(r.Context())
		if !ok {
			t.Fatal("expected authenticated user in context")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	middleware.RequireAuth(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if !tokenService.called {
		t.Fatal("expected token service to be called")
	}

	if tokenService.token != "valid-token" {
		t.Fatalf("expected token %q, got %q", "valid-token", tokenService.token)
	}

	if principal == nil || principal.UserID != userUUID {
		t.Fatalf("expected principal user id %q, got %#v", userUUID, principal)
	}

	if principal.Role != domain.RoleCustomer {
		t.Fatalf("expected role %q, got %q", domain.RoleCustomer, principal.Role)
	}

	if principal.CustomerID == nil || *principal.CustomerID != customerID {
		t.Fatalf("expected customer id %q, got %#v", customerID, principal.CustomerID)
	}
}

func TestJWTMiddleware_RequireAuth_MissingHeader(t *testing.T) {
	middleware := NewJWTMiddleware(&tokenServiceMock{})
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	rec := httptest.NewRecorder()

	middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected next handler not to be called")
	})).ServeHTTP(rec, req)

	assertAuthErrorCode(t, rec, http.StatusUnauthorized, "UNAUTHORIZED")
}

func TestJWTMiddleware_RequireAuth_MalformedHeader(t *testing.T) {
	middleware := NewJWTMiddleware(&tokenServiceMock{})
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Token abc")
	rec := httptest.NewRecorder()

	middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected next handler not to be called")
	})).ServeHTTP(rec, req)

	assertAuthErrorCode(t, rec, http.StatusUnauthorized, "UNAUTHORIZED")
}

func TestJWTMiddleware_RequireAuth_InvalidToken(t *testing.T) {
	tokenService := &tokenServiceMock{err: errors.New("invalid token")}
	middleware := NewJWTMiddleware(tokenService)
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()

	middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected next handler not to be called")
	})).ServeHTTP(rec, req)

	assertAuthErrorCode(t, rec, http.StatusUnauthorized, "INVALID_TOKEN")
}

func TestJWTMiddleware_RequireAuth_ExpiredToken(t *testing.T) {
	tokenService := &tokenServiceMock{err: errors.New("token is expired")}
	middleware := NewJWTMiddleware(tokenService)
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	rec := httptest.NewRecorder()

	middleware.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("expected next handler not to be called")
	})).ServeHTTP(rec, req)

	assertAuthErrorCode(t, rec, http.StatusUnauthorized, "INVALID_TOKEN")
}

func TestJWTMiddleware_OptionalAuth_MissingHeader(t *testing.T) {
	middleware := NewJWTMiddleware(&tokenServiceMock{})
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rec := httptest.NewRecorder()
	called := false

	middleware.OptionalAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		if user, ok := GetAuthenticatedUser(r.Context()); ok || user != nil {
			t.Fatalf("expected no authenticated user, got %#v", user)
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestGetAuthenticatedUser_NoUser(t *testing.T) {
	user, ok := GetAuthenticatedUser(context.Background())
	if ok || user != nil {
		t.Fatalf("expected no authenticated user, got %#v", user)
	}
}

func assertAuthErrorCode(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedCode string) {
	t.Helper()

	if rec.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, rec.Code)
	}

	var got struct {
		Data  interface{} `json:"data"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data != nil {
		t.Fatalf("expected nil data, got %#v", got.Data)
	}

	if got.Error.Code != expectedCode {
		t.Fatalf("expected error code %q, got %q", expectedCode, got.Error.Code)
	}
}

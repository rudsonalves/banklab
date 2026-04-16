package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAppToken_Require_MissingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	AppToken("expected-token")(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("expected next handler not to be called")
	})).ServeHTTP(rec, req)

	assertAppTokenError(t, rec)
}

func TestAppToken_Require_InvalidHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(headerAppToken, "wrong-token")
	rec := httptest.NewRecorder()

	AppToken("expected-token")(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("expected next handler not to be called")
	})).ServeHTTP(rec, req)

	assertAppTokenError(t, rec)
}

func TestAppToken_Require_ValidHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(headerAppToken, "expected-token")
	rec := httptest.NewRecorder()
	called := false

	AppToken("expected-token")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func assertAppTokenError(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	var got struct {
		Data  interface{} `json:"data"`
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data != nil {
		t.Fatalf("expected nil data, got %#v", got.Data)
	}

	if got.Error.Code != "INVALID_APP_TOKEN" {
		t.Fatalf("expected error code %q, got %q", "INVALID_APP_TOKEN", got.Error.Code)
	}

	if got.Error.Message != "invalid application token" {
		t.Fatalf("expected message %q, got %q", "invalid application token", got.Error.Message)
	}
}

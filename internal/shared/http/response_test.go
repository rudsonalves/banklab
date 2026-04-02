package sharedhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteSuccess(rec, http.StatusCreated, map[string]string{"id": "123"})

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected content type %q, got %q", "application/json", contentType)
	}

	var got struct {
		Data  map[string]string `json:"data"`
		Error interface{}       `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data["id"] != "123" {
		t.Fatalf("expected id %q, got %q", "123", got.Data["id"])
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, http.StatusUnauthorized, sharederrors.ErrUnauthorized)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if contentType := rec.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected content type %q, got %q", "application/json", contentType)
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

	if got.Error.Code != "UNAUTHORIZED" {
		t.Fatalf("expected error code %q, got %q", "UNAUTHORIZED", got.Error.Code)
	}

	if got.Error.Message != "Unauthorized" {
		t.Fatalf("expected message %q, got %q", "Unauthorized", got.Error.Message)
	}
}

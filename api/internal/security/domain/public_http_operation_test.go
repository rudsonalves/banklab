package domain

import (
	"errors"
	"testing"
)

func TestNewPublicHTTPOperation_AcceptsValidOperation(t *testing.T) {
	operation, err := NewPublicHTTPOperation("post", "/accounts/internal-transfers")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if operation.Method != "POST" {
		t.Fatalf("expected method to be POST, got %q", operation.Method)
	}

	if operation.Path != "/accounts/internal-transfers" {
		t.Fatalf("unexpected path: %q", operation.Path)
	}
}

func TestNewPublicHTTPOperation_RejectsEmptyMethod(t *testing.T) {
	_, err := NewPublicHTTPOperation("   ", "/accounts/internal-transfers")

	if !errors.Is(err, ErrInvalidStepUpPublicOperationMethod) {
		t.Fatalf("expected ErrInvalidStepUpPublicOperationMethod, got %v", err)
	}
}

func TestNewPublicHTTPOperation_RejectsPathWithoutLeadingSlash(t *testing.T) {
	_, err := NewPublicHTTPOperation("POST", "accounts/internal-transfers")

	if !errors.Is(err, ErrInvalidStepUpPublicOperationPath) {
		t.Fatalf("expected ErrInvalidStepUpPublicOperationPath, got %v", err)
	}
}

func TestNewPublicHTTPOperation_RejectsEmptyPath(t *testing.T) {
	_, err := NewPublicHTTPOperation("POST", "   ")

	if !errors.Is(err, ErrInvalidStepUpPublicOperationPath) {
		t.Fatalf("expected ErrInvalidStepUpPublicOperationPath, got %v", err)
	}
}

func TestNewPublicHTTPOperation_RejectsPathWithScheme(t *testing.T) {
	testCases := []string{
		"http://api.banklab.local/accounts/internal-transfers",
		"https://api.banklab.local/accounts/internal-transfers",
	}

	for _, path := range testCases {
		_, err := NewPublicHTTPOperation("POST", path)
		if !errors.Is(err, ErrInvalidStepUpPublicOperationPath) {
			t.Fatalf("path %q: expected ErrInvalidStepUpPublicOperationPath, got %v", path, err)
		}
	}
}

func TestNewPublicHTTPOperation_RejectsPathWithHost(t *testing.T) {
	_, err := NewPublicHTTPOperation("POST", "//api.banklab.local/accounts/internal-transfers")

	if !errors.Is(err, ErrInvalidStepUpPublicOperationPath) {
		t.Fatalf("expected ErrInvalidStepUpPublicOperationPath, got %v", err)
	}
}

func TestNewPublicHTTPOperation_RejectsPathWithQueryString(t *testing.T) {
	_, err := NewPublicHTTPOperation("POST", "/accounts/internal-transfers?tenant=banklab")

	if !errors.Is(err, ErrInvalidStepUpPublicOperationPath) {
		t.Fatalf("expected ErrInvalidStepUpPublicOperationPath, got %v", err)
	}
}

func TestNewPublicHTTPOperation_RejectsPathWithFragment(t *testing.T) {
	_, err := NewPublicHTTPOperation("POST", "/accounts/internal-transfers#details")

	if !errors.Is(err, ErrInvalidStepUpPublicOperationPath) {
		t.Fatalf("expected ErrInvalidStepUpPublicOperationPath, got %v", err)
	}
}

func TestNewPublicHTTPOperation_AllowsTemplatedPath(t *testing.T) {
	operation, err := NewPublicHTTPOperation("POST", "/accounts/{id}/withdraw")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if operation.Path != "/accounts/{id}/withdraw" {
		t.Fatalf("unexpected path: %q", operation.Path)
	}
}

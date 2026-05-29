package domain

import (
	"errors"
	"testing"
)

func TestDefaultStepUpPublicOperationResolver_ResolvesInternalTransferCreate(t *testing.T) {
	resolver := NewDefaultStepUpPublicOperationResolver()
	operation, err := NewPublicHTTPOperation("POST", "/accounts/internal-transfers")
	if err != nil {
		t.Fatalf("expected valid operation, got %v", err)
	}

	endpointKey, err := resolver.Resolve(operation)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if endpointKey != StepUpEndpointInternalTransferCreate {
		t.Fatalf("expected endpoint key %q, got %q", StepUpEndpointInternalTransferCreate, endpointKey)
	}
}

func TestDefaultStepUpPublicOperationResolver_RejectsDifferentMethodForSamePath(t *testing.T) {
	resolver := NewDefaultStepUpPublicOperationResolver()
	operation, err := NewPublicHTTPOperation("GET", "/accounts/internal-transfers")
	if err != nil {
		t.Fatalf("expected valid operation, got %v", err)
	}

	_, err = resolver.Resolve(operation)
	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

func TestDefaultStepUpPublicOperationResolver_RejectsNonConfiguredPath(t *testing.T) {
	resolver := NewDefaultStepUpPublicOperationResolver()
	operation, err := NewPublicHTTPOperation("POST", "/accounts/pix-transfers")
	if err != nil {
		t.Fatalf("expected valid operation, got %v", err)
	}

	_, err = resolver.Resolve(operation)
	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

func TestWhitelistStepUpPublicOperationResolver_RejectsInvalidOperationInput(t *testing.T) {
	resolver := NewDefaultStepUpPublicOperationResolver()

	_, err := resolver.Resolve(&PublicHTTPOperation{Method: "", Path: "/accounts/internal-transfers"})
	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

func TestWhitelistStepUpPublicOperationResolver_SupportsTemplatedPathMapping(t *testing.T) {
	resolver := NewWhitelistStepUpPublicOperationResolver(PublicStepUpOperationMapping{
		Method:      "POST",
		Path:        "/accounts/{id}/withdraw",
		EndpointKey: "account.withdraw",
	})

	operation, err := NewPublicHTTPOperation("post", "/accounts/{id}/withdraw")
	if err != nil {
		t.Fatalf("expected valid operation, got %v", err)
	}

	endpointKey, err := resolver.Resolve(operation)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if endpointKey != "account.withdraw" {
		t.Fatalf("expected endpoint key account.withdraw, got %q", endpointKey)
	}
}

func TestWhitelistStepUpPublicOperationResolver_NilResolverRejectsOperation(t *testing.T) {
	var resolver *WhitelistStepUpPublicOperationResolver
	operation, err := NewPublicHTTPOperation("POST", "/accounts/internal-transfers")
	if err != nil {
		t.Fatalf("expected valid operation, got %v", err)
	}

	_, err = resolver.Resolve(operation)
	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

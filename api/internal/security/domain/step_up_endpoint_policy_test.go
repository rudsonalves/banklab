package domain

import (
	"errors"
	"testing"
)

func TestDefaultStepUpEndpointPolicy_AllowsInternalTransferCreate(t *testing.T) {
	policy := NewDefaultStepUpEndpointPolicy()

	err := policy.Validate(StepUpEndpointInternalTransferCreate)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultStepUpEndpointPolicy_RejectsUnknownEndpoint(t *testing.T) {
	policy := NewDefaultStepUpEndpointPolicy()

	err := policy.Validate("pix.create")

	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

func TestDefaultStepUpEndpointPolicy_RejectsBlankEndpoint(t *testing.T) {
	policy := NewDefaultStepUpEndpointPolicy()

	err := policy.Validate("   ")

	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

func TestWhitelistStepUpEndpointPolicy_TrimsInput(t *testing.T) {
	policy := NewDefaultStepUpEndpointPolicy()

	err := policy.Validate(" " + StepUpEndpointInternalTransferCreate + " ")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWhitelistStepUpEndpointPolicy_CustomWhitelist(t *testing.T) {
	policy := NewWhitelistStepUpEndpointPolicy("pix.create", " boleto.pay ")

	if err := policy.Validate("pix.create"); err != nil {
		t.Fatalf("expected pix.create to be allowed, got %v", err)
	}

	if err := policy.Validate("boleto.pay"); err != nil {
		t.Fatalf("expected boleto.pay to be allowed, got %v", err)
	}

	err := policy.Validate(StepUpEndpointInternalTransferCreate)
	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

func TestWhitelistStepUpEndpointPolicy_NilPolicyRejectsEndpoint(t *testing.T) {
	var policy *WhitelistStepUpEndpointPolicy

	err := policy.Validate(StepUpEndpointInternalTransferCreate)

	if !errors.Is(err, ErrStepUpEndpointNotAllowed) {
		t.Fatalf("expected ErrStepUpEndpointNotAllowed, got %v", err)
	}
}

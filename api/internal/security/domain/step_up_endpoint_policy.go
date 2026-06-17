package domain

import "strings"

const (
	StepUpEndpointInternalTransferCreate     = "internal_transfer.create"
	StepUpEndpointInstallationRegisterCreate = "installation.register"
)

type WhitelistStepUpEndpointPolicy struct {
	allowed map[string]struct{}
}

// NewWhitelistStepUpEndpointPolicy creates a new instance of
// WhitelistStepUpEndpointPolicy with the provided endpoint keys.
func NewWhitelistStepUpEndpointPolicy(endpointKeys ...string) *WhitelistStepUpEndpointPolicy {
	allowed := make(map[string]struct{}, len(endpointKeys))

	for _, endpointKey := range endpointKeys {
		normalized := strings.TrimSpace(endpointKey)
		if normalized == "" {
			continue
		}

		allowed[normalized] = struct{}{}
	}

	return &WhitelistStepUpEndpointPolicy{
		allowed: allowed,
	}
}

// NewDefaultStepUpEndpointPolicy creates a new instance of WhitelistStepUpEndpointPolicy
// with the default allowed endpoint keys.
func NewDefaultStepUpEndpointPolicy() *WhitelistStepUpEndpointPolicy {
	return NewWhitelistStepUpEndpointPolicy(
		StepUpEndpointInternalTransferCreate,
		StepUpEndpointInstallationRegisterCreate,
	)
}

// Validate checks if the provided endpoint key is allowed by the policy.
func (p *WhitelistStepUpEndpointPolicy) Validate(endpointKey string) error {
	if p == nil {
		return ErrStepUpEndpointNotAllowed
	}

	normalized := strings.TrimSpace(endpointKey)
	if normalized == "" {
		return ErrStepUpEndpointNotAllowed
	}

	if _, ok := p.allowed[normalized]; !ok {
		return ErrStepUpEndpointNotAllowed
	}

	return nil
}

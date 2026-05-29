package domain

import "strings"

const (
	StepUpPublicMethodInternalTransferCreate = "POST"
	StepUpPublicPathInternalTransferCreate   = "/accounts/internal-transfers"
)

type PublicStepUpOperationMapping struct {
	Method      string
	Path        string
	EndpointKey string
}

type WhitelistStepUpPublicOperationResolver struct {
	endpointKeysByOperation map[string]string
}

// NewWhitelistStepUpPublicOperationResolver creates a new instance of
// WhitelistStepUpPublicOperationResolver with the provided operation
// to endpoint key mappings.
func NewWhitelistStepUpPublicOperationResolver(
	mappings ...PublicStepUpOperationMapping,
) *WhitelistStepUpPublicOperationResolver {
	endpointKeysByOperation := make(map[string]string, len(mappings))

	for _, mapping := range mappings {
		operation, err := NewPublicHTTPOperation(mapping.Method, mapping.Path)
		if err != nil {
			continue
		}

		endpointKey := strings.TrimSpace(mapping.EndpointKey)
		if endpointKey == "" {
			continue
		}

		endpointKeysByOperation[operation.LookupKey()] = endpointKey
	}

	return &WhitelistStepUpPublicOperationResolver{
		endpointKeysByOperation: endpointKeysByOperation,
	}
}

// NewDefaultStepUpPublicOperationResolver creates a new instance of
// WhitelistStepUpPublicOperationResolver with the default operation to endpoint
// key mappings.
func NewDefaultStepUpPublicOperationResolver() *WhitelistStepUpPublicOperationResolver {
	return NewWhitelistStepUpPublicOperationResolver(PublicStepUpOperationMapping{
		Method:      StepUpPublicMethodInternalTransferCreate,
		Path:        StepUpPublicPathInternalTransferCreate,
		EndpointKey: StepUpEndpointInternalTransferCreate,
	})
}

// Resolve checks if the provided PublicHTTPOperation matches any of the
// configured operation to endpoint key mappings and returns the corresponding
// endpoint key if a match is found. If the resolver is nil, the operation is
// invalid, or no match is found, it returns an error indicating that the
// endpoint is not allowed.
func (r *WhitelistStepUpPublicOperationResolver) Resolve(operation *PublicHTTPOperation) (string, error) {
	if r == nil || operation == nil {
		return "", ErrStepUpEndpointNotAllowed
	}

	if err := operation.Validate(); err != nil {
		return "", ErrStepUpEndpointNotAllowed
	}

	endpointKey, ok := r.endpointKeysByOperation[operation.LookupKey()]
	if !ok {
		return "", ErrStepUpEndpointNotAllowed
	}

	return endpointKey, nil
}

package domain

import (
	"net/url"
	"strings"
	"unicode"
)

// PublicHTTPOperation represents the public HTTP operation requested for
// step-up authorization. For parameterized routes, clients must send the
// documented public template (example: /accounts/{id}/withdraw).
type PublicHTTPOperation struct {
	Method string
	Path   string
}

// NewPublicHTTPOperation creates a new PublicHTTPOperation with the given
// method and path.
func NewPublicHTTPOperation(method, path string) (*PublicHTTPOperation, error) {
	operation := &PublicHTTPOperation{
		Method: strings.ToUpper(strings.TrimSpace(method)),
		Path:   strings.TrimSpace(path),
	}

	if err := operation.Validate(); err != nil {
		return nil, err
	}

	return operation, nil
}

// Validate checks if the PublicHTTPOperation has valid method and path.
func (o *PublicHTTPOperation) Validate() error {
	if o == nil {
		return ErrInvalidStepUpPublicOperation
	}

	if !isValidHTTPMethodToken(o.Method) {
		return ErrInvalidStepUpPublicOperationMethod
	}

	if !isValidPublicPath(o.Path) {
		return ErrInvalidStepUpPublicOperationPath
	}

	return nil
}

// LookupKey returns the canonical method+path representation used by internal
// operation resolvers.
func (o *PublicHTTPOperation) LookupKey() string {
	if o == nil {
		return ""
	}

	return o.Method + " " + o.Path
}

// isValidHTTPMethodToken checks if the provided method is a valid HTTP
// method token (non-empty, uppercase, no spaces).
func isValidHTTPMethodToken(method string) bool {
	if strings.TrimSpace(method) == "" {
		return false
	}

	for _, r := range method {
		if !unicode.IsUpper(r) {
			return false
		}
	}

	return true
}

// isValidPublicPath checks if the provided path is a valid public HTTP path
// (non-empty, starts with '/', no scheme or host, no query string, no fragment).
func isValidPublicPath(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" || !strings.HasPrefix(path, "/") {
		return false
	}

	if strings.HasPrefix(path, "//") {
		return false
	}

	if strings.Contains(path, "://") {
		return false
	}

	parsed, err := url.Parse(path)
	if err != nil {
		return false
	}

	if parsed.Scheme != "" || parsed.Host != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}

	return parsed.Path == path
}

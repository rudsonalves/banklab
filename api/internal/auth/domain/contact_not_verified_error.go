package domain

// ContactNotVerifiedError is returned when one or more contact methods have
// not yet been verified. It wraps ErrContactNotVerified so errors.Is works.
type ContactNotVerifiedError struct {
	EmailVerified bool
	PhoneVerified bool
}

// NewContactNotVerifiedError returns a ContactNotVerifiedError populated with
// the current verification state of each contact method.
func NewContactNotVerifiedError(emailVerified, phoneVerified bool) error {
	return &ContactNotVerifiedError{
		EmailVerified: emailVerified,
		PhoneVerified: phoneVerified,
	}
}

// Error implements the error interface, returning the sentinel message.
func (e *ContactNotVerifiedError) Error() string {
	return ErrContactNotVerified.Error()
}

// Unwrap returns ErrContactNotVerified, enabling errors.Is and errors.As
// to match against the sentinel.
func (e *ContactNotVerifiedError) Unwrap() error {
	return ErrContactNotVerified
}

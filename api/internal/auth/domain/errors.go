package domain

import "errors"

var (
	ErrEmailAlreadyExists          = errors.New("email already exists")
	ErrPhoneAlreadyExists          = errors.New("phone already exists")
	ErrForbidden                   = errors.New("forbidden")
	ErrInvalidEmail                = errors.New("invalid email")
	ErrInvalidData                 = errors.New("invalid data")
	ErrInvalidPassword             = errors.New("invalid password")
	ErrInvalidCredentials          = errors.New("invalid credentials")
	ErrContactNotVerified          = errors.New("contact not verified")
	ErrAccountApprovalRequired     = errors.New("account approval required")
	ErrUnauthorized                = errors.New("unauthorized")
	ErrInvalidToken                = errors.New("invalid token")
	ErrInvalidUserState            = errors.New("invalid user state: customer role requires customer_id")
	ErrUserNotFound                = errors.New("user not found")
	ErrUserAlreadyActive           = errors.New("user already active")
	ErrSessionNotFound             = errors.New("session not found")
	ErrContactVerificationNotFound = errors.New("contact verification not found")
	ErrInvalidVerificationToken    = errors.New("verification token is invalid")
	ErrContactVerificationExpired  = errors.New("contact verification expired")
)

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

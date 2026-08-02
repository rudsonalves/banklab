package domain

import "errors"

var (
	ErrInvalidInstallationID                = errors.New("invalid installation id")
	ErrInvalidInstallationResourceID        = errors.New("invalid installation resource id")
	ErrInvalidInstallation                  = errors.New("invalid installation")
	ErrInstallationNotFound                 = errors.New("installation not found")
	ErrInstallationMismatch                 = errors.New("installation mismatch")
	ErrInstallationRevoked                  = errors.New("installation revoked")
	ErrInstallationLimitReached             = errors.New("installation limit reached")
	ErrFirstInstallationAlreadyBootstrapped = errors.New("first installation already bootstrapped")

	ErrInvalidRestrictedAuthorization       = errors.New("invalid restricted authorization")
	ErrRestrictedAuthorizationNotFound      = errors.New("restricted authorization not found")
	ErrRestrictedAuthorizationInvalid       = errors.New("restricted authorization invalid")
	ErrRestrictedAuthorizationExpired       = errors.New("restricted authorization expired")
	ErrRestrictedAuthorizationConsumed      = errors.New("restricted authorization already consumed")
	ErrRestrictedAuthorizationRevoked       = errors.New("restricted authorization revoked")
	ErrRestrictedAuthorizationAlreadyActive = errors.New("restricted authorization already active")
)

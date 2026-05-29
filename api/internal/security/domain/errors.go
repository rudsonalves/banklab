package domain

import "errors"

var (
	ErrInvalidTransactionPasswordPIN = errors.New("invalid transaction password pin")
	ErrTransactionPasswordAlreadySet = errors.New("transaction password already set")
	ErrTransactionPasswordNotSet     = errors.New("transaction password not set")
	ErrTransactionPasswordInvalid    = errors.New("transaction password is invalid")
	ErrTransactionPasswordLocked     = errors.New("transaction password is locked due to multiple failed attempts")
	ErrInvalidTransactionPassword    = errors.New("invalid transaction password")

	ErrInvalidStepUpToken  = errors.New("invalid step-up token")
	ErrStepUpTokenExpired  = errors.New("step-up token expired")
	ErrStepUpTokenConsumed = errors.New("step-up token already consumed")

	ErrStepUpEndpointNotAllowed = errors.New("step-up endpoint not allowed")
)

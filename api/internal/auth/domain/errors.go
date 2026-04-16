package domain

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrInvalidData        = errors.New("invalid data")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidUserState   = errors.New("invalid user state: customer role requires customer_id")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyActive  = errors.New("user already active")
	ErrSessionNotFound    = errors.New("session not found")
)

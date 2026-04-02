package domain

import "errors"

var (
	ErrInvalidData      = errors.New("invalid data")
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrAccountNotFound  = errors.New("account not found")
	ErrCustomerNotFound = errors.New("customer not found")
	ErrAccountInactive  = errors.New("account inactive")
)

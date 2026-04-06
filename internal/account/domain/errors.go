package domain

import "errors"

var (
	ErrInvalidData               = errors.New("invalid data")
	ErrInvalidAmount             = errors.New("invalid amount")
	ErrAccountNotFound           = errors.New("account not found")
	ErrInsufficientBalance       = errors.New("insufficient balance")
	ErrSameAccountTransfer       = errors.New("same account transfer")
	ErrCustomerNotFound          = errors.New("customer not found")
	ErrAccountInactive           = errors.New("account inactive")
	ErrForbidden                 = errors.New("forbidden")
	ErrOperationAlreadyProcessed = errors.New("operation already processed")
)

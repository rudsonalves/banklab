package domain

import "errors"

var (
	ErrInvalidData      = errors.New("invalid data")
	ErrCustomerNotFound = errors.New("customer not found")
)

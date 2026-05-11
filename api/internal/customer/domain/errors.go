package domain

import "errors"

var (
	ErrNameRequired = errors.New("name is required")
	ErrCPFRequired  = errors.New("cpf is required")

	ErrCPFInvalid       = errors.New("cpf is invalid")
	ErrCPFAlreadyExists = errors.New("cpf already exists")

	ErrInvalidData = errors.New("invalid data")
	ErrNotFound    = errors.New("customer not found")
)

package domain

import "errors"

var (
	ErrNameRequired = errors.New("name is required")
	ErrCPFRequired  = errors.New("cpf is required")

	ErrCPFAlreadyExists = errors.New("cpf already exists")

	ErrInvalidData = errors.New("invalid data")
	ErrNotFound    = errors.New("customer not found")
)

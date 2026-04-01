package domain

import "errors"

var (
	ErrNameRequired  = errors.New("name is required")
	ErrCPFRequired   = errors.New("cpf is required")
	ErrEmailRequired = errors.New("email is required")

	ErrCPFAlreadyExists   = errors.New("cpf already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrInvalidData = errors.New("invalid data")
)

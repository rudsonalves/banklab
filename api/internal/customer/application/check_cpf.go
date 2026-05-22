package application

import (
	"context"
	"fmt"

	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type CheckCPFUseCase struct {
	documentRepo domain.CustomerDocumentRepository
}

func NewCheckCPFUseCase(documentRepo domain.CustomerDocumentRepository) *CheckCPFUseCase {
	return &CheckCPFUseCase{documentRepo: documentRepo}
}

type CheckCPFInput struct {
	CPF string
}

type CheckCPFOutput struct {
	CPF       string `json:"cpf"`
	Exists    bool   `json:"exists"`
	Available bool   `json:"available"`
}

func (uc *CheckCPFUseCase) Execute(ctx context.Context, input CheckCPFInput) (*CheckCPFOutput, error) {
	cpf := domain.NormalizeCPF(input.CPF)
	if cpf == "" {
		return nil, domain.ErrCPFRequired
	}
	if !domain.ValidateCPF(cpf) {
		return nil, domain.ErrCPFInvalid
	}

	exists, err := uc.documentRepo.ExistsCPF(ctx, cpf)
	if err != nil {
		return nil, fmt.Errorf("check cpf existence: %w", err)
	}

	return &CheckCPFOutput{
		CPF:       cpf,
		Exists:    exists,
		Available: !exists,
	}, nil
}

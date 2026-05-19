package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type checkCPFDocumentRepositoryMock struct {
	existsCPFCalls int
	existsCPFValue string
	existsCPF      bool
	existsCPFErr   error
}

func (m *checkCPFDocumentRepositoryMock) CreateDocument(ctx context.Context, document *domain.CustomerDocument) error {
	return nil
}

func (m *checkCPFDocumentRepositoryMock) ExistsCPF(ctx context.Context, cpf string) (bool, error) {
	m.existsCPFCalls++
	m.existsCPFValue = cpf
	return m.existsCPF, m.existsCPFErr
}

func (m *checkCPFDocumentRepositoryMock) GetPrimaryDocumentByCustomerID(
	ctx context.Context,
	customerID uuid.UUID,
) (*domain.CustomerDocument, error) {
	return nil, nil
}

func (m *checkCPFDocumentRepositoryMock) GetCPFByCustomerID(
	ctx context.Context,
	customerID uuid.UUID,
) (*domain.CustomerDocument, error) {
	return nil, nil
}

func TestCheckCPFUseCase_Execute_Available(t *testing.T) {
	repo := &checkCPFDocumentRepositoryMock{}
	uc := NewCheckCPFUseCase(repo)

	output, err := uc.Execute(context.Background(), CheckCPFInput{CPF: "123.456.789-09"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if output.CPF != "12345678909" {
		t.Fatalf("expected normalized cpf, got %q", output.CPF)
	}
	if output.Exists {
		t.Fatal("expected exists false")
	}
	if !output.Available {
		t.Fatal("expected available true")
	}
	if repo.existsCPFCalls != 1 {
		t.Fatalf("expected ExistsCPF to be called once, got %d", repo.existsCPFCalls)
	}
	if repo.existsCPFValue != "12345678909" {
		t.Fatalf("expected ExistsCPF value 12345678909, got %q", repo.existsCPFValue)
	}
}

func TestCheckCPFUseCase_Execute_AlreadyExists(t *testing.T) {
	repo := &checkCPFDocumentRepositoryMock{existsCPF: true}
	uc := NewCheckCPFUseCase(repo)

	output, err := uc.Execute(context.Background(), CheckCPFInput{CPF: "12345678909"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if output == nil {
		t.Fatal("expected output, got nil")
	}
	if !output.Exists {
		t.Fatal("expected exists true")
	}
	if output.Available {
		t.Fatal("expected available false")
	}
}

func TestCheckCPFUseCase_Execute_InvalidCPF(t *testing.T) {
	repo := &checkCPFDocumentRepositoryMock{}
	uc := NewCheckCPFUseCase(repo)

	output, err := uc.Execute(context.Background(), CheckCPFInput{CPF: "12345678901"})
	if !errors.Is(err, domain.ErrCPFInvalid) {
		t.Fatalf("expected ErrCPFInvalid, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
	if repo.existsCPFCalls != 0 {
		t.Fatalf("expected ExistsCPF not to be called, got %d", repo.existsCPFCalls)
	}
}

func TestCheckCPFUseCase_Execute_RequiredCPF(t *testing.T) {
	repo := &checkCPFDocumentRepositoryMock{}
	uc := NewCheckCPFUseCase(repo)

	output, err := uc.Execute(context.Background(), CheckCPFInput{CPF: "   "})
	if !errors.Is(err, domain.ErrCPFRequired) {
		t.Fatalf("expected ErrCPFRequired, got %v", err)
	}
	if output != nil {
		t.Fatalf("expected nil output, got %+v", output)
	}
}

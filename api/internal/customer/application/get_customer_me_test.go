package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/customer/domain"
)

type customerRepositoryGetByIDMock struct {
	profile *domain.CustomerProfile
	err     error
}

func (m *customerRepositoryGetByIDMock) Create(
	ctx context.Context,
	c *domain.Customer,
) error {
	return nil
}

func (m *customerRepositoryGetByIDMock) Exists(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	return false, nil
}

func (m *customerRepositoryGetByIDMock) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.CustomerProfile, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.profile, nil
}

func TestGetCustomerMe_Execute_Success(t *testing.T) {
	customerID := uuid.New()
	birthDate := time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC)
	repo := &customerRepositoryGetByIDMock{profile: &domain.CustomerProfile{
		Customer: domain.Customer{
			ID:        customerID,
			Name:      "Maria Silva",
			BirthDate: birthDate,
			CreatedAt: time.Now().UTC(),
		},
		Email: "maria@example.com",
		CPF:   "12345678901",
	}}
	uc := NewGetCustomerMe(repo)

	got, err := uc.Execute(
		context.Background(),
		GetCustomerMeInput{CustomerID: customerID},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got == nil {
		t.Fatal("expected customer, got nil")
	}
	if got.Customer.ID != customerID {
		t.Fatalf("expected customer ID %v, got %v", customerID, got.Customer.ID)
	}
	if got.Email != "maria@example.com" {
		t.Fatalf("expected email %q, got %q", "maria@example.com", got.Email)
	}
	if got.CPF != "12345678901" {
		t.Fatalf("expected cpf %q, got %q", "12345678901", got.CPF)
	}
}

func TestGetCustomerMe_Execute_InvalidWhenCustomerIDMissing(t *testing.T) {
	repo := &customerRepositoryGetByIDMock{}
	uc := NewGetCustomerMe(repo)

	got, err := uc.Execute(
		context.Background(),
		GetCustomerMeInput{CustomerID: uuid.Nil},
	)
	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}
	if got != nil {
		t.Fatalf("expected nil customer, got %+v", got)
	}
}

func TestGetCustomerMe_Execute_NotFound(t *testing.T) {
	repo := &customerRepositoryGetByIDMock{}
	uc := NewGetCustomerMe(repo)

	got, err := uc.Execute(
		context.Background(),
		GetCustomerMeInput{CustomerID: uuid.New()},
	)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrNotFound, err)
	}
	if got != nil {
		t.Fatalf("expected nil customer, got %+v", got)
	}
}

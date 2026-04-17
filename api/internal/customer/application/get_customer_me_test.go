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
	customer *domain.Customer
	email    string
	err      error
}

func (m *customerRepositoryGetByIDMock) Create(ctx context.Context, c *domain.Customer) error {
	return nil
}

func (m *customerRepositoryGetByIDMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, string, error) {
	if m.err != nil {
		return nil, "", m.err
	}
	return m.customer, m.email, nil
}

func TestGetCustomerMe_Execute_Success(t *testing.T) {
	customerID := uuid.New()
	repo := &customerRepositoryGetByIDMock{customer: &domain.Customer{
		ID:        customerID,
		Name:      "Maria Silva",
		CPF:       "12345678901",
		CreatedAt: time.Now().UTC(),
	}, email: "maria@example.com"}
	uc := NewGetCustomerMe(repo)

	got, email, err := uc.Execute(context.Background(), GetCustomerMeInput{CustomerID: customerID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got == nil {
		t.Fatal("expected customer, got nil")
	}
	if got.ID != customerID {
		t.Fatalf("expected customer ID %v, got %v", customerID, got.ID)
	}
	if email != "maria@example.com" {
		t.Fatalf("expected email %q, got %q", "maria@example.com", email)
	}
}

func TestGetCustomerMe_Execute_InvalidWhenCustomerIDMissing(t *testing.T) {
	repo := &customerRepositoryGetByIDMock{}
	uc := NewGetCustomerMe(repo)

	got, email, err := uc.Execute(context.Background(), GetCustomerMeInput{CustomerID: uuid.Nil})
	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}
	if got != nil {
		t.Fatalf("expected nil customer, got %+v", got)
	}
	if email != "" {
		t.Fatalf("expected empty email, got %q", email)
	}
}

func TestGetCustomerMe_Execute_NotFound(t *testing.T) {
	repo := &customerRepositoryGetByIDMock{}
	uc := NewGetCustomerMe(repo)

	got, email, err := uc.Execute(context.Background(), GetCustomerMeInput{CustomerID: uuid.New()})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrNotFound, err)
	}
	if got != nil {
		t.Fatalf("expected nil customer, got %+v", got)
	}
	if email != "" {
		t.Fatalf("expected empty email, got %q", email)
	}
}

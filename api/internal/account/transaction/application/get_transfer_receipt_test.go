package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type transferReceiptRepositoryMock struct {
	receipt *domain.TransferReceipt
	err     error
	calls   int
	gotRef  uuid.UUID
}

func (m *transferReceiptRepositoryMock) GetByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (m *transferReceiptRepositoryMock) GetByBranchAndNumber(ctx context.Context, branch, number string) (*domain.Account, error) {
	return nil, nil
}

func (m *transferReceiptRepositoryMock) GetTransferReceiptByReference(ctx context.Context, referenceID uuid.UUID) (*domain.TransferReceipt, error) {
	m.calls++
	m.gotRef = referenceID
	return m.receipt, m.err
}

func (m *transferReceiptRepositoryMock) IncreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *transferReceiptRepositoryMock) DecreaseBalance(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	return 0, nil
}

func (m *transferReceiptRepositoryMock) CreateTransaction(ctx context.Context, tx *domain.Transaction) error {
	return nil
}

func (m *transferReceiptRepositoryMock) GetTransactionByIdempotencyKey(ctx context.Context, accountID uuid.UUID, key string) (*domain.Transaction, error) {
	return nil, nil
}

func (m *transferReceiptRepositoryMock) GetTransactionByReference(ctx context.Context, accountID uuid.UUID, referenceID uuid.UUID, typeName domain.TransactionType) (*domain.Transaction, error) {
	return nil, nil
}

func (m *transferReceiptRepositoryMock) WithTransaction(ctx context.Context, fn func(tx domain.Tx) error) error {
	return errors.New("not implemented")
}

func TestGetTransferReceipt_Execute_SuccessForSourceCustomer(t *testing.T) {
	sourceCustomerID := uuid.New()
	destinationCustomerID := uuid.New()
	referenceID := uuid.New()
	operationDate := time.Date(2026, 5, 6, 10, 30, 0, 0, time.UTC)
	repo := &transferReceiptRepositoryMock{
		receipt: &domain.TransferReceipt{
			OperationType:            domain.TransactionTransferOut,
			Amount:                   2500,
			Status:                   "completed",
			TransactionReference:     referenceID,
			OperationDate:            operationDate,
			SourceCustomerID:         sourceCustomerID,
			SourceBranch:             "0001",
			SourceAccountNumber:      "123456",
			DestinationCustomerID:    destinationCustomerID,
			DestinationBranch:        "0002",
			DestinationAccountNumber: "654321",
			RecipientName:            "Maria Silva",
		},
	}
	uc := NewGetTransferReceipt(repo)

	result, err := uc.Execute(context.Background(), GetTransferReceiptInput{
		User:                 testCustomerUser(sourceCustomerID),
		TransactionReference: referenceID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}

	if result.TransactionReference != referenceID.String() {
		t.Fatalf("expected reference %s, got %s", referenceID, result.TransactionReference)
	}

	if result.OperationType != string(domain.TransactionTransferOut) || result.Amount != 2500 || result.Status != "completed" {
		t.Fatalf("unexpected receipt result: %+v", result)
	}

	if result.OperationDate != operationDate {
		t.Fatalf("expected operation date %s, got %s", operationDate, result.OperationDate)
	}

	if result.RecipientName != "Maria Silva" {
		t.Fatalf("expected recipient name %q, got %q", "Maria Silva", result.RecipientName)
	}

	if repo.calls != 1 || repo.gotRef != referenceID {
		t.Fatalf("expected repository lookup once with %s, got calls=%d ref=%s", referenceID, repo.calls, repo.gotRef)
	}
}

func TestGetTransferReceipt_Execute_SuccessForDestinationCustomer(t *testing.T) {
	sourceCustomerID := uuid.New()
	destinationCustomerID := uuid.New()
	referenceID := uuid.New()
	repo := &transferReceiptRepositoryMock{
		receipt: &domain.TransferReceipt{
			OperationType:            domain.TransactionTransferOut,
			Amount:                   2500,
			Status:                   "completed",
			TransactionReference:     referenceID,
			OperationDate:            time.Now().UTC(),
			SourceCustomerID:         sourceCustomerID,
			DestinationCustomerID:    destinationCustomerID,
			DestinationBranch:        "0002",
			DestinationAccountNumber: "654321",
			RecipientName:            "Maria Silva",
		},
	}
	uc := NewGetTransferReceipt(repo)

	result, err := uc.Execute(context.Background(), GetTransferReceiptInput{
		User:                 testCustomerUser(destinationCustomerID),
		TransactionReference: referenceID,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected result, got nil")
	}
}

func TestGetTransferReceipt_Execute_Forbidden(t *testing.T) {
	repo := &transferReceiptRepositoryMock{
		receipt: &domain.TransferReceipt{
			TransactionReference:  referenceIDForTest(),
			SourceCustomerID:      uuid.New(),
			DestinationCustomerID: uuid.New(),
		},
	}
	uc := NewGetTransferReceipt(repo)

	result, err := uc.Execute(context.Background(), GetTransferReceiptInput{
		User:                 testCustomerUser(uuid.New()),
		TransactionReference: repo.receipt.TransactionReference,
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}
}

func TestGetTransferReceipt_Execute_NotFound(t *testing.T) {
	referenceID := uuid.New()
	repo := &transferReceiptRepositoryMock{err: domain.ErrTransactionNotFound}
	uc := NewGetTransferReceipt(repo)

	result, err := uc.Execute(context.Background(), GetTransferReceiptInput{
		User:                 testCustomerUser(uuid.New()),
		TransactionReference: referenceID,
	})

	if !errors.Is(err, domain.ErrTransactionNotFound) {
		t.Fatalf("expected error %v, got %v", domain.ErrTransactionNotFound, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}
}

func TestGetTransferReceipt_Execute_InvalidReference(t *testing.T) {
	repo := &transferReceiptRepositoryMock{}
	uc := NewGetTransferReceipt(repo)

	result, err := uc.Execute(context.Background(), GetTransferReceiptInput{
		User:                 testCustomerUser(uuid.New()),
		TransactionReference: uuid.Nil,
	})

	if !errors.Is(err, domain.ErrInvalidData) {
		t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
	}

	if result != nil {
		t.Fatalf("expected result to be nil, got %+v", result)
	}

	if repo.calls != 0 {
		t.Fatalf("expected repository not to be called, got %d calls", repo.calls)
	}
}

func referenceIDForTest() uuid.UUID {
	return uuid.New()
}

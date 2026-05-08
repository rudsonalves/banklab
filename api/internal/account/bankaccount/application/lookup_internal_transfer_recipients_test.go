package application

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
)

type lookupRecipientsRepositoryMock struct {
	byAccountCalls  int
	byAccountBranch string
	byAccountNumber string
	byAccountValue  []domain.TransferRecipient
	byAccountErr    error

	byDocumentCalls  int
	byDocumentValue  string
	byDocumentResult []domain.TransferRecipient
	byDocumentErr    error
}

func (m *lookupRecipientsRepositoryMock) Create(ctx context.Context, account *domain.Account) error {
	return nil
}

func (m *lookupRecipientsRepositoryMock) ListByCustomerID(ctx context.Context, customerID uuid.UUID) ([]domain.Account, error) {
	return nil, nil
}

func (m *lookupRecipientsRepositoryMock) FindTransferRecipientsByBranchAndNumber(ctx context.Context, branch, number string) ([]domain.TransferRecipient, error) {
	m.byAccountCalls++
	m.byAccountBranch = branch
	m.byAccountNumber = number
	if m.byAccountErr != nil {
		return nil, m.byAccountErr
	}
	return m.byAccountValue, nil
}

func (m *lookupRecipientsRepositoryMock) FindTransferRecipientsByDocument(ctx context.Context, document string) ([]domain.TransferRecipient, error) {
	m.byDocumentCalls++
	m.byDocumentValue = document
	if m.byDocumentErr != nil {
		return nil, m.byDocumentErr
	}
	return m.byDocumentResult, nil
}

func (m *lookupRecipientsRepositoryMock) ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error) {
	return false, nil
}

func (m *lookupRecipientsRepositoryMock) NextAccountNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (m *lookupRecipientsRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func TestLookupInternalTransferRecipients_Execute_ForbiddenWhenUserMissing(t *testing.T) {
	repo := &lookupRecipientsRepositoryMock{}
	uc := NewLookupInternalTransferRecipients(repo)

	recipients, err := uc.Execute(context.Background(), LookupInternalTransferRecipientsInput{
		Branch:        "0001",
		AccountNumber: "12345678",
	})

	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expected error %v, got %v", domain.ErrForbidden, err)
	}
	if recipients != nil {
		t.Fatalf("expected nil recipients, got %+v", recipients)
	}
	if repo.byAccountCalls != 0 || repo.byDocumentCalls != 0 {
		t.Fatalf("expected repository not to be called")
	}
}

func TestLookupInternalTransferRecipients_Execute_ByBranchAndAccountNumber(t *testing.T) {
	accountID := uuid.New()
	repo := &lookupRecipientsRepositoryMock{
		byAccountValue: []domain.TransferRecipient{
			{
				AccountID:      accountID,
				HolderName:     "Maria Silva",
				MaskedDocument: "***.456.789-**",
				Branch:         "0001",
				AccountNumber:  "12345678",
			},
		},
	}
	uc := NewLookupInternalTransferRecipients(repo)

	recipients, err := uc.Execute(context.Background(), LookupInternalTransferRecipientsInput{
		User:          testCustomerUser(uuid.New()),
		Branch:        " 0001 ",
		AccountNumber: "12345678-9",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.byAccountCalls != 1 {
		t.Fatalf("expected account lookup once, got %d", repo.byAccountCalls)
	}
	if repo.byAccountBranch != "0001" {
		t.Fatalf("expected normalized branch %q, got %q", "0001", repo.byAccountBranch)
	}
	if repo.byAccountNumber != "123456789" {
		t.Fatalf("expected normalized account number %q, got %q", "123456789", repo.byAccountNumber)
	}
	if len(recipients) != 1 || recipients[0].AccountID != accountID {
		t.Fatalf("unexpected recipients: %+v", recipients)
	}
}

func TestLookupInternalTransferRecipients_Execute_ByDocumentWithMultipleAccounts(t *testing.T) {
	firstID := uuid.New()
	secondID := uuid.New()
	repo := &lookupRecipientsRepositoryMock{
		byDocumentResult: []domain.TransferRecipient{
			{AccountID: firstID, HolderName: "Maria Silva"},
			{AccountID: secondID, HolderName: "Maria Silva"},
		},
	}
	uc := NewLookupInternalTransferRecipients(repo)

	recipients, err := uc.Execute(context.Background(), LookupInternalTransferRecipientsInput{
		User:     testCustomerUser(uuid.New()),
		Document: "123.456.789-01",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.byDocumentCalls != 1 {
		t.Fatalf("expected document lookup once, got %d", repo.byDocumentCalls)
	}
	if repo.byDocumentValue != "12345678901" {
		t.Fatalf("expected normalized document %q, got %q", "12345678901", repo.byDocumentValue)
	}
	if len(recipients) != 2 {
		t.Fatalf("expected 2 recipients, got %d", len(recipients))
	}
}

func TestLookupInternalTransferRecipients_Execute_NoResults(t *testing.T) {
	repo := &lookupRecipientsRepositoryMock{byDocumentResult: []domain.TransferRecipient{}}
	uc := NewLookupInternalTransferRecipients(repo)

	recipients, err := uc.Execute(context.Background(), LookupInternalTransferRecipientsInput{
		User:     testCustomerUser(uuid.New()),
		Document: "12345678901",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if recipients == nil {
		t.Fatal("expected empty recipients slice, got nil")
	}
	if len(recipients) != 0 {
		t.Fatalf("expected no recipients, got %d", len(recipients))
	}
}

func TestLookupInternalTransferRecipients_Execute_InvalidQueryModes(t *testing.T) {
	tests := []struct {
		name  string
		input LookupInternalTransferRecipientsInput
	}{
		{
			name:  "missing query",
			input: LookupInternalTransferRecipientsInput{},
		},
		{
			name:  "branch without account number",
			input: LookupInternalTransferRecipientsInput{Branch: "0001"},
		},
		{
			name:  "account number without branch",
			input: LookupInternalTransferRecipientsInput{AccountNumber: "12345678"},
		},
		{
			name: "mixed modes",
			input: LookupInternalTransferRecipientsInput{
				Branch:        "0001",
				AccountNumber: "12345678",
				Document:      "12345678901",
			},
		},
		{
			name:  "invalid document length",
			input: LookupInternalTransferRecipientsInput{Document: "123"},
		},
		{
			name:  "cnpj document is out of scope",
			input: LookupInternalTransferRecipientsInput{Document: "12.345.678/0001-90"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &lookupRecipientsRepositoryMock{}
			uc := NewLookupInternalTransferRecipients(repo)
			tt.input.User = testCustomerUser(uuid.New())

			recipients, err := uc.Execute(context.Background(), tt.input)

			if !errors.Is(err, domain.ErrInvalidData) {
				t.Fatalf("expected error %v, got %v", domain.ErrInvalidData, err)
			}
			if recipients != nil {
				t.Fatalf("expected nil recipients, got %+v", recipients)
			}
			if repo.byAccountCalls != 0 || repo.byDocumentCalls != 0 {
				t.Fatalf("expected repository not to be called")
			}
		})
	}
}

func TestLookupInternalTransferRecipients_Execute_PropagatesRepositoryFailure(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	repo := &lookupRecipientsRepositoryMock{byDocumentErr: expectedErr}
	uc := NewLookupInternalTransferRecipients(repo)

	recipients, err := uc.Execute(context.Background(), LookupInternalTransferRecipientsInput{
		User:     testCustomerUser(uuid.New()),
		Document: "12345678901",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v, got %v", expectedErr, err)
	}
	if recipients != nil {
		t.Fatalf("expected nil recipients, got %+v", recipients)
	}
}

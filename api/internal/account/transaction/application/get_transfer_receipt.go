package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

type GetTransferReceipt struct {
	repo domain.Repository
}

// NewGetTransferReceipt creates a new instance of GetTransferReceipt use case.
func NewGetTransferReceipt(repo domain.Repository) *GetTransferReceipt {
	return &GetTransferReceipt{repo: repo}
}

type GetTransferReceiptInput struct {
	User                 *authdomain.AuthenticatedUser
	TransactionReference uuid.UUID
}

type TransferReceiptResult struct {
	OperationType            string
	Amount                   int64
	Status                   string
	TransactionReference     string
	OperationDate            time.Time
	SourceBranch             string
	SourceAccountNumber      string
	DestinationBranch        string
	DestinationAccountNumber string
	RecipientName            string
	Description              *string
}

// Execute retrieves the transfer receipt based on the provided input.
func (uc *GetTransferReceipt) Execute(ctx context.Context, input GetTransferReceiptInput) (*TransferReceiptResult, error) {
	if input.TransactionReference == uuid.Nil {
		return nil, domain.ErrInvalidData
	}

	receipt, err := uc.repo.GetTransferReceiptByReference(ctx, input.TransactionReference)
	if err != nil {
		return nil, err
	}

	if !canAccessTransferReceipt(input.User, receipt) {
		return nil, domain.ErrForbidden
	}

	operationType := receipt.OperationType
	if input.User != nil &&
		input.User.CustomerID != nil &&
		*input.User.CustomerID == receipt.DestinationCustomerID &&
		*input.User.CustomerID != receipt.SourceCustomerID {
		operationType = domain.TransactionTransferIn
	}

	return &TransferReceiptResult{
		OperationType:            string(operationType),
		Amount:                   receipt.Amount,
		Status:                   receipt.Status,
		TransactionReference:     receipt.TransactionReference.String(),
		OperationDate:            receipt.OperationDate,
		SourceBranch:             receipt.SourceBranch,
		SourceAccountNumber:      receipt.SourceAccountNumber,
		DestinationBranch:        receipt.DestinationBranch,
		DestinationAccountNumber: receipt.DestinationAccountNumber,
		RecipientName:            receipt.RecipientName,
		Description:              receipt.Description,
	}, nil
}

// canAccessTransferReceipt checks if the user has access to the transfer receipt.
func canAccessTransferReceipt(user *authdomain.AuthenticatedUser, receipt *domain.TransferReceipt) bool {
	if user == nil || receipt == nil {
		return false
	}

	if user.Role == authdomain.RoleAdmin {
		return true
	}

	if user.CustomerID == nil {
		return false
	}

	return *user.CustomerID == receipt.SourceCustomerID || *user.CustomerID == receipt.DestinationCustomerID
}

package delivery

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	transactionapp "github.com/seu-usuario/bank-api/internal/account/transaction/application"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type depositUseCase interface {
	Execute(ctx context.Context, input transactionapp.DepositInput) (*domain.Account, error)
}

type withdrawUseCase interface {
	Execute(ctx context.Context, input transactionapp.WithdrawInput) (*domain.Account, error)
}

type transferUseCase interface {
	Execute(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error)
}

type transferReceiptUseCase interface {
	Execute(ctx context.Context, input transactionapp.GetTransferReceiptInput) (*transactionapp.TransferReceiptResult, error)
}

type Handler struct {
	deposit  depositUseCase
	withdraw withdrawUseCase
	transfer transferUseCase
	receipt  transferReceiptUseCase
}

func New(deposit depositUseCase, withdraw withdrawUseCase, transfer transferUseCase, receipt transferReceiptUseCase) *Handler {
	return &Handler{
		deposit:  deposit,
		withdraw: withdraw,
		transfer: transfer,
		receipt:  receipt,
	}
}

// TransferReceipt handles the HTTP request for retrieving a transfer receipt.
func (h *Handler) TransferReceipt(w http.ResponseWriter, r *http.Request) {
	if h.receipt == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := requireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	referenceIDStr := r.PathValue("transaction_reference")
	referenceID, err := uuid.Parse(referenceIDStr)
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	result, err := h.receipt.Execute(r.Context(), transactionapp.GetTransferReceiptInput{
		User:                 user,
		TransactionReference: referenceID,
	})
	if err != nil {
		log.Printf("event=get_transfer_receipt error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, TransferReceiptData{
		OperationType:            result.OperationType,
		Amount:                   result.Amount,
		Status:                   result.Status,
		TransactionReference:     result.TransactionReference,
		OperationDate:            result.OperationDate.Format(time.RFC3339),
		SourceBranch:             result.SourceBranch,
		SourceAccountNumber:      result.SourceAccountNumber,
		DestinationBranch:        result.DestinationBranch,
		DestinationAccountNumber: result.DestinationAccountNumber,
		RecipientName:            result.RecipientName,
		Description:              result.Description,
	})
}

// Deposit handles the HTTP request for depositing funds into an account.
func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	if h.deposit == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := requireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	accountID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	account, err := h.deposit.Execute(r.Context(), transactionapp.DepositInput{
		User:      user,
		AccountID: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		log.Printf("event=deposit error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":      account.ID.String(),
		"balance": account.Balance,
	})
}

// Withdraw handles the HTTP request for withdrawing funds from an account.
func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if h.withdraw == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := requireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	accountID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	account, err := h.withdraw.Execute(r.Context(), transactionapp.WithdrawInput{
		User:      user,
		AccountID: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		log.Printf("event=withdraw error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"id":      account.ID.String(),
		"balance": account.Balance,
	})
}

// Transfer handles the HTTP request for transferring funds between accounts.
func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if h.transfer == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := requireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	var req TransferRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	if req.Amount <= 0 {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidAmount))
		return
	}

	// if req.IdempotencyKey == "" {
	// 	sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
	// 	return
	// }

	if req.FromAccountBranch == "" ||
		req.FromAccountNumber == "" ||
		req.ToAccountBranch == "" ||
		req.ToAccountNumber == "" {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	result, err := h.transfer.Execute(r.Context(), transactionapp.TransferInput{
		User:              user,
		FromAccountBranch: req.FromAccountBranch,
		FromAccountNumber: req.FromAccountNumber,
		ToAccountBranch:   req.ToAccountBranch,
		ToAccountNumber:   req.ToAccountNumber,
		Amount:            req.Amount,
		IdempotencyKey:    req.IdempotencyKey,
	})
	if err != nil {
		log.Printf("event=transfer error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, TransferData{
		FromAccountBranch:    req.FromAccountBranch,
		FromAccountNumber:    req.FromAccountNumber,
		TransactionReference: result.TransactionReference.String(),
		ToAccountBranch:      req.ToAccountBranch,
		ToAccountNumber:      req.ToAccountNumber,
		Amount:               result.Amount,
		FromBalance:          result.FromBalance,
		ToBalance:            result.ToBalance,
	})
}

// requireUser retrieves the authenticated user from the context and returns an
// error if not found.
func requireUser(ctx context.Context) (*authdomain.AuthenticatedUser, error) {
	user, ok := sharedauthctx.GetAuthenticatedUser(ctx)
	if !ok || user == nil {
		return nil, authdomain.ErrUnauthorized
	}

	return user, nil
}

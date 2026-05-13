package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	transactionapp "github.com/seu-usuario/bank-api/internal/account/transaction/application"
	"github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type depositUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input transactionapp.DepositInput) (*domain.Account, error)
}

type withdrawUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input transactionapp.WithdrawInput) (*domain.Account, error)
}

type transferUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error)
}

type transferReceiptUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input transactionapp.GetTransferReceiptInput) (*transactionapp.TransferReceiptResult, error)
}

func (m *depositUseCaseMock) Execute(ctx context.Context, input transactionapp.DepositInput) (*domain.Account, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *withdrawUseCaseMock) Execute(ctx context.Context, input transactionapp.WithdrawInput) (*domain.Account, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *transferUseCaseMock) Execute(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *transferReceiptUseCaseMock) Execute(ctx context.Context, input transactionapp.GetTransferReceiptInput) (*transactionapp.TransferReceiptResult, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func stringPtr(value string) *string {
	return &value
}

func TestHandler_Deposit_MissingAuth(t *testing.T) {
	h := &Handler{deposit: &depositUseCaseMock{}}
	accountID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/terminal/accounts/"+accountID.String()+"/deposit", strings.NewReader(`{"amount":100}`))
	req.SetPathValue("id", accountID.String())
	rec := httptest.NewRecorder()

	h.Deposit(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandler_Deposit_AccountInactive(t *testing.T) {
	customerID := uuid.New()
	depositUC := &depositUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.DepositInput) (*domain.Account, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrAccountInactive
		},
	}
	h := &Handler{deposit: depositUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/terminal/accounts/"+accountID.String()+"/deposit", strings.NewReader(`{"amount":100}`))
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Deposit(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "ACCOUNT_INACTIVE" {
		t.Fatalf("expected error code %q, got %q", "ACCOUNT_INACTIVE", got.Error.Code)
	}

	if depositUC.executeCalls != 1 {
		t.Fatalf("expected use case Execute to be called once, got %d calls", depositUC.executeCalls)
	}
}

func TestHandler_Deposit_Forbidden(t *testing.T) {
	depositUC := &depositUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.DepositInput) (*domain.Account, error) {
			return nil, domain.ErrForbidden
		},
	}
	h := &Handler{deposit: depositUC}
	accountID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/terminal/accounts/"+accountID.String()+"/deposit", strings.NewReader(`{"amount":100}`))
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, uuid.New())
	rec := httptest.NewRecorder()

	h.Deposit(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "FORBIDDEN" {
		t.Fatalf("expected error code %q, got %q", "FORBIDDEN", got.Error.Code)
	}
}

func TestHandler_Withdraw_InsufficientBalance(t *testing.T) {
	customerID := uuid.New()
	withdrawUC := &withdrawUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.WithdrawInput) (*domain.Account, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrInsufficientBalance
		},
	}
	h := &Handler{withdraw: withdrawUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/terminal/accounts/"+accountID.String()+"/withdraw", strings.NewReader(`{"amount":100}`))
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Withdraw(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INSUFFICIENT_FUNDS" {
		t.Fatalf("expected error code %q, got %q", "INSUFFICIENT_FUNDS", got.Error.Code)
	}
}

func TestHandler_Transfer_SameAccount(t *testing.T) {
	customerID := uuid.New()
	accountID := uuid.New()
	transferUC := &transferUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrSameAccountTransfer
		},
	}
	h := &Handler{transfer: transferUC}

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+accountID.String()+`","to_account_id":"`+accountID.String()+`","amount":100,"idempotency_key":"same-account-key"}`))
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "SAME_ACCOUNT_TRANSFER" {
		t.Fatalf("expected error code %q, got %q", "SAME_ACCOUNT_TRANSFER", got.Error.Code)
	}
}

func TestHandler_Transfer_MissingAuth(t *testing.T) {
	transferUC := &transferUseCaseMock{}
	h := &Handler{transfer: transferUC}
	fromAccountID := uuid.New()
	toAccountID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+fromAccountID.String()+`","to_account_id":"`+toAccountID.String()+`","amount":100,"idempotency_key":"missing-auth-key"}`))
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if transferUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", transferUC.executeCalls)
	}
}

func TestHandler_Transfer_UseCaseErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "forbidden", err: domain.ErrForbidden, wantStatus: http.StatusForbidden, wantCode: "FORBIDDEN"},
		{name: "account not found", err: domain.ErrAccountNotFound, wantStatus: http.StatusNotFound, wantCode: "ACCOUNT_NOT_FOUND"},
		{name: "insufficient funds", err: domain.ErrInsufficientBalance, wantStatus: http.StatusUnprocessableEntity, wantCode: "INSUFFICIENT_FUNDS"},
		{name: "inactive account", err: domain.ErrAccountInactive, wantStatus: http.StatusUnprocessableEntity, wantCode: "ACCOUNT_INACTIVE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fromAccountID := uuid.New()
			toAccountID := uuid.New()
			transferUC := &transferUseCaseMock{
				executeFn: func(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error) {
					if input.User == nil {
						return nil, errors.New("missing user")
					}
					return nil, tt.err
				},
			}
			h := &Handler{transfer: transferUC}

			req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+fromAccountID.String()+`","to_account_id":"`+toAccountID.String()+`","amount":100,"idempotency_key":"use-case-error-key"}`))
			req = testAuthenticatedRequest(req, uuid.New())
			rec := httptest.NewRecorder()

			h.Transfer(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rec.Code)
			}

			var got struct {
				Error struct {
					Code string `json:"code"`
				} `json:"error"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if got.Error.Code != tt.wantCode {
				t.Fatalf("expected error code %q, got %q", tt.wantCode, got.Error.Code)
			}

			if transferUC.executeCalls != 1 {
				t.Fatalf("expected use case Execute once, got %d calls", transferUC.executeCalls)
			}
		})
	}
}

func TestHandler_Transfer_SuccessWithNewPayload(t *testing.T) {
	customerID := uuid.New()
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	transactionReference := uuid.New()
	transferUC := &transferUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}

			if input.FromAccountID != fromAccountID {
				return nil, errors.New("unexpected source account")
			}

			if input.ToAccountID != toAccountID {
				return nil, errors.New("unexpected destination account")
			}

			if input.Amount != 100 {
				return nil, errors.New("unexpected amount")
			}

			if input.IdempotencyKey != "success-key" {
				return nil, errors.New("unexpected idempotency key")
			}

			if input.Description == nil || *input.Description != "Aluguel de maio" {
				return nil, errors.New("unexpected description")
			}

			return &transactionapp.TransferResult{
				FromAccountID:        fromAccountID,
				ToAccountID:          toAccountID,
				TransactionReference: transactionReference,
				Amount:               100,
				FromBalance:          900,
				ToBalance:            1100,
			}, nil
		},
	}
	h := &Handler{transfer: transferUC}

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+fromAccountID.String()+`","to_account_id":"`+toAccountID.String()+`","amount":100,"idempotency_key":"success-key","description":"Aluguel de maio"}`))
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got struct {
		Data struct {
			FromAccountID        string `json:"from_account_id"`
			TransactionReference string `json:"transaction_reference"`
			ToAccountID          string `json:"to_account_id"`
			Amount               int64  `json:"amount"`
			FromBalance          int64  `json:"from_balance"`
			ToBalance            int64  `json:"to_balance"`
		} `json:"data"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.FromAccountID != fromAccountID.String() {
		t.Fatalf("unexpected source account in response: %+v", got)
	}

	if got.Data.ToAccountID != toAccountID.String() {
		t.Fatalf("unexpected destination account in response: %+v", got)
	}

	if got.Data.TransactionReference != transactionReference.String() {
		t.Fatalf("expected transaction_reference %q, got %q", transactionReference.String(), got.Data.TransactionReference)
	}

	if got.Data.Amount != 100 || got.Data.FromBalance != 900 || got.Data.ToBalance != 1100 {
		t.Fatalf("unexpected balances/amount in response: %+v", got)
	}

	if transferUC.executeCalls != 1 {
		t.Fatalf("expected use case Execute to be called once, got %d calls", transferUC.executeCalls)
	}
}

func TestHandler_Transfer_InvalidJSON(t *testing.T) {
	customerID := uuid.New()
	transferUC := &transferUseCaseMock{}
	h := &Handler{transfer: transferUC}

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_branch":"0001",`))
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_REQUEST" {
		t.Fatalf("expected error code %q, got %q", "INVALID_REQUEST", got.Error.Code)
	}

	if transferUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", transferUC.executeCalls)
	}
}

func TestHandler_Transfer_InvalidAmount(t *testing.T) {
	customerID := uuid.New()
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	transferUC := &transferUseCaseMock{}
	h := &Handler{transfer: transferUC}

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+fromAccountID.String()+`","to_account_id":"`+toAccountID.String()+`","amount":0,"idempotency_key":"invalid-amount-key"}`))
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_AMOUNT" {
		t.Fatalf("expected error code %q, got %q", "INVALID_AMOUNT", got.Error.Code)
	}

	if transferUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", transferUC.executeCalls)
	}
}

func TestHandler_Transfer_MissingIdempotencyKey(t *testing.T) {
	customerID := uuid.New()
	fromAccountID := uuid.New()
	toAccountID := uuid.New()
	transferUC := &transferUseCaseMock{}
	h := &Handler{transfer: transferUC}

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+fromAccountID.String()+`","to_account_id":"`+toAccountID.String()+`","amount":100}`))
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_DATA" {
		t.Fatalf("expected error code %q, got %q", "INVALID_DATA", got.Error.Code)
	}

	if transferUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", transferUC.executeCalls)
	}
}

func TestHandler_Transfer_EmptyToAccountID(t *testing.T) {
	customerID := uuid.New()
	fromAccountID := uuid.New()
	transferUC := &transferUseCaseMock{}
	h := &Handler{transfer: transferUC}

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+fromAccountID.String()+`","to_account_id":"","amount":100,"idempotency_key":"invalid-data-key"}`))
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_DATA" {
		t.Fatalf("expected error code %q, got %q", "INVALID_DATA", got.Error.Code)
	}

	if transferUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", transferUC.executeCalls)
	}
}

func TestHandler_Transfer_LegacyPayloadRejected(t *testing.T) {
	customerID := uuid.New()
	transferUC := &transferUseCaseMock{}
	h := &Handler{transfer: transferUC}

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_branch":"0001","from_account_number":"123456","to_branch":"0002","to_account_number":"654321","amount":100,"idempotency_key":"legacy-key"}`))
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Transfer(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "INVALID_REQUEST" {
		t.Fatalf("expected error code %q, got %q", "INVALID_REQUEST", got.Error.Code)
	}

	if transferUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", transferUC.executeCalls)
	}
}

func TestHandler_TransferReceipt_Success(t *testing.T) {
	customerID := uuid.New()
	referenceID := uuid.New()
	operationDate := time.Date(2026, 5, 6, 12, 30, 0, 0, time.UTC)
	receiptUC := &transferReceiptUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.GetTransferReceiptInput) (*transactionapp.TransferReceiptResult, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			if input.TransactionReference != referenceID {
				return nil, errors.New("unexpected transaction reference")
			}
			return &transactionapp.TransferReceiptResult{
				OperationType:            string(domain.TransactionTransferOut),
				Amount:                   2500,
				Status:                   "completed",
				TransactionReference:     referenceID.String(),
				OperationDate:            operationDate,
				SourceBranch:             "0001",
				SourceAccountNumber:      "123456",
				DestinationBranch:        "0002",
				DestinationAccountNumber: "654321",
				RecipientName:            "Maria Silva",
				Description:              stringPtr("Aluguel de maio"),
			}, nil
		},
	}
	h := &Handler{receipt: receiptUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts/transfer/"+referenceID.String()+"/receipt", nil)
	req.SetPathValue("transaction_reference", referenceID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.TransferReceipt(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got struct {
		Data TransferReceiptData `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.TransactionReference != referenceID.String() {
		t.Fatalf("expected transaction_reference %q, got %q", referenceID.String(), got.Data.TransactionReference)
	}
	if got.Data.OperationDate != operationDate.Format(time.RFC3339) {
		t.Fatalf("expected operation_date %q, got %q", operationDate.Format(time.RFC3339), got.Data.OperationDate)
	}
	if got.Data.RecipientName != "Maria Silva" {
		t.Fatalf("expected recipient name %q, got %q", "Maria Silva", got.Data.RecipientName)
	}
	if got.Data.Description == nil || *got.Data.Description != "Aluguel de maio" {
		t.Fatalf("expected description %q, got %+v", "Aluguel de maio", got.Data.Description)
	}
	if receiptUC.executeCalls != 1 {
		t.Fatalf("expected use case Execute once, got %d calls", receiptUC.executeCalls)
	}
}

func TestHandler_TransferReceipt_MissingAuth(t *testing.T) {
	referenceID := uuid.New()
	receiptUC := &transferReceiptUseCaseMock{}
	h := &Handler{receipt: receiptUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts/transfer/"+referenceID.String()+"/receipt", nil)
	req.SetPathValue("transaction_reference", referenceID.String())
	rec := httptest.NewRecorder()

	h.TransferReceipt(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
	if receiptUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", receiptUC.executeCalls)
	}
}

func TestHandler_TransferReceipt_InvalidReference(t *testing.T) {
	customerID := uuid.New()
	receiptUC := &transferReceiptUseCaseMock{}
	h := &Handler{receipt: receiptUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts/transfer/not-a-uuid/receipt", nil)
	req.SetPathValue("transaction_reference", "not-a-uuid")
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.TransferReceipt(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if receiptUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", receiptUC.executeCalls)
	}
}

func TestHandler_TransferReceipt_NotFound(t *testing.T) {
	customerID := uuid.New()
	referenceID := uuid.New()
	receiptUC := &transferReceiptUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.GetTransferReceiptInput) (*transactionapp.TransferReceiptResult, error) {
			return nil, domain.ErrTransactionNotFound
		},
	}
	h := &Handler{receipt: receiptUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts/transfer/"+referenceID.String()+"/receipt", nil)
	req.SetPathValue("transaction_reference", referenceID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.TransferReceipt(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandler_TransferReceipt_Forbidden(t *testing.T) {
	customerID := uuid.New()
	referenceID := uuid.New()
	receiptUC := &transferReceiptUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.GetTransferReceiptInput) (*transactionapp.TransferReceiptResult, error) {
			return nil, domain.ErrForbidden
		},
	}
	h := &Handler{receipt: receiptUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts/transfer/"+referenceID.String()+"/receipt", nil)
	req.SetPathValue("transaction_reference", referenceID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.TransferReceipt(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

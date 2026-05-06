package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	transactionapp "github.com/seu-usuario/bank-api/internal/account/transaction/application"
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

func TestHandler_Deposit_MissingAuth(t *testing.T) {
	h := &Handler{deposit: &depositUseCaseMock{}}
	accountID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/accounts/"+accountID.String()+"/deposit", strings.NewReader(`{"amount":100}`))
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

	req := httptest.NewRequest(http.MethodPost, "/accounts/"+accountID.String()+"/deposit", strings.NewReader(`{"amount":100}`))
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
	req := httptest.NewRequest(http.MethodPost, "/accounts/"+accountID.String()+"/deposit", strings.NewReader(`{"amount":100}`))
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

	req := httptest.NewRequest(http.MethodPost, "/accounts/"+accountID.String()+"/withdraw", strings.NewReader(`{"amount":100}`))
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
	transferUC := &transferUseCaseMock{
		executeFn: func(ctx context.Context, input transactionapp.TransferInput) (*transactionapp.TransferResult, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrSameAccountTransfer
		},
	}
	h := &Handler{transfer: transferUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/accounts/transfer", strings.NewReader(`{"from_account_id":"`+accountID.String()+`","to_account_id":"`+accountID.String()+`","amount":100}`))
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

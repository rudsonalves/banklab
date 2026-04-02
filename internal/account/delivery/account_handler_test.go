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
	"github.com/seu-usuario/bank-api/internal/account/application"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

type createAccountUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input application.CreateAccountInput) (*domain.Account, error)
}

type depositUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input application.DepositInput) (*domain.Account, error)
}

func (m *createAccountUseCaseMock) Execute(ctx context.Context, input application.CreateAccountInput) (*domain.Account, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *depositUseCaseMock) Execute(ctx context.Context, input application.DepositInput) (*domain.Account, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func TestHandler_CreateAccount_InvalidJSON(t *testing.T) {
	h := &Handler{createAccount: nil}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader("{"))
	rec := httptest.NewRecorder()

	h.CreateAccount(rec, req)

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
}

func TestHandler_CreateAccount_InvalidCustomerID(t *testing.T) {
	uc := &createAccountUseCaseMock{}
	h := &Handler{createAccount: uc}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"customer_id":"invalid-uuid"}`))
	rec := httptest.NewRecorder()

	h.CreateAccount(rec, req)

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

	if uc.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", uc.executeCalls)
	}
}

func TestHandler_CreateAccount_CustomerNotFound(t *testing.T) {
	uc := &createAccountUseCaseMock{
		executeFn: func(ctx context.Context, input application.CreateAccountInput) (*domain.Account, error) {
			return nil, domain.ErrCustomerNotFound
		},
	}
	h := &Handler{createAccount: uc}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"customer_id":"`+uuid.New().String()+`"}`))
	rec := httptest.NewRecorder()

	h.CreateAccount(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}

	var got struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Error.Code != "CUSTOMER_NOT_FOUND" {
		t.Fatalf("expected error code %q, got %q", "CUSTOMER_NOT_FOUND", got.Error.Code)
	}
}

func TestHandler_CreateAccount_Success(t *testing.T) {
	inputCustomerID := uuid.New()
	returnedAccount := &domain.Account{
		ID:         uuid.New(),
		CustomerID: inputCustomerID,
		Number:     "12345678",
		Branch:     "0001",
		Balance:    0,
		Status:     domain.AccountActive,
	}

	uc := &createAccountUseCaseMock{
		executeFn: func(ctx context.Context, input application.CreateAccountInput) (*domain.Account, error) {
			if input.CustomerID != inputCustomerID {
				return nil, errors.New("unexpected customer id")
			}
			return returnedAccount, nil
		},
	}
	h := &Handler{createAccount: uc}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"customer_id":"`+inputCustomerID.String()+`"}`))
	rec := httptest.NewRecorder()

	h.CreateAccount(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var got struct {
		Data struct {
			ID         string `json:"id"`
			CustomerID string `json:"customer_id"`
			Number     string `json:"number"`
			Branch     string `json:"branch"`
			Balance    int64  `json:"balance"`
			Status     string `json:"status"`
		} `json:"data"`
		Error interface{} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.ID != returnedAccount.ID.String() {
		t.Fatalf("expected id %q, got %q", returnedAccount.ID.String(), got.Data.ID)
	}

	if got.Data.CustomerID != returnedAccount.CustomerID.String() {
		t.Fatalf("expected customer_id %q, got %q", returnedAccount.CustomerID.String(), got.Data.CustomerID)
	}

	if got.Data.Number != returnedAccount.Number {
		t.Fatalf("expected number %q, got %q", returnedAccount.Number, got.Data.Number)
	}

	if got.Data.Branch != returnedAccount.Branch {
		t.Fatalf("expected branch %q, got %q", returnedAccount.Branch, got.Data.Branch)
	}

	if got.Data.Balance != returnedAccount.Balance {
		t.Fatalf("expected balance %d, got %d", returnedAccount.Balance, got.Data.Balance)
	}

	if got.Data.Status != string(returnedAccount.Status) {
		t.Fatalf("expected status %q, got %q", string(returnedAccount.Status), got.Data.Status)
	}
}

func TestHandler_Deposit_AccountInactive(t *testing.T) {
	depositUC := &depositUseCaseMock{
		executeFn: func(ctx context.Context, input application.DepositInput) (*domain.Account, error) {
			return nil, domain.ErrAccountInactive
		},
	}
	h := &Handler{deposit: depositUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/accounts/"+accountID.String()+"/deposit", strings.NewReader(`{"amount":100}`))
	req.SetPathValue("id", accountID.String())
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

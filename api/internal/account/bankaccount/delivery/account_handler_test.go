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
	accountapp "github.com/seu-usuario/bank-api/internal/account/bankaccount/application"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
)

type createAccountUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input accountapp.CreateAccountInput) (*domain.Account, error)
}

type listAccountsUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input accountapp.ListAccountsInput) ([]domain.Account, error)
}

type balanceUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input accountapp.GetAccountBalanceInput) (*accountapp.AccountBalance, error)
}

func (m *createAccountUseCaseMock) Execute(ctx context.Context, input accountapp.CreateAccountInput) (*domain.Account, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *listAccountsUseCaseMock) Execute(ctx context.Context, input accountapp.ListAccountsInput) ([]domain.Account, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *balanceUseCaseMock) Execute(ctx context.Context, input accountapp.GetAccountBalanceInput) (*accountapp.AccountBalance, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func TestHandler_ListAccounts_MissingAuth(t *testing.T) {
	h := &Handler{listAccounts: &listAccountsUseCaseMock{}}
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rec := httptest.NewRecorder()

	h.ListAccounts(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandler_ListAccounts_RejectsQueryParameters(t *testing.T) {
	listUC := &listAccountsUseCaseMock{}
	h := &Handler{listAccounts: listUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts?status=active", nil)
	req = testAuthenticatedRequest(req, uuid.New())
	rec := httptest.NewRecorder()

	h.ListAccounts(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if listUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", listUC.executeCalls)
	}
}

func TestHandler_ListAccounts_Forbidden(t *testing.T) {
	listUC := &listAccountsUseCaseMock{
		executeFn: func(ctx context.Context, input accountapp.ListAccountsInput) ([]domain.Account, error) {
			return nil, domain.ErrForbidden
		},
	}
	h := &Handler{listAccounts: listUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	req = testAuthenticatedRequest(req, uuid.New())
	rec := httptest.NewRecorder()

	h.ListAccounts(rec, req)

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

func TestHandler_ListAccounts_Success(t *testing.T) {
	customerID := uuid.New()
	accountID := uuid.New()
	listUC := &listAccountsUseCaseMock{
		executeFn: func(ctx context.Context, input accountapp.ListAccountsInput) ([]domain.Account, error) {
			if input.User == nil || input.User.CustomerID == nil || *input.User.CustomerID != customerID {
				return nil, errors.New("unexpected user")
			}
			return []domain.Account{
				{
					ID:         accountID,
					CustomerID: customerID,
					Number:     "10000001",
					Branch:     "0001",
					Status:     domain.AccountActive,
				},
			}, nil
		},
	}
	h := &Handler{listAccounts: listUC}
	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.ListAccounts(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got struct {
		Data []struct {
			ID         string `json:"id"`
			CustomerID string `json:"customer_id"`
			Number     string `json:"number"`
			Branch     string `json:"branch"`
			Status     string `json:"status"`
		} `json:"data"`
		Error any `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if len(got.Data) != 1 {
		t.Fatalf("expected 1 account, got %d", len(got.Data))
	}

	if got.Data[0].ID != accountID.String() {
		t.Fatalf("expected id %q, got %q", accountID.String(), got.Data[0].ID)
	}

	if got.Data[0].CustomerID != customerID.String() {
		t.Fatalf("expected customer_id %q, got %q", customerID.String(), got.Data[0].CustomerID)
	}

	if got.Data[0].Number != "10000001" {
		t.Fatalf("expected number %q, got %q", "10000001", got.Data[0].Number)
	}

	if got.Data[0].Branch != "0001" {
		t.Fatalf("expected branch %q, got %q", "0001", got.Data[0].Branch)
	}

	if got.Data[0].Status != string(domain.AccountActive) {
		t.Fatalf("expected status %q, got %q", domain.AccountActive, got.Data[0].Status)
	}

	if got.Error != nil {
		t.Fatalf("expected nil error, got %#v", got.Error)
	}
}

func TestHandler_CreateAccount_InvalidJSON(t *testing.T) {
	h := &Handler{createAccount: nil}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader("{"))
	req = testAuthenticatedRequest(req, uuid.New())
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

func TestHandler_CreateAccount_RejectsUnknownField(t *testing.T) {
	uc := &createAccountUseCaseMock{}
	h := &Handler{createAccount: uc}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"customer_id":"invalid-uuid"}`))
	req = testAuthenticatedRequest(req, uuid.New())
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

	if uc.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", uc.executeCalls)
	}
}

func TestHandler_CreateAccount_CustomerNotFound(t *testing.T) {
	uc := &createAccountUseCaseMock{
		executeFn: func(ctx context.Context, input accountapp.CreateAccountInput) (*domain.Account, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrCustomerNotFound
		},
	}
	h := &Handler{createAccount: uc}
	customerID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{}`))
	req = testAuthenticatedRequest(req, customerID)
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
		executeFn: func(ctx context.Context, input accountapp.CreateAccountInput) (*domain.Account, error) {
			if input.User == nil || input.User.CustomerID == nil || *input.User.CustomerID != inputCustomerID {
				return nil, errors.New("unexpected user")
			}
			return returnedAccount, nil
		},
	}
	h := &Handler{createAccount: uc}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{}`))
	req = testAuthenticatedRequest(req, inputCustomerID)
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

func TestHandler_CreateAccount_SuccessWithEmptyBody(t *testing.T) {
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
		executeFn: func(ctx context.Context, input accountapp.CreateAccountInput) (*domain.Account, error) {
			if input.User == nil || input.User.CustomerID == nil || *input.User.CustomerID != inputCustomerID {
				return nil, errors.New("unexpected user")
			}
			return returnedAccount, nil
		},
	}
	h := &Handler{createAccount: uc}
	req := httptest.NewRequest(http.MethodPost, "/accounts", nil)
	req = testAuthenticatedRequest(req, inputCustomerID)
	rec := httptest.NewRecorder()

	h.CreateAccount(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestHandler_GetBalance_QueryParamsNotAllowed(t *testing.T) {
	balanceUC := &balanceUseCaseMock{}
	h := &Handler{balance: balanceUC}
	accountID := uuid.New()
	customerID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/balance?currency=BRL", nil)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.GetBalance(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if balanceUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", balanceUC.executeCalls)
	}
}

func TestHandler_GetBalance_Unauthorized(t *testing.T) {
	balanceUC := &balanceUseCaseMock{}
	h := &Handler{balance: balanceUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/balance", nil)
	req.SetPathValue("id", accountID.String())
	rec := httptest.NewRecorder()

	h.GetBalance(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}

	if balanceUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", balanceUC.executeCalls)
	}
}

func TestHandler_GetBalance_InvalidUUID(t *testing.T) {
	balanceUC := &balanceUseCaseMock{}
	h := &Handler{balance: balanceUC}
	customerID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/accounts/invalid-uuid/balance", nil)
	req.SetPathValue("id", "invalid-uuid")
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.GetBalance(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if balanceUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", balanceUC.executeCalls)
	}
}

func TestHandler_GetBalance_AccountNotFound(t *testing.T) {
	customerID := uuid.New()
	balanceUC := &balanceUseCaseMock{
		executeFn: func(ctx context.Context, input accountapp.GetAccountBalanceInput) (*accountapp.AccountBalance, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrAccountNotFound
		},
	}
	h := &Handler{balance: balanceUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/balance", nil)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.GetBalance(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandler_GetBalance_Forbidden(t *testing.T) {
	customerID := uuid.New()
	balanceUC := &balanceUseCaseMock{
		executeFn: func(ctx context.Context, input accountapp.GetAccountBalanceInput) (*accountapp.AccountBalance, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrForbidden
		},
	}
	h := &Handler{balance: balanceUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/balance", nil)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.GetBalance(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestHandler_GetBalance_Success(t *testing.T) {
	accountID := uuid.New()
	customerID := uuid.New()
	balanceUC := &balanceUseCaseMock{
		executeFn: func(ctx context.Context, input accountapp.GetAccountBalanceInput) (*accountapp.AccountBalance, error) {
			if input.AccountID != accountID {
				return nil, errors.New("unexpected account id")
			}
			if input.User == nil || input.User.CustomerID == nil || *input.User.CustomerID != customerID {
				return nil, errors.New("unexpected user")
			}

			return &accountapp.AccountBalance{
				AccountID: accountID,
				Balance:   12000,
			}, nil
		},
	}
	h := &Handler{balance: balanceUC}

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/balance", nil)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.GetBalance(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got struct {
		Data struct {
			AccountID string `json:"account_id"`
			Balance   int64  `json:"balance"`
		} `json:"data"`
		Error interface{} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.AccountID != accountID.String() {
		t.Fatalf("expected account_id %q, got %q", accountID.String(), got.Data.AccountID)
	}

	if got.Data.Balance != 12000 {
		t.Fatalf("expected balance %d, got %d", 12000, got.Data.Balance)
	}
}

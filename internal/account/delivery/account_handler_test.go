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

type withdrawUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input application.WithdrawInput) (*domain.Account, error)
}

type transferUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input application.TransferInput) (*application.TransferResult, error)
}

type statementUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input application.GetStatementInput) (*application.Statement, error)
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

func (m *withdrawUseCaseMock) Execute(ctx context.Context, input application.WithdrawInput) (*domain.Account, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *transferUseCaseMock) Execute(ctx context.Context, input application.TransferInput) (*application.TransferResult, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
}

func (m *statementUseCaseMock) Execute(ctx context.Context, input application.GetStatementInput) (*application.Statement, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
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

func TestHandler_CreateAccount_InvalidCustomerID(t *testing.T) {
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
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrCustomerNotFound
		},
	}
	h := &Handler{createAccount: uc}
	customerID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"customer_id":"`+customerID.String()+`"}`))
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
		executeFn: func(ctx context.Context, input application.CreateAccountInput) (*domain.Account, error) {
			if input.CustomerID != inputCustomerID {
				return nil, errors.New("unexpected customer id")
			}
			if input.User == nil || input.User.CustomerID == nil || *input.User.CustomerID != inputCustomerID {
				return nil, errors.New("unexpected user")
			}
			return returnedAccount, nil
		},
	}
	h := &Handler{createAccount: uc}
	req := httptest.NewRequest(http.MethodPost, "/accounts", strings.NewReader(`{"customer_id":"`+inputCustomerID.String()+`"}`))
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
		executeFn: func(ctx context.Context, input application.DepositInput) (*domain.Account, error) {
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
		executeFn: func(ctx context.Context, input application.DepositInput) (*domain.Account, error) {
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
		executeFn: func(ctx context.Context, input application.WithdrawInput) (*domain.Account, error) {
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
		executeFn: func(ctx context.Context, input application.TransferInput) (*application.TransferResult, error) {
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

func TestHandler_Statement_InvalidFromQuery(t *testing.T) {
	statementUC := &statementUseCaseMock{}
	h := &Handler{statement: statementUC}
	accountID := uuid.New()
	customerID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/statement?from=not-a-date", nil)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Statement(rec, req)

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

	if statementUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", statementUC.executeCalls)
	}
}

func TestHandler_Statement_AccountNotFound(t *testing.T) {
	customerID := uuid.New()
	statementUC := &statementUseCaseMock{
		executeFn: func(ctx context.Context, input application.GetStatementInput) (*application.Statement, error) {
			if input.User == nil {
				return nil, errors.New("missing user")
			}
			return nil, domain.ErrAccountNotFound
		},
	}
	h := &Handler{statement: statementUC}
	accountID := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/statement", nil)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Statement(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandler_Statement_CursorWithoutCursorID(t *testing.T) {
	statementUC := &statementUseCaseMock{}
	h := &Handler{statement: statementUC}
	accountID := uuid.New()
	customerID := uuid.New()
	cursor := time.Now().UTC().Truncate(time.Second)

	req := httptest.NewRequest(http.MethodGet, "/accounts/"+accountID.String()+"/statement?cursor="+cursor.Format(time.RFC3339), nil)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, customerID)
	rec := httptest.NewRecorder()

	h.Statement(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	if statementUC.executeCalls != 0 {
		t.Fatalf("expected use case Execute not to be called, got %d calls", statementUC.executeCalls)
	}
}

func TestHandler_Statement_Success(t *testing.T) {
	accountID := uuid.New()
	from := time.Now().UTC().Add(-24 * time.Hour).Truncate(time.Second)
	to := time.Now().UTC().Truncate(time.Second)
	cursorID := uuid.New()
	transactionID := uuid.New()
	referenceID := uuid.New().String()

	statementUC := &statementUseCaseMock{
		executeFn: func(ctx context.Context, input application.GetStatementInput) (*application.Statement, error) {
			if input.AccountID != accountID {
				return nil, errors.New("unexpected account id")
			}
			if input.User == nil || input.User.CustomerID == nil || *input.User.CustomerID != accountID {
				return nil, errors.New("unexpected user")
			}
			if input.Limit != 20 {
				return nil, errors.New("unexpected limit")
			}
			if input.Cursor == nil || !input.Cursor.Equal(to) {
				return nil, errors.New("unexpected cursor")
			}
			if input.CursorID == nil || *input.CursorID != cursorID {
				return nil, errors.New("unexpected cursor id")
			}
			if input.From == nil || !input.From.Equal(from) {
				return nil, errors.New("unexpected from")
			}
			if input.To == nil || !input.To.Equal(to) {
				return nil, errors.New("unexpected to")
			}

			return &application.Statement{
				AccountID: accountID.String(),
				Items: []application.StatementItem{
					{
						TransactionID: transactionID.String(),
						Type:          string(domain.TransactionDeposit),
						Amount:        100,
						BalanceAfter:  500,
						ReferenceID:   &referenceID,
						CreatedAt:     to,
					},
				},
				NextCursor: &application.StatementCursor{CreatedAt: to, ID: transactionID.String()},
			}, nil
		},
	}
	h := &Handler{statement: statementUC}

	req := httptest.NewRequest(
		http.MethodGet,
		"/accounts/"+accountID.String()+"/statement?limit=20&cursor="+to.Format(time.RFC3339)+"&cursor_id="+cursorID.String()+"&from="+from.Format(time.RFC3339)+"&to="+to.Format(time.RFC3339),
		nil,
	)
	req.SetPathValue("id", accountID.String())
	req = testAuthenticatedRequest(req, accountID)
	rec := httptest.NewRecorder()

	h.Statement(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got struct {
		Data struct {
			AccountID string `json:"account_id"`
			Items     []struct {
				TransactionID string  `json:"transaction_id"`
				Type          string  `json:"type"`
				Amount        int64   `json:"amount"`
				BalanceAfter  int64   `json:"balance_after"`
				ReferenceID   *string `json:"reference_id"`
			} `json:"items"`
			NextCursor *struct {
				CreatedAt time.Time `json:"created_at"`
				ID        string    `json:"id"`
			} `json:"next_cursor"`
		} `json:"data"`
		Error interface{} `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if got.Data.AccountID != accountID.String() {
		t.Fatalf("expected account_id %q, got %q", accountID.String(), got.Data.AccountID)
	}

	if len(got.Data.Items) != 1 {
		t.Fatalf("expected 1 statement item, got %d", len(got.Data.Items))
	}

	if got.Data.Items[0].TransactionID != transactionID.String() {
		t.Fatalf("expected transaction_id %q, got %q", transactionID.String(), got.Data.Items[0].TransactionID)
	}

	if got.Data.NextCursor == nil {
		t.Fatal("expected next_cursor to be non-nil")
	}

	if got.Data.NextCursor.ID != transactionID.String() {
		t.Fatalf("expected next_cursor id %q, got %q", transactionID.String(), got.Data.NextCursor.ID)
	}
}

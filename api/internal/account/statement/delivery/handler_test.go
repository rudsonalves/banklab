package delivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
	statementapp "github.com/seu-usuario/bank-api/internal/account/statement/application"
	transactiondomain "github.com/seu-usuario/bank-api/internal/account/transaction/domain"
)

type statementUseCaseMock struct {
	executeCalls int
	executeFn    func(ctx context.Context, input statementapp.GetStatementInput) (*statementapp.Statement, error)
}

func (m *statementUseCaseMock) Execute(ctx context.Context, input statementapp.GetStatementInput) (*statementapp.Statement, error) {
	m.executeCalls++
	if m.executeFn == nil {
		return nil, nil
	}
	return m.executeFn(ctx, input)
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
		executeFn: func(ctx context.Context, input statementapp.GetStatementInput) (*statementapp.Statement, error) {
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
		executeFn: func(ctx context.Context, input statementapp.GetStatementInput) (*statementapp.Statement, error) {
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

			return &statementapp.Statement{
				AccountID: accountID.String(),
				Items: []statementapp.StatementItem{
					{
						TransactionID: transactionID.String(),
						Type:          string(transactiondomain.TransactionDeposit),
						Amount:        100,
						BalanceAfter:  500,
						ReferenceID:   &referenceID,
						CreatedAt:     to,
					},
				},
				NextCursor: &statementapp.StatementCursor{CreatedAt: to, ID: transactionID.String()},
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

package delivery

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	accountapp "github.com/seu-usuario/bank-api/internal/account/application/account"
	statementapp "github.com/seu-usuario/bank-api/internal/account/application/statement"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

// CreateAccount handles the HTTP request for creating a new account.
// It validates the request, checks user authentication, and delegates
// the account creation to the application layer.
func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	var req CreateAccountRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil && !errors.Is(err, io.EOF) {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	input := accountapp.CreateAccountInput{
		User: user,
	}

	account, err := h.createAccount.Execute(r.Context(), input)
	if err != nil {
		log.Printf("event=create_account error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusCreated, AccountData{
		ID:         account.ID.String(),
		CustomerID: account.CustomerID.String(),
		Number:     account.Number,
		Branch:     account.Branch,
		Balance:    account.Balance,
		Status:     string(account.Status),
	})
}

// ListAccounts handles the HTTP request for listing accounts. It validates
// the request, checks user authentication, and retrieves the list of accounts
// accessible to the user from the application layer.
func (h *Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	if len(r.URL.Query()) > 0 {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	accounts, err := h.listAccounts.Execute(r.Context(), accountapp.ListAccountsInput{User: user})
	if err != nil {
		log.Printf("event=list_accounts error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	data := make([]AccountSummaryData, 0, len(accounts))
	for _, account := range accounts {
		data = append(data, AccountSummaryData{
			ID:         account.ID.String(),
			CustomerID: account.CustomerID.String(),
			Number:     account.Number,
			Branch:     account.Branch,
			Status:     string(account.Status),
		})
	}

	sharedhttp.WriteJSON(w, http.StatusOK, data)
}

// Statement handles the HTTP request for retrieving an account statement.
// It validates the request, checks user authentication, and delegates
// the statement retrieval to the application layer.
func (h *Handler) Statement(w http.ResponseWriter, r *http.Request) {
	if h.statement == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	limit, err := parseOptionalInt(r.URL.Query().Get("limit"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	cursor, err := parseOptionalTime(r.URL.Query().Get("cursor"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	cursorID, err := parseOptionalUUID(r.URL.Query().Get("cursor_id"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	if (cursor == nil) != (cursorID == nil) {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	from, err := parseOptionalTime(r.URL.Query().Get("from"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	to, err := parseOptionalTime(r.URL.Query().Get("to"))
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	result, err := h.statement.Execute(r.Context(), statementapp.GetStatementInput{
		User:      user,
		AccountID: accountID,
		Limit:     limit,
		Cursor:    cursor,
		CursorID:  cursorID,
		From:      from,
		To:        to,
	})
	if err != nil {
		log.Printf("event=get_statement error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	items := make([]StatementItemData, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, StatementItemData{
			TransactionID: item.TransactionID,
			Type:          item.Type,
			Amount:        item.Amount,
			BalanceAfter:  item.BalanceAfter,
			ReferenceID:   item.ReferenceID,
			CreatedAt:     item.CreatedAt,
		})
	}

	var nextCursor *StatementCursorData
	if result.NextCursor != nil {
		nextCursor = &StatementCursorData{
			CreatedAt: result.NextCursor.CreatedAt,
			ID:        result.NextCursor.ID,
		}
	}

	sharedhttp.WriteJSON(w, http.StatusOK, StatementData{
		AccountID:  result.AccountID,
		Items:      items,
		NextCursor: nextCursor,
	})
}

// GetBalance handles the HTTP request for retrieving the balance of an account.
// It validates the request, checks user authentication, and delegates
// the balance retrieval to the application layer.
func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if h.balance == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	if len(r.URL.Query()) > 0 {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	result, err := h.balance.Execute(r.Context(), accountapp.GetAccountBalanceInput{
		User:      user,
		AccountID: accountID,
	})
	if err != nil {
		log.Printf("event=get_balance error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, AccountBalanceData{
		AccountID: result.AccountID.String(),
		Balance:   result.Balance,
	})
}

// parseOptionalInt parses an optional integer query parameter. If the
// input string is empty, it returns 0 and no error. Otherwise, it attempts
// to parse the string as an integer and returns the result or an error if
// the parsing fails.
func parseOptionalInt(raw string) (int, error) {
	if raw == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, err
	}

	return value, nil
}

// parseOptionalTime parses an optional time query parameter in RFC3339 format.
// If the input string is empty, it returns nil and no error. Otherwise, it attempts
// to parse the string as a time.Time and returns the result or an error if
// the parsing fails.
func parseOptionalTime(raw string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

// parseOptionalUUID parses an optional UUID query parameter. If the input string
// is empty, it returns nil and no error. Otherwise, it attempts to parse the
// string as a uuid.UUID and returns the result or an error if the parsing fails.
func parseOptionalUUID(raw string) (*uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}

	parsed, err := uuid.Parse(raw)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}

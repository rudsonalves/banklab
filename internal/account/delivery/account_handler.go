package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/application"
	"github.com/seu-usuario/bank-api/internal/account/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	var req CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	input := application.CreateAccountInput{
		User:       user,
		CustomerID: customerID,
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

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	if h.deposit == nil {
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

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	account, err := h.deposit.Execute(r.Context(), application.DepositInput{
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

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if h.withdraw == nil {
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

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	account, err := h.withdraw.Execute(r.Context(), application.WithdrawInput{
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

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if h.transfer == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(sharederrors.ErrInvalidRequest))
		return
	}

	fromAccountID, err := uuid.Parse(req.FromAccountID)
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	toAccountID, err := uuid.Parse(req.ToAccountID)
	if err != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(domain.ErrInvalidData))
		return
	}

	result, err := h.transfer.Execute(r.Context(), application.TransferInput{
		User:          user,
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		log.Printf("event=transfer error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	sharedhttp.WriteJSON(w, http.StatusOK, TransferData{
		FromAccountID: result.FromAccountID.String(),
		ToAccountID:   result.ToAccountID.String(),
		Amount:        result.Amount,
		FromBalance:   result.FromBalance,
		ToBalance:     result.ToBalance,
	})
}

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

	result, err := h.statement.Execute(r.Context(), application.GetStatementInput{
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

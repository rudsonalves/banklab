package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/application"
	"github.com/seu-usuario/bank-api/internal/account/domain"
)

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "customer_id must be a valid UUID")
		return
	}

	input := application.CreateAccountInput{
		CustomerID: customerID,
	}

	account, err := h.createAccount.Execute(r.Context(), input)
	if err != nil {
		log.Printf("event=create_account error=%v", err)

		switch {
		case errors.Is(err, domain.ErrCustomerNotFound):
			writeError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", "customer not found")
			return
		case errors.Is(err, domain.ErrInvalidData):
			writeError(w, http.StatusBadRequest, "INVALID_DATA", "invalid data")
			return
		}

		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusCreated, response{
		Data: AccountData{
			ID:         account.ID.String(),
			CustomerID: account.CustomerID.String(),
			Number:     account.Number,
			Branch:     account.Branch,
			Balance:    account.Balance,
			Status:     string(account.Status),
		},
		Error: nil,
	})
}

func (h *Handler) Deposit(w http.ResponseWriter, r *http.Request) {
	if h.deposit == nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "account id must be a valid UUID")
		return
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	account, err := h.deposit.Execute(r.Context(), application.DepositInput{
		AccountID: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		log.Printf("event=deposit error=%v", err)

		switch {
		case errors.Is(err, domain.ErrInvalidAmount):
			writeError(w, http.StatusBadRequest, "INVALID_AMOUNT", "amount must be greater than zero")
			return
		case errors.Is(err, domain.ErrInvalidData):
			writeError(w, http.StatusBadRequest, "INVALID_DATA", "invalid data")
			return
		case errors.Is(err, domain.ErrAccountNotFound):
			writeError(w, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "account not found")
			return
		case errors.Is(err, domain.ErrAccountInactive):
			writeError(w, http.StatusUnprocessableEntity, "ACCOUNT_INACTIVE", "account is not active")
			return
		}

		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data: map[string]interface{}{
			"id":      account.ID.String(),
			"balance": account.Balance,
		},
		Error: nil,
	})
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if h.withdraw == nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "account id must be a valid UUID")
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	account, err := h.withdraw.Execute(r.Context(), application.WithdrawInput{
		AccountID: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		log.Printf("event=withdraw error=%v", err)

		switch {
		case errors.Is(err, domain.ErrInvalidAmount):
			writeError(w, http.StatusBadRequest, "INVALID_AMOUNT", "amount must be greater than zero")
			return
		case errors.Is(err, domain.ErrInvalidData):
			writeError(w, http.StatusBadRequest, "INVALID_DATA", "invalid data")
			return
		case errors.Is(err, domain.ErrAccountNotFound):
			writeError(w, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "account not found")
			return
		case errors.Is(err, domain.ErrInsufficientBalance):
			writeError(w, http.StatusUnprocessableEntity, "INSUFFICIENT_BALANCE", "insufficient balance")
			return
		case errors.Is(err, domain.ErrAccountInactive):
			writeError(w, http.StatusUnprocessableEntity, "ACCOUNT_INACTIVE", "account is not active")
			return
		}

		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data: map[string]interface{}{
			"id":      account.ID.String(),
			"balance": account.Balance,
		},
		Error: nil,
	})
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if h.transfer == nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "invalid request body")
		return
	}

	fromAccountID, err := uuid.Parse(req.FromAccountID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "from_account_id must be a valid UUID")
		return
	}

	toAccountID, err := uuid.Parse(req.ToAccountID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "to_account_id must be a valid UUID")
		return
	}

	result, err := h.transfer.Execute(r.Context(), application.TransferInput{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        req.Amount,
	})
	if err != nil {
		log.Printf("event=transfer error=%v", err)

		switch {
		case errors.Is(err, domain.ErrInvalidData):
			writeError(w, http.StatusBadRequest, "INVALID_DATA", "invalid data")
			return
		case errors.Is(err, domain.ErrInvalidAmount):
			writeError(w, http.StatusBadRequest, "INVALID_AMOUNT", "amount must be greater than zero")
			return
		case errors.Is(err, domain.ErrSameAccountTransfer):
			writeError(w, http.StatusBadRequest, "SAME_ACCOUNT_TRANSFER", "source and destination accounts must be different")
			return
		case errors.Is(err, domain.ErrAccountNotFound):
			writeError(w, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "account not found")
			return
		case errors.Is(err, domain.ErrInsufficientBalance):
			writeError(w, http.StatusUnprocessableEntity, "INSUFFICIENT_BALANCE", "insufficient balance")
			return
		case errors.Is(err, domain.ErrAccountInactive):
			writeError(w, http.StatusUnprocessableEntity, "ACCOUNT_INACTIVE", "account is not active")
			return
		}

		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data: TransferData{
			FromAccountID: result.FromAccountID.String(),
			ToAccountID:   result.ToAccountID.String(),
			Amount:        result.Amount,
			FromBalance:   result.FromBalance,
			ToBalance:     result.ToBalance,
		},
		Error: nil,
	})
}

func (h *Handler) Statement(w http.ResponseWriter, r *http.Request) {
	if h.statement == nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "account id must be a valid UUID")
		return
	}

	limit, err := parseOptionalInt(r.URL.Query().Get("limit"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "limit must be a valid integer")
		return
	}

	cursor, err := parseOptionalTime(r.URL.Query().Get("cursor"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "cursor must be a valid RFC3339 datetime")
		return
	}

	cursorID, err := parseOptionalUUID(r.URL.Query().Get("cursor_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "cursor_id must be a valid UUID")
		return
	}

	if (cursor == nil) != (cursorID == nil) {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "cursor and cursor_id must be provided together")
		return
	}

	from, err := parseOptionalTime(r.URL.Query().Get("from"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "from must be a valid RFC3339 datetime")
		return
	}

	to, err := parseOptionalTime(r.URL.Query().Get("to"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_DATA", "to must be a valid RFC3339 datetime")
		return
	}

	result, err := h.statement.Execute(r.Context(), application.GetStatementInput{
		AccountID: accountID,
		Limit:     limit,
		Cursor:    cursor,
		CursorID:  cursorID,
		From:      from,
		To:        to,
	})
	if err != nil {
		log.Printf("event=get_statement error=%v", err)

		switch {
		case errors.Is(err, domain.ErrInvalidData):
			writeError(w, http.StatusBadRequest, "INVALID_DATA", "invalid data")
			return
		case errors.Is(err, domain.ErrAccountNotFound):
			writeError(w, http.StatusNotFound, "ACCOUNT_NOT_FOUND", "account not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
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

	writeJSON(w, http.StatusOK, response{
		Data: StatementData{
			AccountID:  result.AccountID,
			Items:      items,
			NextCursor: nextCursor,
		},
		Error: nil,
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

package delivery

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/account/application"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
)

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		writeError(w, http.StatusUnauthorized, authErr)
		return
	}

	var req CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.ErrInvalidRequest)
		return
	}

	customerID, err := uuid.Parse(req.CustomerID)
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "customer_id",
		}))
		return
	}

	input := application.CreateAccountInput{
		User:       user,
		CustomerID: customerID,
	}

	account, err := h.createAccount.Execute(r.Context(), input)
	if err != nil {
		log.Printf("event=create_account error=%v", err)
		appErr, status := mapAccountError(err)
		writeError(w, status, appErr)
		return
	}

	writeSuccess(w, http.StatusCreated, AccountData{
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
		writeError(w, http.StatusInternalServerError, sharederrors.ErrInternal)
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		writeError(w, http.StatusUnauthorized, authErr)
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "account_id",
		}))
		return
	}

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.ErrInvalidRequest)
		return
	}

	account, err := h.deposit.Execute(r.Context(), application.DepositInput{
		User:      user,
		AccountID: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		log.Printf("event=deposit error=%v", err)
		appErr, status := mapAccountError(err)
		writeError(w, status, appErr)
		return
	}

	writeSuccess(w, http.StatusOK, map[string]interface{}{
		"id":      account.ID.String(),
		"balance": account.Balance,
	})
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	if h.withdraw == nil {
		writeError(w, http.StatusInternalServerError, sharederrors.ErrInternal)
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		writeError(w, http.StatusUnauthorized, authErr)
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "account_id",
		}))
		return
	}

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.ErrInvalidRequest)
		return
	}

	account, err := h.withdraw.Execute(r.Context(), application.WithdrawInput{
		User:      user,
		AccountID: accountID,
		Amount:    req.Amount,
	})
	if err != nil {
		log.Printf("event=withdraw error=%v", err)
		appErr, status := mapAccountError(err)
		writeError(w, status, appErr)
		return
	}

	writeSuccess(w, http.StatusOK, map[string]interface{}{
		"id":      account.ID.String(),
		"balance": account.Balance,
	})
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	if h.transfer == nil {
		writeError(w, http.StatusInternalServerError, sharederrors.ErrInternal)
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		writeError(w, http.StatusUnauthorized, authErr)
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.ErrInvalidRequest)
		return
	}

	fromAccountID, err := uuid.Parse(req.FromAccountID)
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "from_account_id",
		}))
		return
	}

	toAccountID, err := uuid.Parse(req.ToAccountID)
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "to_account_id",
		}))
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
		appErr, status := mapAccountError(err)
		writeError(w, status, appErr)
		return
	}

	writeSuccess(w, http.StatusOK, TransferData{
		FromAccountID: result.FromAccountID.String(),
		ToAccountID:   result.ToAccountID.String(),
		Amount:        result.Amount,
		FromBalance:   result.FromBalance,
		ToBalance:     result.ToBalance,
	})
}

func (h *Handler) Statement(w http.ResponseWriter, r *http.Request) {
	if h.statement == nil {
		writeError(w, http.StatusInternalServerError, sharederrors.ErrInternal)
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		writeError(w, http.StatusUnauthorized, authErr)
		return
	}

	accountIDRaw := r.PathValue("id")
	accountID, err := uuid.Parse(accountIDRaw)
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "account_id",
		}))
		return
	}

	limit, err := parseOptionalInt(r.URL.Query().Get("limit"))
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "limit",
		}))
		return
	}

	cursor, err := parseOptionalTime(r.URL.Query().Get("cursor"))
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "cursor",
		}))
		return
	}

	cursorID, err := parseOptionalUUID(r.URL.Query().Get("cursor_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "cursor_id",
		}))
		return
	}

	if (cursor == nil) != (cursorID == nil) {
		writeError(w, http.StatusBadRequest, sharederrors.NewError("INVALID_DATA", "Invalid data"))
		return
	}

	from, err := parseOptionalTime(r.URL.Query().Get("from"))
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "from",
		}))
		return
	}

	to, err := parseOptionalTime(r.URL.Query().Get("to"))
	if err != nil {
		writeError(w, http.StatusBadRequest, sharederrors.NewErrorWithDetails("INVALID_DATA", "Invalid data", map[string]interface{}{
			"field": "to",
		}))
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
		appErr, status := mapAccountError(err)
		writeError(w, status, appErr)
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

	writeSuccess(w, http.StatusOK, StatementData{
		AccountID:  result.AccountID,
		Items:      items,
		NextCursor: nextCursor,
	})
}

func mapAccountError(err error) (*sharederrors.AppError, int) {
	switch application.CategorizeError(err) {
	case application.ErrorCategoryForbidden:
		return sharederrors.ErrForbidden, http.StatusForbidden
	case application.ErrorCategoryInvalidData:
		return sharederrors.ErrInvalidData, http.StatusBadRequest
	case application.ErrorCategoryInvalidAmount:
		return sharederrors.NewError("INVALID_AMOUNT", "Invalid amount"), http.StatusBadRequest
	case application.ErrorCategoryCustomerNotFound:
		return sharederrors.NewError("CUSTOMER_NOT_FOUND", "Customer not found"), http.StatusNotFound
	case application.ErrorCategoryAccountNotFound:
		return sharederrors.NewError("ACCOUNT_NOT_FOUND", "Account not found"), http.StatusNotFound
	case application.ErrorCategoryAccountInactive:
		return sharederrors.NewError("ACCOUNT_INACTIVE", "Account is not active"), http.StatusUnprocessableEntity
	case application.ErrorCategoryInsufficientAmount:
		return sharederrors.NewError("INSUFFICIENT_BALANCE", "Insufficient balance"), http.StatusUnprocessableEntity
	case application.ErrorCategorySameAccount:
		return sharederrors.NewError("SAME_ACCOUNT_TRANSFER", "Source and destination accounts must be different"), http.StatusBadRequest
	default:
		return sharederrors.ErrInternal, http.StatusInternalServerError
	}
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

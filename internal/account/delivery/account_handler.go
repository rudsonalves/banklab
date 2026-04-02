package delivery

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

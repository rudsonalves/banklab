package delivery

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	accountapp "github.com/seu-usuario/bank-api/internal/account/bankaccount/application"
	"github.com/seu-usuario/bank-api/internal/account/bankaccount/domain"
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

// LookupInternalTransferRecipients handles the HTTP request for looking up
// eligible recipient accounts for internal transfers.
func (h *Handler) LookupInternalTransferRecipients(w http.ResponseWriter, r *http.Request) {
	if h.lookupInternalTransferRecipients == nil {
		sharedhttp.WriteError(w, sharederrors.MapError(nil))
		return
	}

	user, authErr := RequireUser(r.Context())
	if authErr != nil {
		sharedhttp.WriteError(w, sharederrors.MapError(authErr))
		return
	}

	query := r.URL.Query()
	recipients, err := h.lookupInternalTransferRecipients.Execute(
		r.Context(),
		accountapp.LookupInternalTransferRecipientsInput{
			User:          user,
			Branch:        query.Get("branch"),
			AccountNumber: query.Get("account_number"),
			Document:      query.Get("document"),
		},
	)
	if err != nil {
		log.Printf("event=lookup_internal_transfer_recipients error=%v", err)
		sharedhttp.WriteError(w, sharederrors.MapError(err))
		return
	}

	data := make([]InternalTransferRecipientData, 0, len(recipients))
	for _, recipient := range recipients {
		data = append(data, InternalTransferRecipientData{
			AccountID:     recipient.AccountID.String(),
			HolderName:    recipient.HolderName,
			Document:      recipient.MaskedDocument,
			Branch:        recipient.Branch,
			AccountNumber: recipient.AccountNumber,
			AccountType:   recipient.AccountType,
		})
	}

	sharedhttp.WriteJSON(w, http.StatusOK, InternalTransferRecipientsData{
		Accounts: data,
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

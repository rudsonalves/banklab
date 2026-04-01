package delivery

import "github.com/seu-usuario/bank-api/internal/account/application"

type Handler struct {
	createAccount *application.CreateAccount
}

func New(createAccount *application.CreateAccount) *Handler {
	return &Handler{
		createAccount: createAccount,
	}
}

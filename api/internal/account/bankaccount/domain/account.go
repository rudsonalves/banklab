package domain

import (
	"github.com/google/uuid"
	accountdomain "github.com/seu-usuario/bank-api/internal/account/domain"
)

type Account = accountdomain.Account
type AccountStatus = accountdomain.AccountStatus
type Transaction = accountdomain.Transaction
type TransactionType = accountdomain.TransactionType
type Tx = accountdomain.Tx

const (
	AccountActive   = accountdomain.AccountActive
	AccountInactive = accountdomain.AccountInactive
	AccountBlocked  = accountdomain.AccountBlocked
)

var (
	ErrInvalidData         = accountdomain.ErrInvalidData
	ErrInvalidAmount       = accountdomain.ErrInvalidAmount
	ErrAccountNotFound     = accountdomain.ErrAccountNotFound
	ErrInsufficientBalance = accountdomain.ErrInsufficientBalance
	ErrCustomerNotFound    = accountdomain.ErrCustomerNotFound
	ErrAccountInactive     = accountdomain.ErrAccountInactive
	ErrForbidden           = accountdomain.ErrForbidden
)

func NewAccount(customerID uuid.UUID, number, branch string) (*Account, error) {
	return accountdomain.NewAccount(customerID, number, branch)
}

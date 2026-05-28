package infrastructure

import (
	"github.com/seu-usuario/bank-api/internal/security/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptTransactionPasswordHasher struct {
	cost int
}

var _ domain.TransactionPasswordHasher = (*BcryptTransactionPasswordHasher)(nil)

func NewBcryptTransactionPasswordHasher(cost int) *BcryptTransactionPasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	return &BcryptTransactionPasswordHasher{cost: cost}
}

func (h *BcryptTransactionPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *BcryptTransactionPasswordHasher) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

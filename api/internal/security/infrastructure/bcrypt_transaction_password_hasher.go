package infrastructure

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/seu-usuario/bank-api/internal/security/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptTransactionPasswordHasher struct {
	cost   int
	pepper string
}

var _ domain.TransactionPasswordHasher = (*BcryptTransactionPasswordHasher)(nil)

func NewBcryptTransactionPasswordHasher(cost int, pepper string) *BcryptTransactionPasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	if strings.TrimSpace(pepper) == "" {
		panic("transaction password pepper is required")
	}

	return &BcryptTransactionPasswordHasher{cost: cost, pepper: pepper}
}

func (h *BcryptTransactionPasswordHasher) Hash(password string) (string, error) {
	peppered := h.pepperPassword(password)

	hash, err := bcrypt.GenerateFromPassword(peppered, h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *BcryptTransactionPasswordHasher) Compare(hash, password string) bool {
	peppered := h.pepperPassword(password)

	return bcrypt.CompareHashAndPassword([]byte(hash), peppered) == nil
}

func (h *BcryptTransactionPasswordHasher) pepperPassword(password string) []byte {
	mac := hmac.New(sha256.New, []byte(h.pepper))
	_, _ = mac.Write([]byte(password))

	encoded := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return []byte(encoded)
}

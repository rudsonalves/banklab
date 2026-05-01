package infrastructure

import (
	"github.com/seu-usuario/bank-api/internal/auth/domain"
	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct {
	cost int
}

var _ domain.PasswordHasher = (*BcryptPasswordHasher)(nil)

// NewBcryptPasswordHasher creates a new BcryptPasswordHasher with the specified cost.
// If the cost is outside the valid range, it defaults to bcrypt.DefaultCost.
func NewBcryptPasswordHasher(cost int) *BcryptPasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}

	return &BcryptPasswordHasher{cost: cost}
}

// Hash generates a bcrypt hash of the given password using the configured cost.
func (h *BcryptPasswordHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Compare checks if the provided password matches the given bcrypt hash.
func (h *BcryptPasswordHasher) Compare(hash string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

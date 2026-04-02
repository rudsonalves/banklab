package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) error
}

type TokenClaims struct {
	UserID     string
	Role       Role
	CustomerID *string
}

type TokenService interface {
	GenerateAccessToken(claims TokenClaims) (string, error)
	ParseAccessToken(token string) (*TokenClaims, error)
}

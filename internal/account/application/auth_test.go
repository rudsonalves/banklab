package application

import (
	"github.com/google/uuid"
	auth "github.com/seu-usuario/bank-api/internal/auth/domain"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

func testCustomerUser(customerID uuid.UUID) *auth.AuthenticatedUser {
	return &auth.AuthenticatedUser{
		UserID:     "user-1",
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	}
}

func testAdminUser() *auth.AuthenticatedUser {
	return &auth.AuthenticatedUser{
		UserID: "admin-1",
		Role:   authdomain.RoleAdmin,
	}
}

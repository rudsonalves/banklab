package delivery

import (
	"net/http"

	"github.com/google/uuid"
	auth "github.com/seu-usuario/bank-api/internal/auth/delivery"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
)

func testAuthenticatedRequest(req *http.Request, customerID uuid.UUID) *http.Request {
	ctx := auth.WithAuthenticatedUser(req.Context(), auth.AuthenticatedUser{
		UserID:     uuid.New(),
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	})

	return req.WithContext(ctx)
}

package delivery

import (
	"net/http"

	"github.com/google/uuid"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
)

func testAuthenticatedRequest(req *http.Request, customerID uuid.UUID) *http.Request {
	ctx := sharedauthctx.WithAuthenticatedUser(req.Context(), sharedauthctx.AuthenticatedUser{
		UserID:     uuid.New(),
		Role:       authdomain.RoleCustomer,
		CustomerID: &customerID,
	})

	return req.WithContext(ctx)
}

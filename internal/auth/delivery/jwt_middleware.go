package delivery

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/seu-usuario/bank-api/internal/auth/application"
	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type JWTMiddleware struct {
	tokenService authdomain.TokenService
}

func NewJWTMiddleware(tokenService authdomain.TokenService) *JWTMiddleware {
	return &JWTMiddleware{tokenService: tokenService}
}

func (m *JWTMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			sharedhttp.WriteError(w, http.StatusUnauthorized, sharederrors.ErrUnauthorized)
			return
		}

		claims, err := m.tokenService.ParseAccessToken(token)
		if err != nil {
			sharedhttp.WriteError(w, http.StatusUnauthorized, sharederrors.ErrInvalidToken)
			return
		}

		customerID, ok := parseNullableCustomerID(claims.CustomerID)
		if !ok {
			sharedhttp.WriteError(w, http.StatusUnauthorized, sharederrors.ErrInvalidToken)
			return
		}

		principal := authdomain.AuthenticatedUser{
			UserID:     claims.UserID,
			Role:       claims.Role,
			CustomerID: customerID,
		}

		ctx := application.WithAuthenticatedUser(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.TrimSpace(r.Header.Get("Authorization"))
		if authorization == "" {
			next.ServeHTTP(w, r)
			return
		}

		token, ok := bearerToken(authorization)
		if !ok {
			sharedhttp.WriteError(w, http.StatusUnauthorized, sharederrors.ErrUnauthorized)
			return
		}

		claims, err := m.tokenService.ParseAccessToken(token)
		if err != nil {
			sharedhttp.WriteError(w, http.StatusUnauthorized, sharederrors.ErrInvalidToken)
			return
		}

		customerID, ok := parseNullableCustomerID(claims.CustomerID)
		if !ok {
			sharedhttp.WriteError(w, http.StatusUnauthorized, sharederrors.ErrInvalidToken)
			return
		}

		principal := authdomain.AuthenticatedUser{
			UserID:     claims.UserID,
			Role:       claims.Role,
			CustomerID: customerID,
		}

		ctx := application.WithAuthenticatedUser(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func bearerToken(authorization string) (string, bool) {
	parts := strings.Split(strings.TrimSpace(authorization), " ")
	if len(parts) != 2 {
		return "", false
	}

	if parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}

	return parts[1], true
}

func parseNullableCustomerID(raw *string) (*uuid.UUID, bool) {
	if raw == nil {
		return nil, true
	}

	parsed, err := uuid.Parse(*raw)
	if err != nil {
		return nil, false
	}

	return &parsed, true
}

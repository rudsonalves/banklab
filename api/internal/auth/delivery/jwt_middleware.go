package delivery

import (
	"net/http"
	"strings"

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
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
			return
		}

		claims, err := m.tokenService.ParseAccessToken(token)
		if err != nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidToken))
			return
		}

		principal := authdomain.AuthenticatedUser{
			UserID:     claims.UserID,
			Role:       claims.Role,
			CustomerID: claims.CustomerID,
		}

		ctx := WithAuthenticatedUser(r.Context(), principal)
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
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
			return
		}

		claims, err := m.tokenService.ParseAccessToken(token)
		if err != nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidToken))
			return
		}

		principal := authdomain.AuthenticatedUser{
			UserID:     claims.UserID,
			Role:       claims.Role,
			CustomerID: claims.CustomerID,
		}

		ctx := WithAuthenticatedUser(r.Context(), principal)
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

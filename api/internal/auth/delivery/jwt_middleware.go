package delivery

import (
	"net/http"
	"strings"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

type JWTMiddleware struct {
	tokenService authdomain.TokenService
}

// NewJWTMiddleware creates a new instance of JWTMiddleware with the provided TokenService.
func NewJWTMiddleware(tokenService authdomain.TokenService) *JWTMiddleware {
	return &JWTMiddleware{tokenService: tokenService}
}

// RequireAuth is a middleware that ensures the request has a valid JWT token.
// If the token is valid, it extracts the user information and adds it to the request context.
// If the token is missing or invalid, it responds with an appropriate error.
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

		ctx := sharedauthctx.WithAuthenticatedUser(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth is a middleware that attempts to authenticate the request using a JWT token if present.
// If a valid token is provided, it extracts the user information and adds it to the request context.
// If no token is provided, it simply passes the request through without authentication.
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

		ctx := sharedauthctx.WithAuthenticatedUser(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// bearerToken extracts the token from the Authorization header if it follows the "
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

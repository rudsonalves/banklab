package delivery

import (
	"net/http"
	"strings"

	authdomain "github.com/seu-usuario/bank-api/internal/auth/domain"
	installationdomain "github.com/seu-usuario/bank-api/internal/installation/domain"
	sharedauthctx "github.com/seu-usuario/bank-api/internal/shared/authctx"
	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
	sharedheaders "github.com/seu-usuario/bank-api/internal/shared/http/headers"
)

type JWTMiddleware struct {
	tokenService          authdomain.TokenService
	restrictedAccessToken installationdomain.RestrictedAccessTokenVerifier
}

// NewJWTMiddleware creates a new instance of JWTMiddleware with the provided TokenService.
func NewJWTMiddleware(tokenService authdomain.TokenService) *JWTMiddleware {
	return &JWTMiddleware{tokenService: tokenService}
}

func (m *JWTMiddleware) WithRestrictedAccessTokenVerifier(
	verifier installationdomain.RestrictedAccessTokenVerifier,
) *JWTMiddleware {
	if verifier != nil {
		m.restrictedAccessToken = verifier
	}

	return m
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
		if claims == nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidToken))
			return
		}

		principal := principalFromClaims(claims)
		ctx := sharedauthctx.WithAuthenticatedUser(r.Context(), principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) RequireOperationalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerInstallationID, err := parseCanonicalInstallationID(r.Header.Get(sharedheaders.InstallationID))
		if err != nil {
			sharedhttp.WriteError(w, sharederrors.MapError(err))
			return
		}

		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
			return
		}

		claims, err := m.tokenService.ParseAccessToken(token)
		if err != nil || claims == nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidToken))
			return
		}
		if claims.InstallationID == nil || *claims.InstallationID != headerInstallationID {
			sharedhttp.WriteError(w, sharederrors.MapError(installationdomain.ErrInstallationMismatch))
			return
		}

		principal := principalFromClaims(claims)
		ctx := sharedauthctx.WithAuthenticatedUser(r.Context(), principal)
		ctx = sharedauthctx.WithOperationalSession(ctx, sharedauthctx.OperationalSession{
			UserID:         claims.UserID,
			Role:           claims.Role,
			CustomerID:     claims.CustomerID,
			InstallationID: claims.InstallationID,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) RequireOperationalOrRestrictedAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerInstallationID, err := parseCanonicalInstallationID(r.Header.Get(sharedheaders.InstallationID))
		if err != nil {
			sharedhttp.WriteError(w, sharederrors.MapError(err))
			return
		}

		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
			return
		}

		if claims, err := m.tokenService.ParseAccessToken(token); err == nil && claims != nil {
			if claims.InstallationID == nil || *claims.InstallationID != headerInstallationID {
				sharedhttp.WriteError(w, sharederrors.MapError(installationdomain.ErrInstallationMismatch))
				return
			}

			principal := principalFromClaims(claims)
			ctx := sharedauthctx.WithAuthenticatedUser(r.Context(), principal)
			ctx = sharedauthctx.WithOperationalSession(ctx, sharedauthctx.OperationalSession{
				UserID:         claims.UserID,
				Role:           claims.Role,
				CustomerID:     claims.CustomerID,
				InstallationID: claims.InstallationID,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		if m.restrictedAccessToken == nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidToken))
			return
		}
		claims, err := m.restrictedAccessToken.VerifyRestrictedAccessToken(r.Context(), token)
		if err != nil {
			sharedhttp.WriteError(w, sharederrors.MapError(err))
			return
		}
		if claims == nil || claims.InstallationID.UUID() != headerInstallationID {
			sharedhttp.WriteError(w, sharederrors.MapError(installationdomain.ErrInstallationMismatch))
			return
		}

		ctx := sharedauthctx.WithRestrictedSession(r.Context(), sharedauthctx.RestrictedSession{
			UserID:         claims.UserID,
			InstallationID: claims.InstallationID.UUID(),
			JTI:            claims.JTI,
			Scope:          claims.Scope,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *JWTMiddleware) RequireRestrictedAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headerInstallationID, err := parseCanonicalInstallationID(r.Header.Get(sharedheaders.InstallationID))
		if err != nil {
			sharedhttp.WriteError(w, sharederrors.MapError(err))
			return
		}
		if m.restrictedAccessToken == nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidToken))
			return
		}

		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrUnauthorized))
			return
		}

		claims, err := m.restrictedAccessToken.VerifyRestrictedAccessToken(r.Context(), token)
		if err != nil {
			sharedhttp.WriteError(w, sharederrors.MapError(err))
			return
		}
		if claims == nil || claims.InstallationID.UUID() != headerInstallationID {
			sharedhttp.WriteError(w, sharederrors.MapError(installationdomain.ErrInstallationMismatch))
			return
		}

		ctx := sharedauthctx.WithRestrictedSession(r.Context(), sharedauthctx.RestrictedSession{
			UserID:         claims.UserID,
			InstallationID: claims.InstallationID.UUID(),
			JTI:            claims.JTI,
			Scope:          claims.Scope,
		})
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

		if claims == nil {
			sharedhttp.WriteError(w, sharederrors.MapError(authdomain.ErrInvalidToken))
			return
		}

		ctx := sharedauthctx.WithAuthenticatedUser(r.Context(), principalFromClaims(claims))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func principalFromClaims(claims *authdomain.TokenClaims) authdomain.AuthenticatedUser {
	return authdomain.AuthenticatedUser{
		UserID:     claims.UserID,
		Role:       claims.Role,
		CustomerID: claims.CustomerID,
	}
}

// bearerToken extracts the token from the Authorization header if it follows
// the "Bearer <token>" format. It returns the token and a boolean indicating
// whether the extraction was successful. If the header does not follow the expected
// format, it returns an empty string and false.
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

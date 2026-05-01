package middleware

import (
	"crypto/subtle"
	"net/http"

	sharederrors "github.com/seu-usuario/bank-api/internal/shared/errors"
	sharedhttp "github.com/seu-usuario/bank-api/internal/shared/http"
)

const headerAppToken = "X-App-Token"

// AppToken is a middleware that checks for a specific application token in the request header.
// It compares the token in a constant time manner to prevent timing attacks.
// If the token is invalid, it responds with an error and does not call the next handler.
func AppToken(expectedToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if subtle.ConstantTimeCompare(
				[]byte(r.Header.Get(headerAppToken)),
				[]byte(expectedToken),
			) != 1 {
				sharedhttp.WriteError(w, sharederrors.ErrInvalidAppToken)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

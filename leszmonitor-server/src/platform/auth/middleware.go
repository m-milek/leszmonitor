package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/platform/log"
)

// RequireAuth wraps a single handler with JwtAuth.
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return JwtAuth(next).ServeHTTP
}

// JwtAuth middleware validates JWT tokens from the Authorization header.
func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpx.WrapResponseWriter(w)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.RespondError(r.Context(), rw, http.StatusUnauthorized, errors.New("no token provided"))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userClaims, err := ValidateJwt(tokenString)
		if err != nil {
			httpx.RespondError(r.Context(), rw, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		// Store the user claims in the request context
		ctx := SetUserInContext(r.Context(), userClaims)

		logger := log.FromContext(ctx).With().Str("username", userClaims.Username).Logger()
		ctx = log.WithContext(ctx, &logger)

		// Call the next handler with the updated context
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

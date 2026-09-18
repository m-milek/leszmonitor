package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/platform/log"
)

// JwtAuth middleware validates JWT tokens from the Authorization header.
func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := httpx.WrapResponseWriter(w)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(rw, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userClaims, err := ValidateJwt(tokenString)
		if err != nil {
			http.Error(rw, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Store the user claims in the request context
		ctx := SetUserInContext(r.Context(), userClaims)

		// Call the next handler with the updated context
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// SetUserContext stores user claims in the request context.
func SetUserContext(ctx context.Context, claims *UserClaims) context.Context {
	logger := log.FromContext(ctx)
	logger.Debug().Msg("Setting user claims in context: " + claims.Username)
	return context.WithValue(ctx, "userClaims", claims)
}

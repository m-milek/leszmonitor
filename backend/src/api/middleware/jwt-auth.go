package middleware

import (
	"context"
	"fmt"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/env"
	"net/http"
	"os"
	"strings"
)

// UserClaims extends standard JWT claims with custom fields
type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JwtAuth middleware validates JWT tokens from the Authorization header
func JwtAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := newResponseWriter(w)

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(rw, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Handle "Bearer <token>" format
		tokenString := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// Get JWT secret from environment
		jwtSecret := os.Getenv(env.JWT_SECRET)
		if jwtSecret == "" {
			http.Error(rw, "Server configuration error", http.StatusInternalServerError)
			return
		}

		// Parse and validate the token with custom claims
		claims := &UserClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Verify signing algorithm
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			http.Error(rw, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(rw, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Verify we have the expected claims type
		userClaims, ok := token.Claims.(*UserClaims)
		if !ok {
			http.Error(rw, "Unauthorized: Invalid claims format", http.StatusUnauthorized)
			return
		}

		// Ensure username is present
		if userClaims.Username == "" {
			http.Error(rw, "Unauthorized: Missing username in token", http.StatusUnauthorized)
			return
		}

		// Store the user claims in the request context
		ctx := SetUserContext(r.Context(), userClaims)

		// Call the next handler with the updated context
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}

// Define a context key type to avoid collisions
type contextKey string

const userClaimsKey contextKey = "userClaims"

// SetUserContext stores user claims in the request context
func SetUserContext(ctx context.Context, claims *UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

// GetUserFromContext retrieves user claims from the request context
func GetUserFromContext(ctx context.Context) (*UserClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*UserClaims)
	return claims, ok
}

package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/appconfig"
)

// JwtClaims represents the claims stored in a Leszmonitor JWT token.
// It includes standard claims and a custom Username field.
type JwtClaims struct {
	jwt.MapClaims
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

// UserClaims extends standard JWT claims with custom fields.
type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// jwtFromRequest extracts the JWT token from the Authorization header of the HTTP request.
func jwtFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	const prefix = "Bearer "
	if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
		return "", nil
	}

	return authHeader[len(prefix):], nil
}

func decodeJwtClaims(jwtString string) (JwtClaims, error) {
	claims := JwtClaims{}
	token, err := jwt.ParseWithClaims(jwtString, &claims, func(_ *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(config.JwtSecret)), nil
	})

	if err != nil {
		return JwtClaims{}, errors.Join(fmt.Errorf("failed to parse JWT token: %w", err))
	}

	if !token.Valid {
		return JwtClaims{}, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}

func ValidateJwt(token string) (*UserClaims, error) {
	jwtSecret := os.Getenv(config.JwtSecret)
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT secret is not configured")
	}

	claims := &UserClaims{}

	parsedJwt, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid JWT token: %w", err)
	}

	if !parsedJwt.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	userClaims, ok := parsedJwt.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("unexpected JWT claims type")
	}

	if userClaims.Username == "" {
		return nil, fmt.Errorf("JWT token is missing username claim")
	}

	return userClaims, nil
}

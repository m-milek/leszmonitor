package api

import (
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logging"
	"net/http"
	"os"
)

// JwtClaims represents the claims stored in a Leszmonitor JWT token.
// It includes standard claims and a custom Username field.
type JwtClaims struct {
	jwt2.MapClaims
	Username string `json:"username"`
	Exp      int64  `json:"exp"`
}

// JwtFromRequest extracts the JWT token from the Authorization header of the HTTP request.
func JwtFromRequest(r *http.Request) (string, error) {
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

func DecodeJwtClaims(jwtString string) (JwtClaims, error) {
	claims := JwtClaims{}
	token, err := jwt2.ParseWithClaims(jwtString, &claims, func(token *jwt2.Token) (interface{}, error) {
		return []byte(os.Getenv(env.JwtSecret)), nil
	})

	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to parse JWT token")
		return JwtClaims{}, err
	}

	if !token.Valid {
		logging.Api.Warn().Msg("Invalid JWT token")
		return JwtClaims{}, fmt.Errorf("invalid JWT token")
	}

	return claims, nil
}
